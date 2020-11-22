[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Version](https://img.shields.io/badge/package-v0.1.2-blue.svg?maxAge=2592000)](https://github.com/nhood-org/nhood-engine-service-tests/releases/tag/v0.1.2)
[![Version](https://img.shields.io/badge/docker-v0.1.2-blue.svg?maxAge=2592000)](https://github.com/nhood-org/repository/packages/199509)
[![CircleCI](https://circleci.com/gh/nhood-org/nhood-engine-service-tests.svg?style=shield)](https://circleci.com/gh/nhood-org/nhood-engine-service-tests)
[![Go Report Card](https://goreportcard.com/badge/github.com/nhood-org/nhood-engine-service-tests)](https://goreportcard.com/report/github.com/nhood-org/nhood-engine-service-tests)

# nhood-engine-service-tests

Acceptance tests for nhood engine service

## Pre-requisites

- Go
- Godog
- Make

## Validate tests

In order to validate tests against a mocked server use the following make command:

```bash
make validate
```

## Usage

In order to run test against a target server, specify TEST_TARGET_HOST and use the following make command:

```bash
export TEST_TARGET_HOST=localhost:8080
make run
```

## CI/CD

Project is continuously integrated within a `circleCi` pipeline that link to which may be found [here](https://circleci.com/gh/nhood-org/workflows/nhood-engine-service-tests)

Pipeline is fairly simple:

1. Validate tests
1. Build docker image with test application

Configuration of CI is implemented in `.circleci/config.yml`.

## Versioning

In order to release new package version, execute the following script:

```bash
export CIRCLE_CI_USER_TOKEN=<CIRCLE_CI_USER_TOKEN>
export NEW_VERSION=<NEW_VERSION>
make trigger-circle-ci-release
```

In order to release new version of docker image, execute the following script:

```bash
export CIRCLE_CI_USER_TOKEN=<CIRCLE_CI_USER_TOKEN>
export NEW_VERSION=<NEW_VERSION>
make trigger-circle-ci-docker-release
```

## License

`nhood-engine-service-tests` is released under the MIT license:
- https://opensource.org/licenses/MIT
