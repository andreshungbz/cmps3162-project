# CMPS3162 Semester Project

## Hotel Room Booking & Housekeeping Management System

| Key               | Value                                          |
| ----------------- | ---------------------------------------------- |
| **Student Name**  | [Andres Hung](https://github.com/andreshungbz) |
| **Student Email** | 2018118240@ub.edu.bz                           |
| **Course**        | CMPS3162 - Advanced Databases                  |
| **Due Date**      | April 7, 2026                                  |

### Test 1 Deliverables

- Entity-Relationship Diagrams (ERD) can be found in the `docs` folder located at https://github.com/andreshungbz/cmps3162-project/tree/main/docs.
- Google Slides presentation can be found at https://docs.google.com/presentation/d/1tE8GPZKjBMg3du9lbhShSSBfoenZ_GDl0b08cgz5pno/edit?usp=sharing.

### Authentication/Authorization/Mail Demo Video

https://youtu.be/CshMb9krVZg

### Running the Application

> [!NOTE]
> You must provide your own credentials for `MAILTRAP_SMTP_USERNAME` and `MAILTRAP_SMTP_PASSWORD` in the `.envrc` file in order to enable email functionality. This is primarily for activation emails when creating an employee.

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

## Attributions

- `/ui/static/media/videos/ocean.mp4` Video by Marianna Sigov: https://www.pexels.com/video/peaceful-ocean-waves-at-sunrise-35446596/
- `/ui/static/media/favicon.ico` - Favicon eagle icon is copyright 2020 Twitter, Inc., and other contributors. The graphics are licensed under [CC-BY 4.0](https://creativecommons.org/licenses/by/4.0/). No modifications were made to the original image.
