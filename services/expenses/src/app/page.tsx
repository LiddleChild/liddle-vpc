"use client";

import { useEffect, useMemo, useRef, useState } from "react";
import dayjs from "dayjs";
import {
  CalendarDays,
  ChevronDown,
  ChevronLeft,
  ChevronRight,
  RefreshCw,
  WifiOff,
} from "lucide-react";
import {
  ExpenseFilters,
  type ExpenseFilters as ExpenseFiltersValue,
} from "@/components/expense-filters";
import { Skeleton } from "@/components/ui/skeleton";
import {
  getExpensesSummary,
  listAccounts,
  listCategories,
  listTags,
} from "@/app/actions/catalog";
import type { MonthlyExpenseSummary } from "@/lib/server/catalog";

const formatMoney = (value: number, currency: string) =>
  new Intl.NumberFormat("th-TH", {
    style: "currency",
    currency,
    maximumFractionDigits: 2,
  }).format(value);

export default function Home() {
  const [monthOffset, setMonthOffset] = useState(0);
  const [showFilters, setShowFilters] = useState(false);
  const [isOffline, setIsOffline] = useState(false);
  const [accountNames, setAccountNames] = useState<string[]>([]);
  const [categories, setCategories] = useState<string[]>([]);
  const [tags, setTags] = useState<string[]>([]);
  const [summary, setSummary] = useState<MonthlyExpenseSummary | null>(null);
  const [summaryError, setSummaryError] = useState(false);
  const [isSummaryLoading, setIsSummaryLoading] = useState(true);
  const screenViewport = useRef<HTMLDivElement>(null);
  const [filters, setFilters] = useState<ExpenseFiltersValue>({
    accounts: [],
    category: "all",
    tags: [],
  });
  const selectedMonth = useMemo(
    () => dayjs("2025-09-01").add(monthOffset, "month").format("YYYY-MM"),
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

    Promise.all([listCategories(), listTags(), listAccounts()]).then(
      ([nextCategories, nextTags, nextAccounts]) => {
        if (!isMounted) return;

        const nextAccountNames = nextAccounts.map((account) => account.name);
        setCategories(nextCategories);
        setTags(nextTags);
        setAccountNames(nextAccountNames);
        setFilters((current) => ({ ...current, accounts: nextAccountNames }));
      },
    );

    return () => {
      isMounted = false;
    };
  }, []);

  const handleScroll = () => {
    const viewport = screenViewport.current;
    if (viewport)
      setShowFilters(viewport.scrollLeft > viewport.clientWidth / 2);
  };
  return isOffline ? (
    <main className="connection-state"><div className="connection-icon"><WifiOff size={26} /></div><h1>Can’t connect</h1><p>Connect to the internet to view your current expenses.</p><button onClick={() => window.location.reload()}>Try again</button></main>
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
                    {isSummaryLoading
                      ? <Skeleton className="expense-amount-skeleton" />
                      : summaryError
                        ? "—"
                        : formatMoney(summary?.expenses ?? 0, summary?.currency ?? "THB")}
                  </strong>
                  <small>this month</small>
                </div>
                <div className="sync-line">
                  <RefreshCw size={13} />
                  {isSummaryLoading ? (
                    <>
                      <span className="sr-only">Loading expenses</span>
                      <Skeleton className="sync-skeleton" />
                    </>
                  ) : summaryError
                    ? "Unable to load expenses"
                    : summary
                      ? `Updated ${dayjs(summary.syncedAt).format("h:mm A")}`
                      : "Loading expenses"}
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
              key={accountNames.join("|")}
              accountNames={accountNames}
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
