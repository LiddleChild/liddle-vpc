"use server";

import {
  getExpenseCatalog,
  getMonthlyExpenseSummary,
  type Account,
  type MonthlyExpenseSummary,
  type ExpenseSummaryFilters,
} from "@/lib/server/catalog";
import dayjs from "dayjs";

export async function listCategories(): Promise<string[]> {
  const catalog = await getExpenseCatalog();
  return catalog.categories;
}

export async function listTags(): Promise<string[]> {
  const catalog = await getExpenseCatalog();
  return catalog.tags;
}

export async function listAccounts(): Promise<Account[]> {
  const catalog = await getExpenseCatalog();
  return catalog.accounts;
}

export async function getExpensesSummary(
  filters: ExpenseSummaryFilters,
): Promise<MonthlyExpenseSummary> {
  if (
    !/^\d{4}-\d{2}$/.test(filters.month) ||
    !dayjs(`${filters.month}-01`).isValid()
  ) {
    throw new Error("Invalid month");
  }

  return getMonthlyExpenseSummary(filters);
}
