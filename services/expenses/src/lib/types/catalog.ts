export type Account = {
  id: string;
  name: string;
};

export type CatalogChoice = {
  value: string;
  textColor?: string;
  fillColor?: string;
};

export type ExpenseCatalog = {
  categories: CatalogChoice[];
  tags: CatalogChoice[];
  accounts: Account[];
};

export type ExpenseSummaryFilters = {
  month: string;
  accounts: string[];
  categories: string[];
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
