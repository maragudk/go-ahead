# go-ahead

A starter template for SaaS web apps. I use this for quickly getting a new web app set up.

## Features

- Starts serving HTTP on an internal and an external interface
- Prometheus metrics exposed at `/metrics` on the internal interface
- A `/health` endpoint for load balancers on the internal interface
- Ready to connect to CockroachDB, set up in development through Docker Compose
- …and of course the building, (integration) testing, linting etc. set up in the `Makefile`

## Usage

- Clone or fork this repository (or use Github's templating feature)
- Make a global search/replace for the keyword `ahead`
- Enjoy with coffee and a biscuit

## About

This project is brought to you by [maragu](https://www.maragu.dk).

If you like and/or use this project, please consider sponsoring it at [github.com/sponsors/maragudk](https://github.com/sponsors/maragudk).
