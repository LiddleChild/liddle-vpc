"use server";

import {
  getExpenseCatalog,
  getMonthlyExpenseSummary,
  type ExpenseCatalog,
  type MonthlyExpenseSummary,
  type ExpenseSummaryFilters,
} from "@/lib/server/catalog";
import dayjs from "dayjs";

export async function listCatalog(): Promise<ExpenseCatalog> {
  return getExpenseCatalog();
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
