#!/bin/bash
set -euxo pipefail

/usr/bin/didimo serve --dir /pb_data --http=0.0.0.0:8090 &
bun run /app/build/index.js
