# This is a simple code show that how to use go-gorm to implement optimistic lock.

## Using packages

- [go-gorm](https://github.com/go-gorm/gorm)
- [go-gorm-optimisticlock](https://github.com/go-gorm/optimisticlock)


## Sample codes

- `gormPlugin` is the sample code that use go-gorm-optimisticlock to implement optimistic lock.

- `pureGormUpdatedAt` is the sample code that use go-gorm and field `UpdatedAt` to implement optimistic lock.

- `pureGormVersion` is the sample code that use go-gorm and field `Version` to implement optimistic lock.

## How to run 

1. create a `.env` file in the sample code directory, and set the database connection string.

   ```toml
    # .env
    POSTGRES_IP="{your_postgres_ip}"
    POSTGRES_PORT="{your_postgres_port}"
    POSTGRES_DB="{your_postgres_db}"
    POSTGRES_USER="{your_postgres_user}"
    POSTGRES_PASSWORD="{your_postgres_password}"
   ```
2. run the sample code. And you will got the same result that the book's count is 20.