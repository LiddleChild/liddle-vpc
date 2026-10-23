import {
  getExpenseCatalog as getBackendExpenseCatalog,
  getMonthlyExpenseSummary as getBackendMonthlyExpenseSummary,
} from "@/lib/server/backend-client";
import type {
  ExpenseCatalog,
  ExpenseSummaryFilters,
  MonthlyExpenseSummary,
} from "@/lib/types/catalog";

export function getExpenseCatalog(): Promise<ExpenseCatalog> {
  return getBackendExpenseCatalog();
}

export function getMonthlyExpenseSummary(
  filters: ExpenseSummaryFilters,
): Promise<MonthlyExpenseSummary> {
  return getBackendMonthlyExpenseSummary(filters);
}
