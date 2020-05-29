![Go Test Workflow](https://github.com/lucasreed/smol/workflows/testing/badge.svg)
# Smol

Take urls and make them [smol](https://www.urbandictionary.com/define.php?term=Smol).

A simple url shortener API server written in go.

## Run locally

To run the dockerized version you only need to have docker installed and run the following:

`make docker`

This will run an ephemeral container that you can run api calls against on port 8080 of your machine.

CTRL-C to exit, any redirect data added to the server will not be persisted.

## API Endpoints

All api endpoints will start with `/api/${VERSION}/`

### v1
`/api/v1/add` - `POST` - add a redirect. Expects json POST data in the following format: `{"Destination":"www.google.com"}`

Example usage:

```shell
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"Destination":"www.google.com"}' \
  http://localhost:8080/api/v1/add
```
