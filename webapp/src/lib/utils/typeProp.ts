// SPDX-FileCopyrightText: 2024 The Forkbomb Company
//
// SPDX-License-Identifier: AGPL-3.0-or-later

export type TypeProp<T> = Array<T>;

export function createTypeProp<T>(): TypeProp<T> {
  return [] as TypeProp<T>;
}
