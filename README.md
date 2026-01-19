# RCD Gopher Social

This is a Go app that simulates a social network. This is for a back-end engineering course with Go
on Udemy. This is not meant to be a production-ready app, but rather a learning exercise. I
highly recommend the course if you want to learn more about building back-end services with Go.

[Backend Engineering with Go Udemy Course](https://udemy.com/course/backend-engineering-with-go/)

## Pre-requisites

- Go 1.25 or later
- Podman or Docker (examples use Podman)
- Homebrew or another package manager to install the tools mentioned below

## Libraries

The following interesting Go libraries are used in this app:

- [Chi](https://github.com/go-chi/chi) a router for building Go services.
- [PGX](https://github.com/jackc/pgx) a PostgreSQL driver for Go.

I used PGX because the older pq library is no longer being maintained. PGX is a modern
PostgreSQL driver that is actively maintained. I also like that I could figure out how to use
the newer library on my own without just blindly following along with the course.

## Flake setup

Create 2 env files. `.envrc` and `.envc`. Put the following in `.envrc`:

```bash
use flake
source .env
```

Put the variables from the Direnv section below in `.env`.

When you run `direnv allow` it will setup everything for you from the flake.

## Tools

### Air

This project uses `air` to automatically reload the server when changes are made. To install it,
run the following command:

[https://github.com/air-verse/air](https://github.com/air-verse/air)

```bash
go install github.com/air-verse/air@latest
```

#### Configuration

Check out `.air.toml` for the configuration. The `air` command will look for this file in the root
of the project.

Assuming you have `$GOPATH/bin` in your `$PATH`, you can run the following command to start the
server:

```bash
air
```

Each time you save a file, the server will automatically reload.

### Taskfile

[https://taskfile.dev](https://taskfile.dev)

Think of this as Go's (much more modern) version of Make. See [./Taskfile.yml](./Taskfile.yml) for
the available commands.

To install it, run the following command:

```bash
brew install go-task/tap/go-task
```

To build:

```bash
task build
```

To run the tests:

```bash
task test
```

### Direnv

[https://github.com/direnv/direnv](https://github.com/direnv/direnv)

This project uses `direnv` to manage environment variables. To use it, create a `.envrc` file in
the root of the project with the following content:

```bash
export ADDR=":3000"
export DATABASE_URL="postgres://admin:adminpassword@localhost/social?sslmode=disable"
```

Then, run the following command to allow the `.envrc` file:

```bash
direnv allow
```

### Docker/Podman Compose

- [https://github.com/containers/podman-compose](https://github.com/containers/podman-compose)
- [https://docs.docker.com/compose/](https://docs.docker.com/compose/)

This project uses either Docker compose or Podman compose to run a PostgreSQL database. Check out
instructions on the website to install it. Let's assume you have `podman-compose` installed.

Because the flake code can break the env file for docker we use `.env` for docker/podman compose.
If you are not using a flake create a symlin from `.envrc` to `.env`.

```bash
ln -s .envrc .env
```

To start the database, run the following command:

```bash
podman-compose up
```

### Migrate

This project uses `migrate` to manage database migrations. To install it, run the following command:

[https://github.com/golang-migrate/migrate](https://github.com/golang-migrate/migrate)

```bash
brew install golang-migrate
```

The [Taskfile.yml](Taskfile.yml) file has some commands to help with migrations. The Taskfile
also contains the commands for reference.

Create migrations (found in ./cmd/migrate/migrations):

```bash
task migrate-create-users
task migrate-create-posts
```

To add a migration other than the base one here is example syntax:

```bash
migrate create -seq -ext sql -dir ./cmd/migrate/migrations/ alter_posts_add_version
```

This will create the up and down SQL files in the migrations directory.

Run migration to upgrade:

```bash
task migrate-up
```

Run migration to downgrade:

```bash
task migrate-down
```

To delete the database and start over:

```bash
task migrate-drop
task migrate-up
task seed-db
```

NOTE: Please add a user or at least set `is_active` to true for any users you want to test with as
you can only do token authentication with users who are active.

### Rainfrog

Rainfrog is a TUI application to view and manage a PostgreSQL database.

[https://github.com/achristmascarl/rainfrog](https://github.com/achristmascarl/rainfrog)

Right now Rust has to be installed first. Also `$HOME/.cargo/bin` has to be in the PATH.

```bash
brew install rustup
rustup-init
```

Then install rainfrog:

```bash
cargo install rainfrog
```

To run Rainfrog, use the following command:

```bash
rainfrog --url $DATABASE_URL
```

### Tredis

Tredis is a Rust based TUI client for Redis.

[https://github.com/huseyinbabal/tredis](https://github.com/huseyinbabal/tredis)

To install:

```bash
brew install huseyinbabal/tap/tredis
```

To run tredis:

```bash
tredis --host localhost --port 6379 --db 1
```

## Seeding the Database for Testing

To seed the database with test data, run the following command:

```bash
task seed-db
```

## Zellij Script

This project is set up to use `zellij` to manage terminal panes. To install it, run the following
command:

```bash
brew install zellij
```

To start the zellij session, run the following command:

```bash
./scripts/zellij.sh
```

This will start a zellij session with the following panes:

- vim - to edit the code
- shell - to run other commands like task build, tests, etc.
- air - to automatically reload the server when changes are made
- posting - TUI app for making API requests
- rainfrog - TUI app for managing the PostgreSQL database
- podman-compose - to manage the PostgreSQL database container

## Generating Self-Signed Certificates for MacOS

Instructions on how to generate the certificate using `KeyChain Access` can be found here:

[https://support.apple.com/en-ca/guide/keychain-access/kyca8916/mac](https://support.apple.com/en-ca/guide/keychain-access/kyca8916/mac)

Then run the following command to sign the binary:

```bash
codesign -f -s "RCD Local" ./bin/gopher-social --deep
```

Add this to the build step in my `Taskfile.yml` file:

```yaml
build:
  deps: [check-deps]
  cmds:
    - go build -o {{.out}} {{.src}}
    - codesign -f -s "RCD Local" {{.out}} --deep
  generates:
    - ./{{.out}}
```

Configure `.air.toml` to call the build task:

```toml
root = "."
testdata_dir = "testdata"
bin_dir = "bin"

[build]
  args_bin = []
  bin = "./bin/gopher-social"
  cmd = "task build"
```

## Notes

- For adding a migration with a foreign key to a database that has existing data please see
  [000013_alter_users_with_roles.up.sql](./cmd/migrate/migrations/000013_alter_users_with_roles.up.sql).
