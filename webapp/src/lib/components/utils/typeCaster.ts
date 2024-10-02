// SPDX-FileCopyrightText: 2024 The Forkbomb Company
//
// SPDX-License-Identifier: AGPL-3.0-or-later

export function typeCaster<T>(): TypeCaster<T> {
	return [] as T[];
}

export type TypeCaster<T> = T[];
