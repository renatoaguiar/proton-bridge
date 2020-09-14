// Copyright (c) 2020 Proton Technologies AG
//
// This file is part of ProtonMail Bridge.
//
// ProtonMail Bridge is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// ProtonMail Bridge is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with ProtonMail Bridge.  If not, see <https://www.gnu.org/licenses/>.

package transfer

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	pkgMessage "github.com/ProtonMail/proton-bridge/pkg/message"
	"github.com/ProtonMail/proton-bridge/pkg/pmapi"
	"github.com/ProtonMail/proton-bridge/pkg/sentry"
	"github.com/pkg/errors"
)

const (
	pmapiImportBatchMaxItems = 10
	pmapiImportBatchMaxSize  = 25 * 1000 * 1000 // 25 MB
)

// DefaultMailboxes returns the default mailboxes for default rules if no other is found.
func (p *PMAPIProvider) DefaultMailboxes(_ Mailbox) []Mailbox {
	return []Mailbox{{
		ID:          pmapi.ArchiveLabel,
		Name:        "Archive",
		IsExclusive: true,
	}}
}

// CreateMailbox creates label in ProtonMail account.
func (p *PMAPIProvider) CreateMailbox(mailbox Mailbox) (Mailbox, error) {
	if mailbox.ID != "" {
		return Mailbox{}, errors.New("mailbox is already created")
	}

	exclusive := 0
	if mailbox.IsExclusive {
		exclusive = 1
	}

	label, err := p.client().CreateLabel(&pmapi.Label{
		Name:      mailbox.Name,
		Color:     mailbox.Color,
		Exclusive: exclusive,
		Type:      pmapi.LabelTypeMailbox,
	})
	if err != nil {
		return Mailbox{}, errors.Wrap(err, fmt.Sprintf("failed to create mailbox %s", mailbox.Name))
	}
	mailbox.ID = label.ID
	return mailbox, nil
}

// TransferFrom imports messages from channel.
func (p *PMAPIProvider) TransferFrom(rules transferRules, progress *Progress, ch <-chan Message) {
	log.Info("Started transfer from channel to PMAPI")
	defer log.Info("Finished transfer from channel to PMAPI")

	// Cache has to be cleared before each transfer to not contain
	// old stuff from previous cancelled run.
	p.importMsgReqMap = map[string]*pmapi.ImportMsgReq{}
	p.importMsgReqSize = 0

	for msg := range ch {
		if progress.shouldStop() {
			break
		}

		if p.isMessageDraft(msg) {
			p.transferDraft(rules, progress, msg)
		} else {
			p.transferMessage(rules, progress, msg)
		}
	}

	if len(p.importMsgReqMap) > 0 {
		p.importMessages(progress)
	}
}

func (p *PMAPIProvider) isMessageDraft(msg Message) bool {
	for _, target := range msg.Targets {
		if target.ID == pmapi.DraftLabel {
			return true
		}
	}
	return false
}

func (p *PMAPIProvider) transferDraft(rules transferRules, progress *Progress, msg Message) {
	importedID, err := p.importDraft(msg, rules.globalMailbox)
	progress.messageImported(msg.ID, importedID, err)
}

func (p *PMAPIProvider) importDraft(msg Message, globalMailbox *Mailbox) (string, error) {
	message, attachmentReaders, err := p.parseMessage(msg)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse message")
	}

	if err := message.Encrypt(p.keyRing, nil); err != nil {
		return "", errors.Wrap(err, "failed to encrypt draft")
	}

	if globalMailbox != nil {
		message.LabelIDs = append(message.LabelIDs, globalMailbox.ID)
	}

	attachments := message.Attachments
	message.Attachments = nil

	draft, err := p.createDraft(message, "", pmapi.DraftActionReply)
	if err != nil {
		return "", errors.Wrap(err, "failed to create draft")
	}

	for idx, attachment := range attachments {
		attachment.MessageID = draft.ID
		attachmentBody, _ := ioutil.ReadAll(attachmentReaders[idx])

		r := bytes.NewReader(attachmentBody)
		sigReader, err := attachment.DetachedSign(p.keyRing, r)
		if err != nil {
			return "", errors.Wrap(err, "failed to sign attachment")
		}

		r = bytes.NewReader(attachmentBody)
		encReader, err := attachment.Encrypt(p.keyRing, r)
		if err != nil {
			return "", errors.Wrap(err, "failed to encrypt attachment")
		}

		_, err = p.createAttachment(attachment, encReader, sigReader)
		if err != nil {
			return "", errors.Wrap(err, "failed to create attachment")
		}
	}

	return draft.ID, nil
}

func (p *PMAPIProvider) transferMessage(rules transferRules, progress *Progress, msg Message) {
	importMsgReq, err := p.generateImportMsgReq(msg, rules.globalMailbox)
	if err != nil {
		progress.messageImported(msg.ID, "", err)
		return
	}

	importMsgReqSize := len(importMsgReq.Body)
	if p.importMsgReqSize+importMsgReqSize > pmapiImportBatchMaxSize || len(p.importMsgReqMap) == pmapiImportBatchMaxItems {
		p.importMessages(progress)
	}
	p.importMsgReqMap[msg.ID] = importMsgReq
	p.importMsgReqSize += importMsgReqSize
}

