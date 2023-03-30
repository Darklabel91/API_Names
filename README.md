# API_Names

The API_Names project is designed to find a name and its possible variations from a given name, and correct any misspellings. It uses the Metaphone (br) algorithm to search the database and the Levenshtein distance method for correction.

## How to Run

Before running the API, make sure you have MySQL installed on your machine. Then, create a .env file and set the following environment variables:
 ```
DB_USERNAME
DB_PASSWORD
DB_NAME
DB_HOST
DB_PORT
SECRET
```
Finally, run the API using the command:
```go
go run main.go
```
on the first run the prompt will return:
```
go run main.go
-       Create Database
-       Created first user
-       Upload data start
-       Upload data finished time_uploading
-       Listening and serving
```


## API Endpoints
The main endpoint for the API is ```http://localhost:8080/metaphone/:name```. You need to log in to get an access token before you can access any other endpoint. We use JWT to generate access tokens.

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
        "CreatedAt": "2023-03-28T23:18:23.624-03:00",
        "UpdatedAt": "2023-03-28T23:18:23.624-03:00",
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

- GET - ```http://localhost:8080/3```
```json
{
  "ID": 3,
  "CreatedAt": "0001-01-01T00:00:00Z",
  "UpdatedAt": "0001-01-01T00:00:00Z",
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
  "CreatedAt": "0001-01-01T00:00:00Z",
  "UpdatedAt": "0001-01-01T00:00:00Z",
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
  "CreatedAt": "0001-01-01T00:00:00Z",
  "UpdatedAt": "0001-01-01T00:00:00Z",
  "DeletedAt": null,
  "Name": "ARON",
  "Classification": "M",
  "Metaphone": "ARM",
  "NameVariations": [
    "ARON",
    "AARON",
    "AHARON",
    "AROM",
    "ARYON",
    "HARON",
    "HARNON",
    "AIROM",
    "AIRON",
    "AIRYON",
    "AYRON",
    "HAIRON",
    "HAYRON",
    "IARON",
    "YARON",
    "ARLON",
    "ARILON",
    "ARLOM",
    "HARLON",
    "ARION",
    "ARNON",
    "ARNOM",
    "ARONE",
    "ARONI",
    "ARONY",
    "ARTON",
    "JARON",
    "JAROM",
    "KARON",
    "CARON",
    "MARON",
    "MAROM",
    "MARRON",
    "MARYON",
    "NARON",
    "RARON",
    "SARON",
    "SAROM"
  ]
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
