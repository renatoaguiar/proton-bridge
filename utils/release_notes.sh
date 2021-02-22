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


# Generate HTML release notes
# hosted at https://protonmail.com/download/{ie,bridge}/{stable,early}_releases.html
INFILE=$1
OUTFILE=${INFILE//.md/.html}

# Load props
APP_NAME="Import-Export app"
if [[ "$INFILE" =~ bridge ]]; then 
    APP_NAME="Bridge"
fi

CHANNEL=early
if [[ "$INFILE" =~ stable ]]; then 
    CHANNEL=stable
fi

# Check dependencies
if ! which pandoc; then
  echo "PANDOC NOT FOUND!\nPlease install pandoc in order to build release notes."
  exit 1
fi

# Build release notes
pandoc $INFILE -f markdown -t html -s -o $OUTFILE -c utils/release_notes.css --self-contained --section-divs --metadata title="Release notes - ProtonMail $APP_NAME - $CHANNEL"
