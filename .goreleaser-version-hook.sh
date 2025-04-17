#!/bin/env bash

# SPDX-FileCopyrightText: 2025 Forkbomb BV
#
# SPDX-License-Identifier: AGPL-3.0-or-later
set -euo pipefail

echo "$GORELEASER_VERSION" > VERSION

# Optional: configure git if running in CI (especially GitHub Actions)
git config user.name "forkboteu"
git config user.email "apps@forkbomb.eu"

# Add and commit the file
git add VERSION
if git diff --cached --quiet; then
    echo "No changes to commit."
else
    git commit -m "chore(pre-release): 🚀 update VERSION to ${GORELEASER_VERSION} [ci skip]"
    git push origin main
fi
