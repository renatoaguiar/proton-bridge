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

// Package constants contains variables that are set via ldflags during build.
package constants

// nolint[gochecknoglobals]
var (
	// Version of the build.
	Version = ""

	// Revision is current hash of the build.
	Revision = ""

	// BuildTime stamp of the build.
	BuildTime = ""

	// DSNSentry client keys to be able to report crashes to Sentry.
	DSNSentry = ""

	// LongVersion is derived from Version and Revision.
	LongVersion = Version + " (" + Revision + ")"

	// BuildVersion is derived from LongVersion and BuildTime.
	BuildVersion = LongVersion + " " + BuildTime
)
