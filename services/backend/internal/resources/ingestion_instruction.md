All data is attached at the end of the instruction

Maps given notifications into a list of dictionary which has the following shape

```json
{
    "transactions": [
        {
            "version": "1",
            "notification": {},
            "transaction": {
                "app": "string",
                "type": "string",
                "date": "string",
                "amount": "number",
                "currency": "string",
            }
        }
    ]
}
```

Rules:
- version is a string which is always "1"
- notification is an object holding the raw notification dictionary
- transaction (object, has the following keys)
    - app is a string holding app from notification
    - amount is a number holding absolute value of transaction amount
    - currency is a string holding iso currency, defaults to THB
    - type is a string from this list: INCOME, EXPENSE, or UNKNOWN. use UNKNOWN when you cannot precisely decide
    - date is a string holding transaction date in YYYY-MM-DD format. prioritize date from notification content over date field.

Data:
