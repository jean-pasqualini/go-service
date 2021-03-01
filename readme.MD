# Go service

##### How to test authentication + authorization
```
export TOKEN=$(make genkeys)                                            
curl -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/readiness
```