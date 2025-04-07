#!/bin/bash

# SPDX-FileCopyrightText: 2025 Forkbomb BV
#
# SPDX-License-Identifier: AGPL-3.0-or-later

set -euxo pipefail
/app/didimo serve --dir /pb_data --http=0.0.0.0:8090 &
bun run /app/webapp/build/index.js
