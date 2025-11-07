# Go Auth Microservice

![Go CI](https://github.com/isubhampadhi56/go-auth-microservice/actions/workflows/ci.yml/badge.svg)
![Go Version](https://img.shields.io/badge/go-1.25-blue.svg)
![License](https://img.shields.io/github/license/isubhampadhi56/go-auth-microservice.svg)

# API Documentation

This document provides detailed instructions on how to interact with the application API for user authentication and token management.

## Getting Started

The main application file is located under the `cmd` directory. To start the application, run the following command:

```bash
go run cmd/main.go
```

The application uses SQLite as the database to store user data. All configuration files are already present in the `pkg/config` directory.

---

## 1. User Signup

The signup request requires a JSON object in the body with `email` and `password` attributes. The email must be in a valid format, and the password must be at least 8 characters long. 

### Endpoint: `POST /api/v1/auth/signup`

**Request Body**:
```json
{
    "email": "user1@mail.com",
    "password": "123456789"
}
```

**Validation**:
- The email should be in a valid email format.
- The password should be at least 8 characters long.

### Example Request (using `curl`):
```bash
curl --location 'http://localhost:8080/api/v1/auth/signup' --header 'Content-Type: application/json' --data-raw '{
    "email": "user1@mail.com",
    "password": "123456789"
}'
```

---

## 2. User Login

The login request requires a JSON object in the body with `email` and `password` attributes. Upon successful login, the response will include both `accessToken` and `refreshToken`.

### Endpoint: `POST /api/v1/auth/login`

**Request Body**:
```json
{
    "email": "user1@mail.com",
    "password": "123456789"
}
```

### Example Request (using `curl`):
```bash
curl --location 'http://localhost:8080/api/v1/auth/login' --header 'Content-Type: application/json' --data-raw '{
    "email": "user1@mail.com",
    "password": "123456789"
}'
```

**Response**:
```json
{
    "accessToken": "<access_token_here>",
    "refreshToken": "<refresh_token_here>"
}
```

---

## 3. Authorization Token and Protected Routes

To access protected routes, the request should include a valid `accessToken` in the `Authorization` header.

### Protected Routes:

- `GET /api/v1/me`: Checks if the user's `accessToken` is still valid.
- `GET /api/v1/user`: Returns the details of the logged-in user.
- `PATCH /api/v1/user/deactivate`: Deactivates the user account.
- `PATCH /api/v1/user/changePassword`: Change the user password.

### Example Request (Get User Details):
```bash
curl --location 'http://localhost:8080/api/v1/me' --header 'Authorization: <access_token_here>'
```

```bash
curl --location 'http://localhost:8080/api/v1/user' --header 'Authorization: <access_token_here>'
```

---

## 4. Revocation of Token

The `PATCH /api/v1/deactivate` endpoint will deactivate the logged-in user account. 

- After deactivation, the `accessToken` expires within 5 minutes. For instant invalidation it's value stored in a in-memory cache to blacklist the access token.
- The user will no longer be able to generate new tokens or refresh them using the `refreshToken`.
- Any existing authorization tokens will be invalidated.

### Endpoint: `PATCH /api/v1/deactivate`

### Example Request (using `curl`):
```bash
curl --location --request PATCH 'http://localhost:8080/api/v1/deactivate' --header 'Authorization: <access_token_here>'
```
---

The `PATCH /api/v1/changePassword` endpoint will change password of the logged-in user account. 

- After changing password, the `accessToken` expires within 5 minutes. For instant invalidation it's value stored in a in-memory cache to blacklist the access token.
- The user will no longer be able to generate new tokens or refresh them using the `refreshToken`.
- Any existing authorization tokens will be invalidated.

### Endpoint: `PATCH /api/v1/changePassword`

### Example Request (using `curl`):
```bash
curl --location --request PATCH 'http://localhost:8080/api/v1/changePassword' \
--header 'Authorization: <access_token_here>' \
--header 'Content-Type: application/json' \
--data '{
    "password": "123456789"
}'
```
---

## 5. Refresh Access Token

Access tokens expire in 5 minutes, while refresh tokens expire in 72 hours (configured in the `pkg/config`). To renew an `accessToken`, send a GET request to the `api/v1/auth/token` route with the valid `refreshToken`.

### Endpoint: `GET /api/v1/auth/token`

**Request Header**:
```plaintext
RefreshToken: <refresh_token_here>
```

**Response**:
```json
{
    "accessToken": "<new_access_token_here>"
}
```

If the refresh token is valid, and the user is active, a new `accessToken` will be generated. If the user's account has been deactivated or updated after the `refreshToken` was issued, the request will fail, prompting the user to log in again.

### Example Request (using `curl`):
```bash
curl --location 'http://localhost:8080/api/v1/auth/token' --header 'RefreshToken: <refresh_token_here>'
```
---