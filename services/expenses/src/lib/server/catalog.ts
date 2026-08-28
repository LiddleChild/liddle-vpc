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

/**
 * The provider boundary keeps the eventual Grist table and column names out
 * of the application. Replace this implementation with the Grist adapter
 * when the private document is configured.
 */
const mockCatalog: ExpenseCatalog = {
  categories: ["Everyday", "Bills", "Other"],
  tags: ["Needs", "Card", "Recurring", "Shared", "Cash"],
  accounts: [
    { id: "daily-spending", name: "Daily spending" },
    { id: "home-bills", name: "Home & bills" },
    { id: "shared-expenses", name: "Shared expenses" },
    { id: "cash-other", name: "Cash & other" },
  ],
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
  // Keep this asynchronous so the mock and live providers share the same API.
  return mockCatalog;
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
