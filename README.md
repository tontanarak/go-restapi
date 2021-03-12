#  REST API for a spa booking system.
- Customers can book a spa session with a predefined time slot.
- Spa owners can increase/decrease the number of spa sessions.

## Header
```bash
Authorization: Bearer <jwt token>

# sample <jwt token>
# Header 
{
  "alg": "HS256",
  "typ": "JWT"
}

# Payload
{
  "id" : 1,
  "name" : "Warunthorn",
  "admin" : true
}



```

## Endpoints

### Increase session

```bash
POST /spasessions

# sample
{
    "time" : "2021-03-12 05:12:48"
}

```
### Decrease session

```bash
DELETE /spasessions/{id}
```
### Book session

```bash
PATCH /spasessions/{id}

# sample

{
    "CUSTOMER": "Tanarak Chunsanit"
}

```
