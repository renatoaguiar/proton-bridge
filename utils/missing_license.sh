#!/bin/bash

# Copyright (c) 2021 Proton Technologies AG
#
# This file is part of ProtonMail Bridge.
#
# ProtonMail Bridge is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# ProtonMail Bridge is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with ProtonMail Bridge.  If not, see <https://www.gnu.org/licenses/>.

YEAR=`date +%Y`
MISSING_FILES=`find . -not -path "./vendor/*" -not -path "./vendor-cache/*" -not -path "./.cache/*" -not -name "*mock*.go" -regextype posix-egrep -regex ".*\.go|.*\.qml|.*\.sh|.*\.py" -exec grep -L "Copyright (c) ${YEAR} Proton Technologies AG" {} \;`

for f in ${MISSING_FILES}
do
    echo -n "MISSING LICENSE or WRONG YEAR in $f"
    if [[ $1 == "add" ]]
    then
        cat ./utils/license_header.txt $f > tmp
        mv tmp $f
        echo -n "... license added"
    fi
    if [[ $1 == "change-year" ]]
    then
        sed -i "s/Copyright (c) [0-9]\\{4\\} Proton Technologies AG/Copyright (c) ${YEAR} Proton Technologies AG/" $f || exit 3
        echo -n "... replaced copyright year"
    fi
    echo
done

[[ "$1" == "check" ]] && [[ -n ${MISSING_FILES} ]] && exit 1
exit 0
