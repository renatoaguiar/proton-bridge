// Copyright (c) 2021 Proton Technologies AG
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

package imap

import "github.com/sirupsen/logrus"

const (
	fetchMessagesWorkers    = 5 // In how many workers to fetch message (group list on IMAP).
	fetchAttachmentsWorkers = 5 // In how many workers to fetch attachments (for one message).
)

var (
	log = logrus.WithField("pkg", "imap") //nolint[gochecknoglobals]
)
