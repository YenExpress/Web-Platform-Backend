# YenExpress Online Pharmacy and Telemedicine API

This web server provides the YenExpress web application  with authentication capabilities, relevant data for different user operations and handles all data processing and storage for seamless user experience.


The project is written in Go 1.18.2 and uses the
[Gin framework](https://gin-gonic.com/) to build a blazing fast, yet
efficient web API following the micro service design pattern.

Furthermore, it is capable of caching the latest data, to on the one hand
reduce the amount of requests (and outgoing traffic) to the upstream APIs, and
on the other hand reduce the response time to a minimum, especially for
subsequent requests. This works really well due to the asynchronous model
Gin provides.


## Install

    To download all relevant packages use go mod download


## Update go.mod file

    use go mod tidy

## Run the app

    If you just want to run your code, use go build or go run . - your dependencies will be downloaded and built automatically

# REST API Documentation

The REST API documentation is hosted [here](https://documenter.getpostman.com/view/22798352/2s93eSZaV8)