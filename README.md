# go-ahead

[![GoDoc](https://godoc.org/github.com/maragudk/go-ahead?status.svg)](https://godoc.org/github.com/maragudk/go-ahead)

A SaaS web app starter template written in Go.

## Features

- Web app with best practice project layout and single-binary deployment
- Stateless app ready for working behind a load balancer, such as [Caddy](https://caddyserver.com)
- Storage using [Postgres](https://www.postgresql.org) (but can be easily changed to another database system)
- Metrics using [Prometheus](https://prometheus.io)
- Transactional emails using [Postmark](https://postmarkapp.com)
- Unit and integration tests using [Docker Compose](https://docs.docker.com/compose/),
  with CI using Github Actions and [CircleCI](https://circleci.com)
- Views using [gomponents](https://github.com/maragudk/gomponents) and [TailwindCSS](https://tailwindcss.com)

### Roadmap

- [Stripe](https://stripe.com) integration for subscription payments
- Authentication and authorization without third-party dependencies
- Admininistration panel
- [Sentry](https://sentry.io/) integration

## Usage

- Clone or fork this repository (or use Github's templating feature)
- Make a global search/replace for the keyword `ahead`
- Enjoy with coffee and a biscuit

## About

This project is brought to you by [maragu](https://www.maragu.dk).

If you like and/or use this project, [please consider sponsoring it](https://github.com/sponsors/maragudk).
