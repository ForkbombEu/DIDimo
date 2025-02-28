#!/bin/bash
set -euxo pipefail
/app/didimo serve --dir /pb_data --http=0.0.0.0:8090 &
bun run /app/webapp/build/index.js
