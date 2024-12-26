# Pocketbase backend

Starter for a backend that includes pocketbase and [keypairoom](https://github.com/LedgerProject/keypairoom) key generation.

## Run Locally

```
docker-compose build
docker-compose up
```

## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

`SALT`: 32 bytes base64 encoded

`RESTROOM_URL` the URL of a [zenflows-crypto](https://github.com/interfacerproject/zenflows-crypto) instance

## Recipes

To protect a collection using the `authorizations` collection:

1. Add an `owner` attribute to the model you want to protect
2. Paste this string in the API rules section, where needed

```
owner.id = @request.auth.id || (@collection.authorizations.users.id ?= @request.auth.id && @collection.authorizations.record_id ?= id)
```

## API Reference

#### Keypairoom HMAC generation

```http
  GET /api/keypairoom-server
```

| Parameter | Type         | Description                          |
| :-------- | :----------- | :----------------------------------- |
| `data`    | `dictionary` | **Required**. Identifier of the user |

## Authors

-   [@albertolerda](https://github.com/albertolerda)
-   [@bbtgnn](https://github.com/bbtgnn)
-   [@puria](https://github.com/puria)

## License

```
Pocketbase starter - Zenflows economic model API's
Copyleft (ɔ) 2023 Dyne.org foundation

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.

```
