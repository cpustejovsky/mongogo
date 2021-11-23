# Mongogo

## Description

Service to determine if a email domain is a catch-all domain. Pulled and modified from my [catchall](https://github.com/cpustejovsky/catchall) project

## Set Up
Move to the `app` directory and run `go build` and run `./app` to start

To connect to a MongoDB database, either pass in a `-uri` flag or set a `.env` file with `MONGO_URI` as the property
To specify a API port, pass in a `-addr` flag
To specify a pprof port, pass in a `-pprof` flag

## Usage

Starter kit for a Go service using MongoDB as database

## To-Dos
* Separate Database actions from domain
  * Wrap `primitive.ObjectIDFromHex(id)` in a helper function in this package
* ~~Add CRUD functionality for users~~
* Add integration tests for domain methods
* Add authentication
* Add metrics
* Make better use of the `context` package

### For `idiomatic_structure`
* Change helpers to a better name
  * Potentially handlerHelpers?
  * Generalize handlers package, include these helpers, and have a specific user handlers repo?
* TBD