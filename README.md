# Pursuit Documentation
Follow the instructions below to use the pursuit backend setup

## Setup
First, install the latest version of Go from golang.org

Next, install PostgreSQL for database
Then, navigate to postgresql.org and find the tutorial for installing and setting up Postgres

Make sure you have everything set up, and then connect to your database by running the connection string.
Example connection string: `"postgres://username:password@host:port/database`

If everything works, the database should be set up correctly.

### Other dependencies
Install goose for migrations
`go install github.com/pressly/goose/v3/cmd/goose@latest`

install SQLC
`go install github.com/kyleconroy/sqlc/cmd/sqlc@latest`

install Go dependencies
`go mod tidy`

### Env Setup
Set up an env file that looks something like this:
```
DB_URL="postgres://username:password@localhost:5432/yourdb?sslmode=disable"
PLATFORM="dev"
JWT_SECRET="your_jwt_secret"
```

### Database Migrations
Run the database migrations
`goose -dir sql/schema postgres "$DB_URL" up`

If something unexpected happens, drop the migration tables:
`goose -dir sql/schema postgres "$DB_URL" down`

to generate Go code from SQL with sqlc:
`sqlc generate`


### Run the Program
To run the program, navigate to /cmd/server and run:
`go run main.go`