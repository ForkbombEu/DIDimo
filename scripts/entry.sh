#!/bin/bash
set -euxo pipefail

didimo serve --http=0.0.0.0 &
wait-for-it.sh localhost:8090 && didimo-ui &
