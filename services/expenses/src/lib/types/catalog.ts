export type Account = {
  id: number;
  name: string;
};

export type CatalogChoice = {
  value: string;
  textColor?: string;
  fillColor?: string;
};

export type ExpenseCatalog = {
  tags: CatalogChoice[];
  accounts: Account[];
};

export type ExpenseSummaryFilters = {
  month: string;
  accounts: number[];
  tags: string[];
};

export type MonthlyExpenseSummary = {
  month: string;
  expenses: number;
  transactionCount: number;
  accountCount: number;
  accountExpenses: AccountExpense[];
  syncedAt: string;
};

export type AccountExpense = {
  accountId: number;
  accountName: string;
  expenses: number;
  transactionCount: number;
};
