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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIncrementRestartFlag(t *testing.T) {
	var tests = []struct {
		in  []string
		out []string
	}{
		{[]string{"./bridge", "--restart", "1"}, []string{"./bridge", "--restart", "2"}},
		{[]string{"./bridge", "--restart", "2"}, []string{"./bridge", "--restart", "3"}},
		{[]string{"./bridge", "--other", "--restart", "2"}, []string{"./bridge", "--other", "--restart", "3"}},
		{[]string{"./bridge", "--restart", "2", "--other"}, []string{"./bridge", "--restart", "3", "--other"}},
		{[]string{"./bridge", "--restart", "2", "--other", "2"}, []string{"./bridge", "--restart", "3", "--other", "2"}},
		{[]string{"./bridge"}, []string{"./bridge", "--restart", "1"}},
		{[]string{"./bridge", "--something"}, []string{"./bridge", "--something", "--restart", "1"}},
		{[]string{"./bridge", "--something", "--else"}, []string{"./bridge", "--something", "--else", "--restart", "1"}},
		{[]string{"./bridge", "--restart", "bad"}, []string{"./bridge", "--restart", "1"}},
		{[]string{"./bridge", "--restart", "bad", "--other"}, []string{"./bridge", "--restart", "1", "--other"}},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.in, " "), func(t *testing.T) {
			assert.Equal(t, tt.out, incrementRestartFlag(tt.in))
		})
	}
}
