"use client"

import { useState, useEffect } from "react"
import { Card } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { ScrollArea } from "@/components/ui/scroll-area"
import { TileConfirmationDialog } from "./tile-confirmation-dialog"
import type { BingoTile } from "@/lib/bingoUtils"
import { Lock, CheckCircle2 } from "lucide-react"

interface TileLock {
  tileId: number
  lockedBy: string
  expiresAt: number
}

export function TileConfirmationPanel() {
  const [tiles, setTiles] = useState<BingoTile[]>([])
  const [loading, setLoading] = useState(true)
  const [selectedTile, setSelectedTile] = useState<BingoTile | null>(null)
  const [locks, setLocks] = useState<Map<number, TileLock>>(new Map())
  const [confirmedTiles, setConfirmedTiles] = useState<Set<number>>(new Set())

  useEffect(() => {
    fetchTiles()
    const interval = setInterval(cleanupExpiredLocks, 1000)
    return () => clearInterval(interval)
  }, [])

  const fetchTiles = async () => {
    try {
      const response = await fetch("http://localhost:8080/tiles")
      const data = await response.json()
      setTiles(data)
    } catch (error) {
      console.error("Failed to fetch tiles:", error)
    } finally {
      setLoading(false)
    }
  }

  const cleanupExpiredLocks = () => {
    const now = Date.now()
    setLocks((prev) => {
      const updated = new Map(prev)
      for (const [tileId, lock] of updated.entries()) {
        if (lock.expiresAt < now) {
          updated.delete(tileId)
        }
      }
      return updated
    })
  }

  const handleTileClick = (tile: BingoTile) => {
    const lock = locks.get(tile.id)
    if (lock) {
      alert(`This tile is currently locked by ${lock.lockedBy}`)
      return
    }

    // Create a 5-second lock
    const newLock: TileLock = {
      tileId: tile.id,
      lockedBy: "You",
      expiresAt: Date.now() + 5000,
    }
    setLocks((prev) => new Map(prev).set(tile.id, newLock))
    setSelectedTile(tile)
  }

  const handleConfirm = async (context: string) => {
    if (!selectedTile) return

    try {
      await fetch("http://localhost:8080/api/host/confirm-tile", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          tileId: selectedTile.id,
          tileName: selectedTile.text,
          context,
        }),
      })

      setConfirmedTiles((prev) => new Set(prev).add(selectedTile.id))
      setLocks((prev) => {
        const updated = new Map(prev)
        updated.delete(selectedTile.id)
        return updated
      })
    } catch (error) {
      console.error("Failed to confirm tile:", error)
    }

    setSelectedTile(null)
  }

  const handleCancel = () => {
    if (selectedTile) {
      setLocks((prev) => {
        const updated = new Map(prev)
        updated.delete(selectedTile.id)
        return updated
      })
    }
    setSelectedTile(null)
  }

  // Group tiles by category
  const tilesByCategory = tiles.reduce(
    (acc, tile) => {
      const category = tile.category || "General"
      if (!acc[category]) acc[category] = []
      acc[category].push(tile)
      return acc
    },
    {} as Record<string, BingoTile[]>,
  )

  const categories = Object.keys(tilesByCategory).sort()

  if (loading) {
    return (
      <Card className="flex items-center justify-center p-8">
        <p className="text-muted-foreground">Loading tiles...</p>
      </Card>
    )
  }

  return (
    <>
      <ScrollArea className="h-[calc(100vh-12rem)]">
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {categories.map((category) => (
            <Card key={category} className="p-4">
              <h3 className="mb-3 font-semibold text-foreground">{category}</h3>
              <div className="space-y-2">
                {tilesByCategory[category].map((tile) => {
                  const isLocked = locks.has(tile.id)
                  const isConfirmed = confirmedTiles.has(tile.id)

                  return (
                    <Button
                      key={tile.id}
                      variant={isConfirmed ? "secondary" : "outline"}
                      className="w-full justify-start text-left"
                      onClick={() => handleTileClick(tile)}
                      disabled={isLocked && locks.get(tile.id)?.lockedBy !== "You"}
                    >
                      <span className="flex-1 truncate text-sm">{tile.text}</span>
                      {isConfirmed && <CheckCircle2 className="ml-2 h-4 w-4 text-primary" />}
                      {isLocked && <Lock className="ml-2 h-4 w-4 text-muted-foreground" />}
                    </Button>
                  )
                })}
              </div>
            </Card>
          ))}
        </div>
      </ScrollArea>

      <TileConfirmationDialog
        tile={selectedTile}
        open={!!selectedTile}
        onConfirm={handleConfirm}
        onCancel={handleCancel}
      />
    </>
  )
}
