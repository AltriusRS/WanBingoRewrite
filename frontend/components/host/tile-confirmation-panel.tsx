"use client"

import {useEffect, useState} from "react"
import {Card} from "@/components/ui/card"
import {Button} from "@/components/ui/button"
import {ScrollArea} from "@/components/ui/scroll-area"
import {TileConfirmationDialog} from "./tile-confirmation-dialog"
import type {BingoTile} from "@/lib/bingoUtils"
import {CheckCircle2, Clock, Hash, Lock, Pause, Play, RotateCcw, Trash2} from "lucide-react"
import {getApiRoot} from "@/lib/auth";
import {toast} from "sonner"
import {useHost} from "./host-context"
import {useAuth} from "@/components/auth"


interface Timer {
    id: string
    title: string
    duration: number
    created_by?: string
    show_id?: string
    starts_at?: string
    expires_at?: string
    is_active: boolean
    settings: any
    created_at: string
    updated_at: string
}

interface TileConfirmationPanelProps {
  showLateButton?: boolean
}

export function TileConfirmationPanel({ showLateButton }: TileConfirmationPanelProps = {}) {
    const {confirmedTiles, locks} = useHost()
    const {user} = useAuth()

    const getHostPanelTextSize = (): string => {
        if (!user?.settings) return 'medium'
        const settings = user.settings as any
        if (settings.appearance?.hostPanel?.textSize) {
            return settings.appearance.hostPanel.textSize
        }
        return 'medium'
    }

    const getTextSizeClass = (size?: string) => {
        switch (size) {
            case 'small': return 'text-xs'
            case 'large': return 'text-base'
            default: return 'text-sm'
        }
    }
    const [tiles, setTiles] = useState<BingoTile[]>([])
    // const [stats, setStats] = useState<Map<string, TileStats>>(new Map())
    const [loading, setLoading] = useState(true)
    const [selectedTile, setSelectedTile] = useState<BingoTile | null>(null)
    const [showTileIds, setShowTileIds] = useState<Set<string>>(new Set())
    const [revokeMode, setRevokeMode] = useState(false)
    const [timers, setTimers] = useState<Timer[]>([])
    const [currentTime, setCurrentTime] = useState(new Date())
    const [selectedTimer, setSelectedTimer] = useState<Timer | null>(null)
    const [lateTile, setLateTile] = useState<BingoTile | null>(null)

    useEffect(() => {
        fetchTiles()
        fetchShowTiles()
        fetchTimers()
        if (showLateButton) {
            fetchLateTile()
        }
        // Update current time every second for countdown
        const interval = setInterval(() => {
            setCurrentTime(new Date())
        }, 1000)
        return () => clearInterval(interval)
    }, [showLateButton])

    const formatTime = (expiresAt?: string) => {
        if (!expiresAt) return "00:00"
        const expires = new Date(expiresAt)
        const diff = Math.max(0, Math.floor((expires.getTime() - currentTime.getTime()) / 1000))
        const mins = Math.floor(diff / 60)
        const secs = diff % 60
        return `${mins}:${secs.toString().padStart(2, "0")}`
    }

    const handleTimerClick = (timer: Timer) => {
        setSelectedTimer(timer)
    }

    const toggleTimer = async (id: string) => {
        try {
            const timer = timers.find(t => t.id === id)
            if (!timer) return

            const action = timer.is_active ? "stop" : "start"
            await fetch(`${getApiRoot()}/timers/${id}/${action}`, {
                method: "POST",
                credentials: "include",
            })
            fetchTimers()
        } catch (error) {
            console.error("Failed to toggle timer:", error)
        }
    }

    const resetTimer = async (id: string) => {
        try {
            await fetch(`${getApiRoot()}/timers/${id}/reset`, {
                method: "POST",
                credentials: "include",
            })
            fetchTimers()
        } catch (error) {
            console.error("Failed to reset timer:", error)
        }
    }

    const deleteTimer = async (id: string) => {
        try {
            await fetch(`${getApiRoot()}/timers/${id}`, {
                method: "DELETE",
                credentials: "include",
            })
            fetchTimers()
        } catch (error) {
            console.error("Failed to delete timer:", error)
        }
    }

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

    const fetchTimers = async () => {
        try {
            const response = await fetch(`${getApiRoot()}/timers?is_active=true`)
            const data = await response.json()
            setTimers(data.timers)
        } catch (error) {
            console.error("Failed to fetch timers:", error)
        }
    }

    const fetchLateTile = async () => {
        try {
            const response = await fetch(`${getApiRoot()}/tiles?category=Late`)
            const data = await response.json()
            const tiles: BingoTile[] = data.tiles
            const late = tiles.find(t => t.title === "Show Is Late")
            if (late) setLateTile(late)
        } catch (error) {
            console.error("Failed to fetch late tile:", error)
        }
    }

    const handleConfirmLate = async () => {
        if (!lateTile) return

        console.log(`[TileConfirmationPanel] Confirming late tile ${lateTile.id}`)
        console.log(`[TileConfirmationPanel] Current confirmed tiles:`, Array.from(confirmedTiles))

        try {
            const response = await fetch(`${getApiRoot()}/tiles/confirmations`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
                body: JSON.stringify({
                    tile_id: lateTile.id,
                    context: "",
                }),
            })

            if (!response.ok) {
                const errorText = await response.text()
                throw new Error(`Failed to confirm tile: ${response.status} ${errorText}`)
            }

            console.log(`[TileConfirmationPanel] Late tile confirmation successful`)
            toast.success("Show Is Late has been confirmed successfully.")
        } catch (error) {
            console.error("Failed to confirm late tile:", error)
            toast.error(error instanceof Error ? error.message : "An error occurred")
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
                body: JSON.stringify({tile_id: tile.id}),
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

        try {
            if (revokeMode) {
                console.log(`[TileConfirmationPanel] Revoking tile ${selectedTile.id}`)
                // Revoke tile confirmation
                const response = await fetch(`${getApiRoot()}/host/confirmed-tiles/${selectedTile.id}`, {
                    method: "DELETE",
                    credentials: "include",
                })

                if (!response.ok) {
                    const errorText = await response.text()
                    throw new Error(`Failed to revoke tile: ${response.status} ${errorText}`)
                }

                console.log(`[TileConfirmationPanel] Revocation successful for tile ${selectedTile.id}`)
                toast.success(`${selectedTile.title} confirmation has been revoked.`)
            } else {
                console.log(`[TileConfirmationPanel] Confirming tile ${selectedTile.id} with context: "${context}"`)
                // Confirm tile
                const response = await fetch(`${getApiRoot()}/host/show-tiles`, {
                    method: "POST",
                    headers: {"Content-Type": "application/json"},
                    credentials: "include",
                    body: JSON.stringify({
                        tile_id: selectedTile.id,
                        context: context || undefined,
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
                    body: JSON.stringify({tile_id: selectedTile.id}),
                })
            } catch (unlockError) {
                console.warn("Failed to unlock tile after confirmation:", unlockError)
            }

            setSelectedTile(null)
        } catch (error) {
            console.error("Failed to process tile:", error)
            toast.error(error instanceof Error ? error.message : "Failed to process tile")
        }
    }

    const handleStartTimer = async () => {
        if (!selectedTile) return

        const tileSettings = selectedTile.settings as any
        const timerInfo = tileSettings?.timer
        if (!timerInfo) return

        try {
            // First, try to delete any existing timer with the same name
            const existingTimer = timers.find(t => t.title === timerInfo.name)
            if (existingTimer) {
                await fetch(`${getApiRoot()}/timers/${existingTimer.id}`, {
                    method: "DELETE",
                    credentials: "include",
                })
            }

            // Create new timer
            const response = await fetch(`${getApiRoot()}/timers`, {
                method: "POST",
                headers: {"Content-Type": "application/json"},
                credentials: "include",
                body: JSON.stringify({
                    title: timerInfo.name,
                    duration: timerInfo.duration,
                }),
            })

            if (!response.ok) {
                const errorText = await response.text()
                throw new Error(`Failed to create timer: ${response.status} ${errorText}`)
            }

            // Refresh timers list
            fetchTimers()
        } catch (error) {
            console.error("Failed to start timer:", error)
            throw error
        }
    }

    const handleCancel = async () => {
        if (selectedTile) {
            try {
                // Call backend to release lock
                await fetch(`${getApiRoot()}/host/tile-unlocks`, {
                    method: "POST",
                    headers: {"Content-Type": "application/json"},
                    credentials: "include",
                    body: JSON.stringify({tile_id: selectedTile.id}),
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
                    {showLateButton && lateTile && (
                        <div className="bg-muted p-4 rounded-lg">
                            <Button
                                variant="outline"
                                size="sm"
                                className="w-full gap-2"
                                onClick={handleConfirmLate}
                                disabled={confirmedTiles.has(lateTile.id)}
                            >
                                <Clock className="h-4 w-4" />
                                {confirmedTiles.has(lateTile.id) ? "Show Is Late Confirmed" : "Confirm Show Is Late"}
                            </Button>
                        </div>
                    )}
                    {timers.length > 0 && (
                        <div className="mb-4">
                            <h3 className="font-semibold text-foreground mb-2">Ongoing Timers</h3>
                            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-2">
                                 {timers.map((timer) => (
                                     <div key={timer.id}
                                          className="flex items-center justify-between p-2 bg-muted rounded border cursor-pointer hover:bg-muted/80 transition-colors"
                                          onClick={() => handleTimerClick(timer)}>
                                         <div className="flex-1">
                                             <span className="font-medium text-sm">{timer.title}</span>
                                             <div className="text-xs text-muted-foreground">
                                                 <span className="font-mono">{formatTime(timer.expires_at)}</span>
                                                 <span className="ml-1">/ {timer.duration}s</span>
                                             </div>
                                         </div>
                                         {timer.is_active && (
                                             <div className="w-2 h-2 bg-green-500 rounded-full flex-shrink-0"></div>
                                         )}
                                     </div>
                                 ))}
                            </div>
                        </div>
                    )}
                <div className="flex flex-row flex-wrap gap-4">
                    {categories.map((category) => (
                        <div key={category} className="max-w-[40dvw] bg-muted p-4 rounded-lg">
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
                            <div className="flex flex-row flex-wrap gap-2">
                                {tilesByCategory[category]
                                    .sort((a, b) => a.title.localeCompare(b.title))
                                    .map((tile) => {
                                        const isLocked = locks.has(tile.id)
                                        const isConfirmed = confirmedTiles.has(tile.id)
                                        const isInPlay = showTileIds.has(tile.id)

                                         const tileSettings = tile.settings as any
                                         const requiresTimer = tileSettings?.requiresTimer

                                         return (
                                             <Button
                                                 key={tile.id}
                                                 variant="outline"
                                                 className={`justify-start text-left h-auto py-2 ${isConfirmed ? "border-primary" : ""}`}
                                                 onClick={() => handleTileClick(tile)}
                                                 disabled={isLocked && locks.get(tile.id)?.lockedBy !== "You"}
                                             >
                                                  <span className={`flex-1 truncate ${getTextSizeClass(getHostPanelTextSize())}`}>{tile.title}</span>
                                                 {requiresTimer && <Clock
                                                     className="ml-2 h-4 w-4 text-orange-500 flex-shrink-0"/>}
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
                </div>
             </ScrollArea>

            <TileConfirmationDialog
                tile={selectedTile}
                open={!!selectedTile}
                revokeMode={revokeMode}
                onConfirm={handleConfirm}
                onCancel={handleCancel}
                onStartTimer={handleStartTimer}
            />

            {/* Timer Control Dialog */}
            {selectedTimer && (
                <TimerControlDialog
                    timer={selectedTimer}
                    open={!!selectedTimer}
                    onClose={() => setSelectedTimer(null)}
                    onToggle={toggleTimer}
                    onReset={resetTimer}
                    onDelete={deleteTimer}
                    currentTime={currentTime}
                />
            )}
        </>
    )
}

// Timer Control Dialog Component
interface TimerControlDialogProps {
    timer: Timer
    open: boolean
    onClose: () => void
    onToggle: (id: string) => void
    onReset: (id: string) => void
    onDelete: (id: string) => void
    currentTime: Date
}

function TimerControlDialog({ timer, open, onClose, onToggle, onReset, onDelete, currentTime }: TimerControlDialogProps) {
    const formatTime = (expiresAt?: string) => {
        if (!expiresAt) return "00:00"
        const expires = new Date(expiresAt)
        const diff = Math.max(0, Math.floor((expires.getTime() - currentTime.getTime()) / 1000))
        const mins = Math.floor(diff / 60)
        const secs = diff % 60
        return `${mins}:${secs.toString().padStart(2, "0")}`
    }

    const formatDateTime = (dateString: string) => {
        return new Date(dateString).toLocaleString()
    }

    if (!open) return null

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-background/80">
            <Card className="w-full max-w-md p-6">
                <div className="space-y-4">
                    <div className="flex items-center justify-between">
                        <h3 className="text-lg font-semibold">Timer Controls</h3>
                        <Button variant="ghost" size="icon" onClick={onClose}>
                            ×
                        </Button>
                    </div>

                    <div className="space-y-3">
                        <div className="rounded-lg border border-border bg-muted p-4">
                            <div className="flex items-center gap-2 mb-2">
                                <Clock className="h-5 w-5" />
                                <span className="font-medium">{timer.title}</span>
                            </div>

                            <div className="grid grid-cols-2 gap-4 text-sm">
                                <div>
                                    <span className="text-muted-foreground">Status:</span>
                                    <div className="flex items-center gap-2 mt-1">
                                        {timer.is_active ? (
                                            <>
                                                <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                                                <span>Active</span>
                                            </>
                                        ) : (
                                            <>
                                                <div className="w-2 h-2 bg-gray-400 rounded-full"></div>
                                                <span>Stopped</span>
                                            </>
                                        )}
                                    </div>
                                </div>

                                <div>
                                    <span className="text-muted-foreground">Time Left:</span>
                                    <div className="font-mono text-lg mt-1">
                                        {formatTime(timer.expires_at)}
                                    </div>
                                </div>

                                <div>
                                    <span className="text-muted-foreground">Duration:</span>
                                    <div className="mt-1">{timer.duration}s</div>
                                </div>

                                <div>
                                    <span className="text-muted-foreground">Created:</span>
                                    <div className="text-xs mt-1">
                                        {formatDateTime(timer.created_at)}
                                    </div>
                                </div>

                                {timer.starts_at && (
                                    <div>
                                        <span className="text-muted-foreground">Started:</span>
                                        <div className="text-xs mt-1">
                                            {formatDateTime(timer.starts_at)}
                                        </div>
                                    </div>
                                )}

                                {timer.expires_at && (
                                    <div>
                                        <span className="text-muted-foreground">Expires:</span>
                                        <div className="text-xs mt-1">
                                            {formatDateTime(timer.expires_at)}
                                        </div>
                                    </div>
                                )}
                            </div>
                        </div>

                        <div className="flex gap-2">
                            <Button
                                variant="outline"
                                onClick={() => onToggle(timer.id)}
                                className="flex-1"
                            >
                                {timer.is_active ? (
                                    <>
                                        <Pause className="h-4 w-4 mr-2" />
                                        Stop
                                    </>
                                ) : (
                                    <>
                                        <Play className="h-4 w-4 mr-2" />
                                        Start
                                    </>
                                )}
                            </Button>

                            <Button
                                variant="outline"
                                onClick={() => onReset(timer.id)}
                                className="flex-1"
                            >
                                <RotateCcw className="h-4 w-4 mr-2" />
                                Reset
                            </Button>

                            <Button
                                variant="destructive"
                                onClick={() => {
                                    onDelete(timer.id)
                                    onClose()
                                }}
                            >
                                <Trash2 className="h-4 w-4" />
                            </Button>
                        </div>
                    </div>
                </div>
            </Card>
        </div>
    )
}
