## Microservice outsourcing-auth

## 📖 API Documentation

### Login into account

Authorization for companies and clients. When the user is logged in, they will have a token that represents a valid session.

**URL** : `v1/login`

**Method** : `POST`

**Auth required** : NO

**Permissions required** : None

**Attributes:**
* email
* password

```json
{
  "user": {
    "login": {
      "email": "myemail@gmail.com",
      "password": "hash"
    }
  }
}
```

### Success Response

**Code** : `200 OK`

**Content examples**

For a user with an email address myemail@gmail.com and with some password, 
the server service will send a text status code, user ID, session token, and account type:  

```json
{
  "status": "success",
  "user": {
    "id": 123,
    "token": "token",
    "type": "client"
  }
}
```

----------
### Client registration

This API provides client registration, 
and the server service sends a text status code, user ID, session token, and account type:  

**URL** : `v1/register/client`

**Method** : `POST`

**Auth required** : NO

**Permissions required** : None

**Attributes:**
* Full_name
* Email
* Phone
* Password
* Photo (not necessary)
* Account type

```json
{
  "user": {
    "register": {
      "full_name": "name surname patronymic",
      "email": "myemail@gmail.com",
      "phone": "89017335432",
      "password": "hash",
      "photo": null,
      "type": "client"
    }
  }
}
```

### Success Response

**Code** : `200 OK`

**Content examples**

For a client with attributes, 
the server service will send a text status code, user ID, session token, and account type:  

```json
{
  "status": "success",
  "user": {
    "id": 123,
    "token": "token",
    "type": "client"
  }
}
```

----------
### Company registration

This API provides company registration,
and the server service sends a text status code, user ID, session token, and account type:

**URL** : `v1/register/company`

**Method** : `POST`

**Auth required** : NO

**Permissions required** : None

**Attributes:**
* Company_name
* Email
* Phone
* Full_name
* Position_agent
* ID_company (INN/OGRN)
* Address
* Type_service
* Password
* Photo (not necessary)
* Account type

```json
{
  "user": {
    "register": {
      "company_name": "Apple",
      "email": "myemail@gmail.com",
      "phone": "89017335432",
      "full_name": "name surname patronymic",
      "position_agent": "CEO",
      "id_company": "00000000",
      "address": "Moscow, Leninskiy prospekt, dom 7",
      "type_service": "Car dealership",
      "password": "hash",
      "photo": null,
      "type": "company"
    }
  }
}
```

### Success Response

**Code** : `200 OK`

**Content examples**

For a company with attributes,
the server service will send a text status code, user ID, session token, and account type:

```json
{
  "status": "success",
  "user": {
    "id": 123,
    "token": "token",
    "type": "company"
  }
}
```

----------
### Client/Company account

This API provides company and client account by session token,
and the server service sends a text status code, user ID, session token, and account type:

**URL** : `v1/account`

**Method** : `POST`

**Auth required** : Yes

**Permissions required** : None

**Attributes:**
* Token

```json
{
  "user": {
    "account": {
      "token": "token"
    }
  }
}
```

### Success Response

**Code** : `200 OK`

**Content examples**

For a **client**,
the server service will send a text status code, user ID, full name, phone, email, photo url, session token, and account type:

```json
{
  "status": "success",
  "user": {
    "account": {
      "id": 123,
      "full_name": "name surname patronymic",
      "phone": "89017335432",
      "email": "myemail@gmail.com",
      "photo": "photo_url",
      "token": "token",
      "type": "client"
    }
  }
}
```
----------
**Code** : `200 OK`

**Content examples**

For a **company**,
the server service will send a text status code, user ID, company name, stars, email, photo url, session token, and account type:

```json
{
  "status": "success",
  "user": {
    "account": {
      "id": 123,
      "company_name": "Apple",
      "stars": 4.7,
      "email": "myemail@gmail.com",
      "photo": "photo_url",
      "token": "token",
      "type": "company"
    }
  }
}
```

----------

## Notes

* This API created by UI design reference.
* `Stars` - this is a rating system (now it's a mock)
* `Session token` - this is an authentication token (which we can use to identify the user)