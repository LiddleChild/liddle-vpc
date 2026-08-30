"use client";

import { useEffect, useMemo, useRef, useState } from "react";
import dayjs from "dayjs";
import {
  CalendarDays,
  ChevronDown,
  ChevronLeft,
  ChevronRight,
  WifiOff,
} from "lucide-react";
import {
  ExpenseFilters,
  type ExpenseFilters as ExpenseFiltersValue,
} from "@/components/expense-filters";
import { Skeleton } from "@/components/ui/skeleton";
import { getExpensesSummary, listCatalog } from "@/app/actions/catalog";
import type {
  Account,
  CatalogChoice,
  MonthlyExpenseSummary,
} from "@/lib/types/catalog";

const formatMoney = (value: number) =>
  new Intl.NumberFormat("th-TH", {
    maximumFractionDigits: 2,
  }).format(value);

const FILTER_STORAGE_KEY = "expenses:selected-filters";

const readCachedFilters = (): ExpenseFiltersValue | null => {
  try {
    const cached = window.localStorage.getItem(FILTER_STORAGE_KEY);
    if (!cached) return null;

    const parsed: unknown = JSON.parse(cached);
    if (!parsed || typeof parsed !== "object") return null;

    const filters = parsed as Record<string, unknown>;
    if (
      !Array.isArray(filters.accounts) ||
      !filters.accounts.every((account) => Number.isInteger(account)) ||
      !Array.isArray(filters.categories) ||
      !filters.categories.every((category) => typeof category === "string") ||
      !Array.isArray(filters.tags) ||
      !filters.tags.every((tag) => typeof tag === "string")
    ) {
      return null;
    }

    return {
      accounts: filters.accounts as number[],
      categories: filters.categories as string[],
      tags: filters.tags as string[],
    };
  } catch {
    return null;
  }
};

const writeCachedFilters = (filters: ExpenseFiltersValue) => {
  try {
    window.localStorage.setItem(FILTER_STORAGE_KEY, JSON.stringify(filters));
  } catch {
    // Storage may be unavailable in private browsing or when quota is exceeded.
  }
};

const applyCatalogToCachedFilters = (
  cached: ExpenseFiltersValue | null,
  catalog: { accounts: Account[]; categories: CatalogChoice[]; tags: CatalogChoice[] },
): ExpenseFiltersValue => ({
  accounts: cached
    ? cached.accounts.filter((id) => catalog.accounts.some((account) => account.id === id))
    : catalog.accounts.map(({ id }) => id),
  categories: cached
    ? cached.categories.filter((value) => catalog.categories.some((category) => category.value === value))
    : catalog.categories.map(({ value }) => value),
  tags: cached
    ? cached.tags.filter((value) => catalog.tags.some((tag) => tag.value === value))
    : catalog.tags.map(({ value }) => value),
});

