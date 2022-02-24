# Product and subscription API

Go REST API with GORM ORM (POSTGRESQL), Mux and JWT

## Installation
To simplify the installation process, there is a dev docker file and a docker compose file that can be used to run the project locally, so the only dependency is a running [Docker](https://docs.docker.com/engine/install/) installation  
Once the Docker client is running, use the following command:

```bash
docker-compose up
```

Please note that the build described in `docker-compose.yml` file is purely a dev tool. We are running [air](https://github.com/cosmtrek/air) to have a hot reload API that monitors the code and rebuilds on the go.

*NOTE*: If you are using Windows you might run into issues with the hot reload because of how virtualization works on Windows (it is a WSL2 limitation, documented [here](https://github.com/cosmtrek/air/issues/190))

## Dockerization
On startup, Docker creates three containers, one for the actual API, one for Postgres and an additional container for PGAdmin4 (this is really useful for development and for general DB monitoring).
In general you fill find this:
| **Container** |      **Purpose**      | **Local port** |
|:-------------:|:---------------------:|:--------------:|
|      API      |   Hosting the Go app  |      8080      |
|   Postgresql  |    Database system    |      5432      |
|    PgAdmin4   | Database UI Navigator |      5050      |

# The `.env`

For development purposes I have setup a `.env` file on the root of the folder
 It contains the database credentials and also the pgAdmin credentials, you can log in into PgAdmin on your browser by visiting [http://localhost:5050/browser/](http://localhost:5050/browser/)

# Endpoints and API

To make the testing process easier, I created a Postman Collection and Environment. They are located at the root of the project on `go-api.postman_collection.json` and `go-api-test.postman_environment.json` respectively. Please import them into Postman to be able to check the endpoints smoothly.

## Auth
To actually authenticate requests, I created two endpoints `/auth/register` where you can create your user like this:
```bash
curl --location --request POST 'http://localhost:8080/auth/register' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username": "juan",
    "password": "admin123"
}'
```

And to actually authenticate requests we are using JWT, so to get a new token, use the `auth/token` endpoint with your newly created user, like this:
```bash
curl --location --request POST 'http://localhost:8080/auth/token' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username": "juan",
    "password": "admin123"
}'
```

That should be it. You now have a bearer token ready to authenticate all your requests to the API.

The rest of the endpoints are located on `api/...`, please remember to always send your bearer token with your request, or you will get rejected by the API:
```bash
curl --location --request GET 'http://localhost:8080/api/products' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2NDU2ODYyMjUsInVzZXJfaWQiOjJ9.n-RRL1rKCfXaFAV8-4R53o63rFzBnQI--Akbm2d3lDk'
```

# Project structure
```
go-api
├─ .air.toml
├─ .env
├─ .gitignore
├─ docker-compose.yml
├─ Dockerfile
├─ go-api-test.postman_environment.json
├─ go-api.postman_collection.json
├─ go.mod
├─ go.sum
├─ init.sql
├─ LICENSE
├─ main.go
├─ pkg
│  ├─ api
│  │  ├─ middleware
│  │  │  ├─ auth.go
│  │  │  └─ logs.go
│  │  └─ service
│  │     ├─ products
│  │     │  ├─ repository.go
│  │     │  └─ service.go
│  │     ├─ subscriptions
│  │     │  ├─ repository.go
│  │     │  └─ service.go
│  │     ├─ user
│  │     │  ├─ repository.go
│  │     │  └─ service.go
│  │     └─ voucher
│  │        └─ repository.go
│  ├─ auth
│  │  ├─ jwt.go
│  │  └─ passwords.go
│  ├─ migrations
│  │  └─ automigrate.go
│  ├─ models
│  │  ├─ products.go
│  │  ├─ shared
│  │  │  ├─ pagination.go
│  │  │  └─ query.go
│  │  ├─ subscription.go
│  │  ├─ users.go
│  │  └─ voucher.go
│  └─ shared
│     └─ utils
│        ├─ queryParser.go
│        └─ writeError.go
└─ README.md
```
The project was setup using a Microservice oriented structure, we have interface defined repositories and services for all entities that register routes with a central mux Router. I also include Middleware to actually get the user id from the bearer token and another one to log the request to stout. Because all repos and services are defined as interfaces they are highly testable both on the unit level and on the service integration level. You could spin up a server with mocked repositories to test all services end to end. Unfortunately, I did not have enough time to add unit tests to the project :(

I hope you have as much fun reviewing the challenge as I had building it :)
## License
[MIT](https://choosealicense.com/licenses/mit/)