project-xm
===========

## Prerequisites
- Docker (with compose support)
- Go 1.18+
- GNU Make

## Running locally
The following `make` command is provided for easy deployment on the local machine:
```
make docker/build docker/run
```

This should automatically build the project and launch all the necessary containers in Docker.

You can see all available commands along with a short help text simply by typing `make`.

For testing, you can use any valid HS256 JWT created with the secret key from `config.json` (`supersecret` is the default), such as this one:
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1MTYyMzkwMjJ9.rf0rI9lvdgd91-qfOjv4ChS0JHRWK2qeZMqTiMvvF5Q
```

## Running in production
1. Build the Docker image using `make docker/build`
2. Optionally, tag the image with appropriate name for your container registry
3. Deploy to Compose, Kubernetes, or plain Docker

## Tests
1. Run unit tests using `make unittest`
2. Run integration tests using `make integration-test`. These tests will run in a docker container using postgres database image
3. Run all tests using `make test`

## TODOs
- [ ] Better error handling & logging (i.e. return proper HTTP error codes for low-level errors)
- [ ] Stricter JWT validation
- [ ] GRPC delivery layer
- [ ] Atomic mutations & events (outbox table?)
- [ ] Add actual Kafka support
