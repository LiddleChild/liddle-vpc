"use client";

import { useMemo, useRef, useState } from "react";
import dayjs from "dayjs";
import {
  CalendarDays,
  ChevronDown,
  ChevronLeft,
  ChevronRight,
  RefreshCw,
} from "lucide-react";
import {
  ExpenseFilters,
  type ExpenseFilters as ExpenseFiltersValue,
} from "@/components/expense-filters";

type Account = {
  name: string;
  color: string;
  amount: number;
  category: string;
  tags: string[];
};
const accounts: Account[] = [
  {
    name: "Daily spending",
    color: "#e9a34d",
    amount: 428.35,
    category: "Everyday",
    tags: ["Needs", "Card"],
  },
  {
    name: "Home & bills",
    color: "#8c80e7",
    amount: 1360,
    category: "Bills",
    tags: ["Needs", "Recurring"],
  },
  {
    name: "Shared expenses",
    color: "#75b9a8",
    amount: 225.9,
    category: "Everyday",
    tags: ["Shared"],
  },
  {
    name: "Cash & other",
    color: "#de7d6d",
    amount: 86.2,
    category: "Other",
    tags: ["Cash"],
  },
];
const formatMoney = (value: number) =>
  new Intl.NumberFormat("th-TH", {
    style: "currency",
    currency: "THB",
    maximumFractionDigits: 2,
  }).format(value);

export default function Home() {
  const [monthOffset, setMonthOffset] = useState(0);
  const [showFilters, setShowFilters] = useState(false);
  const screenViewport = useRef<HTMLDivElement>(null);
  const [filters, setFilters] = useState<ExpenseFiltersValue>({
    accounts: accounts.map((account) => account.name),
    category: "all",
    tags: [],
  });
  const month = useMemo(
    () => dayjs("2025-09-01").add(monthOffset, "month").format("MMMM YYYY"),
    [monthOffset],
  );
  const visibleAccounts = accounts.filter(
    (account) =>
      filters.accounts.includes(account.name) &&
      (filters.category === "all" || account.category === filters.category) &&
      filters.tags.every((tag) => account.tags.includes(tag)),
  );
  const totalExpenses = visibleAccounts.reduce(
    (sum, account) => sum + account.amount,
    0,
  );
  const categories = [...new Set(accounts.map((account) => account.category))];
  const tags = [...new Set(accounts.flatMap((account) => account.tags))];

  const handleScroll = () => {
    const viewport = screenViewport.current;
    if (viewport)
      setShowFilters(viewport.scrollLeft > viewport.clientWidth / 2);
  };
  return (
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
              <section className="expense-hero" aria-label="Current expenses">
                <div className="expense-circle">
                  <span>Expenses</span>
                  <strong>{formatMoney(totalExpenses)}</strong>
                  <small>this month</small>
                </div>
                <div className="sync-line">
                  <RefreshCw size={13} /> Updated 8 min ago
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
              accountNames={accounts.map((account) => account.name)}
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
