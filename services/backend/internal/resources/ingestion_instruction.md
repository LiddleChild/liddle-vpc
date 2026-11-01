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
                "createdAt": "string",
                "amount": string",
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
    - amount is a string holding absolute value of transaction amount
    - currency is a string holding iso currency, defaults to THB
    - type is a string from this list: INFLOW, OUTFLOW, or UNKNOWN. it is decided from notification title or subtitle or body
        - INFLOW usually contains "รายการเงินเข้า", ได้รับเงิน"
        - OUTFLOW usually contains "รายการเงินออก", "รายการโอน", "ถอน", "โอน"
        - UNKNOWN if it cannot be precisely decided
    - date is a string holding transaction date in YYYY-MM-DD format. prioritize date from notification content over date field.
    - createdAt is a string holding notification timestamp in RFC3339 format
- unescaped all backslash escaped characters to their original characters

Data:
