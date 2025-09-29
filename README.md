# RCD Gopher Social

This is a Go app that simulates a social network. This is for a back-end engineering with Go
course on Udemy.

[https://udemy.com/course/backend-engineering-with-go/](https://udemy.com/course/backend-engineering-with-go/)

## Libraries

The following interesting Go libraries are used in this app:

- [Chi](https://github.com/go-chi/chi) a router for building Go services.
- [PGX](https://github.com/jackc/pgx) a PostgreSQL driver for Go.

## Tools

### Air

[https://github.com/air-verse/air](https://github.com/air-verse/air)

This project uses `air` to automatically reload the server when changes are made. To install it,
run the following command:

```bash
go install github.com/air-verse/air@latest
```

#### Configuration

Checkout `.air.toml` for the configuration. The `air` command will look for this file in the root
of the project.

Assuming you have `$GOPATH/bin` in your `$PATH`, you can run the following command to start the
server:

```bash
air
```

Each time you save a file, the server will automatically reload.

### Taskfile

[https://taskfile.dev](https://taskfile.dev)

Think of this as Go's (much more modern) version of make. See [./Taskfile.yml](./Taskfile.yml) for
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

If I symlink `.envrc` to `.env` I can also pull in the environment into `docker-compose.yml`:

```yaml
services:
  app:
    env_file:
        - .env
```

### Docker Compose

[https://docs.docker.com/compose/](https://docs.docker.com/compose/)

This project uses `docker-compose` to run the PostgreSQL database. To install it, run the following
command:

```bash
brew install docker-compose
```

To start the database, run the following command:

```bash
docker-compose up
```

### Migrate

[https://github.com/golang-migrate/migrate](https://github.com/golang-migrate/migrate)

This project uses `migrate` to manage database migrations. To install it, run the following command:

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

Run migration to upgrade:

```bash
task migrate-up
```

Run migration to downgrade:

```bash
task migrate-down
```

### Rainfrog

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
rainfrog --url $DB_ADDR
```

## Generating Self-Signed Certificates for MacOS

Instructions on how to generate the certificate using `KeyChain Access` can be found here:

[https://support.apple.com/en-ca/guide/keychain-access/kyca8916/mac](https://support.apple.com/en-ca/guide/keychain-access/kyca8916/mac)

Then run the following command to sign the binary:

```bash
codesign -f -s "RCD Local" ./bin/gopher-social --deep
```

I added this to the build step in my `Taskdev.yaml` file:

```yaml
  build:
    deps: [check-deps]
    cmds:
      - go build -o {{.out}} {{.src}}
      - codesign -f -s "RCD Local" {{.out}} --deep
    generates:
      - ./{{.out}}
```

Then I configured `.air.toml` to call the build task:

```toml
root = "."
testdata_dir = "testdata"
bin_dir = "bin"

[build]
  args_bin = []
  bin = "./bin/gopher-social"
  cmd = "task build"
```

This worked.
