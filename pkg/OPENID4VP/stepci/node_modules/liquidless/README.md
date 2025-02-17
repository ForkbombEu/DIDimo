# liquidless

Shopify's Liquid template engine, but less powerful. Perfect for configuration files

## Get Started

Install the dependency from npm:

```sh
npm i liquidless
```

Import the `renderString` function

```js
import { renderString } from 'liquidless'
```

## Usage

**Example: Rendering a string**

```js
renderString('Hello, {{ world }}', { world: 'world!' })
```

Outputs:

```
Hello, world!
```

**Example: Using filters**

```js
renderString('Hello, {{ world | upcase }}', { world: 'world!' })
```

Outputs:

```
Hello, WORLD!
```

**Example: Supplying custom filters**

```js
renderString('Hello, {{ world | something: 1, 2, 3 }}', { world: 'world!' }, {
  filters: {
    something: (value, args, variable) => `${value} ${args.join(', ')} (${variable})`
  }
})
```

Outputs:

```
Hello, world! 1, 2, 3 (world)
```

**Example: Rending values in an object**

```js
import { renderObject } from 'liquidless'
renderObject([{hello: {world: '{{ world }}'}}], { world: 'world!' })
```

Outputs

```js
[{hello: {world: 'world!'}}]
```

## Filters

- `append` - Adds a given string to the end of a string.
- `base64_decode` - Decodes a string in Base64 format.
- `base64_encode` - Encodes a string to Base64 format.
- `camelize` - Converts a string to CamelCase.
- `capitalize` - Capitalizes the first word in a string and downcases the remaining characters.
- `downcase` - Converts a string to all lowercase characters.
- `escape` - Escapes special characters in HTML, such as `<>`, ', and `&`, and converts characters into escape sequences.
- `hmac_sha1` - Converts a string into an SHA-1 hash using a hash message authentication code (HMAC).
- `hmac_sha256` - Converts a string into an SHA-256 hash using a hash message authentication code (HMAC).
- `lstrip` - Strips all whitespace from the left of a string.
- `md5` - Converts a string into an MD5 hash.
- `newline_to_br` - Converts newlines (`\n`) in a string to HTML line breaks (`<br>`).
- `pluralize` - Outputs the singular or plural version of a string based on a given number.
- `prepend` - Adds a given string to the beginning of a string.
- `remove` - Removes any instance of a substring inside a string.
- `remove_first` - Removes the first instance of a substring inside a string.
- `remove_last` - Removes the last instance of a substring inside a string.
- `replace` - Replaces any instance of a substring inside a string with a given string.
- `replace_first` - Replaces the first instance of a substring inside a string with a given string.
- `replace_last` - Replaces the last instance of a substring inside a string with a given string.
- `rstrip` - Strips all whitespace from the right of a string.
- `sha1` - Converts a string into an SHA-1 hash using a hash message authentication code (HMAC).
- `sha256` - Converts a string into an SHA-256 hash using a hash message authentication code (HMAC).
- `slice` - Returns a substring or series of array items, starting at a given 0-based index.
- `split` - Splits a string into an array of substrings based on a given separator.
- `strip` - Strips all whitespace from the left and right of a string.
- `strip_html` - Strips all HTML tags from a string.
- `strip_newlines` - Strips all newline characters (line breaks) from a string.
- `toInt` - Converts a value to Int
- `toFloat` - Converts a value to Float
- `toString` - Converts a value to String
- `upcase` - Converts a string to all uppercase characters.
- `url_decode` - Decodes a string to URL-safe format by converting percent-encoded characters to special characters.
- `url_encode` - Encodes a string to URL-safe format by converting special characters to percent-encoded characters.
