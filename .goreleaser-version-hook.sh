#!/bin/env bash
set -euo pipefail

echo "$GORELEASER_VERSION" > VERSION

# Optional: configure git if running in CI (especially GitHub Actions)
git config user.name "goreleaser"
git config user.email "goreleaser@forkbomb.eu"

# Add and commit the file
git add VERSION
if git diff --cached --quiet; then
    echo "No changes to commit."
else
    git commit -m "chore(pre-release): 🚀 update VERSION to ${GORELEASER_VERSION} [ci skip]"
    git push origin main
fi
