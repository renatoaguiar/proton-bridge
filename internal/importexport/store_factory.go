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

package importexport

import (
	"github.com/ProtonMail/proton-bridge/internal/store"
)

// storeFactory implements dummy factory creating no store (not needed by Import-Export).
type storeFactory struct{}

// New does nothing.
func (f *storeFactory) New(user store.BridgeUser) (*store.Store, error) {
	return nil, nil
}

// Remove does nothing.
func (f *storeFactory) Remove(userID string) error {
	return nil
}
