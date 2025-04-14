# WT app (backend) 

Description here.

## Preparing

Create a Postgres database on your local PG server.

## Development

Start the app using the command:
```bash
go run .\cmd\api -pgLogin "your_pg_login" -pgPassword "your_pg_password" -dbName "your_db_name" -port 4000
```
Required parameters: pgLogin, pgPassword, dbName.
Port is an optional parameter.
