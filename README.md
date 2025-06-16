# Chirpy Backend

Chirpy is a lightweight social media backend service built in Go with PostgreSQL and JWT-based authentication. This project supports user authentication, chirp posting, token refresh/revocation, and a subscription upgrade mechanism.

---

## Table of Contents

* [Setup](#setup)
* [Environment Variables](#environment-variables)
* [Endpoints](#endpoints)

  * [User Management](#user-management)
  * [Authentication](#authentication)
  * [Chirp Management](#chirp-management)
  * [Admin Utilities](#admin-utilities)
  * [Chirpy Red Webhook](#chirpy-red-webhook)
* [Types](#types)

---

## Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/hursty1/chirpy
   cd chirpy
   ```

2. Set up a `.env` file with the required keys.

3. Run the server:

   ```bash
   go run main.go
   ```

---

## Environment Variables

```env
DB_URL=postgres://user:password@localhost:5432/chirpy\PLATFORM=dev
TOKENSECRET=your_secret_key
POLKA_KEY=your_polka_api_key
```

---

## Endpoints

### User Management

#### `POST /api/users`

Creates a new user.

**Request:**

```json
{
  "email": "user@example.com",
  "password": "strongpassword"
}
```

**Response:** `201 Created`

```json
{
  "id": "uuid",
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "email": "user@example.com",
  "is_chirpy_red": false
}
```

#### `PUT /api/users`

Updates user email/password.

**Auth Required:** Bearer token

**Request:**

```json
{
  "email": "new@example.com",
  "password": "newpassword"
}
```

**Response:** `200 OK`
User object (same as above).

---

### Authentication

#### `POST /api/login`

Authenticates a user.

**Request:**

```json
{
  "email": "user@example.com",
  "password": "password"
}
```

**Response:** `200 OK`

```json
{
  "token": "JWT",
  "refresh_token": "token",
  "id": "uuid",
  "email": "user@example.com",
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "is_chirpy_red": false
}
```

#### `POST /api/refresh`

Renews an access token.

**Auth:** `Authorization: Bearer <refresh_token>`

**Response:** `200 OK`

```json
{
  "token": "new_jwt"
}
```

#### `POST /api/revoke`

Revokes a refresh token.

**Auth:** `Authorization: Bearer <refresh_token>`

**Response:** `204 No Content`

---

### Chirp Management

#### `GET /api/chirps`

Returns all chirps.

**Optional Query Params:**

* `author_id=uuid` (filter by author)
* `sort=desc|asc` (sort by date)

**Response:** `200 OK`

```json
[
  {
    "id": "uuid",
    "created_at": "timestamp",
    "updated_at": "timestamp",
    "body": "chirp text",
    "user_id": "uuid"
  }
]
```

#### `POST /api/chirps`

Create a new chirp.

**Auth Required:** Bearer token

**Request:**

```json
{
  "body": "hello chirpy"
}
```

**Response:** `201 Created`
Chirp object.

#### `GET /api/chirps/{chirpID}`

Fetch chirp by ID.

**Response:** `200 OK`
Chirp object.

#### `DELETE /api/chirps/{chirpID}`

Delete chirp by ID (must be owner).

**Auth Required:** Bearer token

**Response:** `204 No Content`

---

### Admin Utilities

#### `GET /api/healthz`

Returns `OK` to indicate health.

#### `POST /admin/reset`

**Dev only.** Resets metrics and deletes users.

#### `GET /admin/metrics`

Returns HTML page showing number of hits to `/app/*` routes.

---

### Chirpy Red Webhook

#### `POST /api/polka/webhooks`

Webhook endpoint to upgrade user to `is_chirpy_red=true`.

**Auth Required:** `Authorization: ApiKey <key>`

**Request Body:**

```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "uuid"
  }
}
```

**Response:** `204 No Content`

---

## Types

### `UserResponse`

```json
{
  "id": "uuid",
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "email": "user@example.com",
  "is_chirpy_red": false
}
```

### `ResponseChirp`

```json
{
  "id": "uuid",
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "body": "chirp text",
  "user_id": "uuid"
}
```

---

## License

This project is provided as-is for educational or personal use.
