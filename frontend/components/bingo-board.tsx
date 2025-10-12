"use client"

import {useEffect, useState} from "react"
import {Card} from "@/components/ui/card"
import {Button} from "@/components/ui/button"
import {RotateCcw} from "lucide-react"
import confetti from "canvas-confetti"
import {BingoBoardProps, BingoTile, randomizeTiles} from "@/lib/bingoUtils";
import {useChat} from "@/components/chat/chat-context";
import {Tooltip, TooltipContent, TooltipTrigger,} from "@/components/ui/tooltip"


export function BingoBoard({onWin}: BingoBoardProps) {
    const ctx = useChat();
    const [tiles, setTiles] = useState<BingoTile[]>([])
    const [hasWon, setHasWon] = useState(false)

    const resetBoard = () => {
        let newTiles: BingoTile[] = randomizeTiles(tiles, ctx)

        setTiles(newTiles)
        setHasWon(false)
    }

    const checkWin = (currentTiles: BingoTile[]) => {
        // Check rows
        for (let i = 0; i < 5; i++) {
            const row = currentTiles.slice(i * 5, i * 5 + 5)
            if (row.every((tile) => tile.marked)) return true
        }

        // Check columns
        for (let i = 0; i < 5; i++) {
            const column = [
                currentTiles[i],
                currentTiles[i + 5],
                currentTiles[i + 10],
                currentTiles[i + 15],
                currentTiles[i + 20],
            ]
            if (column.every((tile) => tile.marked)) return true
        }

        // Check diagonals
        const diagonal1 = [currentTiles[0], currentTiles[6], currentTiles[12], currentTiles[18], currentTiles[24]]
        const diagonal2 = [currentTiles[4], currentTiles[8], currentTiles[12], currentTiles[16], currentTiles[20]]

        if (diagonal1.every((tile) => tile.marked) || diagonal2.every((tile) => tile.marked)) return true

        return false
    }

    const triggerConfetti = () => {
        const duration = 3000
        const animationEnd = Date.now() + duration
        const defaults = {startVelocity: 30, spread: 360, ticks: 60, zIndex: 0}

        const randomInRange = (min: number, max: number) => Math.random() * (max - min) + min

        const interval = window.setInterval(() => {
            const timeLeft = animationEnd - Date.now()

            if (timeLeft <= 0) {
                return clearInterval(interval)
            }

            const particleCount = 50 * (timeLeft / duration)

            confetti({
                ...defaults,
                particleCount,
                origin: {x: randomInRange(0.1, 0.3), y: Math.random() - 0.2},
            })
            confetti({
                ...defaults,
                particleCount,
                origin: {x: randomInRange(0.7, 0.9), y: Math.random() - 0.2},
            })
        }, 250)
    }

    const toggleTile = (id: number) => {
        let mappedTile = tiles.filter((f) => f.id === id)[0];
        if (!mappedTile) return console.error("Unable to find tile matching id", id);
        if (mappedTile.id === -1) return

        setTiles((prev) => {
            const newTiles = prev.map((tile) => (tile.id === id ? {...tile, marked: !tile.marked} : tile))

            if (!hasWon && checkWin(newTiles)) {
                setHasWon(true)
                triggerConfetti()
                onWin?.()
            }

            return newTiles
        })
    }

    useEffect(() => {
        resetBoard()
    }, [])

    return (
        <div className="flex h-full max-h-[calc(100vh-8rem)] flex-col space-y-4 md:max-h-[calc(100vh-6rem)]">
            {/* Header with New Board button and stats */}
            <div className="flex shrink-0 items-center justify-between">
                <div className="flex items-center gap-3">
                    <div className="flex gap-3 text-xs text-muted-foreground sm:text-sm">
                        <div>
                            <span
                                className="font-medium text-foreground">{tiles.filter((t) => t.marked).length}</span> /
                            25
                        </div>
                        {hasWon && (
                            <>
                                <div className="h-4 w-px bg-border"/>
                                <div className="font-medium text-primary">🎉 BINGO!</div>
                            </>
                        )}
                    </div>
                    <Tooltip>
                        <TooltipTrigger asChild>
                            <Button onClick={(e) => {
                                if (ctx.episode.isLive) {
                                    return;
                                }
                                resetBoard(e);
                            }} variant="outline" size="sm" className="gap-2 bg-transparent"
                            >
                                <RotateCcw className="h-4 w-4"/>
                                <span className="hidden sm:inline">New Board</span>
                                <span className="sm:hidden">New</span>
                            </Button>
                        </TooltipTrigger>
                        <TooltipContent>
                            {ctx.episode.isLive ? "You can't reset the board while the show is live" : "Reset the board to a new randomized set of tiles"}
                        </TooltipContent>
                    </Tooltip>

                </div>
            </div>

            <Card className="flex-1 overflow-auto p-3 sm:p-4 lg:p-6">
                <div className="grid h-full grid-cols-5 gap-1.5 sm:gap-2 lg:gap-3">
                    {tiles.map((tile) => (
                        <button
                            key={tile.id}
                            onClick={() => toggleTile(tile.id)}
                            className={`
                relative aspect-square rounded-lg border-2 p-1.5 text-center text-[10px] font-medium transition-all
                hover:scale-105 active:scale-95 sm:p-2 sm:text-xs lg:text-sm
                ${
                                tile.marked
                                    ? "border-primary bg-primary/20 text-foreground shadow-lg shadow-primary/20"
                                    : "border-border bg-card text-card-foreground hover:border-primary/50 hover:bg-secondary"
                            }
                ${tile.text === "FREE" ? "cursor-default border-primary bg-primary text-primary-foreground" : "cursor-pointer"}
              `}
                        >
                            <div className="flex h-full items-center justify-center">
                                <span className="text-balance leading-tight">{tile.text}</span>
                            </div>
                            {tile.marked && tile.text !== "FREE" && (
                                <div className="absolute inset-0 flex items-center justify-center">
                                    <div className="h-1 w-full rotate-45 bg-primary/40"/>
                                    <div className="absolute h-1 w-full -rotate-45 bg-primary/40"/>
                                </div>
                            )}
                        </button>
                    ))}
                </div>
            </Card>
        </div>
    )
}
