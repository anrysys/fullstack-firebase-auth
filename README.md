# Fullstack on Next.js (frontend) / Fiber Go-lang (backend) with Firebase Authentication and App Check and ReCaptcha v3

## Table of Contents

1. [Introduction](#introduction)
2. [Prerequisites](#prerequisites)
3. [Setup](#setup)
4. [Running the Application](#running-the-application)
5. [Database Migrations (PostgreSQL)](#database-migrations-postgresql)
6. [Start Backend Server (Fiber Go-lang)](#start-backend-server-fiber-go-lang)
7. [Start Frontend Server (NextJs)](#start-frontend-server-nextjs)
8. [Technologies Used](#technologies-used)
9. [Contact](#contact)

## Introduction

`Next Firebase Authentication` is a project that demonstrates how to implement user authentication in a Next.js application using `Firebase Authentication`. It provides a complete user authentication flow, including user `registration`, `login`, `reset password`, `profile management`, `password change` and `forgot password`.

## Features

- **User registration and login**: Users can register and login with email and password.
- **Profile management**: Users can update their profile information, including name and photo URL.
- **Password change**: Users can change their password.
- **Firebase integration**: The application uses Firebase for authentication and data storage.
- **TypeScript**: The codebase is written in TypeScript for type safety and better maintainability.
- **Next.js**: A React framework for building JavaScript applications with server-side rendering and static site generation.
- **Firebase Authentication**: A service that provides backend services, easy-to-use SDKs, and ready-made UI libraries to authenticate users to your app.
- **Firebase App-Check**:
  - Enter in your Firebase <https://console.firebase.google.com> (enter your project ID in section `App Check`) exp: `https://console.firebase.google.com/project/_YOUUR-PROJECT-ID_/appcheck/products`
  - Create a new ReCaptcha v3: <https://www.google.com/recaptcha/>
- **React**: A JavaScript library for building user interfaces.

## Prerequisites

Before you begin, ensure you have met the following requirements:

- You have installed Docker and Docker Compose.

- You have installed Go.
- You have a `<Windows/Linux/Mac>` machine.
- IMPORTANT!!! Make sure your computer is not using ports 3000, 5432, 6379, and 8000

## Setup

To setup the project, follow these steps:

1. Clone the repository

```shell
git clone https://github.com/anrysys/next-firebase-authentication.git
```

2. Change file name from example.env.local to .env.local:

```shell
mv example.env.local .env.local
```

3. Change a .env file(s) to match your project:

```shell
cd next-firebase-authentication
nano ./backend/.env
nano .env.local
```

Attention! It is mandatory to change in the `./backend/.env` file all the values for the keys:

- `ACCESS_TOKEN_PRIVATE_KEY=INSERT NEW PRIVATE KEY HERE`
- `ACCESS_TOKEN_PUBLIC_KEY=INSERT NEW PUBLIC KEY HERE`
- `REFRESH_TOKEN_PRIVATE_KEY=INSERT NEW PRIVATE KEY HERE`
- `REFRESH_TOKEN_PUBLIC_KEY=INSERT NEW PUBLIC KEY HERE`

To do this, you need to generate an RSA key pair (4096 is recommended, but you can use 1024 or 2048, which will decrease security but save a bit on traffic). You can generate new keys using utilities such as OpenSSL or ssh-keygen, or online at: <https://it-tools.tech/rsa-key-pair-generator>

4. Navigate to the project directory:

```shell
cd next-firebase-authentication
npm install
```

5. Go to ./backend and run:

```shell
go init
go mod init github.com/anrysys/next-firebase-authentication
go mod tidy
```

## Running the Application

To run the application, execute the following command:

Start the service`s container (Postgresql and Redis):

```shell
docker-compose up
```

## Database Migrations (PostgreSQL)

We use [golang-migrate](https://github.com/golang-migrate/migrate) for handling database migrations.

To create a new migration, run:

```shell
docker-compose run --rm migrate create -ext sql -dir /migrations create_customer_table
```

To apply migrations, run:

```shell
make migrate.run
```

To rollback migrations, run:

```shell
make migrate.down
```

## Start Backend Server (Fiber Go-lang)

To start the backend server, execute the following command:

```shell
cd ./backend
go run main.go
-- OR --
air
```

## Start Frontend Server (NextJs)

To start the frontend server, execute the following command:

```shell
cd next-firebase-authentication
npm run dev
```

## Technologies Used

This project uses the following technologies:

- [Next.js](https://nextjs.org)
- [Docker](https://www.docker.com)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Go](https://golang.org)
- [PostgreSQL](https://www.postgresql.org)
- [Redis](https://redis.io)
- [Firebase](https://firebase.google.com)
- [TypeScript](https://www.typescriptlang.org)
- [React](https://reactjs.org)
- [Fiber Go-lang](https://github.com/gofiber/fiber)

## Contact

If you want to contact me you can reach me at `<anrysys@gmail.com>`.
