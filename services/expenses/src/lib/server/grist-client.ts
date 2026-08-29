export type GristColumnFields = {
  label?: string;
  type?: string;
  formula?: string;
  isFormula?: boolean;
  visibleCol?: number;
  widgetOptions?: string;
  [key: string]: unknown;
};

export type GristColumn = {
  id: string;
  fields: GristColumnFields;
};

type ListColumnsResponse = {
  columns: GristColumn[];
};

export type GristClientOptions = {
  baseUrl: string;
  docId: string;
  apiKey: string;
};

export type ListTableColumnsOptions = {
  includeHidden?: boolean;
};

export class GristApiError extends Error {
  readonly status: number;

  constructor(message: string, status: number) {
    super(message);
    this.name = "GristApiError";
    this.status = status;
  }
}

const isRecord = (value: unknown): value is Record<string, unknown> =>
  typeof value === "object" && value !== null;

const parseColumnsResponse = (value: unknown): ListColumnsResponse => {
  if (!isRecord(value) || !Array.isArray(value.columns)) {
    throw new Error("Grist returned an invalid columns response");
  }

  const columns = value.columns.map((column) => {
    if (
      !isRecord(column) ||
      typeof column.id !== "string" ||
      !isRecord(column.fields)
    ) {
      throw new Error("Grist returned an invalid column");
    }

    return {
      id: column.id,
      fields: column.fields as GristColumnFields,
    };
  });

  return { columns };
};

export class GristClient {
  private readonly baseUrl: string;
  private readonly docId: string;
  private readonly apiKey: string;

  constructor(options: GristClientOptions) {
    if (!options.baseUrl || !options.docId || !options.apiKey) {
      throw new Error("Grist client requires baseUrl, docId, and apiKey");
    }

    this.baseUrl = options.baseUrl.replace(/\/+$/, "");
    this.docId = options.docId;
    this.apiKey = options.apiKey;
  }

  async listTableColumns(
    tableId: string,
    options: ListTableColumnsOptions = {},
  ): Promise<GristColumn[]> {
    if (!tableId) {
      throw new Error("A table ID is required");
    }

    const url = new URL(
      `${this.baseUrl}/api/docs/${encodeURIComponent(this.docId)}/tables/${encodeURIComponent(tableId)}/columns`,
    );
    if (options.includeHidden) {
      url.searchParams.set("hidden", "true");
    }

    const response = await fetch(url, {
      headers: {
        Accept: "application/json",
        Authorization: `Bearer ${this.apiKey}`,
      },
      cache: "no-store",
    });

    if (!response.ok) {
      throw new GristApiError(
        `Grist columns request failed with status ${response.status}`,
        response.status,
      );
    }

    let payload: unknown;
    try {
      payload = await response.json();
    } catch {
      throw new GristApiError("Grist returned invalid JSON", response.status);
    }

    try {
      return parseColumnsResponse(payload).columns;
    } catch (error) {
      throw new GristApiError(
        error instanceof Error ? error.message : "Grist returned invalid data",
        response.status,
      );
    }
  }
}
