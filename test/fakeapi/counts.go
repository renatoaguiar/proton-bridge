// Copyright (c) 2021 Proton Technologies AG
//
// This file is part of ProtonMail Bridge.Bridge.
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

package fakeapi

import "github.com/ProtonMail/proton-bridge/pkg/pmapi"

func (api *FakePMAPI) CountMessages(addressID string) ([]*pmapi.MessagesCount, error) {
	if err := api.checkAndRecordCall(GET, "/mail/v4/messages/count?AddressID="+addressID, nil); err != nil {
		return nil, err
	}
	return api.getCounts(addressID), nil
}

func (api *FakePMAPI) getAllCounts() []*pmapi.MessagesCount {
	return api.getCounts("")
}

func (api *FakePMAPI) getCounts(addressID string) []*pmapi.MessagesCount {
	allCounts := map[string]*pmapi.MessagesCount{}
	for _, message := range api.messages {
		if addressID != "" && message.AddressID != addressID {
			continue
		}
		for _, labelID := range message.LabelIDs {
			if counts, ok := allCounts[labelID]; ok {
				counts.Total++
				if message.Unread == 1 {
					counts.Unread++
				}
			} else {
				allCounts[labelID] = &pmapi.MessagesCount{
					LabelID: labelID,
					Total:   1,
					Unread:  message.Unread,
				}
			}
		}
	}

	res := []*pmapi.MessagesCount{}
	for _, counts := range allCounts {
		res = append(res, counts)
	}
	return res
}
