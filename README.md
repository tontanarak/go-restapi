#  REST API for a spa booking system.
- Customers can book a spa session with a predefined time slot.
- Spa owners can increase/decrease the number of spa sessions.


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
