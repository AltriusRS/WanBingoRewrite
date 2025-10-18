"use client"

import { SuggestionsPanel } from "@/components/host/suggestions-panel"

export default function SuggestionsPage() {
  return (
    <div className="container mx-auto p-6">
      <h1 className="text-2xl font-bold mb-6">Tile Suggestions</h1>
      <SuggestionsPanel />
    </div>
  )
}