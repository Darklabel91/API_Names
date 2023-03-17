# API_Names

API developed to search for a name by it's metaphone (br) on database and returns all similar names that it can found based on metaphone code and the levenshtein algorithm (0.8 of similarity).

## Setup
- first run the API ```go run main.go```
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

- GET - /:id 
```json
{
  "ID": 1,
  "CreatedAt": "0001-01-01T00:00:00Z",
  "UpdatedAt": "0001-01-01T00:00:00Z",
  "DeletedAt": null,
  "Name": "ALINE",
  "Classification": "F",
  "Metaphone": "A",
  "NameVariations": "|AALINE|AILINE|ALEINE|ALIINE|ALINE|ALINER|ALINHE|ALINNE|ALYNE|ALYNNE|AYLINE|EALINE|ELEINE|ELINE|ELINER|ELINNE|ELYNE|EULINE|HALINE|HALYNE|HELEINE|HELINE|HELYNE|IALINE|ILEINE|ILINE|LEINE|LEINER|LEYNE|LINE|LINER|LUEINE|LUINE|LUYNE|LYNE|LYNNE|OLINE|UELINE|"
}
```

- GET - /name/:name 
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
```

- GET - /metaphone/:name 
```json
{
  "Name": "TIAGO",
  "metaphone": "TG",
  "name_variations": [
    "TAGO",
    "TAGOR",
    "TAIGA",
    "TAIGO",
    "TAIGOR",
    "TAYGO",
    "TAYGOR",
    "TEAGO",
    "TEAGO",
    "THAGO",
    "THAGO",
    "THAIGO",
    "THAIGOR",
    "THAYGO",
    "THAYGOR",
    "THEAGO",
    "THIAGA",
    "THIAGO",
    "THIEGO",
    "THIIAGO",
    "THIOGO",
    "THYAGO",
    "THYEGO",
    "TIAGA",
    "TIAGA",
    "TIAGHO",
    "TIAGO",
    "TIAGO",
    "TIAHGO",
    "TIEGO",
    "TIEGO",
    "TIHAGO",
    "TIIAGO",
    "TIOGO",
    "TIOGO",
    "TYAGO"
  ]
}
```
