# User Service
User service is a microservice which is providing user CRUD operations via RestAPI having MONGODB as storage and a simple(more like a placeholder) auth mechanism.

## **Folder Descriptions**

- **`cmd/`**: The entry point of the application where the `main.go` file resides.

- **`internal/`**: The main application code, organized in layers.

- **`pkg/`**: Utility packages that can be reused across different parts of the application (e.g., logging, middlewares, httpserver init, error handling utilities, db initialization).

- **`docker-compose.yml`**: Docker Compose file to set up and run the application with its dependencies (e.g., MongoDB).

- **`Dockerfile`**: The Dockerfile to build a container image for the application.

## How to run

### With Docker
#### Prerequisities:
- Docker should be installed in your machine

UserAPI can run via docker simply by running `make compose-up` command in terminal. By running this command service and it's dependencies(mongodb) will be up and ready to use from `localhost:8080`

If you want to use your own configuration instead of default values, you can simply create `.env` file in project root with below params for example:

```
HOST_ADDRESS=0.0.0.0
PORT=3000
DB_NAME=users
MONGODB_URI=mongodb://127.0.0.1:27017
LOG_LEVEL=DEBUG
READ_TIMEOUT_IN_SECONDS=10
WRITE_TIMEOUT_IN_SECONDS=10
IDLE_TIMEOUT_IN_SECONDS=10
```

**NOTE**: After running you can run a healthcheck by manually calling `GET localhost:8080/health` or you can check docker logs since it is automatically running every 30 seconds.

### Alternative Run
---
#### Prerequisities:
- Mongodb instance should run.

Run `make run` in command line to run with default parameters.

## Make Http Requests:
- Create User: `curl -X POST localhost:8080/users --header "authorization: Bearer valid-token" -d '{"firstName":"John", "lastName":"Doe", "nickName":"johndoe", "email":"johndoe@email.com", "country":"TR"}'`

- Update User: `curl -X PUT localhost:8080/users/{id} --header "authorization: Bearer valid-token" -d '{"firstName":"Jane"}'`

- Delete User: `curl -X DELETE localhost:8080/users/{id} --header "authorization: Bearer valid-token"`

- List Users: `curl -X POST 'localhost:8080/users/filter?limit=5&offset=0' --header "authorization: Bearer valid-token" -d '{"firstName":"John", "country":"TR"}'`

## Data seeding
For data seeding you can use user-service-automation after running user-service app. There is a test method `TestUserCreate` under `tests/user_create_test` to create many user as defined in `CreateUserAmount` const.


### Response examples
#### Success response
---
- List users response (`POST /users/filter?limit=20&offset=0`)
```json
{
    "totalRecords": 1,
    "limit": 20,
    "offset": 0,
    "hasNext": false,
    "hasPrevious": false,
    "items": [
        {
            "id": "f84aaec4-f894-4797-8152-aa71710ab303",
            "firstName": "John",
            "lastName": "Doe",
            "nickName": "john.doe",
            "email": "johndoe@email.com",
            "country": "UK",
            "status": 1,
            "createdAt": "2024-11-28T09:14:41.523Z",
            "updatedAt": "2024-11-28T09:14:41.523Z",
            "version": 0
        },
        {
            "id": "22e04a63-bd56-450f-a2d9-627bbc79aa79",
            ..
            ..
        }
    ]
}
```

- Create user response(`POST /users`)
```json
{
    "id": "f26110a0-33cd-4142-95c4-687e82ef7d30",
    "firstName": "John",
    "lastName": "Doe",
    "nickName": "john.doe",
    "email": "john.doe@email.com",
    "country": "TR",
    "status": 1,
    "createdAt": "2024-12-01T19:49:13.090670886Z",
    "updatedAt": "2024-12-01T19:49:13.090670886Z",
    "version": 0
}
```

#### Error response
Same format for all errors. Returns HttpStatus code in response header according to error.
```json
{
    "Message": "firstName can't be empty",
    "Code": "400"
}
```


## MongoDB
- From the terminal connect mongoDB with "`mongosh`" and list values of `users` collection with these commands:
```sh
    #connect to MongoDB with 
    mongosh --port 27017

    #list databases 
    show dbs

    #use users database
    use users
    
    #run a mongo command
    db.users.find()
```

## Unit tests
In Unit tests `testify` lib used for assertion and mocking. For mocking the interfacer `mockery` tool has been used.

For repository layer unit test not written because it would be better to test it with integration tests. It's also challenging to write unit tests for repository.

**NOTE**: Unit test coverage can be increased by writing missing unit tests for handlers and codes under pkg folder