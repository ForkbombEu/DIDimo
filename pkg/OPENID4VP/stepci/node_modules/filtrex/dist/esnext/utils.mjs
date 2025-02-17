import { UnexpectedTypeError } from "./errors.mjs";

/**
 * Determines whether an object has a property with the specified name.
 * @param {object} obj the object to be checked
 * @param {string|number} prop property name
 */
export function hasOwnProperty(obj, prop) {
  if (typeof obj === "object" || typeof obj === "function") {
    return Object.prototype.hasOwnProperty.call(obj, prop);
  }

  return false;
}

/**
 * Mathematically correct modulo
 * @param {number} a
 * @param {number} b
 * @returns {number}
 */
export function mod(a, b) {
  return ((a % b) + b) % b;
}

/**
 * Converts instances of Number, String and Boolean to primitives
 */
export function unbox(value) {
  if (typeof value !== "object") return value;

  if (
    value instanceof Number ||
    value instanceof String ||
    value instanceof Boolean
  )
    return value.valueOf();
}

/**
 * Unboxes value and unwraps it from a single-element array
 */
export function unwrap(value) {
  if (Array.isArray(value) && value.length === 1) value = value[0];

  return unbox(value);
}

/**
 * Returns the type of a value in a neat, user-readable way
 */
export function prettyType(value) {
  value = unwrap(value);

  if (value === undefined) return "undefined";
  if (value === null) return "null";
  if (value === true) return "true";
  if (value === false) return "false";

  if (typeof value === "number") return "number";
  if (typeof value === "string") return "text";

  if (typeof value !== "object" && typeof value !== "function")
    return "unknown type";

  if (Array.isArray(value)) return "list";

  return "object";
}

// Type assertions/coertions

export function num(value) {
  value = unwrap(value);
  if (typeof value === "number") return value;

  throw new UnexpectedTypeError("number", prettyType(value));
}

export function str(value) {
  value = unwrap(value);
  if (typeof value === "string") return value;

  throw new UnexpectedTypeError("text", prettyType(value));
}

export function numstr(value) {
  value = unwrap(value);
  if (typeof value === "string" || typeof value === "number") return value;

  throw new UnexpectedTypeError("text or number", prettyType(value));
}

export function bool(value) {
  value = unwrap(value);
  if (typeof value === "boolean") return value;

  throw new UnexpectedTypeError(
    "logical value (“true” or “false”)",
    prettyType(value),
  );
}

export function arr(value) {
  if (value === undefined || value === null) {
    throw new UnexpectedTypeError("list", prettyType(value));
  }

  if (Array.isArray(value)) {
    return value;
  } else {
    return [value];
  }
}

/**
 * Array.flat polyfill from MDN
 */
export function flatten(input) {
  const stack = [...input];
  const res = [];
  while (stack.length) {
    // pop value from stack
    const next = stack.pop();
    if (Array.isArray(next)) {
      // push back array items, won't modify the original input
      stack.push(...next);
    } else {
      res.push(next);
    }
  }
  // reverse to restore input order
  return res.reverse();
}

/**
 * Template string noop tag from
 * https://github.com/lleaff/tagged-template-noop/blob/master/index.js
 */
export function defaultTag(strings, ...keys) {
  const lastIndex = strings.length - 1;
  return (
    strings.slice(0, lastIndex).reduce((p, s, i) => p + s + keys[i], "") +
    strings[lastIndex]
  );
}

function _code(fragments, params, skipParentheses) {
  const args = [];

  for (let i = 0; i < fragments.length - 1; i++) {
    args.push(fragments[i]);
    args.push(params[i]);
  }

  args.push(fragments[fragments.length - 1]);

  const argsJs = args
    .map(function (a) {
      return typeof a == "number" ? "$" + a : JSON.stringify(a);
    })
    .join(",");

  return skipParentheses
    ? "$$ = [" + argsJs + "];"
    : '$$ = ["(", ' + argsJs + ', ")"];';
}

export function code(fragments, ...params) {
  return _code(fragments, params, false);
}

export function parenless(fragments, ...params) {
  return _code(fragments, params, true);
}

export function noopTag(...args) {
  return [defaultTag(...args), parenless(...args)];
}
