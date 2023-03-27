# API_Names

API developed to search for a name by it's metaphone (br) on database and returns all similar names that it can found based on metaphone code and the levenshtein algorithm (0.8 of similarity).

## Setup
- create a MySql database, you can use [this script](https://github.com/Darklabel91/API_Names/blob/main/database/create_database.txt) if needed
- set a .env file on the project with ```DB_USERNAME```, ```DB_PASSWORD```, ```DB_NAME```, ```DB_HOST``` and ```DB_PORT```
- run the API ```go run main.go```
- use import wizzard on mysql workbench to upload the [.csv file](https://github.com/Darklabel91/API_Names/blob/main/database/name_types.csv)
- have fun

## API
This project suports a simple CRUD. the main endpoint is  ```/name``` where you can search for a name and it returns the name, metaphone code and all it's variations

Every method expect Status:200 and JSON content-type as show bellow:

| Req    | Endpoint                               | Description                         | Success           | Error              |
|--------|----------------------------------------|-------------------------------------|-------------------|--------------------|
| POST   | /name                                  | Create a name in the database       | Status:200 - JSON | Status: 400 - JSON |
| DELETE | /:id                                   | Delete a name by given id           | Status:200 - JSON | Status: 404 - JSON |
| PUT    | /:id                                   | Update a name by given id           | Status:200 - JSON | Status: 500 - JSON |
| GET    | /:id                                   | Read name with given id             | Status:200 - JSON | Status: 400 - JSON |
| GET    | /name/:name                            | Read name with given name           | Status:200 - JSON | Status: 404 - JSON |
| GET    | /metaphone/:name                       | Read metaphones of given name       | Status:200 - JSON | Status: 404 - JSON |


## Endpoint Examples

- GET - /3 
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

- GET - /name/aron 
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

- GET - /metaphone/haron
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
## Dependacy
[metaphone-br](https://github.com/Darklabel91/metaphone-br)
