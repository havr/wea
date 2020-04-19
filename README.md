## What is it

This is a test task for FastBill. 
It's a web server that serves a single endpoint `/GET city-information?name=<location-name>`. 
It responds with the Wikipedia brief introduction for the given location joined with it's current temperature and a short weather description.

## API Spec

The server responds in the JSON format.

### Error Response

The error response use the same structure:

| field    | type   | description |
|----------|--------|-------------|
| `error`  | string | a brief description of what happened |

### Endpoints

`GET /city-information?name=<city-name>`

*Response*

| field    | type | description |
|----------|------|-------|
| `description` | string | the city description from Wikipedia. |
| `weatherSituation` | string | the current weather situation description (example: a few clouds). |
| `temperature` | number | the city temperature in celsius.|

*Error codes*

| code | meaning |
|------|-------|
| 404 | either the wikipedia entry or weater data not found |
| 500 | an unexpected server error |

## Assumptions

The app doesn't check if a user queries exactly a city, so it returns a result if there is a city in OpenWeather registry and the corresponding wikipedia article.

Moreover, it doesn't resolve ambiguities (e.g. Washington D.C. and Washington the State).

## Envs

```bash
export WE_SERVE_AT=:7070 # Required, a host to serve at.
export WE_OW_API_KEY=d623c13ca45eae0da784b3e6f8d6b17d # Required, an OpenWeather API KEY.
```

## Running

Set the environment variables up and then simply run.
```
make run
# or alternatively
go run cmd/wea/wea.go
```

## Linters

Install the tools:
 - golangci-lint: use https://github.com/golangci/golangci-lint#install
 - golint: `go get -u golang.org/x/lint/golint`

Run linters:
```
make lint
```

## Testing

HTTP server and business logic are partially covered with unit tests.
Third-party clients have basic integration tests to ensure they work properly.
In order to test OpenWeather, one should provide the `WE_OW_API_KEY` env to the test.
For example:
```
WE_OW_API_KEY=d623c13ca45eae0da784b3e6f8d6b17d make test
```

