// SPDX-FileCopyrightText: 2024 The Forkbomb Company
//
// SPDX-License-Identifier: AGPL-3.0-or-later

export function log(...args: unknown[]) {
	if (process.env.NODE_ENV === 'development') {
		console.log(...args);
	}
}
