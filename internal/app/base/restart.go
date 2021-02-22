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

package base

import (
	"os"
	"os/exec"
	"strconv"

	"github.com/sirupsen/logrus"
)

// maxAllowedRestarts controls after how many crashes the app will give up restarting.
const maxAllowedRestarts = 10

func (b *Base) restartApp(crash bool) error {
	var args []string

	if crash {
		args = incrementRestartFlag(os.Args)[1:]
	} else {
		args = os.Args[1:]
	}

	logrus.
		WithField("command", b.command).
		WithField("args", args).
		Warn("Restarting")

	return exec.Command(b.command, args...).Start() // nolint[gosec]
}

// incrementRestartFlag increments the value of the restart flag.
// If no such flag is present, it is added with initial value 1.
func incrementRestartFlag(args []string) []string {
	res := append([]string{}, args...)

	hasFlag := false

	for k, v := range res {
		if v != "--restart" {
			continue
		}

		hasFlag = true

		if k+1 >= len(res) {
			continue
		}

		n, err := strconv.Atoi(res[k+1])
		if err != nil {
			res[k+1] = "1"
		} else {
			res[k+1] = strconv.Itoa(n + 1)
		}
	}

	if !hasFlag {
		res = append(res, "--restart", "1")
	}

	return res
}
