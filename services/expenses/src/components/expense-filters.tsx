"use client"

import { useState } from "react"
import { RotateCcw } from "lucide-react"
import type { CatalogChoice } from "@/lib/server/catalog"

export type ExpenseFilters = {
  accounts: string[]
  categories: string[]
  tags: string[]
}

type ExpenseFiltersProps = {
  accountNames: string[]
  categories: CatalogChoice[]
  tags: CatalogChoice[]
  value: ExpenseFilters
  onApply: (filters: ExpenseFilters) => void
}

export function ExpenseFilters({ accountNames, categories, tags, value, onApply }: ExpenseFiltersProps) {
  const [draft, setDraft] = useState(value)

  const toggle = (items: string[], item: string) => items.includes(item) ? items.filter((entry) => entry !== item) : [...items, item]
  const updateFilters = (next: ExpenseFilters) => {
    setDraft(next)
    onApply(next)
  }

  const reset = () => updateFilters({ accounts: accountNames, categories: categories.map(({ value }) => value), tags: tags.map(({ value }) => value) })

  return <section className="filter-page" role="dialog" aria-modal="true">
        <div className="filter-settings-header"><h2>Settings</h2><button className="reset-button" onClick={reset}><RotateCcw size={14} /> Clear all</button></div>
        <div className="filter-section"><span className="filter-label">Accounts</span><div className="account-options">{accountNames.map((account) => { const selected = draft.accounts.includes(account); return <button className="filter-account" key={account} onClick={() => updateFilters({ ...draft, accounts: toggle(draft.accounts, account) })} aria-pressed={selected}><span className={`filter-checkbox ${selected ? "selected" : ""}`} />{account}</button> })}</div></div>
        <div className="filter-section"><span className="filter-label">Categories</span><div className="tag-options">{categories.map((category) => { const selected = draft.categories.includes(category.value); return <button className={`tag-option ${selected ? "selected" : ""}`} style={selected ? { color: category.textColor, backgroundColor: category.fillColor, borderColor: category.fillColor } : undefined} key={category.value} onClick={() => updateFilters({ ...draft, categories: toggle(draft.categories, category.value) })} aria-pressed={selected}>{category.value}</button> })}</div></div>
        <div className="filter-section"><span className="filter-label">Tags</span><div className="tag-options">{tags.map((tag) => { const selected = draft.tags.includes(tag.value); return <button className={`tag-option ${selected ? "selected" : ""}`} style={selected ? { color: tag.textColor, backgroundColor: tag.fillColor, borderColor: tag.fillColor } : undefined} key={tag.value} onClick={() => updateFilters({ ...draft, tags: toggle(draft.tags, tag.value) })} aria-pressed={selected}>{tag.value}</button> })}</div></div>
      </section>
}
