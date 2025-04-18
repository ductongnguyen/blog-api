# blog-api
Golang Clean Architecture REST API with gin example. Code is inspired from and based on [Golang Clean Architecture REST API example](https://github.com/AleksK1NG/blog-apiitecture-REST-API)

#### Full list what has been used:
* [gin](https://github.com/gin-gonic/gin) - Web framework
* [gorm](https://gorm.io/) - Extensions to database/sql.
* [caarlos0](https://github.com/caarlos0/env) - Go configuration with fangs
* [go-redis](https://github.com/go-redis/redis) - Type-safe Redis client for Golang
* [zap](https://github.com/uber-go/zap) - Logger
* [validator](https://github.com/go-playground/validator) - Go Struct and Field validation
* [jwt-go](https://github.com/golang-jwt/jwt) - JSON Web Tokens (JWT)
* [uuid](https://github.com/google/uuid) - UUID
* [testify](https://github.com/stretchr/testify) - Testing toolkit
* [gomock](https://github.com/golang/mock) - Mocking framework

This project has 4 layer :

- Models Layer
- Repository Layer
- Usecase Layer
- Delivery Layer

### Requirement
- Golang 1.20 or highest
- MySQL 8.0 or highest
- Redis 6 or highest 

### How To Run This Project

> Make Sure you have run the migrations/01_create_initial_tables.up.sql in your mysql

#### Run the Applications

Here is the steps to run it with `Makefile`


```bash
# Create a .env file from .env.dist and ensure that the configurations inside it are correct according to your environment settings
$ cp .env.example .env

# Run the application
$ make run
```