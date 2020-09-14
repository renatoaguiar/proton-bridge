// Copyright (c) 2020 Proton Technologies AG
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

package context

import (
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/ProtonMail/proton-bridge/pkg/config"
	"github.com/ProtonMail/proton-bridge/pkg/constants"
	"github.com/ProtonMail/proton-bridge/pkg/pmapi"
	"github.com/sirupsen/logrus"
)

type fakeConfig struct {
	dir string
}

// newFakeConfig creates a temporary folder for files.
// It's expected the test calls `ClearData` before finish to remove it from the file system.
func newFakeConfig() *fakeConfig {
	dir, err := ioutil.TempDir("", "example")
	if err != nil {
		panic(err)
	}

	cfg := &fakeConfig{
		dir: dir,
	}

	// We must generate cert.pem and key.pem to prevent errors when attempting to open them.
	if _, err = config.GenerateTLSConfig(cfg.GetTLSCertPath(), cfg.GetTLSKeyPath()); err != nil {
		logrus.WithError(err).Fatal()
	}

	return cfg
}

func (c *fakeConfig) ClearData() error {
	return os.RemoveAll(c.dir)
}
func (c *fakeConfig) GetAPIConfig() *pmapi.ClientConfig {
	return &pmapi.ClientConfig{
		AppVersion: "Bridge_" + constants.Version,
		ClientID:   "bridge",
	}
}
func (c *fakeConfig) GetDBDir() string {
	return c.dir
}
func (c *fakeConfig) GetVersion() string {
	return constants.Version
}
func (c *fakeConfig) GetLogDir() string {
	return c.dir
}
func (c *fakeConfig) GetLogPrefix() string {
	return "test"
}
func (c *fakeConfig) GetPreferencesPath() string {
	return filepath.Join(c.dir, "prefs.json")
}
func (c *fakeConfig) GetTransferDir() string {
	return c.dir
}
func (c *fakeConfig) GetTLSCertPath() string {
	return filepath.Join(c.dir, "cert.pem")
}
func (c *fakeConfig) GetTLSKeyPath() string {
	return filepath.Join(c.dir, "key.pem")
}
func (c *fakeConfig) GetEventsPath() string {
	return filepath.Join(c.dir, "events.json")
}
func (c *fakeConfig) GetIMAPCachePath() string {
	return filepath.Join(c.dir, "user_info.json")
}
func (c *fakeConfig) GetDefaultAPIPort() int {
	return 21042
}
func (c *fakeConfig) GetDefaultIMAPPort() int {
	return 21100 + rand.Intn(100)
}
func (c *fakeConfig) GetDefaultSMTPPort() int {
	return 21200 + rand.Intn(100)
}
