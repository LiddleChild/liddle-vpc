import { GristClient, type GristColumn } from "@/lib/server/grist-client";
import dayjs from "dayjs";
import type {
  Account,
  CatalogChoice,
  ExpenseCatalog,
  AccountExpense,
  ExpenseSummaryFilters,
  MonthlyExpenseSummary,
} from "@/lib/types/catalog";

const getGristClient = () =>
  new GristClient({
    baseUrl: process.env.GRIST_URL ?? "",
    docId: process.env.GRIST_DOC_ID ?? "",
    apiKey: process.env.GRIST_API_KEY ?? "",
  });

const getChoiceValues = (column: GristColumn): CatalogChoice[] => {
  if (!column.fields.widgetOptions) {
    throw new Error(`Column ${column.id} has no choice options`);
  }

  let widgetOptions: unknown;
  try {
    widgetOptions = JSON.parse(column.fields.widgetOptions);
  } catch {
    throw new Error(`Column ${column.id} has invalid choice options`);
  }

  if (typeof widgetOptions !== "object" || widgetOptions === null) {
    throw new Error(`Column ${column.id} has invalid choice options`);
  }

  const choices = (widgetOptions as { choices?: unknown }).choices;
  const choiceOptions = (widgetOptions as { choiceOptions?: unknown })
    .choiceOptions;
  if (
    !Array.isArray(choices) ||
    !choices.every((choice) => typeof choice === "string") ||
    (choiceOptions !== undefined &&
      (typeof choiceOptions !== "object" || choiceOptions === null))
  ) {
    throw new Error(`Column ${column.id} has invalid choice options`);
  }

  return choices.map((choice) => {
    const option = (choiceOptions as Record<string, unknown> | undefined)?.[
      choice
    ];
    if (option === undefined) return { value: choice };
    if (typeof option !== "object" || option === null) {
      throw new Error(`Column ${column.id} has invalid choice styling`);
    }

    const { textColor, fillColor } = option as {
      textColor?: unknown;
      fillColor?: unknown;
    };
    if (
      (textColor !== undefined && typeof textColor !== "string") ||
      (fillColor !== undefined && typeof fillColor !== "string")
    ) {
      throw new Error(`Column ${column.id} has invalid choice styling`);
    }

    return {
      value: choice,
      ...(textColor === undefined ? {} : { textColor }),
      ...(fillColor === undefined ? {} : { fillColor }),
    };
  });
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

const getAccountRecords = async (client: GristClient): Promise<Account[]> => {
  const records = await client.listTableRecords("Accounts");
  return records.map((record) => {
    const name = record.fields.Name;
    if (typeof name !== "string" || !name.trim()) {
      throw new Error("Accounts table contains a record without a valid Name");
    }
    if (typeof record.id !== "number" || !Number.isInteger(record.id)) {
      throw new Error("Accounts table contains a record with an invalid ID");
    }
    return { id: record.id, name };
  });
};

const quoteIdentifier = (identifier: string) => {
  if (!/^[A-Za-z][A-Za-z0-9_]*$/.test(identifier)) {
    throw new Error(`Invalid Grist identifier: ${identifier}`);
  }
  return `"${identifier}"`;
};

export async function getExpenseCatalog(): Promise<ExpenseCatalog> {
  const client = getGristClient();
  const [columns, accountRecords] = await Promise.all([
    client.listTableColumns("Transactions"),
    getAccountRecords(client),
  ]);

  return {
    categories: getChoiceValues(getColumn(columns, "Category")),
    tags: getChoiceValues(getColumn(columns, "Tags")),
    accounts: accountRecords,
  };
}

export async function getMonthlyExpenseSummary(
  filters: ExpenseSummaryFilters,
): Promise<MonthlyExpenseSummary> {
  const client = getGristClient();
  const accountRecordsPromise = getAccountRecords(client);
  const monthStart = dayjs(`${filters.month}-01`);
  const nextMonthStart = monthStart.add(1, "month");
  const transactionsTable = quoteIdentifier("Transactions");
  const sql = `
    SELECT
      t."Account" AS accountId,
      COALESCE(SUM(t."Amount_In_Account_Currency" * t."Inflow_Direction"), 0) AS expenses,
      COUNT(*) AS transactionCount
    FROM ${transactionsTable} t
    WHERE t."Date" >= ?
      AND t."Date" < ?
      AND t."Account" IN (SELECT value FROM json_each(?))
      AND (
        t."Category" = ""
        OR t."Category" IN (SELECT value FROM json_each(?))
      )
      AND (
        t."Tags" is null
        OR EXISTS (
          SELECT 1
          FROM json_each(t."Tags") transaction_tag
          JOIN json_each(?) selected_tag
          ON selected_tag.value = transaction_tag.value
        )
      )
    GROUP BY t."Account"
  `;

  const [accountRecords, result] = await Promise.all([
    accountRecordsPromise,
    client.runSql(sql, [
    monthStart.unix(),
    nextMonthStart.unix(),
    JSON.stringify(filters.accounts),
    JSON.stringify(filters.categories),
    JSON.stringify(filters.tags),
    ]),
  ]);

  const groupedExpenses = result.map(({ fields }): AccountExpense => {
    const accountId = fields.accountId;
    const expenses = fields.expenses;
    const transactionCount = fields.transactionCount;
    if (
      typeof accountId !== "number" ||
      !Number.isInteger(accountId) ||
      typeof expenses !== "number" ||
      !Number.isFinite(expenses) ||
      typeof transactionCount !== "number" ||
      !Number.isInteger(transactionCount)
    ) {
      throw new Error("Grist returned invalid account expense totals");
    }

    const account = accountRecords.find((candidate) => candidate.id === accountId);
    if (!account) {
      throw new Error("Grist returned an expense for an unknown account");
    }

    return {
      accountId,
      accountName: account.name,
      expenses,
      transactionCount,
    };
  });

  const accountExpenses = accountRecords
    .filter((account) => filters.accounts.includes(account.id))
    .map((account) => groupedExpenses.find((expense) => expense.accountId === account.id) ?? {
      accountId: account.id,
      accountName: account.name,
      expenses: 0,
      transactionCount: 0,
    });
  const expenses = groupedExpenses.reduce((total, account) => total + account.expenses, 0);
  const transactionCount = groupedExpenses.reduce((total, account) => total + account.transactionCount, 0);

  return {
    month: filters.month,
    expenses,
    transactionCount,
    accountCount: groupedExpenses.length,
    accountExpenses,
    syncedAt: new Date().toISOString(),
  };
}
