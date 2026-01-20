# Altcha Server

Altcha server for generating challenges and validating solutions.

## Usage

### Run the server
From the binary:
```sh
$ ALTCHA_HMAC_KEY="HMAC_KEY" bin/altcha run
```

From the Docker image:
```sh
$ docker run -e ALTCHA_HMAC_KEY="HMAC_KEY" -p 3333:3333 ghcr.io/scottbass3/altcha-server:latest
```

From source:
```sh
$ ALTCHA_HMAC_KEY="HMAC_KEY" go run ./cmd/altcha run
```

Once the server is running, visit `http://localhost:3333/request` to request a challenge.
Post the solution to `http://localhost:3333/verify`.

Altcha documentation: https://altcha.org/fr/docs/get-started/

### Other commands
Generate a challenge:
```sh
$ ALTCHA_HMAC_KEY="HMAC_KEY" bin/altcha generate
```

Solve a challenge:
```sh
$ ALTCHA_HMAC_KEY="HMAC_KEY" bin/altcha solve [CHALLENGE] [SALT]
```

Verify a solution:
```sh
$ ALTCHA_HMAC_KEY="HMAC_KEY" bin/altcha verify [CHALLENGE] [SALT] [SIGNATURE] [SOLUTION]
```

## Environment variables
| Name                | Description                                                      | Default   | Required |
|---------------------|------------------------------------------------------------------|-----------|----------|
| ALTCHA_BASE_URL     | Base URL prefix for endpoints                                   |           | No       |
| ALTCHA_PORT         | Server listen port                                               | 3333      | No       |
| ALTCHA_HMAC_KEY     | HMAC key used to sign challenges                                 |           | Yes      |
| ALTCHA_MAX_NUMBER   | Max iterations for solving a challenge (difficulty)              | 1000000   | No       |
| ALTCHA_ALGORITHM    | Hash algorithm (SHA-1, SHA-256, SHA-512)                         | SHA-256   | No       |
| ALTCHA_SALT         | Force a fixed salt for challenges                                |           | No       |
| ALTCHA_EXPIRE       | Challenge expiration (seconds)                                   | 600       | No       |
| ALTCHA_CHECK_EXPIRE | Whether to check challenge expiration                            | 1         | No       |

## Build the binary
```sh
$ make build
```

## Build the Docker image
```sh
$ make build-image
```

## Publish the Docker image
```sh
$ make release-image
```
