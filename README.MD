# Go service

##### What I learned

- Create your web framework over the http.Server
- create several middlewares (panics, metrics, logger, errors, auth)
- metrics are exposed by expvar (build, requests count, goroutines counts, errors count, memstats.Alloc count)
- profile via pprof
- tracing your application with opentelemetry and send your spans to zipking
- manage the graceful shutdown of your application
- create units tests in golang
- create integrations tests in golang
- build the docker image of your golang service
- create a kubernetes configuration and use kustomize to generate the differents version (dev, preprod, prod, ...)
- deploy a local k8s with kind
- generate a private key, public key in the programatic way in golang
- create, start, wait for, stop, remove a container in the programmatic way in golang
- create database migrations with darwin
- seed your database
- split your code in clean architecture in golang with function, business, and app layers

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
_curl -v -H "Authorization: Bearer $TOKEN" http://localhost:3000/users/1/10_
```

###### Loading server capacity

```
sudo apt install -y hey
go get github.com/divan/expvarmon
hey -m GET -c 100 -n 10000 -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/users/1/10 
expvarmon -ports=":4000" -vars="build,requests,goroutines,errors,mem:memstats.Alloc"
```