# Finance Tracker - Backend
This repository serves as backend for the Finance Tracker project, which offers users a simplified platform for effective financial management.  From recording expenses and income to managing lent and borrowed amounts, every transaction is easily tracked and categorised. Users can customise monthly budgets, receive notifications for upcoming payments, and set thresholds for expense categories.

## Features
- Manage income and expenses
- Track lent and borrowed amounts
- Set and monitor monthly budgets
- Receive payment notifications
- Set expense category thresholds


## Requirements :
* Go
* Makefile
* Postgres
* Lint

## Run server
```
git clone https://github.com/Thivyasree-Rajaraman/finance-tracker-backend.git
cd backend
go mod tidy
make run
```

## Run lint
```
make lint
```


## Installation
1. Clone the repository and navigate into the project directory:
    ```bash
    git clone https://github.com/Thivyasree-Rajaraman/finance-tracker-backend.git
    cd backend
    go mod tidy
    ```

## Configuration
2. Create a `.env` file in the root directory and add the necessary environment variables:
    ```makefile
    DATABASE_URL=your_db_url
    GOOGLE_CLIENT_ID=your_client_id
    GOOGLE_CLIENT_SECRET=your_client_secret
    GOOGLE_REDIRECT_URL=http://localhost:3000
    GOOGLE_USER_INFO_ENDPOINT=https://www.googleapis.com/oauth2/v2/userinfo
    SECRET_KEY=your_secret_key
    ```

## Database Setup
3. Ensure Postgres is running and execute the following commands to set up the database:
    ```sql
    CREATE DATABASE finance_tracker;
    ```

## Run Server
4. Start the server using:
    ```bash
    make run
    ```

## Run Lint
5. To check the code quality and linting issues, run:
    ```bash
    make lint
    ```

## Run Tests
6. To run tests:
    ```bash
    make test
    ```
