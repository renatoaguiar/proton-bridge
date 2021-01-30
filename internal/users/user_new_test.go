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

package users

import (
	"errors"
	"testing"

	"github.com/ProtonMail/proton-bridge/internal/events"
	"github.com/ProtonMail/proton-bridge/internal/users/credentials"
	"github.com/ProtonMail/proton-bridge/pkg/pmapi"
	gomock "github.com/golang/mock/gomock"
	a "github.com/stretchr/testify/assert"
)

func TestNewUserNoCredentialsStore(t *testing.T) {
	m := initMocks(t)
	defer m.ctrl.Finish()

	m.credentialsStore.EXPECT().Get("user").Return(nil, errors.New("fail"))

	_, err := newUser(m.PanicHandler, "user", m.eventListener, m.credentialsStore, m.clientManager, m.storeMaker)
	a.Error(t, err)
}

func TestNewUserAppOutdated(t *testing.T) {
	m := initMocks(t)
	defer m.ctrl.Finish()

	m.clientManager.EXPECT().GetClient("user").Return(m.pmapiClient).MinTimes(1)

	gomock.InOrder(
		m.credentialsStore.EXPECT().Get("user").Return(testCredentials, nil),
		m.credentialsStore.EXPECT().Get("user").Return(testCredentials, nil),
		m.pmapiClient.EXPECT().AuthRefresh("token").Return(nil, pmapi.ErrUpgradeApplication),
		m.eventListener.EXPECT().Emit(events.UpgradeApplicationEvent, ""),
		m.pmapiClient.EXPECT().ListLabels().Return(nil, pmapi.ErrUpgradeApplication),
		m.pmapiClient.EXPECT().Addresses().Return(nil),
	)

	checkNewUserHasCredentials(testCredentials, m)
}

func TestNewUserNoInternetConnection(t *testing.T) {
	m := initMocks(t)
	defer m.ctrl.Finish()

	m.clientManager.EXPECT().GetClient("user").Return(m.pmapiClient).MinTimes(1)

	gomock.InOrder(
		m.credentialsStore.EXPECT().Get("user").Return(testCredentials, nil),
		m.credentialsStore.EXPECT().Get("user").Return(testCredentials, nil),
		m.pmapiClient.EXPECT().AuthRefresh("token").Return(nil, pmapi.ErrAPINotReachable),
		m.eventListener.EXPECT().Emit(events.InternetOffEvent, ""),

		m.pmapiClient.EXPECT().ListLabels().Return(nil, pmapi.ErrAPINotReachable),
		m.pmapiClient.EXPECT().Addresses().Return(nil),
		m.pmapiClient.EXPECT().GetEvent("").Return(nil, pmapi.ErrAPINotReachable).AnyTimes(),
	)

	checkNewUserHasCredentials(testCredentials, m)
}

func TestNewUserAuthRefreshFails(t *testing.T) {
	m := initMocks(t)
	defer m.ctrl.Finish()

	m.clientManager.EXPECT().GetClient("user").Return(m.pmapiClient).MinTimes(1)
	m.eventListener.EXPECT().Emit(events.LogoutEvent, "user")
	m.eventListener.EXPECT().Emit(events.UserRefreshEvent, "user")
	m.eventListener.EXPECT().Emit(events.CloseConnectionEvent, "user@pm.me")

	gomock.InOrder(
		m.credentialsStore.EXPECT().Get("user").Return(testCredentials, nil),
		m.credentialsStore.EXPECT().Get("user").Return(testCredentials, nil),
		m.pmapiClient.EXPECT().AuthRefresh("token").Return(nil, errors.New("bad token")),
		m.credentialsStore.EXPECT().Logout("user").Return(nil),

		m.pmapiClient.EXPECT().Logout(),
		m.credentialsStore.EXPECT().Logout("user").Return(nil),
		m.credentialsStore.EXPECT().Get("user").Return(testCredentialsDisconnected, nil),
	)

	checkNewUserHasCredentials(testCredentialsDisconnected, m)
}

func TestNewUserUnlockFails(t *testing.T) {
	m := initMocks(t)
	defer m.ctrl.Finish()

	m.clientManager.EXPECT().GetClient("user").Return(m.pmapiClient).MinTimes(1)

	m.eventListener.EXPECT().Emit(events.LogoutEvent, "user")
	m.eventListener.EXPECT().Emit(events.UserRefreshEvent, "user")
	m.eventListener.EXPECT().Emit(events.CloseConnectionEvent, "user@pm.me")

	gomock.InOrder(
		m.credentialsStore.EXPECT().Get("user").Return(testCredentials, nil),
		m.credentialsStore.EXPECT().Get("user").Return(testCredentials, nil),
		m.pmapiClient.EXPECT().AuthRefresh("token").Return(testAuthRefresh, nil),

		m.pmapiClient.EXPECT().Unlock([]byte("pass")).Return(errors.New("bad password")),
		m.credentialsStore.EXPECT().Logout("user").Return(nil),
		m.pmapiClient.EXPECT().Logout(),
		m.credentialsStore.EXPECT().Logout("user").Return(nil),
		m.credentialsStore.EXPECT().Get("user").Return(testCredentialsDisconnected, nil),
	)

	checkNewUserHasCredentials(testCredentialsDisconnected, m)
}

func TestNewUser(t *testing.T) {
	m := initMocks(t)
	defer m.ctrl.Finish()

	m.clientManager.EXPECT().GetClient("user").Return(m.pmapiClient).MinTimes(1)
	mockConnectedUser(m)
	mockEventLoopNoAction(m)

	checkNewUserHasCredentials(testCredentials, m)
}

func checkNewUserHasCredentials(creds *credentials.Credentials, m mocks) {
	user, _ := newUser(m.PanicHandler, "user", m.eventListener, m.credentialsStore, m.clientManager, m.storeMaker)
	defer cleanUpUserData(user)

	_ = user.init()

	waitForEvents()

	a.Equal(m.t, creds, user.creds)
}

func _TestUserEventRefreshUpdatesAddresses(t *testing.T) { // nolint[funlen]
	a.Fail(t, "not implemented")
}
