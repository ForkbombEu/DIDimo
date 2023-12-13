export function typeCaster<T>(): TypeCaster<T> {
	return [] as T[];
}

export type TypeCaster<T> = T[];
