# Go service

##### Routes
- Check
    - GET    /readiness
    - GET    /liveness
- Users    
    - GET    /users/:page/:rows
    - GET    /users/:id
    - GET    /users/token/:kid
    - POST   /users
    - PUT    /users/:id
    - DELETE /users/:id

##### How to build the docker image
```
make sales-api
```

##### How to run in local k8s
```
make sales-api
make kind-up
make kind-load
make kind-services
make kind-status
```

##### How to test automatically (unit tests + integration tests)
```
go get -u honnef.co/go/tools/cmd/staticcheck
make test
```

##### How to test manually

###### Authentication
```
curl --user "admin@example.com:gophers" http://localhost:3000/users/token/54bb2165-71e1-41a6-af3e-7da4a0e1e2c1
export TOKEN=$("the token")                                            
```

###### List users
```
curl -v -H "Authorization: Bearer $TOKEN" http://localhost:3000/users/1/10
```

###### Loading server capacity

```
sudo apt install -y hey
go get github.com/divan/expvarmon
hey -m GET -c 100 -n 10000 -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/users/1/10 
expvarmon -ports=":4000" -vars="build,requests,goroutines,errors,mem:memstats.Alloc"
```