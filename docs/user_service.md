# User Service API

## Register

### Request
Endpoint: `/api/user/register`<br>
Method: `POST`<br>
Request Body: 
```javascript
{
    "email": "pac@pacn.in",
    "username": "pac",
    "password": "pakku"
}
```
### Response

#### Success
```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 10 Sep 2018 05:45:52 GMT
Content-Length: 58
```
```javascript
{
	"code": 200,
	"message": "Registered!",
	"data": null
}
```

#### Failure
```
HTTP/1.1 409 Conflict
Content-Type: application/json
Date: Mon, 10 Sep 2018 05:31:35 GMT
Content-Length: 72
```
```javaScript
{
	"code": 409,
	"message": "User Already Exists",
	"data": null
}
```

## Login

### Request
Endpoint: `/api/user/login`<br>
Method: `POST`<br>
Request Body: 
```javascript
{
    email: "pac@pacn.in",
    password: "pakku"
}
```

### Response

#### Success
```
HTTP/1.1 200 OK
Content-Type: application/json
Set-Cookie: session=mm2vXMGtOyS92KD2hn6M9sLIjtKd9JL0jKcamB6GO14BqcSPyCX80AGzzM9PZ1q7xYlYJTFBhcLFdd7f1jfJQ1nmSz5mOdOVweFnoIS1luyT9lnV-SvTIonUIWigrVtS7GkP57fkDxJDjS-5Y_IDtA; Path=/; HttpOnly
Date: Mon, 10 Sep 2018 05:33:39 GMT
Content-Length: 57
```
```javaScript
{
	"code": 200,
	"message": "Logged in!",
	"data": null
}
```

#### Failure
```
HTTP/1.1 401 Unauthorized
Content-Type: application/json
Date: Mon, 10 Sep 2018 05:31:35 GMT
Content-Length: 72
```
```javaScript
{
	"code": 401,
	"message": "Invalid login credentials",
	"data": null
}
```
