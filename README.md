# Stan Coding Challenge

Implementation of the [Stan coding challenge](https://challengeaccepted.streamco.com.au/)

## Running

Runs with basic config in [cmd/stan.webhook/config.go](cmd/stan.webhook/config.go).
A config file can be passed via `-config` flag

```sh
go build ./cmd/stan.webhook/ && ./stan.webhook
```

## Build with Docker

```sh
make build
```
