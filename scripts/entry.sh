#!/bin/bash
set -euxo pipefail
/app/didimo serve --dir /pb_data "$COOLIFY_URL" &
bun run /app/webapp/build/index.js
