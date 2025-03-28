ui: ./scripts/wait-for-it.sh -s -t 0 localhost:8090 && bun run /app/webapp/build/index.js
pocketbase: didimo serve --dir /pb_data --http 0.0.0.0:8090
temporal: temporal server start-dev --ui-port 8080 --db-filename /pb_data/temporal.db --ui-public-path /workflows
workflow-cred: go run pkg/credential_issuer/worker/worker.go
workflow-oid4: go run pkg/OpenID4VP/worker/worker.go
