export type TypeProp<T> = Array<T>;

export function createTypeProp<T>(): TypeProp<T> {
  return [] as TypeProp<T>;
}