func (p *PMAPIProvider) generateImportMsgReq(msg Message, globalMailbox *Mailbox) (*pmapi.ImportMsgReq, error) {
	message, attachmentReaders, err := p.parseMessage(msg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse message")
	}

	body, err := p.encryptMessage(message, attachmentReaders)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encrypt message")
	}

	unread := 0
	if msg.Unread {
		unread = 1
	}

	labelIDs := []string{}
	for _, target := range msg.Targets {
		// Frontend should not set All Mail to Rules, but to be sure...
		if target.ID != pmapi.AllMailLabel {
			labelIDs = append(labelIDs, target.ID)
		}
	}
	if globalMailbox != nil {
		labelIDs = append(labelIDs, globalMailbox.ID)
	}

	return &pmapi.ImportMsgReq{
		AddressID: p.addressID,
		Body:      body,
		Unread:    unread,
		Time:      message.Time,
		Flags:     computeMessageFlags(labelIDs),
		LabelIDs:  labelIDs,
	}, nil
}

func (p *PMAPIProvider) parseMessage(msg Message) (m *pmapi.Message, r []io.Reader, err error) {
	// Old message parser is panicking in some cases.
	// Instead of crashing we try to convert to regular error.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic while parse: %v", r)
			if sentryErr := sentry.ReportSentryCrash(
				p.clientConfig.ClientID,
				p.clientConfig.AppVersion,
				p.clientConfig.UserAgent,
				err,
			); sentryErr != nil {
				log.Error("Sentry crash report failed: ", sentryErr)
			}
		}
	}()
	message, _, _, attachmentReaders, err := pkgMessage.Parse(bytes.NewBuffer(msg.Body), "", "")
	return message, attachmentReaders, err
}

func (p *PMAPIProvider) encryptMessage(msg *pmapi.Message, attachmentReaders []io.Reader) ([]byte, error) {
	if msg.MIMEType == pmapi.ContentTypeMultipartEncrypted {
		return []byte(msg.Body), nil
	}
	return pkgMessage.BuildEncrypted(msg, attachmentReaders, p.keyRing)
}

func computeMessageFlags(labels []string) (flag int64) {
	for _, labelID := range labels {
		switch labelID {
		case pmapi.SentLabel:
			flag = (flag | pmapi.FlagSent)
		case pmapi.ArchiveLabel, pmapi.InboxLabel:
			flag = (flag | pmapi.FlagReceived)
		case pmapi.DraftLabel:
			log.Error("Found draft target in non-draft import")
		}
	}

	// NOTE: if the labels are custom only
	if flag == 0 {
		flag = pmapi.FlagReceived
	}

	return flag
}

func (p *PMAPIProvider) importMessages(progress *Progress) {
	if progress.shouldStop() {
		return
	}

	importMsgIDs := []string{}
	importMsgRequests := []*pmapi.ImportMsgReq{}
	for msgID, req := range p.importMsgReqMap {
		importMsgIDs = append(importMsgIDs, msgID)
		importMsgRequests = append(importMsgRequests, req)
	}

	log.WithField("msgIDs", importMsgIDs).WithField("size", p.importMsgReqSize).Debug("Importing messages")
	results, err := p.importRequest(importMsgRequests)

	// In case the whole request failed, try to import every message one by one.
	if err != nil || len(results) == 0 {
		log.WithError(err).Warning("Importing messages failed, trying one by one")
		for msgID, req := range p.importMsgReqMap {
			importedID, err := p.importMessage(progress, req)
			progress.messageImported(msgID, importedID, err)
		}
		return
	}

	// In case request passed but some messages failed, try to import the failed ones alone.
	for index, result := range results {
		msgID := importMsgIDs[index]
		if result.Error != nil {
			log.WithError(result.Error).WithField("msg", msgID).Warning("Importing message failed, trying alone")
			req := importMsgRequests[index]
			importedID, err := p.importMessage(progress, req)
			progress.messageImported(msgID, importedID, err)
		} else {
			progress.messageImported(msgID, result.MessageID, nil)
		}
	}

	p.importMsgReqMap = map[string]*pmapi.ImportMsgReq{}
	p.importMsgReqSize = 0
}

func (p *PMAPIProvider) importMessage(progress *Progress, req *pmapi.ImportMsgReq) (importedID string, importedErr error) {
	progress.callWrap(func() error {
		results, err := p.importRequest([]*pmapi.ImportMsgReq{req})
		if err != nil {
			return errors.Wrap(err, "failed to import messages")
		}
		if len(results) == 0 {
			importedErr = errors.New("import ended with no result")
			return nil // This should not happen, only when there is bug which means we should skip this one.
		}
		if results[0].Error != nil {
			importedErr = errors.Wrap(results[0].Error, "failed to import message")
			return nil // Call passed but API refused this message, skip this one.
		}
		importedID = results[0].MessageID
		return nil
	})
	return
}
