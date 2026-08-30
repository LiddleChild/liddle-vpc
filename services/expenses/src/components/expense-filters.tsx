"use client"

import { useState } from "react"
import { RotateCcw } from "lucide-react"
import type { Account, CatalogChoice } from "@/lib/types/catalog"

export type ExpenseFilters = {
  accounts: number[]
  categories: string[]
  tags: string[]
}

type ExpenseFiltersProps = {
  accounts: Account[]
  categories: CatalogChoice[]
  tags: CatalogChoice[]
  value: ExpenseFilters
  onApply: (filters: ExpenseFilters) => void
}

export function ExpenseFilters({ accounts, categories, tags, value, onApply }: ExpenseFiltersProps) {
  const [draft, setDraft] = useState(value)

  const toggle = <T,>(items: T[], item: T) => items.includes(item) ? items.filter((entry) => entry !== item) : [...items, item]
  const updateFilters = (next: ExpenseFilters) => {
    setDraft(next)
    onApply(next)
  }

  const reset = () => updateFilters({ accounts: accounts.map(({ id }) => id), categories: categories.map(({ value }) => value), tags: tags.map(({ value }) => value) })

  return <section className="filter-page" role="dialog" aria-modal="true">
        <div className="filter-settings-header"><h2>Settings</h2><button className="reset-button" onClick={reset}><RotateCcw size={14} /> Clear all</button></div>
        <div className="filter-section"><span className="filter-label">Accounts</span><div className="tag-options">{accounts.map((account) => { const selected = draft.accounts.includes(account.id); return <button className={`tag-option ${selected ? "selected" : ""}`} key={account.id} onClick={() => updateFilters({ ...draft, accounts: toggle(draft.accounts, account.id) })} aria-pressed={selected}>{account.name}</button> })}</div></div>
        <div className="filter-section"><span className="filter-label">Categories</span><div className="tag-options">{categories.map((category) => { const selected = draft.categories.includes(category.value); return <button className={`tag-option ${selected ? "selected" : ""}`} style={selected ? { color: category.textColor, backgroundColor: category.fillColor, borderColor: category.fillColor } : undefined} key={category.value} onClick={() => updateFilters({ ...draft, categories: toggle(draft.categories, category.value) })} aria-pressed={selected}>{category.value}</button> })}</div></div>
        <div className="filter-section"><span className="filter-label">Tags</span><div className="tag-options">{tags.map((tag) => { const selected = draft.tags.includes(tag.value); return <button className={`tag-option ${selected ? "selected" : ""}`} style={selected ? { color: tag.textColor, backgroundColor: tag.fillColor, borderColor: tag.fillColor } : undefined} key={tag.value} onClick={() => updateFilters({ ...draft, tags: toggle(draft.tags, tag.value) })} aria-pressed={selected}>{tag.value}</button> })}</div></div>
      </section>
}
