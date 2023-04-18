# API_Names

The API_Names project is an API designed to find a name and its possible variations from a given name, and correct any misspellings. It uses the Metaphone (br) algorithm to search the database and the Levenshtein distance method for correction.

## Features
- JWT Authentication
- Limited access by Token
- Sign-up
- Login
- Log local and imported to the database from time to time
- Middleware

## Requirements
Before running the API, make sure you have MySQL installed on your machine.

## Installation and Setup
1. Clone the repository 
2. Create a .env file at the root of your project and set the following environment variables:
  ```
  DB_USERNAME=<your_username>
  DB_PASSWORD=<your_password>
  DB_NAME=<your_database_name>
  DB_HOST=<your_database_host>
  DB_PORT=<your_database_port>
  SECRET=<your_jwt_secret>
  ```
  Replace the values with your own database credentials and a secret for JWT token generation. Worht to mention that the DB_NAME does not require an existing database.
  
3. Finally, run the API using the following command:
  ```go
  go run main.go
  ```
  On the first run, the prompt will return:
  ```bash
  go run main.go
  2023/04/12 18:48:48 -   Created Database
  2023/04/12 18:48:48 -   Upload data start
  2023/04/12 18:49:21 -   Upload data finished 33.113701109s
  2023/04/12 18:49:21 -   Created first user
  2023/04/12 18:49:21 -   Listening and serving...
  ```

## API Endpoints
The main endpoint for the API is ```http://localhost:8080/metaphone/:name```. You need to log in to get an access token before you can access any other endpoint.

The following table shows the available endpoints, their corresponding HTTP methods, and a brief description:
| Req    | Endpoint                               | Description                         | Success           | Error                  |
|--------|----------------------------------------|-------------------------------------|-------------------|------------------------|
| POST   | /signup                                | Create a new user                   | Status:200 - JSON | Status: 400/401 - JSON |
| POST   | /login                                 | Login user on API                   | Status:200 - JSON | Status: 400/401 - JSON |
| POST   | /name                                  | Create a name in the database       | Status:200 - JSON | Status: 400/401 - JSON |
| DELETE | /:id                                   | Delete a name by given id           | Status:200 - JSON | Status: 404/401 - JSON |
| PUT    | /:id                                   | Update a name by given id           | Status:200 - JSON | Status: 500/401 - JSON |
| GET    | /:id                                   | Read name with given id             | Status:200 - JSON | Status: 400/401 - JSON |
| GET    | /name/:name                            | Read name with given name           | Status:200 - JSON | Status: 404/401 - JSON |
| GET    | /metaphone/:name                       | Read metaphones of given name       | Status:200 - JSON | Status: 404/401 - JSON |


## Endpoint Examples

- POST - ```http://localhost:8080/signup```
```json
{
    "Email": "user@user.com",
    "Password": "123456"
}
```
Return:
```json
{
    "Message": "User created",
    "User": {
        "ID": 2,
        "CreatedAt": "2023-04-12T18:48:48.475-03:00",
        "UpdatedAt": "2023-04-12T18:48:48.475-03:00",
        "DeletedAt": null,
        "Email": "user@user.com",
        "Password": "$2a$10$crIN3KKScm.HafCl9qQkzeehuK5XUfnGrAxCyymyMPnNHkwDwHBVS"
    }
}
```

- POST - ```http://localhost:8080/login```
```json
{
    "Email": "user@user.com",
    "Password": "123456"
}
```
Return: ```status 200```
```json
{
    "Message": "Login successful"
}
```

- GET - ```http://localhost:8080/3```
```json
{
  "ID": 3,
  "CreatedAt": "2023-04-12T18:48:48.475-03:00",
  "UpdatedAt": "2023-04-12T18:48:48.475-03:00",
  "DeletedAt": null,
  "Name": "ARON",
  "Classification": "M",
  "Metaphone": "ARM",
  "NameVariations": "|AARON|AHARON|AROM|ARON|ARYON|HARON|"
}
```

- GET - ```http://localhost:8080/name/aron```
```json
{
  "ID": 3,
  "CreatedAt": "2023-04-12T18:48:48.475-03:00",
  "UpdatedAt": "2023-04-12T18:48:48.475-03:00",
  "DeletedAt": null,
  "Name": "ARON",
  "Classification": "M",
  "Metaphone": "ARM",
  "NameVariations": "|AARON|AHARON|AROM|ARON|ARYON|HARON|"
}
```

- GET - ```http://localhost:8080/metaphone/haron```
```json
{
    "ID": 3,
    "CreatedAt": "2023-04-12T18:48:48.475-03:00",
    "UpdatedAt": "2023-04-12T18:48:48.475-03:00",
    "DeletedAt": null,
    "Name": "ARON",
    "Classification": "M",
    "Metaphone": "ARM",
    "NameVariations": "ARON | AROM | AARON | ARYON | HARON | AHARON | "
}
```
## Dependencies
- [METAPHONE - BR](https://github.com/DanielFillol/metaphone-br)
- [GIN](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io)
- [MySQL - GORM](https://github.com/go-gorm/mysql)
- [GO.ENV](https://github.com/joho/godotenv)
- [JWT](https://github.com/golang-jwt/jwt)

## Extra
If you're interested in checking out my API caller, you can find it by clicking on this [link](https://github.com/DanielFillol/API_Caller)
