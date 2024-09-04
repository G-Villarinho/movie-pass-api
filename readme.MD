# **🎬 Movie Pass API 🎟️**

 ###### Welcome to the Movie Pass API! This API was created to simplify movie ticket purchasing, manage cinemas, theaters, movies, and users. Everything is ready for you to deliver an amazing seat reservation and ticket buying experience! 🍿

 [![Go Report Card](https://goreportcard.com/badge/github.com/GSVillas/movie-pass-api)](https://goreportcard.com/report/github.com/GSVillas/movie-pass-api) [![Go Version](https://img.shields.io/github/go-mod/go-version/GSVillas/movie-pass-api)](https://golang.org) [![Issues](https://img.shields.io/github/issues/GSVillas/movie-pass-api)](https://github.com/GSVillas/movie-pass-api/issues)

> [!WARNING]
>
>API is currently under development. Expect potential changes and unstable behavior.



## 🚀 Main Features

- 🎟️ **Ticket Purchase**: Simplify the online ticket buying process.
- 🏙️ **Cinema Management**: Add and update cinemas and their respective theaters with custom seat arrangements.
- 📽️ **Movie Catalog**: Manage the catalog of available movies.
- 📅 **Real-Time Seat Reservation**: View and reserve available seats in real time.
- 👤 **User Authentication**: Complete system for registration, sign-in, and authentication, including session management with Redis.
- 🛡️ **JWT Authentication**: All routes are secured with JWT-based authentication for enhanced security.

## 📋 Installation

To get started with the **Movie Pass API**, follow these steps:

### Requirements

Ensure you have the following installed:

- **Go** (version 1.22 or higher)
- **MySQL** (or any other compatible database)
- **Redis** (for caching and session management)
- **Docker** (optional, but recommended for deployment)

### Steps

1. **Clone the repository**:
   ```bash 
   git clone https://github.com/GSVillas/movie-pass-api.git
   ```
    
2. **Install the dependencies**:
   ```bash 
   go mod tidy
   ```
3. **Set up the environment variables in the .env file**:
   ```bash 
   Create a .env based on .en.example
   ```

### ⚙️ Makefile Instructions
 > To simplify your development process, we have provided a Makefile with useful commands to manage keys, linting, running the application, and migrations. 

 #### Requirements
* Generate Keys: Generates private and public ECDSA keys:
     ```bash 
   make generate-keys
   ```
* Lint the Code: Runs the linter to check for code issues:
     ```bash 
   make lint
   ```
* Clean Keys: Deletes the generated private and public keys:
     ```bash 
   make clean
   ```
* Run the Application: Runs the Go application:
     ```bash 
   make rub
   ```
* Run Migrations: Executes database migrations:
     ```bash 
   make migration
   ```

* Run Test: Executes unity tests:
     ```bash 
   make test
   ```

* Run e2e Test: Executes end to end:
     ```bash 
   make e2e
   ```
