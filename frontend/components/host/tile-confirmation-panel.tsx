"use client"

import {useEffect, useState} from "react"
import {Card} from "@/components/ui/card"
import {Button} from "@/components/ui/button"
import {ScrollArea} from "@/components/ui/scroll-area"
import {TileConfirmationDialog} from "./tile-confirmation-dialog"
import type {BingoTile} from "@/lib/bingoUtils"
import {CheckCircle2, Lock, Play, Hash} from "lucide-react"
import {getApiRoot} from "@/lib/auth";

interface TileLock {
    tileId: number
    lockedBy: string
    expiresAt: number
}

export function TileConfirmationPanel() {
    const [tiles, setTiles] = useState<BingoTile[]>([])
    const [loading, setLoading] = useState(true)
    const [selectedTile, setSelectedTile] = useState<BingoTile | null>(null)
    const [locks, setLocks] = useState<Map<string, TileLock>>(new Map())
    const [confirmedTiles, setConfirmedTiles] = useState<Set<string>>(new Set())
    const [showTileIds, setShowTileIds] = useState<Set<string>>(new Set())

    useEffect(() => {
        fetchTiles()
        fetchShowTiles()
        const interval = setInterval(cleanupExpiredLocks, 1000)
        return () => clearInterval(interval)
    }, [])

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
            const response = await fetch(`${getApiRoot()}/tiles/confirmations`, {
                method: "POST",
                headers: {"Content-Type": "application/json"},
                credentials: "include",
                body: JSON.stringify({
                    tile_id: selectedTile.id,
                    context,
                }),
            })

            if (!response.ok) {
                throw new Error("Failed to confirm tile")
            }

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

    // Calculate category stats
    const categoryStats = Object.keys(tilesByCategory).reduce((acc, cat) => {
        const tilesInCat = tilesByCategory[cat]
        acc[cat] = {
            total: tilesInCat.length,
            inPlay: tilesInCat.filter(t => showTileIds.has(t.id)).length
        }
        return acc
    }, {} as Record<string, { total: number; inPlay: number }>)

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
                <div className="space-y-4">
                    {categories.map((category) => (
                        <div key={category} className="bg-muted p-4 rounded-lg">
                            <h3 className="mb-3 font-semibold text-foreground text-lg flex items-center gap-2">
                                {category}
                                <span className="text-sm text-muted-foreground flex items-center gap-1">
                                    <Play className="h-4 w-4 text-green-500" />
                                    {categoryStats[category].inPlay}
                                    <span className="mx-1">/</span>
                                    <Hash className="h-4 w-4" />
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
                                                variant={isConfirmed ? "secondary" : "outline"}
                                                className="justify-start text-left h-auto py-2"
                                                onClick={() => handleTileClick(tile)}
                                                disabled={isLocked && locks.get(tile.id)?.lockedBy !== "You"}
                                            >
                                                <span className="flex-1 truncate text-sm">{tile.title}</span>
                                                {isInPlay && <Play className="ml-2 h-4 w-4 text-green-500 flex-shrink-0" />}
                                                {isConfirmed && <CheckCircle2 className="ml-2 h-4 w-4 text-primary flex-shrink-0" />}
                                                {isLocked && <Lock className="ml-2 h-4 w-4 text-muted-foreground flex-shrink-0" />}
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
                onConfirm={handleConfirm}
                onCancel={handleCancel}
            />
        </>
    )
}
