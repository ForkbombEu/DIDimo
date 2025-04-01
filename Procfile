FE: ./scripts/wait-for-it.sh -s -t 0 localhost:8090 && bun run /app/webapp/build/index.js
BE: didimo serve --http 0.0.0.0:8090
TEMPORAL: temporal server start-dev --ui-port 8280 --db-filename /usr/local/bin/pb_data/temporal.db --ui-public-path /workflows
ISSUER: ./scripts/wait-for-it.sh -s -t 0 localhost:8280 && go run pkg/credential_issuer/worker/worker.go
WALLET: ./scripts/wait-for-it.sh -s -t 0 localhost:8280 && go run pkg/OpenID4VP/worker/worker.go
