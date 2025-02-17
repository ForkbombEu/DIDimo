# phasic

![phasic Demo](https://i.imgur.com/VUV3sdS.gif)

`phasic` is a low-level library for running asynchronous code repeatedly, in phases. This is useful for load and performance testing, where you want to start at a lower request per second rate and slowly increase the rate

## Get Started

**Install from npm**

```
npm i phasic
```

**Import in your project**

```js
import { runPhases } from 'phasic'
```

**Define phases**

```js
const phases = [{
  duration: 10,
  arrivalRate: 1
},
{
  duration: 10,
  arrivalRate: 100
}]
```

Where `duration` is the duration of the phase and `arrivalRate` is how many times the function will be called per second

In this example we have two phases, where the first phase lasts 10 seconds and a function is called once every second (10 calls) and a second phase, where a phase lasts 10 seconds, but the function is called 100 times per second (1000 calls). In total there would be 1010 function calls

> Note: Phase duration only defines the time window in which the functions should be called (not finished!). You should take that into account, if you're intending to run asynchronous tasks like HTTP requests or database reads

**Run your code in phases**

```js
import { runPhases } from '../src'
import fetch, { Response } from 'node-fetch'

const phases = [{
  duration: 10,
  arrivalRate: 1
},
{
  duration: 10,
  arrivalRate: 100
}]

let errorCount = 0

runPhases<Response>(phases, () => fetch('http://localhost:3000').catch(() => errorCount++))
.then((resultList) => {
  console.log('Status Codes:', resultList.map((result) => (result as PromiseFulfilledResult<Response>).value.status))
  console.log('Error rate:', (errorCount / resultList.length * 100).toFixed(2) + '%')
})
```

In this example, we're making a HTTP request to `http://localhost:3000` and reporting the error rate

The runPhases function returns a `Promise<PromiseSettledResult<T>[]>`. You can learn more about it in the [`Promise.allSettled()`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise/allSettled) MDN article
