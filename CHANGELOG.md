# Changelog

All notable changes to this project will be documented in this file.

<!-- The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html). -->

## v0.1.2

### Changed

- Use application/json content-type between core and sub runtime [#30](https://github.com/scaleway/serverless-functions-go/pull/30)
- Remove TriggerType (only HTTP) and use APIGatewayProxyRequest type for Event [#29](https://github.com/scaleway/serverless-functions-go/pull/29)

## v0.1.1

### Fixed

- Headers added by the function handler response were lost #26
- `FunctionInvoker` now construct a request with a path that matches the one of the incoming event #22

## v0.1.0

### Added

- Initial project setup
- Local testing utils
- Repository setup