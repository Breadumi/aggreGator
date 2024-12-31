This is a learning project from Boot.dev. It requires [Go](https://go.dev/doc/install) and [Postgres](https://www.postgresql.org/download/) to be installed. At this point, it will only run in UNIX-like systems.

The main purpose of this project is to understand database creation, migration, and queries. Queries are written in PostgreSQL and translated into Go using SQLC. Migrations are handled with Goose. 

Functionality includes:
- Registering users, RSS feeds, and follows in a database
- Scraping RSS feeds at set intervals and saving new posts
- Browsing saved posts from followed RSS feeds

This program can be installed from the terminal:
```go install github.com/Breadumi/aggreGator/cmd/gator@latest```

A JSON config file `.gatorconfig` must be placed in your system root folder (found using `echo $HOME`) with the following contents:

```
{
  "db_url": "<connection string>"
}
```
The connection string will have the form `postgres://[username]:[password]@localhost:5432/[database_name]`. The database must be created manually using Postgres.

Use `gator help` for a full list of commands.
