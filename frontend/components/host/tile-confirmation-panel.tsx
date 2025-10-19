"use client"

import {useEffect, useState} from "react"
import {Card} from "@/components/ui/card"
import {Button} from "@/components/ui/button"
import {ScrollArea} from "@/components/ui/scroll-area"
import {TileConfirmationDialog} from "./tile-confirmation-dialog"
import type {BingoTile} from "@/lib/bingoUtils"
import {CheckCircle2, Hash, Lock, Play} from "lucide-react"
import {getApiRoot} from "@/lib/auth";
import { toast } from "sonner"
import { useHost } from "./host-context"



export function TileConfirmationPanel() {
  const { confirmedTiles, locks } = useHost()
  const [tiles, setTiles] = useState<BingoTile[]>([])
  // const [stats, setStats] = useState<Map<string, TileStats>>(new Map())
  const [loading, setLoading] = useState(true)
  const [selectedTile, setSelectedTile] = useState<BingoTile | null>(null)
  const [showTileIds, setShowTileIds] = useState<Set<string>>(new Set())
  const [revokeMode, setRevokeMode] = useState(false)

  useEffect(() => {
    fetchTiles()
    fetchShowTiles()
  }, [])

  // Debug logging for state changes
  useEffect(() => {
    console.log("[TileConfirmationPanel] Locks updated:", Array.from(locks.entries()))
  }, [locks])

  useEffect(() => {
    console.log("[TileConfirmationPanel] Confirmed tiles updated:", Array.from(confirmedTiles))
  }, [confirmedTiles])

    const fetchTiles = async () => {
        try {
            const allTiles: BingoTile[] = []
            let page = 1
            let hasNext = true

            while (hasNext) {
                const response = await fetch(`${getApiRoot()}/tiles?page=${page}&limit=100`)
                const data = await response.json()
                allTiles.push(...data.tiles)
                hasNext = data.pagination.has_next
                page++
            }

            setTiles(allTiles)
        } catch (error) {
            console.error("Failed to fetch tiles:", error)
        } finally {
            setLoading(false)
        }
    }

  const fetchShowTiles = async () => {
    try {
      const response = await fetch(`${getApiRoot()}/tiles/show`)
      const data = await response.json()
      setShowTileIds(new Set(data.tile_ids))
    } catch (error) {
      console.error("Failed to fetch show tiles:", error)
    }
  }

  const handleTileClick = async (tile: BingoTile) => {
    console.log(`[TileConfirmationPanel] handleTileClick called for tile: ${tile.id} (${tile.title})`)
    const lock = locks.get(tile.id)
    console.log(`[TileConfirmationPanel] Current lock state for ${tile.id}:`, lock)
    console.log(`[TileConfirmationPanel] All locks:`, Array.from(locks.entries()))

    if (lock) {
      console.log(`[TileConfirmationPanel] Tile ${tile.id} is locked by ${lock.lockedBy}`)
      alert(`This tile is currently locked by ${lock.lockedBy}`)
      return
    }

    try {
      console.log(`[TileConfirmationPanel] Attempting to lock tile ${tile.id}`)
      // Call backend to acquire lock
      const response = await fetch(`${getApiRoot()}/host/tile-locks`, {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        credentials: "include",
        body: JSON.stringify({ tile_id: tile.id }),
      })

      if (!response.ok) {
        const errorText = await response.text()
        throw new Error(`Failed to lock tile: ${response.status} ${errorText}`)
      }

      console.log(`[TileConfirmationPanel] Lock request successful for tile ${tile.id}`)
      // The lock will be updated via SSE event from the host context
      setSelectedTile(tile)
      setRevokeMode(confirmedTiles.has(tile.id))
    } catch (error) {
      console.error("Failed to lock tile:", error)
      toast.error(error instanceof Error ? error.message : "Failed to lock tile")
    }
  }

  const handleConfirm = async (context: string) => {
    if (!selectedTile) return

    console.log(`[TileConfirmationPanel] handleConfirm called for tile ${selectedTile.id}, revokeMode: ${revokeMode}`)

    try {
      if (revokeMode) {
        console.log(`[TileConfirmationPanel] Revoking confirmation for tile ${selectedTile.id}`)
        // Revoke confirmation
        const response = await fetch(`${getApiRoot()}/host/confirmed-tiles/${selectedTile.id}`, {
          method: "DELETE",
          credentials: "include",
        })

        if (!response.ok) {
          const errorText = await response.text()
          throw new Error(`Failed to revoke confirmation: ${response.status} ${errorText}`)
        }

        console.log(`[TileConfirmationPanel] Revocation successful for tile ${selectedTile.id}`)
        toast.success(`${selectedTile.title} confirmation has been revoked.`)
      } else {
        console.log(`[TileConfirmationPanel] Confirming tile ${selectedTile.id} with context: "${context}"`)
        // Confirm tile
        const response = await fetch(`${getApiRoot()}/tiles/confirmations`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          credentials: "include",
          body: JSON.stringify({
            tile_id: selectedTile.id,
            context,
          }),
        })

        if (!response.ok) {
          const errorText = await response.text()
          throw new Error(`Failed to confirm tile: ${response.status} ${errorText}`)
        }

        console.log(`[TileConfirmationPanel] Confirmation successful for tile ${selectedTile.id}`)
        toast.success(`${selectedTile.title} has been confirmed successfully.`)
      }

      // Unlock the tile
      console.log(`[TileConfirmationPanel] Unlocking tile ${selectedTile.id}`)
      try {
        await fetch(`${getApiRoot()}/host/tile-unlocks`, {
          method: "POST",
          headers: {"Content-Type": "application/json"},
          credentials: "include",
          body: JSON.stringify({ tile_id: selectedTile.id }),
        })
        console.log(`[TileConfirmationPanel] Unlock request sent for tile ${selectedTile.id}`)
      } catch (error) {
        console.error("Failed to unlock tile:", error)
      }

      // The unlock will be updated via SSE event from the host context
    } catch (error) {
      console.error("Failed to confirm/revoke tile:", error)
      toast.error(error instanceof Error ? error.message : "An error occurred")
    }

    setSelectedTile(null)
    setRevokeMode(false)
  }

  const handleCancel = async () => {
    if (selectedTile) {
      try {
        // Call backend to release lock
        await fetch(`${getApiRoot()}/host/tile-unlocks`, {
          method: "POST",
          headers: {"Content-Type": "application/json"},
          credentials: "include",
          body: JSON.stringify({ tile_id: selectedTile.id }),
        })
      } catch (error) {
        console.error("Failed to unlock tile:", error)
      }

      // The unlock will be updated via SSE event from the host context
    }
    setSelectedTile(null)
    setRevokeMode(false)
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

    // Calculate category stats
    const categoryStats = Object.keys(tilesByCategory).reduce((acc, cat) => {
        const tilesInCat = tilesByCategory[cat]
        acc[cat] = {
            total: tilesInCat.length,
            inPlay: tilesInCat.filter(t => showTileIds.has(t.id)).length
        }
        return acc
    }, {} as Record<string, { total: number; inPlay: number }>)

    const categories = Object.keys(tilesByCategory).sort().filter(cat => cat !== 'Late')

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
                <div className="space-y-4">
                    {categories.map((category) => (
                        <div key={category} className="bg-muted p-4 rounded-lg">
                                    <h3 className="mb-3 font-semibold text-foreground text-lg flex items-center gap-2">
                                        {category}
                                        <span className="text-sm text-muted-foreground flex items-center gap-1">
                                    <Play className="h-4 w-4 text-green-500"/>
                                            {categoryStats[category].inPlay}
                                            <span className="mx-1">/</span>
                                    <Hash className="h-4 w-4"/>
                                            {categoryStats[category].total}
                                </span>
                                    </h3>
                                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-2">
                                        {tilesByCategory[category]
                                            .sort((a, b) => a.title.localeCompare(b.title))
                                             .map((tile) => {
                                                 const isLocked = locks.has(tile.id)
                                                 const isConfirmed = confirmedTiles.has(tile.id)
                                                 const isInPlay = showTileIds.has(tile.id)

                                                 return (
                                  <Button
                                    key={tile.id}
                                    variant="outline"
                                    className={`justify-start text-left h-auto py-2 ${isConfirmed ? "border-primary" : ""}`}
                                    onClick={() => handleTileClick(tile)}
                                    disabled={isLocked && locks.get(tile.id)?.lockedBy !== "You"}
                                  >
                                                         <span className="flex-1 truncate text-sm">{tile.title}</span>
                                                         {isInPlay && <Play
                                                             className="ml-2 h-4 w-4 text-green-500 flex-shrink-0"/>}
                                                         {isConfirmed && <CheckCircle2
                                                             className="ml-2 h-4 w-4 text-primary flex-shrink-0"/>}
                                                         {isLocked && <Lock
                                                             className="ml-2 h-4 w-4 text-muted-foreground flex-shrink-0"/>}
                                                      </Button>
                                                 )
                                             })}
                                     </div>
                                 </div>
                     ))}
                </div>
            </ScrollArea>

      <TileConfirmationDialog
        tile={selectedTile}
        open={!!selectedTile}
        revokeMode={revokeMode}
        onConfirm={handleConfirm}
        onCancel={handleCancel}
      />
        </>
    )
}
