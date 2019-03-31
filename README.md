[![Go Report Card](https://goreportcard.com/badge/github.com/tupass/tupass-backend)](https://goreportcard.com/report/github.com/tupass/tupass-backend) 
[![Build Status](https://travis-ci.org/tupass/tupass-backend.svg?branch=master)](https://travis-ci.org/tupass/tupass-backend) 
[![Documentation](https://godoc.org/github.com/tupass/tupass-backend?status.svg)](https://godoc.org/github.com/tupass/tupass-backend) 
[![GPLv3 License](https://img.shields.io/badge/License-GPLv3-brightgreen.svg)](https://github.com/tupass/tupass-frontend/LICENSE) 
[![Release](https://img.shields.io/github/release/tupass/tupass-backend.svg?label=Release)](https://github.com/tupass/tupass-backend/releases)

# TUPass Backend

Backend part of the TUPass Password Strength Meter.  
For use with a web interface also see [TUPass Frontend](https://github.com/tupass/tupass-frontend).  
For a live version see [tupass.pw](https://tupass.pw).

## Cloning

This project uses Go version 1.12.1. If you have not already installed Go, please [go get it](https://golang.org/dl/).

Run `go get github.com/tupass/tupass-backend` to pull this repository.

## Installing Dependencies

Run `make dep` to fetch dependencies for this project.

## Building

Run `make build` to build the project.

## Deployment

For local-only usage there is also a bundled version (containing frontend & backend) available as Debian package, Windows executable or Linux binary (see [Releases](https://github.com/tupass/tupass-backend/releases)). 

### For Local Usage

Run `make` to start a dev server. The API will be available at `http://localhost:8000`.

### For Remote / External usage

Run `make run-prod` to start a staging/production server. The API will be available at `http://localhost:8001`.

## Testing

Run `make test` to execute tests.

## License

This project is licensed under the GPLv3 License - see the [LICENSE](LICENSE) file for details.
