PIOT Server
===========

Authentication part inspired by https://medium.com/@theShiva5/creating-simple-login-api-using-go-and-mongodb-9b3c1c775d2f

User registration
-----------------
```
curl -v -X POST localhost:9096/register -d '{"email": "hello@example.com", "password": "hello"}'
```

User authentication - get token
-------------------------------
```
curl -v -X POST localhost:9096/login -d '{"email": "hello@example.com", "password": "hello"}'
```


