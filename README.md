# go-generic-dao

### Reducing `database/sql` boilerplate using Go generics without the need for added dependencies. 

## Medium Article

For a full explanation of the code, please see the accompanying blog post.

- [Abstracting database/sql boilerplate with Go generics](https://medium.com/eureka-engineering/abstracting-database-sql-boilerplate-with-go-generics-1471d969beed)

## Setup

### Environment variables 

Use [direnv](https://direnv.net/) as below, 
or manually set contents of `.env` in your terminal or editor. 

```bash
direnv allow
```

### Database 

```bash
docker-compose up -d
make db-init
```

## Run tests

```bash
make test
```

## Code Overview

### `config` package

* Initializes environment variables and stores them in a singleton
* Initialize with `config.Init()`
* Access with `config.Instance()`

### `db` package

* Initializes a database instance and stores it in a singleton
* Initialize with `db.Init()`
* Access with `db.Instance()`

### `test` package

* Provides helpers for testing
* Call `test.InitConfig()` instead of `config.Init()` from `TestMain` functions

### `like` and `user` packages

* Contain entities, DAOs and tests
