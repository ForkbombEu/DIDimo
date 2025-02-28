#!/bin/bash
set -euxo pipefail
/app/didimo serve --dir /pb_data --https=0.0.0.0:8090 "$COOLIFY_FQDN" &
bun run /app/webapp/build/index.js