export default function Home() {
  const [monthOffset, setMonthOffset] = useState(0);
  const [showFilters, setShowFilters] = useState(false);
  const [isOffline, setIsOffline] = useState(false);
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [categories, setCategories] = useState<CatalogChoice[]>([]);
  const [tags, setTags] = useState<CatalogChoice[]>([]);
  const [summary, setSummary] = useState<MonthlyExpenseSummary | null>(null);
  const [summaryError, setSummaryError] = useState(false);
  const [isSummaryLoading, setIsSummaryLoading] = useState(true);
  const screenViewport = useRef<HTMLDivElement>(null);
  const hasLoadedCatalog = useRef(false);
  const [filters, setFilters] = useState<ExpenseFiltersValue>({
    accounts: [],
    categories: [],
    tags: [],
  });
  const selectedMonth = useMemo(
    () => dayjs().startOf("month").add(monthOffset, "month").format("YYYY-MM"),
    [monthOffset],
  );
  const month = dayjs(`${selectedMonth}-01`).format("MMMM YYYY");
  useEffect(() => {
    const updateConnection = () => setIsOffline(!navigator.onLine);
    updateConnection();
    window.addEventListener("online", updateConnection);
    window.addEventListener("offline", updateConnection);
    if (screenViewport.current) screenViewport.current.scrollLeft = 0;
    return () => {
      window.removeEventListener("online", updateConnection);
      window.removeEventListener("offline", updateConnection);
    };
  }, []);

  useEffect(() => {
    let isMounted = true;

    void (async () => {
      await Promise.resolve();
      if (!isMounted) return;

      setIsSummaryLoading(true);
      setSummaryError(false);

      try {
        const nextSummary = await getExpensesSummary({
          month: selectedMonth,
          ...filters,
        });
        if (!isMounted) return;
        setSummary(nextSummary);
        setIsSummaryLoading(false);
      } catch {
        if (!isMounted) return;
        setSummaryError(true);
        setIsSummaryLoading(false);
      }
    })();

    return () => {
      isMounted = false;
    };
  }, [filters, selectedMonth]);

  useEffect(() => {
    let isMounted = true;
    const cachedFilters = readCachedFilters();

    listCatalog().then((catalog) => {
      if (!isMounted) return;

      setCategories(catalog.categories);
      setTags(catalog.tags);
      setAccounts(catalog.accounts);
      setFilters(applyCatalogToCachedFilters(cachedFilters, catalog));
      hasLoadedCatalog.current = true;
    });

    return () => {
      isMounted = false;
    };
  }, []);

  useEffect(() => {
    if (!hasLoadedCatalog.current) return;
    writeCachedFilters(filters);
  }, [filters]);

  const handleScroll = () => {
    const viewport = screenViewport.current;
    if (viewport)
      setShowFilters(viewport.scrollLeft > viewport.clientWidth / 2);
  };
  return isOffline ? (
    <main className="connection-state">
      <div className="connection-icon">
        <WifiOff size={26} />
      </div>
      <h1>Can’t connect</h1>
      <p>Connect to the internet to view your current expenses.</p>
      <button onClick={() => window.location.reload()}>Try again</button>
    </main>
  ) : (
    <main className="app-shell">
      <div className="ambient-shape ambient-shape-one" />
      <div className="ambient-shape ambient-shape-two" />
      <div
        className="screen-viewport"
        ref={screenViewport}
        onScroll={handleScroll}
      >
        <div className="page-track">
          <div className="page-screen">
            <div className="dashboard">
              <section
                className="expense-hero"
                aria-label="Current expenses"
                aria-busy={isSummaryLoading}
              >
                <div className="expense-circle">
                  <span>Expenses</span>
                  <strong aria-live="polite">
                    {isSummaryLoading ? (
                      <Skeleton className="expense-amount-skeleton" />
                    ) : summaryError ? (
                      "—"
                    ) : (
                      formatMoney(summary?.expenses ?? 0)
                    )}
                  </strong>
                  <small>this month</small>
                </div>
              </section>
              <div className="month-picker" aria-label="Choose month">
                <button
                  className="month-arrow"
                  onClick={() => setMonthOffset((value) => value - 1)}
                  aria-label="Previous month"
                >
                  <ChevronLeft size={18} />
                </button>
                <div className="month-label">
                  <CalendarDays size={16} />
                  <span>{month}</span>
                  <ChevronDown size={15} />
                </div>
                <button
                  className="month-arrow"
                  onClick={() => setMonthOffset((value) => value + 1)}
                  aria-label="Next month"
                >
                  <ChevronRight size={18} />
                </button>
              </div>
            </div>
          </div>
          <div className="page-screen">
            <ExpenseFilters
              key={accounts.map((account) => account.id).join("|")}
              accounts={accounts}
              categories={categories}
              tags={tags}
              value={filters}
              onApply={setFilters}
            />
          </div>
        </div>
      </div>
      <div
        className="page-dots"
        role="status"
        aria-label={`Page ${showFilters ? 2 : 1} of 2`}
      >
        <span className={!showFilters ? "active" : ""} />
        <span className={showFilters ? "active" : ""} />
      </div>
    </main>
  );
}
