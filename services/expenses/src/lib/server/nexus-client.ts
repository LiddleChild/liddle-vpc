import type {
  AccountExpense,
  CatalogChoice,
  ExpenseCatalog,
  ExpenseSummaryFilters,
  MonthlyExpenseSummary,
} from "@/lib/types/catalog";

const isRecord = (value: unknown): value is Record<string, unknown> =>
  typeof value === "object" && value !== null;

const getNexusURL = () => {
  const baseURL = process.env.NEXUS_URL?.replace(/\/+$/, "");
  if (!baseURL) {
    throw new Error("Nexus URL is not configured");
  }
  return baseURL;
};

const requestJSON = async (url: URL): Promise<unknown> => {
  const response = await fetch(url, {
    headers: { Accept: "application/json" },
    cache: "no-store",
  });
  if (!response.ok) {
    throw new Error("Nexus request failed");
  }

  try {
    return await response.json();
  } catch {
    throw new Error("Nexus returned invalid JSON");
  }
};

const parseChoice = (value: unknown): CatalogChoice => {
  if (!isRecord(value) || typeof value.value !== "string") {
    throw new Error("Nexus returned an invalid catalog choice");
  }
  if (
    (value.textColor !== undefined && typeof value.textColor !== "string") ||
    (value.fillColor !== undefined && typeof value.fillColor !== "string")
  ) {
    throw new Error("Nexus returned invalid catalog styling");
  }

  return {
    value: value.value,
    ...(value.textColor === undefined ? {} : { textColor: value.textColor }),
    ...(value.fillColor === undefined ? {} : { fillColor: value.fillColor }),
  };
};

const parseCatalog = (value: unknown): ExpenseCatalog => {
  if (
    !isRecord(value) ||
    !Array.isArray(value.categories) ||
    !Array.isArray(value.tags) ||
    !Array.isArray(value.accounts)
  ) {
    throw new Error("Nexus returned an invalid expense catalog");
  }

  const accounts = value.accounts.map((account) => {
    if (
      !isRecord(account) ||
      typeof account.id !== "number" ||
      !Number.isInteger(account.id) ||
      account.id < 1 ||
      typeof account.name !== "string" ||
      !account.name.trim()
    ) {
      throw new Error("Nexus returned an invalid account");
    }
    return { id: account.id, name: account.name };
  });

  return {
    categories: value.categories.map(parseChoice),
    tags: value.tags.map(parseChoice),
    accounts,
  };
};

const parseAccountExpense = (value: unknown): AccountExpense => {
  if (
    !isRecord(value) ||
    typeof value.accountId !== "number" ||
    !Number.isInteger(value.accountId) ||
    value.accountId < 1 ||
    typeof value.accountName !== "string" ||
    typeof value.expenses !== "number" ||
    !Number.isFinite(value.expenses) ||
    typeof value.transactionCount !== "number" ||
    !Number.isInteger(value.transactionCount) ||
    value.transactionCount < 0
  ) {
    throw new Error("Nexus returned invalid account expense totals");
  }

  return {
    accountId: value.accountId,
    accountName: value.accountName,
    expenses: value.expenses,
    transactionCount: value.transactionCount,
  };
};

const parseMonthlySummary = (value: unknown): MonthlyExpenseSummary => {
  if (
    !isRecord(value) ||
    typeof value.month !== "string" ||
    typeof value.expenses !== "number" ||
    !Number.isFinite(value.expenses) ||
    typeof value.transactionCount !== "number" ||
    !Number.isInteger(value.transactionCount) ||
    value.transactionCount < 0 ||
    typeof value.accountCount !== "number" ||
    !Number.isInteger(value.accountCount) ||
    value.accountCount < 0 ||
    !Array.isArray(value.accountExpenses) ||
    typeof value.syncedAt !== "string" ||
    !Number.isFinite(Date.parse(value.syncedAt))
  ) {
    throw new Error("Nexus returned an invalid monthly summary");
  }

  return {
    month: value.month,
    expenses: value.expenses,
    transactionCount: value.transactionCount,
    accountCount: value.accountCount,
    accountExpenses: value.accountExpenses.map(parseAccountExpense),
    syncedAt: value.syncedAt,
  };
};

const toCSV = (values: Array<string | number>) =>
  values
    .map((value) => {
      const text = String(value);
      return /[",\r\n]/.test(text) ? `"${text.replace(/"/g, '""')}"` : text;
    })
    .join(",");

export async function getExpenseCatalog(): Promise<ExpenseCatalog> {
  const url = new URL(`${getNexusURL()}/expenses/v1/catalog`);
  return parseCatalog(await requestJSON(url));
}

export async function getMonthlyExpenseSummary(
  filters: ExpenseSummaryFilters,
): Promise<MonthlyExpenseSummary> {
  const url = new URL(`${getNexusURL()}/expenses/v1/monthly-summary`);
  url.searchParams.set("month", filters.month);
  url.searchParams.set("accounts", toCSV(filters.accounts));
  url.searchParams.set("categories", toCSV(filters.categories));
  url.searchParams.set("tags", toCSV(filters.tags));

  const summary = parseMonthlySummary(await requestJSON(url));
  if (summary.month !== filters.month) {
    throw new Error("Nexus returned a summary for a different month");
  }
  return summary;
}
