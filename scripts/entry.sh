#!/bin/bash
set -euxo pipefail
/app/didimo serve --dir /pb_data "$COOLIFY_FQDN" &
bun run /app/webapp/build/index.js
