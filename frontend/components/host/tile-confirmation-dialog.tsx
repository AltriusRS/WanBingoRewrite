"use client"

import { useState } from "react"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import type { BingoTile } from "@/lib/bingoUtils"
import { Clock } from "lucide-react"
import { getApiRoot } from "@/lib/auth"
import { toast } from "sonner"

interface TileConfirmationDialogProps {
  tile: BingoTile | null
  open: boolean
  revokeMode: boolean
  onConfirm: (context: string) => void
  onCancel: () => void
  onStartTimer?: () => void
}

export function TileConfirmationDialog({ tile, open, revokeMode, onConfirm, onCancel, onStartTimer }: TileConfirmationDialogProps) {
  const [context, setContext] = useState("")
  const [isStartingTimer, setIsStartingTimer] = useState(false)

  const handleConfirm = () => {
    onConfirm(context)
    setContext("")
  }

  const handleCancel = () => {
    setContext("")
    onCancel()
  }

  const handleStartTimer = async () => {
    if (!onStartTimer) return
    setIsStartingTimer(true)
    try {
      await onStartTimer()
      toast.success("Timer started!")
    } catch (error) {
      toast.error("Failed to start timer")
    } finally {
      setIsStartingTimer(false)
    }
  }

  if (!tile) return null

  const tileSettings = tile.settings as any
  const requiresTimer = tileSettings?.requiresTimer
  const timerInfo = tileSettings?.timer

  return (
    <Dialog open={open} onOpenChange={(open) => !open && handleCancel()}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{revokeMode ? "Revoke Confirmation" : "Confirm Tile"}</DialogTitle>
        </DialogHeader>

        <div className="space-y-4 py-4">
          <div className="rounded-lg border border-border bg-muted p-4">
            <p className="font-medium text-foreground">{tile.title}</p>
            {requiresTimer && timerInfo && (
              <div className="flex items-center mt-2 text-sm text-muted-foreground">
                <Clock className="h-4 w-4 mr-1" />
                Requires {timerInfo.duration}s timer: {timerInfo.name}
              </div>
            )}
          </div>

          {!revokeMode && (
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
          )}

          {revokeMode && (
            <p className="text-sm text-muted-foreground">
              Are you sure you want to revoke the confirmation for this tile?
            </p>
          )}
        </div>

        <DialogFooter className="flex-col sm:flex-row gap-2">
          <div className="flex gap-2 w-full sm:w-auto">
            <Button variant="outline" onClick={handleCancel}>
              Cancel
            </Button>
            {requiresTimer && !revokeMode && (
              <Button
                variant="outline"
                onClick={handleStartTimer}
                disabled={isStartingTimer}
                className="flex items-center gap-2"
              >
                <Clock className="h-4 w-4" />
                {isStartingTimer ? "Starting..." : "Start Timer"}
              </Button>
            )}
          </div>
          <Button variant={revokeMode ? "destructive" : "default"} onClick={handleConfirm}>
            {revokeMode ? "Revoke Confirmation" : "Confirm Tile"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
