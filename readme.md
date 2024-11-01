# This is a simple code show that how to use go-gorm to implement optimistic lock.

## Using packages

- [go-gorm](https://github.com/go-gorm/gorm)
- [go-gorm-optimisticlock](https://github.com/go-gorm/optimisticlock)


## Example codes

- `gormPlugin` is the example code that use go-gorm-optimisticlock to implement optimistic lock.

- `pureGormUpdatedAt` is the example code that use go-gorm and field `UpdatedAt` to implement optimistic lock.

- `pureGormVersion` is the example code that use go-gorm and field `Version` to implement optimistic lock.
- `pureGormTransaction` is the example code that use go-gorm and transaction to implement AddBookCount method.

## How to run examples

1. create a `.env` file in the example code directory, and set the database connection string.

   ```toml
    # .env
    POSTGRES_IP="{your_postgres_ip}"
    POSTGRES_PORT="{your_postgres_port}"
    POSTGRES_DB="{your_postgres_db}"
    POSTGRES_USER="{your_postgres_user}"
    POSTGRES_PASSWORD="{your_postgres_password}"
   ```
2. run the example code. And you will got the same result that the book's count is 20.