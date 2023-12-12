# Billing system

## Description

This repository is for managing batches, orders and charges for IoT products
## Features
* Batch Management: Ability to create, update, and manage batches of IoT products.
* Order Processing: Handles the complete order process from placement to delivery.
* Charge Calculation: Accurately calculates charges for each IoT product or service.
* Cron tasks for monthly billing
* Async email notification
* Swagger documentation
* Rate limiter to protect server from DOS attacks
* Logging and performance monitoring using Loki, Prometheus & Grafana stack

## Installation
To run the application

```bash
docker compose up
```

## How to generate code

- Generate SQL CRUD with sqlc:

    ```bash
    make sqlc
    ```
- Migrate to latest database schema
    ```bash
    make migrate
    ```

- Generate DB mock with gomock:

    ```bash
    make mock
    ```

- Create a new db migration:

    ```bash
    make new_migration name=<migration_name>
    ```

## Swagger documentation
```
http://localhost:8080/swagger/index.html#/
```
## Testing
To run all unit testing and see code coverage
- Start redis container:
```bash
make redis
```
- Start postgres container:
```bash
make postgres
```
- Run all unit tests

```bash
make test
```
