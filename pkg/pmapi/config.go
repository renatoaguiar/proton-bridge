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

package pmapi

import (
	"net/http"
)

// rootURL is the API root URL.
//
// This can be changed using build flags: pmapi_local for "localhost/api", pmapi_dev or pmapi_prod.
// Default is pmapi_prod.
//
// It must not contain the protocol! The protocol should be in rootScheme.
var rootURL = "api.protonmail.ch" //nolint[gochecknoglobals]
var rootScheme = "https"          //nolint[gochecknoglobals]

// The HTTP transport to use by default.
var defaultTransport = &http.Transport{ //nolint[gochecknoglobals]
	Proxy: http.ProxyFromEnvironment,
}

// checkTLSCerts controls whether TLS certs are checked against known fingerprints.
// The default is for this to always be done.
var checkTLSCerts = true //nolint[gochecknoglobals]
