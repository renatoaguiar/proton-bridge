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

// +build !darwin

package versioner

import (
	"os"

	"github.com/Masterminds/semver/v3"
	"github.com/sirupsen/logrus"
)

// RemoveOldVersions removes all but the latest app version.
func (v *Versioner) RemoveOldVersions() error {
	versions, err := v.ListVersions()
	if err != nil {
		return err
	}

	// darwin does not currently use the versioner.
	if len(versions) == 0 {
		return nil
	}

	for _, version := range versions[1:] {
		if err := os.RemoveAll(version.path); err != nil {
			logrus.WithError(err).Error("Failed to remove old app version")
		}
	}

	return nil
}

// RemoveOtherVersions removes all but the specific provided app version.
func (v *Versioner) RemoveOtherVersions(versionToKeep *semver.Version) error {
	versions, err := v.ListVersions()
	if err != nil {
		return err
	}

	for _, version := range versions {
		if version.Equal(versionToKeep) {
			continue
		}
		if err := os.RemoveAll(version.path); err != nil {
			logrus.WithError(err).Error("Failed to remove old app version")
		}
	}

	return nil
}
