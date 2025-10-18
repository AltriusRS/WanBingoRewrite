"use client"

import { useState, useEffect } from "react"
import { Card } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { getApiRoot } from "@/lib/auth"

interface Suggestion {
  id: string
  name: string
  tile_name: string
  reason: string
  status: string
  reviewed_by?: string
  reviewed_at?: string
  created_at: string
}

interface TileSettings {
  requiresTimer?: boolean
  requiresContext?: boolean
  description?: string
  confirmationRules?: string
  timer?: {
    duration: number
    name: string
    description: string
  }
  contextExample?: string
}

interface SuggestionAcceptanceModalProps {
  suggestion: Suggestion
  onClose: () => void
  onAccept: () => void
}

export function SuggestionAcceptanceModal({ suggestion, onClose, onAccept }: SuggestionAcceptanceModalProps) {
  const [text, setText] = useState(suggestion.tile_name)
  const [category, setCategory] = useState("General")
  const [weight, setWeight] = useState("1")
  const [score, setScore] = useState("5")
  const [description, setDescription] = useState("")
  const [confirmationRules, setConfirmationRules] = useState("")
  const [requiresTimer, setRequiresTimer] = useState(false)
  const [requiresContext, setRequiresContext] = useState(false)
  const [timerName, setTimerName] = useState("")
  const [timerDuration, setTimerDuration] = useState("60")
  const [timerDescription, setTimerDescription] = useState("")
  const [contextExample, setContextExample] = useState("")
  const [saving, setSaving] = useState(false)

  const existingCategories = ["Linus", "Luke", "Dan", "Late", "Sponsors", "Topics", "Set/Production", "Events"]

  useEffect(() => {
    if (requiresTimer && !timerName) {
      setTimerName(text)
    }
  }, [text, requiresTimer])

  const handleSave = async () => {
    setSaving(true)
    const newSettings: TileSettings = {
      description: description || undefined,
      confirmationRules: confirmationRules || undefined,
      requiresTimer,
      requiresContext,
      contextExample: requiresContext ? contextExample : undefined,
      timer: requiresTimer ? {
        name: timerName,
        duration: Number.parseInt(timerDuration),
        description: timerDescription,
      } : undefined,
    }

    const data = {
      text,
      category: category === "General" ? null : category,
      weight: Number.parseFloat(weight),
      score: Number.parseFloat(score),
      settings: newSettings,
    }

    try {
      // Create the tile
      const response = await fetch(`${getApiRoot()}/host/tiles`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify(data),
      })

      if (!response.ok) {
        console.error("Failed to create tile:", response.status)
        return
      }

      // Call onAccept to update suggestion status
      onAccept()
    } catch (error) {
      console.error("Failed to accept suggestion:", error)
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-background/80">
      <Card className="w-full max-w-md p-6 max-h-[90vh] overflow-y-auto">
        <h3 className="mb-4 text-lg font-semibold">Accept Tile Suggestion</h3>
        <p className="mb-4 text-sm text-muted-foreground">
          Suggested by: {suggestion.name}<br />
          Reason: {suggestion.reason}
        </p>

        <div className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="tile-text">Tile Title</Label>
            <Input id="tile-text" value={text} onChange={(e) => setText(e.target.value)} />
          </div>

          <div className="space-y-2">
            <Label htmlFor="tile-category">Category</Label>
            <Input
              id="tile-category"
              value={category}
              onChange={(e) => setCategory(e.target.value)}
              list="categories"
            />
            <datalist id="categories">
              {existingCategories.map((cat) => (
                <option key={cat} value={cat} />
              ))}
            </datalist>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="tile-weight">Weight (difficulty)</Label>
              <Input
                id="tile-weight"
                type="number"
                step="0.1"
                value={weight}
                onChange={(e) => setWeight(e.target.value)}
              />
              <p className="text-xs text-muted-foreground">Higher weight = more likely</p>
            </div>
            <div className="space-y-2">
              <Label htmlFor="tile-score">Score</Label>
              <Input
                id="tile-score"
                type="number"
                step="1"
                value={score}
                onChange={(e) => setScore(e.target.value)}
              />
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="tile-description">Description</Label>
            <Input id="tile-description" value={description} onChange={(e) => setDescription(e.target.value)} />
          </div>

          <div className="space-y-2">
            <Label htmlFor="confirmation-rules">Confirmation Rules</Label>
            <Input id="confirmation-rules" value={confirmationRules} onChange={(e) => setConfirmationRules(e.target.value)} />
          </div>

          <div className="flex gap-4">
            <div className="flex items-center space-x-2">
              <input type="checkbox" id="requires-timer" checked={requiresTimer} onChange={(e) => setRequiresTimer(e.target.checked)} />
              <Label htmlFor="requires-timer">Requires Timer</Label>
            </div>
            <div className="flex items-center space-x-2">
              <input type="checkbox" id="requires-context" checked={requiresContext} onChange={(e) => setRequiresContext(e.target.checked)} />
              <Label htmlFor="requires-context">Requires Context</Label>
            </div>
          </div>

          {requiresTimer && (
            <div className="space-y-4 p-4 border rounded">
              <h4 className="font-medium">Timer Configuration</h4>
              <div className="space-y-2">
                <Label htmlFor="timer-name">Timer Name</Label>
                <Input id="timer-name" value={timerName} onChange={(e) => setTimerName(e.target.value)} />
              </div>
              <div className="space-y-2">
                <Label htmlFor="timer-duration">Duration (seconds)</Label>
                <Input id="timer-duration" type="number" value={timerDuration} onChange={(e) => setTimerDuration(e.target.value)} />
              </div>
              <div className="space-y-2">
                <Label htmlFor="timer-description">Timer Description</Label>
                <Input id="timer-description" value={timerDescription} onChange={(e) => setTimerDescription(e.target.value)} />
              </div>
            </div>
          )}

          {requiresContext && (
            <div className="space-y-4 p-4 border rounded">
              <h4 className="font-medium">Context Configuration</h4>
              <div className="space-y-2">
                <Label htmlFor="context-example">Context Example</Label>
                <Input id="context-example" value={contextExample} onChange={(e) => setContextExample(e.target.value)} />
              </div>
            </div>
          )}
        </div>

        <div className="mt-6 flex justify-end gap-2">
          <Button variant="outline" onClick={onClose} disabled={saving}>
            Cancel
          </Button>
          <Button onClick={handleSave} disabled={saving}>
            {saving ? "Creating..." : "Accept & Create Tile"}
          </Button>
        </div>
      </Card>
    </div>
  )
}