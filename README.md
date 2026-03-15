# CMPS3162 Semester Project

## Hotel Room Booking & Housekeeping Management System

| Key               | Value                                          |
| ----------------- | ---------------------------------------------- |
| **Student Name**  | [Andres Hung](https://github.com/andreshungbz) |
| **Student Email** | 2018118240@ub.edu.bz                           |
| **Course**        | CMPS3162 - Advanced Databases                  |
| **Due Date**      | March 19, 2026                                 |

### Test 1 Deliverables

- Entity-Relationship Diagrams (ERD) can be found in the `docs` folder located at https://github.com/andreshungbz/cmps3162-project/tree/main/docs.
- Google Slides presentation can be found at {TODO}.

### Running the Application

#### Docker Compose

```
docker compose up
```

#### Manual Method

##### Pre-requisites

- make
- curl
- golang-migrate

##### Database Setup

```
CREATE role hotel_user WITH LOGIN PASSWORD 'hotel_password';
CREATE DATABASE hotel;
ALTER DATABASE hotel OWNER TO hotel_user;
```

##### Application Setup

```
cp .envrc.example .envrc
make db/migrations/up
make run
```
