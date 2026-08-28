"use client"

import { useState } from "react"
import { Check, RotateCcw } from "lucide-react"

export type ExpenseFilters = {
  accounts: string[]
  category: string
  tags: string[]
}

type ExpenseFiltersProps = {
  accountNames: string[]
  categories: string[]
  tags: string[]
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

  const reset = () => updateFilters({ accounts: accountNames, category: "all", tags: [] })

  return <section className="filter-page" role="dialog" aria-modal="true">
        <div className="filter-settings-header"><h2>Settings</h2><button className="reset-button" onClick={reset}><RotateCcw size={14} /> Clear all</button></div>
        <div className="filter-section"><label className="filter-label" htmlFor="category">Category</label><select id="category" value={draft.category} onChange={(event) => updateFilters({ ...draft, category: event.target.value })}><option value="all">All categories</option>{categories.map((category) => <option value={category} key={category}>{category}</option>)}</select></div>
        <div className="filter-section"><span className="filter-label">Tags</span><div className="tag-options">{tags.map((tag) => { const selected = draft.tags.includes(tag); return <button className={`tag-option ${selected ? "selected" : ""}`} key={tag} onClick={() => updateFilters({ ...draft, tags: toggle(draft.tags, tag) })} aria-pressed={selected}>{selected && <Check size={13} />}{tag}</button> })}</div></div>
        <div className="filter-section"><span className="filter-label">Accounts</span><div className="account-options">{accountNames.map((account) => { const selected = draft.accounts.includes(account); return <button className="filter-account" key={account} onClick={() => updateFilters({ ...draft, accounts: toggle(draft.accounts, account) })} aria-pressed={selected}><span className={`filter-checkbox ${selected ? "selected" : ""}`}>{selected && <Check size={12} strokeWidth={3} />}</span>{account}</button> })}</div></div>
      </section>
}
