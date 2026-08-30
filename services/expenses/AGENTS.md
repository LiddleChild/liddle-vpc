# Expenses project guidance

## Product goal

Expenses is a mobile-first PWA optimized for modern iOS Safari. It consumes private Grist data and answers one question quickly:

> How much have I spent across all of my accounts this month?

Prioritize glanceability, minimal interaction, clear financial states, and accessible mobile presentation.

The initial dashboard should focus on:

- Current expenses as the dominant visual.
- Selected month.
- Total expenses across included accounts.
- Last successful sync time.
- Account filtering with toggle controls so users can include or exclude accounts from the summary. Account names are display-only; filters use numeric Grist row IDs.

Expense entry, editing, category analytics, and complex charts are out of scope unless explicitly requested.

## Architecture

- Use Next.js with the App Router, React, and TypeScript.
- Use server components by default; add client components only for browser interaction or state.
- Use Next.js route handlers for the same-origin read-only API.
- Expose the monthly summary through `GET /api/monthly-summary?month=YYYY-MM`.
- Keep Grist credentials and private tokens on the server; never expose them in client bundles.
- Send all Grist API requests through the Next.js backend only; browser and client-side code must not call Grist directly.
- Keep Grist table and column names isolated in a server-side data adapter. The current fixed tables are `Transactions` and `Accounts`.
- Normalize and validate proxy responses before rendering them.
- Keep data access, calculation logic, and presentation in separate modules.
- Use shadcn/ui components and Tailwind CSS for interface primitives and styling.
- Prefer small typed modules and pure calculation functions.
- Do not add dependencies without a concrete need.
- Use `dayjs` for date and time parsing, calculations, and formatting throughout the project.
- This is a view-only app: it must not create, edit, delete, or otherwise mutate expense or account data. Local display filters and month navigation are allowed.

The current application uses a server-side Grist client in `src/lib/server/grist-client.ts`. Credentials come from `GRIST_URL`, `GRIST_DOC_ID`, and `GRIST_API_KEY`; do not add a table-name environment variable for the fixed `Transactions` and `Accounts` tables.

`getExpenseCatalog` reads `Category` and `Tags` choice metadata from `Transactions`, including `choiceOptions` colors, and reads account IDs and `Name` values from `Accounts`. `getMonthlyExpenseSummary` uses Grist's parameterized SQL endpoint to aggregate transactions.

## Data contract

The Grist adapter should map upstream records into these application-level shapes:

```ts
type ExpenseSourceRecord = {
  month: string // YYYY-MM
  amount: number // total expenses from this source for the month
  currency: string
}

type ExpenseRecord = {
  date: string
  amount: number
  status: 'approved' | 'pending' | 'canceled' | 'rejected'
  refundAmount?: number
  excluded?: boolean
}

type MonthlyExpenseSummary = {
  month: string
  expenses: number
  transactionCount: number
  accountCount: number
  syncedAt: string
}
```

Shared application types live in `src/lib/types/catalog.ts`. `Account.id` is a number. `ExpenseSummaryFilters.accounts` is a `number[]`; category and tag filters are string arrays. `CatalogChoice` carries the choice value plus optional `textColor` and `fillColor`.

Actual Grist schema names must be confined to the server-side adapter rather than spread through UI components or calculation code.

## API route requirements

The route handler at `app/api/monthly-summary/route.ts` must:

- Accept only the `month=YYYY-MM` query parameter.
- Validate the month format before querying a provider.
- Return a typed `MonthlyExpenseSummary` response when data is available.
- Return safe, generic error responses without exposing Grist details or credentials.
- Remain read-only.
- Keep the data-access functions isolated so the UI does not call Grist directly.

Do not put Grist API keys in `NEXT_PUBLIC_*` variables or any other browser-exposed configuration.

## Expense aggregation

Use the configured/local user timezone for month boundaries. The current SQL query uses the selected month's Unix-second start (inclusive) and the next month's start (exclusive).

```text
expenses = SUM(Amount_In_Account_Currency * Inflow_Direction)
```

- The current query assumes all transactions are in THB; currency is not part of the response.
- Accounts are matched by `Transactions.Account` row ID.
- Categories use scalar `IN` matching.
- Tags are stored as a JSON array and use a join through `json_each`; a transaction matches when it contains any selected tag.
- Empty filter arrays do not bypass SQL predicates; the UI defaults to selecting every available value.

## UX and iOS PWA requirements

- Use semantic HTML and accessible labels.
- Ensure controls have comfortable touch targets.
- Support keyboard navigation and visible focus states.
- Account for iOS safe-area insets and viewport sizing.
- Configure the Next.js manifest with Expenses name, short name, start URL, theme color, and standalone display mode.
- Include iOS-compatible viewport and Apple web-app metadata.
- Keep app icons in Next.js metadata-compatible locations.
- Do not rely on hover interactions.
- Allow users to filter accounts using accessible toggle controls; display account names but store numeric IDs, make active and inactive states clear, support keyboard and touch interaction, and update the displayed summary when account filters change.
- Respect `prefers-reduced-motion`.
- Keep the first screen focused on the selected month and its current expense total.
- Keep the first screen minimal: show the current expense total and month picker, without account or expense breakdowns.
- Open Settings by swiping left using horizontal CSS scroll snap; swipe right to return to the main page. Do not add a filter button.
- Settings filters apply immediately. Categories, tags, and accounts are multi-select. All catalog values are selected by default. Use a compact “Clear all” action in the top-right.
- Support both portrait and landscape layouts. Do not force orientation or show an orientation-rotation prompt.
- Do not reset the selected page on `pagehide`; preserve normal browser/PWA page lifecycle behavior.
- Show explicit loading, empty-month, API failure, malformed-data, offline/connectivity, stale-data, and overspending states.
- Offline editing and offline snapshots are not supported. When the API cannot be reached, explain that connectivity is required.
- Do not add offline caching or service-worker behavior unless the offline policy is explicitly expanded.

## Security and privacy

- Never place Grist API keys, private tokens, or server credentials in `NEXT_PUBLIC_*` variables.
- Always proxy Grist requests through the Next.js backend; do not expose Grist endpoints, credentials, or direct client access paths to the browser.
- The proxy must be read-only unless write access is explicitly requested.
- Do not return raw upstream errors or sensitive Grist details to the browser.
- Validate response shape, numeric values, and dates before use. Currency is assumed to be THB and is not queried or returned.
- Avoid logging private financial records or credentials.

## Development workflow

Use pnpm for all package operations and scripts:

```bash
pnpm install
pnpm dev
pnpm run lint
pnpm run build
```

Use Next.js App Router conventions and preserve server/client boundaries. Run lint and the production build before completing implementation work.

## Testing requirements

Add tests for pure expense aggregation logic covering:

- Approved and pending expenses being aggregated.
- Canceled and rejected expenses being excluded.
- Explicit exclusions.
- Refunds and adjustments.
- Month boundaries and configured timezone behavior.
- Exact zero expense totals.
- Refunds and adjustments reducing the expense total.
- Empty and malformed API responses.
- Invalid and missing `month` query parameters.
- Safe route-handler errors without credential or upstream-detail leakage.

Manually verify the responsive UI on a current iOS Safari viewport, including safe-area layout, standalone PWA metadata, touch interaction, loading, API failure, and connectivity states.
