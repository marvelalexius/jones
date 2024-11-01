# Jones Backend Service
Backend API for jones (Yes, exactly like you're thinking, if you know, you know). This API service can be used to (hopefully) get you out from that status (yes, i'm talking about *THAT* status).

## Features

- Register & Login
- JWT authentication token (including the expiry for the bearer key and the refresh token itself).
- Find Users
- Create Reaction (swipe left or right)
- Subscription using stripe (management, create, update, and cancel)
- Premium features to unlock swipe limit and to see who's been liking you

## Clone

First clone this repo by run:

```sh
$ git clone git@github.com:marvelalexius/jones.git
```

## Run jones in local

### Initialize

Firstly, run:

```sh
$ go mod tidy
```

### Environment
- The sample environments are provided in root folder
  - If you run jones in local, use `.env.example` to be `.env` file.
- This service also fully integrated with Stripe, so you'll need to add stripe keys if you want to test the subscription system.
  - If you want to use Stripe, you'll need to create a price in stripe dashboard and add the price id to newly seeded subscription plan in `stripe_product_id`
- For easier testing if you don't want to use stripe, I've also added a feature flag to toggle the subscription system.

### Database Migration

- Ensure you have already created the database. To migrate tables, run:

```sh
$ go run main.go migrate --direction=up
```

- To seed necessary data:

```sh
$ go run main.go migrate seed
```

## Running app

- To run HTTP server, hit:

```sh
$ go run main.go serve
```

## Unit Test & Lint

- Test command includes linters which you need to install `golangci-lint`.

```
<!-- Ubuntu -->
# binary will be $(go env GOPATH)/bin/golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.61.0
```

Then run `golangci-lint run`.

- To test the unit tests, you can run `go test ./... -v --cover`

## API Documentations

- To test the API, import `postman collection` from folder `api-docs/`. All the API is available there.
- To start testing the API, you need to add your base url padded with `/api/v1`add:
  - `host` to collection variable, represents the base url value
- For authentication tokens, it will be automatically included when you successfully ran register or login, so you don't need to add it manually

## Architecture

- jones uses a clean architecture with different naming. It has 4 domain layers such

  - Models Layer provides domain models
  - Repository Layer communicates to persistence layer
  - Service Layer stores action of process
  - HTTP Layer exchanges data between client & system

- Here's some reason why do I use clean architecture
  - Scalable:
    - As the application grows, it'll be easier adding new features or modify existing ones without widespreading negative impacts.
    - Flexible to changes, either it's database, API versions, or the responses, changes are contained without affecting one layer to another
  - Reliable:
    - Services are easily and independently testable.
    - Adapting Domain driven design & SOLID principles.
  - Maintainable
    - A way more organized code. With separated layers, we can easily see what does this folder specifically told to do.
    - Readable code, new developers can understand the code's purpose just by looking to the specific folders without getting lost to a technical details.
- this service could run process separately using `cobra` as commands in one `main` function. For the example, `serve` is the command for running http service, likewise `migration` to run migration. This also good approach if we intend to use worker or exposing another port in the future.

## Stacks

- `golang` as the programming language
- `gin` as the HTTP Framework
- `cobra` runs command
- `postgres` as the RDBMS
- `gorm` as the ORM
- `dbmate` as db migration
- `mockery` for mocking
- `golangci-lint` as linters
- `stripe` for payment and subscription management

## Author

```
Marvel Alexius
```
