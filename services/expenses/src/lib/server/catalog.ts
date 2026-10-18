import {
  getExpenseCatalog as getNexusExpenseCatalog,
  getMonthlyExpenseSummary as getNexusMonthlyExpenseSummary,
} from "@/lib/server/nexus-client";
import type {
  ExpenseCatalog,
  ExpenseSummaryFilters,
  MonthlyExpenseSummary,
} from "@/lib/types/catalog";

export function getExpenseCatalog(): Promise<ExpenseCatalog> {
  return getNexusExpenseCatalog();
}

export function getMonthlyExpenseSummary(
  filters: ExpenseSummaryFilters,
): Promise<MonthlyExpenseSummary> {
  return getNexusMonthlyExpenseSummary(filters);
}
