pocketbase: make build && ./didimo serve $DOMAIN
ui: ./scripts/wait-for-it.sh $DOMAIN:8090 && cd webapp && bun run serve
docs: cd docs && bun i && bun run docs:build && bun run docs:preview
