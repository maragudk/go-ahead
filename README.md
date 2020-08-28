# go-ahead

A SaaS web app starter template written in Go.

## Features

- Web app with best practice project layout and single-binary deployment
- Stateless app ready for working behind a load balancer, such as [Caddy](https://caddyserver.com)
- Storage using [CockroachDB](https://www.cockroachlabs.com) (but can be easily changed to another database system)
- Metrics using [Prometheus](https://prometheus.io)
- Transactional emails using [Postmark](https://postmarkapp.com)
- Unit and integration tests using [Docker Compose](https://docs.docker.com/compose/),
  with CI using Github Actions and [CircleCI](https://circleci.com)

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
