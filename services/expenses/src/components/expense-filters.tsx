"use client"

import { Check, ListFilter, RotateCcw } from "lucide-react"
import { useState, type CSSProperties } from "react"
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

type FilterSectionProps<T> = {
  title: string
  items: T[]
  selected: Array<string | number>
  getKey: (item: T) => string | number
  getLabel: (item: T) => string
  getStyle?: (item: T) => CSSProperties | undefined
  onToggle: (item: T) => void
  onSelectAll: () => void
  onClear: () => void
}

function FilterSection<T>({ title, items, selected, getKey, getLabel, getStyle, onToggle, onSelectAll, onClear }: FilterSectionProps<T>) {
  const allSelected = items.length > 0 && selected.length === items.length

  return (
    <section className="filter-section" aria-labelledby={`${title.toLowerCase()}-filter-label`}>
      <div className="filter-section-heading">
        <div>
          <h3 id={`${title.toLowerCase()}-filter-label`}>{title}</h3>
          <p>{selected.length === 0 ? "Nothing selected" : allSelected ? "Everything included" : `${selected.length} of ${items.length} selected`}</p>
        </div>
        <button className="filter-section-action" type="button" onClick={allSelected ? onClear : onSelectAll}>
          {allSelected ? "Clear" : "Select all"}
        </button>
      </div>
      {items.length === 0 ? <p className="filter-empty">No {title.toLowerCase()} available.</p> : (
        <div className="filter-options" role="group" aria-label={`${title} filters`}>
          {items.map((item) => {
            const isSelected = selected.includes(getKey(item))
            return <button className={`filter-option ${isSelected ? "selected" : ""}`} key={getKey(item)} type="button" style={isSelected ? getStyle?.(item) : undefined} onClick={() => onToggle(item)} aria-pressed={isSelected}>
              <span className="filter-option-check" aria-hidden="true">{isSelected && <Check size={13} strokeWidth={3} />}</span>
              <span>{getLabel(item)}</span>
            </button>
          })}
        </div>
      )}
    </section>
  )
}

export function ExpenseFilters({ accounts, categories, tags, value, onApply }: ExpenseFiltersProps) {
  const [draft, setDraft] = useState(value)

  const update = <K extends keyof ExpenseFilters>(key: K, item: ExpenseFilters[K][number]) => {
    const current = draft[key]
    const next = current.includes(item as never) ? current.filter((entry) => entry !== item) : [...current, item]
    const nextFilters = { ...draft, [key]: next } as ExpenseFilters
    setDraft(nextFilters)
    onApply(nextFilters)
  }

  const setAll = <K extends keyof ExpenseFilters>(key: K, items: ExpenseFilters[K]) => {
    const nextFilters = { ...draft, [key]: items } as ExpenseFilters
    setDraft(nextFilters)
    onApply(nextFilters)
  }

  const clearAll = () => {
    const nextFilters = { accounts: [], categories: [], tags: [] }
    setDraft(nextFilters)
    onApply(nextFilters)
  }

  return (
    <section className="filter-page" role="dialog" aria-modal="true" aria-labelledby="filters-title">
      <header className="filter-settings-header">
        <div className="filter-title-row"><span className="filter-title-icon"><ListFilter size={17} /></span><div><p className="filter-kicker">Refine your view</p><h2 id="filters-title">Filters</h2></div></div>
        <button className="reset-button" type="button" onClick={clearAll}><RotateCcw size={14} /> Clear all</button>
      </header>
      <p className="filter-intro">Choose what to include in your monthly total. Changes apply instantly.</p>
      <FilterSection title="Accounts" items={accounts} selected={draft.accounts} getKey={(account) => account.id} getLabel={(account) => account.name} onToggle={(account) => update("accounts", account.id)} onSelectAll={() => setAll<"accounts">("accounts", accounts.map(({ id }) => id))} onClear={() => setAll<"accounts">("accounts", [])} />
      <FilterSection title="Categories" items={categories} selected={draft.categories} getKey={(category) => category.value} getLabel={(category) => category.value} getStyle={(category) => ({ color: category.textColor, backgroundColor: category.fillColor, borderColor: category.fillColor })} onToggle={(category) => update("categories", category.value)} onSelectAll={() => setAll<"categories">("categories", categories.map(({ value }) => value))} onClear={() => setAll<"categories">("categories", [])} />
      <FilterSection title="Tags" items={tags} selected={draft.tags} getKey={(tag) => tag.value} getLabel={(tag) => tag.value} getStyle={(tag) => ({ color: tag.textColor, backgroundColor: tag.fillColor, borderColor: tag.fillColor })} onToggle={(tag) => update("tags", tag.value)} onSelectAll={() => setAll<"tags">("tags", tags.map(({ value }) => value))} onClear={() => setAll<"tags">("tags", [])} />
    </section>
  )
}
