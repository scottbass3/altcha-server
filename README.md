# Altcha Server

Self-hosted server for generating [ALTCHA](https://altcha.org/docs/get-started/) challenges and verifying solutions.

## Usage

### Run the server

From the binary:
```sh
ALTCHA_HMAC_KEY="your-secret-key-min-16-chars" bin/altcha run
```

From Docker:
```sh
docker run -e ALTCHA_HMAC_KEY="your-secret-key-min-16-chars" -p 3333:3333 ghcr.io/scottbass3/altcha-server:latest
```

From source:
```sh
ALTCHA_HMAC_KEY="your-secret-key-min-16-chars" go run ./cmd/altcha run
```

### Other commands

Generate a challenge:
```sh
ALTCHA_HMAC_KEY="..." bin/altcha generate
```

Solve a challenge (no HMAC key needed):
```sh
bin/altcha solve [CHALLENGE] [SALT]
```

Verify a solution:
```sh
ALTCHA_HMAC_KEY="..." bin/altcha verify [CHALLENGE] [SALT] [SIGNATURE] [SOLUTION]
```

## API

### `GET /health`

Returns `200 OK` when the server is up.

---

### `GET /request`

Returns a challenge for the client to solve.

**Response:**
```json
{
  "algorithm": "SHA-256",
  "challenge": "a3f1...",
  "maxNumber": 1000000,
  "salt": "abc123?expires=1748000000&",
  "signature": "d4e5..."
}
```

---

### `POST /verify`

Verifies a solved challenge. Each challenge can only be submitted once (replay protection).

**Request:**
```json
{
  "algorithm": "SHA-256",
  "challenge": "a3f1...",
  "number": 482910,
  "salt": "abc123?expires=1748000000&",
  "signature": "d4e5..."
}
```

**Response:**
```json
{ "success": true, "data": { "algorithm": "SHA-256", "challenge": "a3f1...", "number": 482910, "salt": "abc123?expires=1748000000&", "signature": "d4e5..." } }
```

Returns `400 Bad Request` for an invalid or expired solution, and `409 Conflict` if the challenge was already used.

> **Replay protection** is only complete with `ALTCHA_CHECK_EXPIRE=true`. When it is `false`, expired challenges stay valid and a used challenge is only remembered for `ALTCHA_EXPIRE` after its submission.

> **Multi-instance deployments:** the default `memory` store is per-process. Use `redis` or `memcached` to share nonce state across instances. Keys are prefixed with `altcha:nonce:`.

---

### `POST /verify-fields`

Verifies that specific form fields haven't been tampered with, using a hash committed at challenge time.

**Request:**
```json
{
  "formData": {
    "email": ["user@example.com"],
    "message": ["Hello"]
  },
  "fields": ["email", "message"],
  "fieldsHash": "e3b0c4..."
}
```

**Response:**
```json
{ "success": true }
```

---

### `POST /verify-server-signature`

Verifies a signed payload issued by ALTCHA's verification service. The signature is checked against `ALTCHA_HMAC_KEY`, so it must be set to the secret used to sign these payloads.

**Request:**
```json
{
  "algorithm": "SHA-256",
  "verificationData": "verified=true&expire=1748000000&...",
  "signature": "f9a2...",
  "verified": true
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "verified": true,
    "expire": 1748000000,
    "score": 0.9,
    "classification": "human",
    "email": "",
    "fields": ["email", "message"],
    "fieldsHash": "e3b0c4..."
  }
}
```

All `POST` endpoints reject bodies larger than 64 KiB with `413 Request Entity Too Large`.

## Environment variables

| Name                      | Description                                              | Default  | Required |
|---------------------------|----------------------------------------------------------|----------|----------|
| `ALTCHA_HMAC_KEY`         | HMAC key used to sign challenges (min 16 characters)     |          | Yes      |
| `ALTCHA_PORT`             | Server listen port                                       | `3333`   | No       |
| `ALTCHA_BASE_URL`         | URL prefix for all endpoints (e.g. `/altcha`)            |          | No       |
| `ALTCHA_ALGORITHM`        | Hash algorithm: `SHA-1`, `SHA-256`, or `SHA-512`         | `SHA-256`| No       |
| `ALTCHA_MAX_NUMBER`       | Max iterations for proof-of-work (controls difficulty)   | `1000000`| No       |
| `ALTCHA_SALT`             | Fixed salt for challenges (random per-challenge if unset)|          | No       |
| `ALTCHA_SALT_LENGTH`      | Random salt length in bytes                              | `12`     | No       |
| `ALTCHA_EXPIRE`           | Challenge TTL as a Go duration (`10m`, `1h`) or seconds (`600`) | `600s`   | No       |
| `ALTCHA_CHECK_EXPIRE`     | Reject expired challenges (see replay protection note)   | `true`   | No       |
| `ALTCHA_CORS_ORIGINS`     | `*` or comma-separated list of allowed origins           | `*`      | No       |
| `ALTCHA_STORE`            | Nonce store backend: `memory`, `redis`, or `memcached`  | `memory` | No       |
| `ALTCHA_REDIS_URL`        | Redis connection URL (required when `ALTCHA_STORE=redis`) |        | No       |
| `ALTCHA_MEMCACHED_SERVERS` | Comma-separated `host:port` list (required when `ALTCHA_STORE=memcached`) | | No |
| `ALTCHA_DEBUG`            | Enable debug logging                                     | `false`  | No       |
| `ALTCHA_DISABLE_VALIDATION` | Skip solution verification (development only)          | `false`  | No       |

## Build

```sh
# Binary
make build

# Docker image
make build-image

# Publish Docker image (requires VERSION=vX.Y.Z)
make release-image VERSION=v1.0.0
```
