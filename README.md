# Trainings API

## Description

Training application developed with [Gin](https://gin-gonic.com/), [GORM](https://gorm.io/index.html), [Go Redis](https://redis.uptrace.dev/) and [Gorilla WebSocket](https://pkg.go.dev/github.com/gorilla/websocket).
Work in progress...

## Environment Variables

Create a `.env` file in the root directory. And add these default values:

```
DATABASE_URL=postgres://postgres:password@postgres:5432/trainings
REDIS_URL=redis://:password@redis:6379/0
API_SECRET=SecretSecretSecret
TOKEN_HOUR_LIFESPAN=12
ENVIRONMENT=dev
PORT=8080
CORS_ORIGIN=*

RATE_LIMIT=30
RATE_WINDOW=1s
LOGIN_ATTEMPTS=5
LOGIN_WINDOW=10m
OTP_EXPIRY=10m
```

If you want to use SMTP for one time password emails. Add your SMTP credentials:

```
SENDER_NAME=Only Ada 💬
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
EMAIL=mail@gmail.com
EMAIL_PASSWORD=password
```

## Run Application with Docker

More information about [Docker](https://www.docker.com/).
To run the application type this command in the root folder.

```bash
$ docker compose up
```

You might have to run this command twice if it doesn't work the first time :)

## API Endpoints

For authentication are used bearer tokens.

- **User** `/api/users`

  - [GET] `/` - Get all users

  - [GET] `/:id` - Get user by ID

  - [GET] `/current` - Get current user

  - [POST] `/send-otp` - Send OTP email

    ```json
    {
      "email": "firstlast@mail.com"
    }
    ```

  - [POST] `/register` - Register user

    ```json
    {
      "firstName": "First",
      "lastName": "Last",
      "email": "firstlast@mail.com",
      "password": "password"
    }
    ```

  - [POST] `/register-otp` - Register user with OTP

    ```json
    {
      "firstName": "First",
      "lastName": "Last",
      "email": "firstlast@mail.com",
      "password": "password",
      "otp": "123456"
    }
    ```

  - [POST] `/login` - Login user

    ```json
    {
      "email": "firstlast@mail.com",
      "password": "password"
    }
    ```

  - [POST] `/login-otp` - Login user with OTP

    ```json
    {
      "email": "firstlast@mail.com",
      "password": "password",
      "otp": "123456"
    }
    ```

  - [GET] `/email/:email` - Search user by email

  - [PUT] `/update` - Update user

    ```json
    {
      "firstName": "First",
      "lastName": "Last",
      "email": "firstlast@mail.com"
    }
    ```

  - [PUT] `/update-otp` - Update user with OTP

    ```json
    {
      "firstName": "First",
      "lastName": "Last",
      "email": "firstlast@mail.com",
      "otp": "123456"
    }
    ```

  - [POST] `/:id/roles/:role` - Add role to user

  - [DELETE] `/:id/roles/:role` - Remove role from user

- **Training** `/api/trainings`

  - [GET] `/` - Get all trainings

  - [GET] `/:id` - Get training by ID

  - [POST] `/` - Create training

    ```json
    {
      "name": "training"
    }
    ```

  - [PUT] `/:id` - Update training by ID

    ```json
    {
      "name": "updated training"
    }
    ```

  - [DELETE] `/:id` - Delete training by ID

  - [POST] `/:training_id/users/:user_id` - Add user to training

  - [DELETE] `/:training_id/users/:user_id` - Remove user from training

- **Post** `/api/posts`

  - [GET] `/:training_id?page=1&size=10` - Get paginated posts by training ID

  - [POST] `/:training_id` - Create post

    ```json
    {
      "text": "Hello World!"
    }
    ```

  - [PUT] `/:id` - Update post by ID

    ```json
    {
      "text": "Goodbye World!"
    }
    ```

  - [DELETE] `/:id` - Delete post by ID
