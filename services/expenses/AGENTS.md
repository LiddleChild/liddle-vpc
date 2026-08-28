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
- Account-name filtering with toggle controls so users can include or exclude accounts from the summary.

Expense entry, editing, category analytics, and complex charts are out of scope unless explicitly requested.

## Architecture

- Use Next.js with the App Router, React, and TypeScript.
- Use server components by default; add client components only for browser interaction or state.
- Use Next.js route handlers for the same-origin read-only API.
- Expose the monthly summary through `GET /api/monthly-summary?month=YYYY-MM`.
- Keep Grist credentials and private tokens on the server; never expose them in client bundles.
- Send all Grist API requests through the Next.js backend only; browser and client-side code must not call Grist directly.
- Keep Grist table and column names isolated in a server-side data adapter.
- Normalize and validate proxy responses before rendering them.
- Keep data access, calculation logic, and presentation in separate modules.
- Use shadcn/ui components and Tailwind CSS for interface primitives and styling.
- Prefer small typed modules and pure calculation functions.
- Do not add dependencies without a concrete need.
- Use `dayjs` for date and time parsing, calculations, and formatting throughout the project.
- This is a view-only app: it must not create, edit, delete, or otherwise mutate expense or account data. Local display filters and month navigation are allowed.

The current application is a minimal starter. The Grist route-handler boundary should support a mock/provider implementation before live credentials are available.

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
  currency: string
  syncedAt: string
}
```

Actual Grist schema names must be confined to the server-side adapter rather than spread through UI components or calculation code.

## API route requirements

The route handler at `app/api/monthly-summary/route.ts` must:

- Accept only the `month=YYYY-MM` query parameter.
- Validate the month format before querying a provider.
- Return a typed `MonthlyExpenseSummary` response when data is available.
- Return safe, generic error responses without exposing Grist details or credentials.
- Remain read-only.
- Keep the provider mockable so the UI and calculation logic do not require a live Grist document during development.

Do not put Grist API keys in `NEXT_PUBLIC_*` variables or any other browser-exposed configuration.

## Expense aggregation

Use the configured/local user timezone for month boundaries.

```text
expenses = approved expenses + pending expenses - refunds/adjustments
```

- Canceled and rejected expenses are excluded.
- Explicitly excluded records are excluded.
- Refunds and adjustments reduce the displayed expense total.
- Preserve the source currency and format amounts consistently.
- Treat missing, malformed, or contradictory financial values as data errors rather than silently inventing values.

## UX and iOS PWA requirements

- Use semantic HTML and accessible labels.
- Ensure controls have comfortable touch targets.
- Support keyboard navigation and visible focus states.
- Account for iOS safe-area insets and viewport sizing.
- Configure the Next.js manifest with Expenses name, short name, start URL, theme color, and standalone display mode.
- Include iOS-compatible viewport and Apple web-app metadata.
- Keep app icons in Next.js metadata-compatible locations.
- Do not rely on hover interactions.
- Allow users to filter accounts by name using accessible toggle controls; make active and inactive states clear, support keyboard and touch interaction, and update the displayed summary when account filters change.
- Respect `prefers-reduced-motion`.
- Keep the first screen focused on the selected month and its current expense total.
- Keep the first screen minimal: show the current expense total and month picker, without account or expense breakdowns.
- Open Settings by swiping left using horizontal CSS scroll snap; swipe right to return to the main page. Do not add a filter button.
- Settings filters apply immediately. Category is single-select; tags and accounts are multi-select. Use a compact “Clear all” action in the top-right.
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
- Validate response shape, numeric values, dates, statuses, and currency before use.
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
