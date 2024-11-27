# user-service

## Make Http Requests:
- Create User: `curl -X POST localhost:3000/users -d '{"firstName":"John", "lastName":"Doe", "nickName":"johndoe", email:"johndoe@email.com", "country":"TR"}'`

- Update User: `curl -X PUT localhost:3000/users/{id} -d '{"firstName":"Jane"}'`

- Delete User: `curl -X DELETE localhost:3000/users/{id}`

- List Users: `curl -X POST localhost:3000/users/filter?limit=5&offset=0 -d '{"firstName":"John", "country":"TR"}'`

## MongoDB
- From the terminal connect mongoDB with "`mongosh`" and list values of `users` collection with these commands:
    - connect to DB with `mongosh --port 27017`
    - `show dbs`
    - `use users`
    - `db.users.find()`
