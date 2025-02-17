# liquidless-faker

[Faker.js](https://github.com/faker-js/faker) filter for liquidless templates

## Get Started

Install the dependency from npm:

```sh
npm i liquidless-faker
```

Import the `fake` filter

```js
import { fake } from 'liquidless-faker'
```

## Usage

**Example: Rendering a string**

```js
renderTemplate('{{ lorem.lines | fake: 1 }}', {}, {
  filters: {
    fake
  }
})
```

Outputs:

```
Saepe nesciunt officiis.
```
