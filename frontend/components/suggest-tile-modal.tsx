"use client"

import type React from "react"

import { useState } from "react"
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import {getApiRoot} from "@/lib/auth";

interface SuggestTileModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onSubmit?: (data: { name: string; tileName: string; reason: string }) => void
}

export function SuggestTileModal({ open, onOpenChange, onSubmit }: SuggestTileModalProps) {
  const [name, setName] = useState("")
  const [tileName, setTileName] = useState("")
  const [reason, setReason] = useState("")
  const [submitting, setSubmitting] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (name.trim() && tileName.trim() && reason.trim()) {
      setSubmitting(true)
      const startTime = Date.now()

      try {
        await fetch(`${getApiRoot()}/suggestions`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ name, tileName, reason }),
        })
      } catch (error) {
        console.error("Failed to submit suggestion:", error)
      }

      // Call optional callback
      onSubmit?.({ name, tileName, reason })

      // Reset form
      setName("")
      setTileName("")
      setReason("")
      // Ensure minimum loading time of 200ms
      const elapsed = Date.now() - startTime
      const remaining = Math.max(0, 200 - elapsed)
      setTimeout(() => {
        setSubmitting(false)
        onOpenChange(false)
      }, remaining)
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>Suggest a New Tile</DialogTitle>
          <DialogDescription>
            Have an idea for a new bingo tile? Share it with the community! Your suggestion will be reviewed by the
            hosts.
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="name">Your Name</Label>
            <Input
              id="name"
              placeholder="Enter your name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="tileName">Tile Text</Label>
            <Input
              id="tileName"
              placeholder="e.g., 'Linus spills his drink'"
              value={tileName}
              onChange={(e) => setTileName(e.target.value)}
              required
              maxLength={50}
            />
            <p className="text-xs text-muted-foreground">{tileName.length}/50 characters</p>
          </div>
          <div className="space-y-2">
            <Label htmlFor="reason">Why should this be a tile?</Label>
            <Textarea
              id="reason"
              placeholder="Explain why this would make a great bingo tile..."
              value={reason}
              onChange={(e) => setReason(e.target.value)}
              required
              rows={4}
              maxLength={200}
            />
            <p className="text-xs text-muted-foreground">{reason.length}/200 characters</p>
          </div>
          <div className="flex justify-end gap-3">
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)} disabled={submitting}>
              Cancel
            </Button>
            <Button type="submit" disabled={submitting}>
              {submitting ? "Submitting..." : "Submit Suggestion"}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  )
}
