# Customer Service

This is a microservice project for managing customer information.

## Table of Contents

- [Introduction](#customer-service)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Contributing](#contributing)
- [License](#license)

## Getting Started

### Prerequisites

- Go installed
- [BoltDB](https://github.com/boltdb/bolt) (embedded key/value database for Go)

### Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/sahitya-pavurala/customer-service.git
    ```

2. Change into the project directory:

    ```bash
    cd customer-service
    ```

3. Run the application:

    ```bash
    go run main.go
    ```

## Usage

The microservice exposes REST API endpoints for managing customers and accounts. Below are some example API endpoints:

- **Register as a user to use the API:**
  ```http
  POST /register

- **Get access token to use the API:**
  ```http
  POST /login
  
- **Get All Customers:**
  ```http
  GET /customers
  
- **Get Customer by id:**
  ```http
  GET /customers/:id

- **Create Customer:**
  ```http
  POST /customers

- **Create Account for custimer:**
  ```http
  POST /customers/:id/accounts
