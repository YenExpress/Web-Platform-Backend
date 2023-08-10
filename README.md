# Ignite.dev Platform Authentication Microservice

This web server provides the Ignite.dev platform  with authentication capabilities including user account creation, login and session management, and all authentication and authorization related operations.


The project is written in Go 1.18.2 and leverages the
[Echo framework](https://echo.labstack.com/) for a fast, performant and highly scalable service.


### Directory Structure Explanation

- `cmd`: This directory contains executable entry point for the microservice.

- `internal`: This directory holds internal packages that are localized to the microservice. This is a convention to prevent other projects from importing these internal packages without proper context and relationship.
  - `api`: This package contains all API-related logic.
    - `common` : Contains common constants, error and custom datatypes used throughout the application.
    - `pkg` : Contains sub-modules for authentication specific to user category e.g `developer`
    - `config`: Configuration management related code resides   here.
    - `database`: Database-related code, such as database connections.
    - `shared` : Contains shared authentication resources generic across user categories e.g generic `dto` (Data Transfer Objects), `interactors`, `handlers` and `middlewares`.
  - `utils` : Helper functions not closely tied to the low level API implementation logic but relevant to achieve tasks.
- `prestart` : Contains code to be run before application starts.
- `test` : All test files reside here.
- `web` : Contains public folders and files including static files and templates.




## To run all tests, build and preview application 

    use command `make all`

## To run all tests alone

    use command `make test`

## To build application executable

    use command `make build`

## To run application in debug mode

    use command `make run`

## To preview application by executing binary

    use command `make preview`