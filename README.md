# goapi

A simple Go API server

## Usage
`$ ./web -addr=:8080 -dbhost=postgres -dbuser=postgres -dbpass=postgres -dbname=postgres -secret=my_server_secret_key`

## Setup
For init a test env of the App:
`$ docker-compose up -d --build`

## API endpoints

| Endpoint              | Method |
|-----------------------|:------:|
| /ping                 |  POST  |
| /auth                 |  POST  |
| /user                 |  POST  |
| /user/{uuid}         |  GET   |

## Example requests:
```
# Check if server is up
curl -s -XPOST localhost:8080/ping

Response:
OK


# Create user
curl -s -XPOST localhost:8080/user -H 'content-type: application/json' -d'
{
  "firstname": "cloud",
  "lastname": "strife",
  "cpf": "61849763623",
  "email": "cloud@strife.com",
  "password": "my_pass"
}'

# Authenticate
curl -s -XPOST localhost:8080/auth -H 'content-type: application/json' -d'
{
  "login": "cloud@strife.com",
  "password": "my_pass"
}'

Response:
{
  "expires": "2019-06-25T17:01:19Z",
  "login": {
    "uuid": "56689a6b-179d-4388-a7ba-be6955ed441d",
    "firstname": "cloud",
    "lastname": "strife",
    "cpf": "61849763623",
    "email": "cloud@strife.com",
    "created": "2019-06-24T17:00:05.348706Z",
    "active": true
  },
  "message": "User found and token generated",
  "status": "success",
  "tokenjwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dpbiI6ImNsb3VkQHN0cmlmZS5jb20iLCJleHAiOjE1NjE0ODIwNzl9.q05sLlPVUlc4Bzws5sB7Pzy0lDsfddwuDXtvbXFIHz0"
}

# Get user info
curl -s -XGET localhost:8080/user/56689a6b-179d-4388-a7ba-be6955ed441d \
-H 'content-type: application/json' \
-H 'Authorization: bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dpbiI6ImNsb3VkQHN0cmlmZS5jb20iLCJleHAiOjE1NjE0ODIwNzl9.q05sLlPVUlc4Bzws5sB7Pzy0lDsfddwuDXtvbXFIHz0'

Response:
{
  "uuid": "56689a6b-179d-4388-a7ba-be6955ed441d",
  "firstname": "cloud",
  "lastname": "strife",
  "cpf": "61849763623",
  "email": "cloud@strife.com",
  "created": "2019-06-24T17:00:05.348706Z",
  "active": true
}

