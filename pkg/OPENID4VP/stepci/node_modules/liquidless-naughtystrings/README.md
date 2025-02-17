# liquidless-naughtystrings

[Big List of Naughty Strings](https://github.com/minimaxir/big-list-of-naughty-strings) filter for liquidless templates

## Get Started

Install the dependency from npm:

```sh
npm i liquidless-naughtystrings
```

Import the `naughtystring` filter

```js
import { naughtystring } from 'liquidless-naughtystrings'
```

## Usage

**Example: Rendering a string**

```js
renderTemplate('{{ | naughtystring }}', {}, {
  filters: {
    naughtystring
  }
})
```

Outputs:

```
<a href="\xE1\x9A\x80javascript:javascript:alert(1)" id="fuzzelement1">test</a>
```
