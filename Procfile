# SPDX-FileCopyrightText: 2025 Forkbomb BV
#
# SPDX-License-Identifier: AGPL-3.0-or-later
FE: ./scripts/wait-for-it.sh -s -t 0 localhost:8090 && bun run /app/webapp/build/index.js
BE: didimo serve --http=0.0.0.0:8090
TEMPORAL: temporal server start-dev --db-filename /app/pb_data/temporal.db
ISSUER: ./scripts/wait-for-it.sh -s -t 0 localhost:8280 && go run pkg/credential_issuer/worker/worker.go
WALLET: ./scripts/wait-for-it.sh -s -t 0 localhost:8280 && go run pkg/OpenID4VP/worker/worker.go
