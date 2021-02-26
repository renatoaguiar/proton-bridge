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

// Package frontend provides all interfaces of the Bridge.
package frontend

import (
	"github.com/ProtonMail/go-autostart"
	"github.com/ProtonMail/proton-bridge/internal/bridge"
	"github.com/ProtonMail/proton-bridge/internal/config/settings"
	"github.com/ProtonMail/proton-bridge/internal/config/useragent"
	"github.com/ProtonMail/proton-bridge/internal/frontend/cli"
	cliie "github.com/ProtonMail/proton-bridge/internal/frontend/cli-ie"
	"github.com/ProtonMail/proton-bridge/internal/frontend/qt"
	qtie "github.com/ProtonMail/proton-bridge/internal/frontend/qt-ie"
	"github.com/ProtonMail/proton-bridge/internal/frontend/types"
	"github.com/ProtonMail/proton-bridge/internal/importexport"
	"github.com/ProtonMail/proton-bridge/internal/locations"
	"github.com/ProtonMail/proton-bridge/internal/updater"
	"github.com/ProtonMail/proton-bridge/pkg/listener"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.WithField("pkg", "frontend") // nolint[unused]
)

// Frontend is an interface to be implemented by each frontend type (cli, gui, html).
type Frontend interface {
	Loop() error
	NotifyManualUpdate(update updater.VersionInfo, canInstall bool)
	SetVersion(update updater.VersionInfo)
	NotifySilentUpdateInstalled()
	NotifySilentUpdateError(error)
	WaitUntilFrontendIsReady()
}

// New returns initialized frontend based on `frontendType`, which can be `cli` or `qt`.
func New(
	version,
	buildVersion,
	programName,
	frontendType string,
	showWindowOnStart bool,
	panicHandler types.PanicHandler,
	locations *locations.Locations,
	settings *settings.Settings,
	eventListener listener.Listener,
	updater types.Updater,
	userAgent *useragent.UserAgent,
	bridge *bridge.Bridge,
	noEncConfirmator types.NoEncConfirmator,
	autostart *autostart.App,
	restarter types.Restarter,
) Frontend {
	bridgeWrap := types.NewBridgeWrap(bridge)
	return newBridgeFrontend(
		version,
		buildVersion,
		programName,
		frontendType,
		showWindowOnStart,
		panicHandler,
		locations,
		settings,
		eventListener,
		updater,
		userAgent,
		bridgeWrap,
		noEncConfirmator,
		autostart,
		restarter,
	)
}

func newBridgeFrontend(
	version,
	buildVersion,
	programName,
	frontendType string,
	showWindowOnStart bool,
	panicHandler types.PanicHandler,
	locations *locations.Locations,
	settings *settings.Settings,
	eventListener listener.Listener,
	updater types.Updater,
	userAgent *useragent.UserAgent,
	bridge types.Bridger,
	noEncConfirmator types.NoEncConfirmator,
	autostart *autostart.App,
	restarter types.Restarter,
) Frontend {
	switch frontendType {
	case "cli":
		return cli.New(
			panicHandler,
			locations,
			settings,
			eventListener,
			updater,
			bridge,
			restarter,
		)
	default:
		return qt.New(
			version,
			buildVersion,
			programName,
			showWindowOnStart,
			panicHandler,
			locations,
			settings,
			eventListener,
			updater,
			userAgent,
			bridge,
			noEncConfirmator,
			autostart,
			restarter,
		)
	}
}

// NewImportExport returns initialized frontend based on `frontendType`, which can be `cli` or `qt`.
func NewImportExport(
	version,
	buildVersion,
	programName,
	frontendType string,
	panicHandler types.PanicHandler,
	locations *locations.Locations,
	settings *settings.Settings,
	eventListener listener.Listener,
	updater types.Updater,
	ie *importexport.ImportExport,
	restarter types.Restarter,
) Frontend {
	ieWrap := types.NewImportExportWrap(ie)
	return newIEFrontend(
		version,
		buildVersion,
		programName,
		frontendType,
		panicHandler,
		locations,
		settings,
		eventListener,
		updater,
		ieWrap,
		restarter,
	)
}

func newIEFrontend(
	version,
	buildVersion,
	programName,
	frontendType string,
	panicHandler types.PanicHandler,
	locations *locations.Locations,
	settings *settings.Settings,
	eventListener listener.Listener,
	updater types.Updater,
	ie types.ImportExporter,
	restarter types.Restarter,
) Frontend {
	switch frontendType {
	case "cli":
		return cliie.New(
			panicHandler,
			locations,
			eventListener,
			updater,
			ie,
			restarter,
		)
	default:
		return qtie.New(
			version,
			buildVersion,
			programName,
			panicHandler,
			locations,
			settings,
			eventListener,
			updater,
			ie,
			restarter,
		)
	}
}
