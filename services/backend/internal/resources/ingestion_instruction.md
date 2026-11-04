Convert the given notifications into a list of extracted transactions.

Return valid JSON only, without Markdown fences or explanations.

### Output schema

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
        "amount": "string",
        "currency": "string"
      }
    }
  ]
}
```

### Rules

- Return one item per input notification, preserving order.
- `version` is always `"1"`.
- `notification` contains the original notification dictionary.
- `app`: Application name from `notification.app`.
- `amount`: Absolute transaction amount as a string, without currency symbols. Use `"0"` if unavailable.
- `currency`: ISO 4217 currency code. Default to `"THB"`.
- `type`: One of `INFLOW`, `OUTFLOW`, or `UNKNOWN`.
  - `INFLOW`: Money received, incoming transfers, deposits, refunds, and reversals of previous transactions. Examples: "รายการเงินเข้า", "ได้รับเงิน", "รายการยกเลิก". A cancelled or voided credit card transaction (`รายการยกเลิก`) is an `INFLOW` because it restores the available credit card balance.
  - `OUTFLOW`: Money sent, outgoing transfers, withdrawals, card payments, and purchases. Examples: "รายการเงินออก", "รายการโอน", "ถอน", "โอน", "คุณโอน", "รายการใช้จ่าย", "รายการใช้บัตร".
  - `UNKNOWN`: Non-transactional notifications, ambiguous messages, and investment activities.
  - Exclude investment activities entirely, including fund purchases, redemptions, and fund transfers. Always classify them as `UNKNOWN`.
  - Ignore promotional messages, even if they mention monetary amounts or cashback.
  - Classify based on the actual financial activity, not merely keywords or the presence of an amount.
- `date`: Transaction date in `YYYY-MM-DD` format. Prefer the date in notification content over `notification.date`. Convert Thai Buddhist Era years to Gregorian years by subtracting 543. If unavailable, use the date from `notification.date`.
- `createdAt`: Notification timestamp from `notification.date`, converted to RFC3339. Preserve timezone when available; otherwise use `+07:00`.
- Unescape backslash-escaped characters in notification strings (e.g. `\\/` to `/`).

### Input

All input data is provided at the end of this instruction. Each record contains a notification dictionary.
