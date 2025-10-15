"use client"

import { useState } from "react"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import type { BingoTile } from "@/lib/bingoUtils"

interface TileConfirmationDialogProps {
  tile: BingoTile | null
  open: boolean
  onConfirm: (context: string) => void
  onCancel: () => void
}

export function TileConfirmationDialog({ tile, open, onConfirm, onCancel }: TileConfirmationDialogProps) {
  const [context, setContext] = useState("")

  const handleConfirm = () => {
    onConfirm(context)
    setContext("")
  }

  const handleCancel = () => {
    setContext("")
    onCancel()
  }

  if (!tile) return null

  return (
    <Dialog open={open} onOpenChange={(open) => !open && handleCancel()}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Confirm Tile</DialogTitle>
        </DialogHeader>

        <div className="space-y-4 py-4">
          <div className="rounded-lg border border-border bg-muted p-4">
            <p className="font-medium text-foreground">{tile.text}</p>
          </div>

          <div className="space-y-2">
            <Label htmlFor="context">Context (optional)</Label>
            <Input
              id="context"
              placeholder="e.g., during sponsor segment"
              value={context}
              onChange={(e) => setContext(e.target.value)}
              maxLength={100}
            />
            <p className="text-xs text-muted-foreground">Add optional context about when this occurred</p>
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={handleCancel}>
            Cancel
          </Button>
          <Button onClick={handleConfirm}>Confirm Tile</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
