"use client"

import { useState } from "react"
import { Check } from "lucide-react"

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
  onClose: () => void
}

export function ExpenseFilters({ accountNames, categories, tags, value, onApply, onClose }: ExpenseFiltersProps) {
  const [draft, setDraft] = useState(value)

  const toggle = (items: string[], item: string) => items.includes(item) ? items.filter((entry) => entry !== item) : [...items, item]
  const apply = () => {
    onApply(draft)
    onClose()
  }

  const reset = () => setDraft({ accounts: accountNames, category: "all", tags: [] })

  return <section className="filter-page" role="dialog" aria-modal="true" aria-labelledby="filter-title">
        <div className="filter-panel-header"><div><p className="eyebrow">Swipe right to return</p><h2 id="filter-title">Filters</h2></div></div>
        <div className="filter-section"><label className="filter-label" htmlFor="category">Category</label><select id="category" value={draft.category} onChange={(event) => setDraft({ ...draft, category: event.target.value })}><option value="all">All categories</option>{categories.map((category) => <option value={category} key={category}>{category}</option>)}</select></div>
        <div className="filter-section"><span className="filter-label">Tags</span><div className="tag-options">{tags.map((tag) => { const selected = draft.tags.includes(tag); return <button className={`tag-option ${selected ? "selected" : ""}`} key={tag} onClick={() => setDraft({ ...draft, tags: toggle(draft.tags, tag) })} aria-pressed={selected}>{selected && <Check size={13} />}{tag}</button> })}</div></div>
        <div className="filter-section"><span className="filter-label">Accounts</span><div className="account-options">{accountNames.map((account) => { const selected = draft.accounts.includes(account); return <button className="filter-account" key={account} onClick={() => setDraft({ ...draft, accounts: toggle(draft.accounts, account) })} aria-pressed={selected}><span className={`filter-checkbox ${selected ? "selected" : ""}`}>{selected && <Check size={12} strokeWidth={3} />}</span>{account}</button> })}</div></div>
        <div className="filter-panel-footer"><button className="reset-button" onClick={reset}>Reset</button><button className="apply-button" onClick={apply}>Apply filters</button></div>
      </section>
}
