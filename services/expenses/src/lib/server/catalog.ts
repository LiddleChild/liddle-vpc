import { GristClient, type GristColumn } from "@/lib/server/grist-client";

export type Account = {
  id: string;
  name: string;
};

export type ExpenseCatalog = {
  categories: string[];
  tags: string[];
  accounts: Account[];
};

export type ExpenseSummaryFilters = {
  month: string;
  accounts: string[];
  category: string;
  tags: string[];
};

export type MonthlyExpenseSummary = {
  month: string;
  expenses: number;
  transactionCount: number;
  accountCount: number;
  currency: string;
  syncedAt: string;
};

type MockExpense = {
  account: string;
  month: string;
  amount: number;
  category: string;
  tags: string[];
  transactionCount: number;
  currency: string;
};

const mockAccounts: Account[] = [
  { id: "daily-spending", name: "Daily spending" },
  { id: "home-bills", name: "Home & bills" },
  { id: "shared-expenses", name: "Shared expenses" },
  { id: "cash-other", name: "Cash & other" },
];

const getGristClient = () =>
  new GristClient({
    baseUrl: process.env.GRIST_URL ?? "",
    docId: process.env.GRIST_DOC_ID ?? "",
    apiKey: process.env.GRIST_API_KEY ?? "",
  });

const getChoiceValues = (column: GristColumn): string[] => {
  if (!column.fields.widgetOptions) {
    throw new Error(`Column ${column.id} has no choice options`);
  }

  let widgetOptions: unknown;
  try {
    widgetOptions = JSON.parse(column.fields.widgetOptions);
  } catch {
    throw new Error(`Column ${column.id} has invalid choice options`);
  }

  if (
    typeof widgetOptions !== "object" ||
    widgetOptions === null ||
    !Array.isArray((widgetOptions as { choices?: unknown }).choices) ||
    !(widgetOptions as { choices: unknown[] }).choices.every(
      (choice) => typeof choice === "string",
    )
  ) {
    throw new Error(`Column ${column.id} has invalid choice options`);
  }

  return (widgetOptions as { choices: string[] }).choices;
};

const getColumn = (columns: GristColumn[], columnId: string) => {
  const column = columns.find(
    (candidate) =>
      candidate.id === columnId || candidate.fields.label === columnId,
  );
  if (!column) {
    throw new Error(`Column ${columnId} was not found`);
  }
  return column;
};

const mockExpenses: MockExpense[] = [
  {
    account: "Daily spending",
    month: "2025-09",
    amount: 428.35,
    category: "Everyday",
    tags: ["Needs", "Card"],
    transactionCount: 12,
    currency: "THB",
  },
  {
    account: "Home & bills",
    month: "2025-09",
    amount: 1360,
    category: "Bills",
    tags: ["Needs", "Recurring"],
    transactionCount: 4,
    currency: "THB",
  },
  {
    account: "Shared expenses",
    month: "2025-09",
    amount: 225.9,
    category: "Everyday",
    tags: ["Shared"],
    transactionCount: 5,
    currency: "THB",
  },
  {
    account: "Cash & other",
    month: "2025-09",
    amount: 86.2,
    category: "Other",
    tags: ["Cash"],
    transactionCount: 2,
    currency: "THB",
  },
];

export async function getExpenseCatalog(): Promise<ExpenseCatalog> {
  const client = getGristClient();
  const columns = await client.listTableColumns(
    process.env.GRIST_TRANSACTIONS_TABLE ?? "Transactions",
  );

  return {
    categories: getChoiceValues(getColumn(columns, "Category")),
    tags: getChoiceValues(getColumn(columns, "Tags")),
    accounts: mockAccounts,
  };
}

export async function getMonthlyExpenseSummary(
  filters: ExpenseSummaryFilters,
): Promise<MonthlyExpenseSummary> {
  const matchingExpenses = mockExpenses.filter(
    (expense) =>
      expense.month === filters.month &&
      filters.accounts.includes(expense.account) &&
      (filters.category === "all" || expense.category === filters.category) &&
      filters.tags.every((tag) => expense.tags.includes(tag)),
  );

  const currencies = new Set(
    matchingExpenses.map((expense) => expense.currency),
  );
  if (currencies.size > 1) {
    throw new Error("Mixed currencies are not supported");
  }

  return {
    month: filters.month,
    expenses: matchingExpenses.reduce(
      (total, expense) => total + expense.amount,
      0,
    ),
    transactionCount: matchingExpenses.reduce(
      (total, expense) => total + expense.transactionCount,
      0,
    ),
    accountCount: new Set(matchingExpenses.map((expense) => expense.account))
      .size,
    currency: matchingExpenses[0]?.currency ?? "THB",
    syncedAt: new Date().toISOString(),
  };
}
