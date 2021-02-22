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

package tls

import "os/exec"

func addTrustedCert(certPath string) error {
	return exec.Command( // nolint[gosec]
		"/usr/bin/security",
		"execute-with-privileges",
		"/usr/bin/security",
		"add-trusted-cert",
		"-d",
		"-r", "trustRoot",
		"-p", "ssl",
		"-k", "/Library/Keychains/System.keychain",
		certPath,
	).Run()
}

func removeTrustedCert(certPath string) error {
	return exec.Command( // nolint[gosec]
		"/usr/bin/security",
		"execute-with-privileges",
		"/usr/bin/security",
		"remove-trusted-cert",
		"-d",
		certPath,
	).Run()
}

func (t *TLS) InstallCerts() error {
	return addTrustedCert(t.getTLSCertPath())
}

func (t *TLS) UninstallCerts() error {
	return removeTrustedCert(t.getTLSCertPath())
}
