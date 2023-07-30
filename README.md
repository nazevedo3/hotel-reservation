# Hotel reservation backend

# Project enviornment variables

```
HTTP_LISTEN_ADDRESS=:3000
JWT_SECRET=somethingsupersecretthatNOBODYKNOWS
MONGO_DB_NAME=hotel-reservation
MONGO_DB_NAME_TEST=hotel-reservation-test
MONGO_DB_URL=mongodb://localhost:27017
MONGO_DB_URL_TEST=mongodb://localhost:27017
```

## Project outline

- users -> book from a hotel
- admins -> going to check reservation/bookings
- Authentication and authorization -> JWT tokens
- Hotels -> CRUD API -> JSON
- Rooms -> CRUD API -> JSON
- Scripts -> database management -> seeding, migration

## Resources

### Mongodb driver

### gofiber

## Docker

### Installing mongodb as a Docker container

```
docker run --name mongodb -d mongo:latest -p 27017:27017
```
