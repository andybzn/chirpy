# API Documentation

## User Resource

### POST /api/users

Creates a new user

Request body:

```json
{
  "email": "example@domain.tld",
  "password": "hunter2"
}
```

Response body:

```json
{
  "id": "925dcff8-d2f8-4c65-b95e-cbb567b3e204",
  "created_at": "2025-01-01T12:00:00Z",
  "updated_at": "2025-01-01T12:00:00Z",
  "email": "example@domain.tld",
  "is_chirpy_red": false
}
```

Possible status codes:

- `201 Created` if the user was created successfully
- `500 Internal Server Error` if there was an error creating the user

### PUT /api/users

Updates the email and password of the current user

Request body:

```json
{
  "email": "exampler@domain.tld",
  "password": "hunter3"
}
```

Response body:

```json
{
  "id": "925dcff8-d2f8-4c65-b95e-cbb567b3e204",
  "created_at": "2025-01-01T12:00:00Z",
  "updated_at": "2025-03-25T16:55:00Z",
  "email": "example@domain.tld",
  "is_chirpy_red": false
}
```

Possible status codes:

- `200 OK` if the user was updated successfully
- `400 Bad Request` if the request body is invalid
- `401 Unauthorized` if the user is not authenticated
- `500 Internal Server Error` if there was an error updating the user

## Chirps Resource

### GET /api/chirps

Retrieves an un-paginated list of chirps

Example request:

```http request
GET /api/chirps
```

Response body:

```json
[
  {
    "id": "69661993-a6c7-4269-b9db-1447d32efff7",
    "created_at": "2025-01-01T12:00:00Z",
    "updated_at": "2025-01-01T12:00:00Z",
    "body": "The internet is a series of tubes...",
    "user_id": "925dcff8-d2f8-4c65-b95e-cbb567b3e204"
  }
]
```

Query parameters:

- `author_id` (optional): Filter chirps by the author ID
- `sort` (optional): Sort chirps by `created_at` date, either `asc` or `desc`

Possible status codes:

- `200 OK` if the chirps were found
- `400 Bad Request` if the `author_id` is invalid
- `500 Internal Server Error` if there was an error returning the chirps

### POST /api/chirps

Request body:

```json
{
  "body": "The internet is a series of tubes..."
}
```

Response body:

```json
{
  "id": "69661993-a6c7-4269-b9db-1447d32efff7",
  "created_at": "2025-01-01T12:00:00Z",
  "updated_at": "2025-01-01T12:00:00Z",
  "body": "The internet is a series of tubes...",
  "user_id": "925dcff8-d2f8-4c65-b95e-cbb567b3e204"
}
```

Possible status codes:

- `201 Created` if the chirp was created successfully
- `400 Bad Request` if the chirp body is invalid
- `401 Unauthorized` if the user is not authenticated
- `500 Internal Server Error` if there was an error creating the chirp

### GET /api/chirps/{chirp_id}

Retrieves a specific chirp by its ID

Example request:

```http request
GET /api/chirps/69661993-a6c7-4269-b9db-1447d32efff7
```

Response body:

```json
{
  "id": "69661993-a6c7-4269-b9db-1447d32efff7",
  "created_at": "2025-01-01T12:00:00Z",
  "updated_at": "2025-01-01T12:00:00Z",
  "body": "The internet is a series of tubes...",
  "user_id": "925dcff8-d2f8-4c65-b95e-cbb567b3e204"
}
```

Possible status codes:

- `200 OK` if the chirp was found
- `400 Bad Request` if the chirp ID is invalid
- `404 Not Found` if the chirp does not exist
- `500 Internal Server Error` if there was an error returning the chirp

### DELETE /api/chirps/{chirp_id}

Deletes a specific chirp by its ID

Example request:

```http request
DELETE /api/chirps/69661993-a6c7-4269-b9db-1447d32efff7
```

Response:

```http request
204 No Content
```

Possible status codes:

- `204 No Content` if the chirp was deleted successfully
- `400 Bad Request` if the chirp ID is invalid
- `401 Unauthorized` if the user is not authenticated
- `403 Forbidden` if the chirp does not belong to the current user
- `404 Not Found` if the chirp does not exist

## Authentication

### POST /api/login

Logs in a user and returns an access token & refresh token

Request body:

```json
{
  "email": "example@domain.tld",
  "password": "hunter2"
}
```

Response body:

```json
{
  "id": "925dcff8-d2f8-4c65-b95e-cbb567b3e204",
  "created_at": "2025-01-01T12:00:00Z",
  "updated_at": "2025-01-01T12:00:00Z",
  "email": "example@domain.tld",
  "is_chirpy_red": false,
  "access_token": "eyJzdWIiOiIxMjM0NTY3ODkwIiwi...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cC0..."
}
```

Possible status codes:

- `200 OK` if the user was logged in successfully
- `401 Unauthorized` if the email or password is incorrect, or if the user is not found
- `500 Internal Server Error` if there was an error logging in the user

### POST /api/refresh

Generates a new access token using the refresh token

Example request:

```http request
POST /api/refresh
```

Response body:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cC0..."
}
```

Possible status codes:

- `200 OK` if the access token was generated successfully
- `401 Unauthorized` if the refresh token is invalid or expired
- `500 Internal Server Error` if there was an error generating the access token

### POST /api/revoke

Revokes the refresh token for the current user

Example request:

```http request
POST /api/revoke
```

Response:

```http request
204 No Content
```

Possible status codes:

- `204 No Content` if the refresh token was revoked successfully
- `401 Unauthorized` if the refresh token is invalid or expired
- `500 Internal Server Error` if there was an error revoking the refresh token