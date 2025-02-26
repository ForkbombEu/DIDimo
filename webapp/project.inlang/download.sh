#!/bin/bash

urls=(
	"https://cdn.jsdelivr.net/npm/@inlang/message-lint-rule-empty-pattern@1/dist/index.js"
	"https://cdn.jsdelivr.net/npm/@inlang/message-lint-rule-identical-pattern@1/dist/index.js"
	"https://cdn.jsdelivr.net/npm/@inlang/message-lint-rule-missing-translation@1/dist/index.js"
	"https://cdn.jsdelivr.net/npm/@inlang/message-lint-rule-without-source@1/dist/index.js"
	"https://cdn.jsdelivr.net/npm/@inlang/message-lint-rule-valid-js-identifier@1/dist/index.js"
	"https://cdn.jsdelivr.net/npm/@inlang/plugin-message-format@2/dist/index.js"
	"https://cdn.jsdelivr.net/npm/@inlang/plugin-m-function-matcher@0/dist/index.js"
)

for url in "${urls[@]}"; do
	# Extract the package name with version from the URL using sed.
	# This grabs the text between "/npm/" and "/dist/"
	package_with_version=$(echo "$url" | sed -n 's|.*/npm/\(.*\)/dist/.*|\1|p')

	# Remove the version (everything after the last "@")
	package_name="${package_with_version%@*}"

	# Create the target directory (e.g. "./@inlang/plugin-m-function-matcher/dist/")
	mkdir -p "./${package_name}/dist"

	# Download the file into the created directory as index.js
	echo "Downloading $url to ./${package_name}/dist/index.js"
	curl -L "$url" -o "./${package_name}/dist/index.js"
done
