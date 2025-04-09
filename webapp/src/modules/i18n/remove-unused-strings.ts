// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

import fs from 'node:fs';
import path from 'node:path';

/**
 * Helper function to recursively find all files in a directory
 */
function getAllFilesInFolder(dirPath: string, exclude: string[] = [], arrayOfFiles: string[] = []) {
	const files = fs.readdirSync(dirPath);

	files.forEach(function (file) {
		const fullPath = path.join(dirPath, file);
		if (exclude.includes(fullPath)) return;
		if (fs.statSync(fullPath).isDirectory()) {
			arrayOfFiles = getAllFilesInFolder(fullPath, exclude, arrayOfFiles);
		} else {
			arrayOfFiles.push(fullPath);
		}
	});

	return arrayOfFiles;
}

/**
 * Extract used keys from the files based on the "m." pattern
 * */
function getKeysInFiles(files: string[], keys: string[]) {
	const usedKeys = new Set();

	files.forEach((file) => {
		const fileContent = fs.readFileSync(file, 'utf-8');
		keys.filter((k) => fileContent.includes(k)).forEach((k) => usedKeys.add(k));
	});

	return Array.from(usedKeys) as string[];
}

/**
 * Filter the JSON file by removing keys not present in the valid keys list
 */
function filterJsonByKeys(json: Record<string, string>, validKeys: string[]) {
	const newJson: Record<string, string> = {};

	const SCHEMA_KEY = '$schema';
	if (SCHEMA_KEY in json) {
		newJson[SCHEMA_KEY] = json[SCHEMA_KEY];
	}

	Object.keys(json).forEach((key) => {
		if (validKeys.includes(key)) {
			newJson[key] = json[key];
		}
	});

	return newJson;
}

/**
 * Main function to extract keys and filter the JSON
 */
function main(searchFolder: string, exclude: string[], jsonFilePath: string) {
	const allFiles = getAllFilesInFolder(searchFolder, exclude);
	const json = JSON.parse(fs.readFileSync(jsonFilePath, 'utf-8')) as Record<string, string>;
	const keys = Object.keys(json);
	const validKeys = getKeysInFiles(allFiles, keys);
	const filteredJson = filterJsonByKeys(json, validKeys);
	fs.writeFileSync(jsonFilePath, JSON.stringify(filteredJson, null, 2));
	console.log();
	console.log(`Removed unused strings ✨`);
}

// // Run the script with folder and JSON file arguments
// const searchFolder = process.argv[2];
// const jsonFilePath = process.argv[3];

// if (!searchFolder || !jsonFilePath) {
//   console.log('Usage: node script.js <searchFolder> <jsonFilePath>');
//   process.exit(1);
// }

// TODO - Optimize also for other languages
main('./src/', ['src/modules/i18n/paraglide'], `./messages/en.json`);
