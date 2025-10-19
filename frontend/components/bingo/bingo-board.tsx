"use client"

import {useCallback, useEffect, useState, useMemo} from "react"
import {Card} from "@/components/ui/card"
import {BingoBoardProps, BingoTile as IBingoTile, fetchBoardFromAPI, regenerateBoardAPI, BoardData, checkBingoWin, fetchConfirmedTiles, isValidWin} from "@/lib/bingoUtils";
import {useChat} from "@/components/chat/chat-context";
import {getApiRoot} from "@/lib/auth";
import {useAuth} from "@/components/auth";
import {BingoTile} from "./tile"

import {BingoStatusBar} from "@/components/bingo/status-bar";

export function BingoBoard({onWin}: BingoBoardProps) {
     const ctx = useChat()
     const { user } = useAuth()
     const [tiles, setTiles] = useState<IBingoTile[]>([])
     const [hasWon, setHasWon] = useState(false)
     const [regenerationCount, setRegenerationCount] = useState(0)
     const [regenerationDiminisher, setRegenerationDiminisher] = useState(1)
     const [confirmedTiles, setConfirmedTiles] = useState<Set<string>>(new Set())

    const shouldHighlightConfirmedTiles = (): boolean => {
         if (!user?.settings) return true // Default to true if no settings
         const settings = user.settings as any
         if (settings.gameplay) {
             return settings.gameplay.highlightConfirmedTiles !== false
         }
         // Fallback for old flat structure
         return settings.highlightConfirmedTiles !== false
     }

     const shouldShowTileScores = (): boolean => {
         // Hide scores for anonymous users
         if (!user || user.id === 'anonymous') return false

         if (!user.settings) return true // Default to true if no settings
         const settings = user.settings as any
         if (settings.gameplay) {
             return settings.gameplay.showTileScores !== false
         }
         return true // Default to true for new setting
     }

     const shouldShowMaxScore = (): boolean => {
         // Hide max score for anonymous users
         if (!user || user.id === 'anonymous') return false

         if (!user.settings) return true // Default to true if no settings
         const settings = user.settings as any
         if (settings.gameplay) {
             return settings.gameplay.showMaxScore !== false
         }
         return true // Default to true for new setting
     }

     const shouldShowMultiplier = (): boolean => {
         // Hide multiplier for anonymous users
         if (!user || user.id === 'anonymous') return false

         if (!user.settings) return true // Default to true if no settings
         const settings = user.settings as any
         if (settings.gameplay) {
             return settings.gameplay.showMultiplier !== false
         }
         return true // Default to true for new setting
     }

      const shouldShowRegenerations = (): boolean => {
          // Hide regeneration counter for anonymous users
          if (!user || user.id === 'anonymous') return false
          return true
      }

      const getBoardTextSize = (): string => {
          if (!user?.settings) return 'medium' // Default
          const settings = user.settings as any
          if (settings.appearance?.board?.textSize) {
              return settings.appearance.board.textSize
          }
          return 'medium'
      }

    const fetchConfirmed = useCallback(async () => {
        const confirmed = await fetchConfirmedTiles()
        setConfirmedTiles(confirmed)
    }, [])

     const resetBoard = useCallback(async (e: unknown) => {
         const boardData = await fetchBoardFromAPI()
         setTiles(boardData.tiles)
         setRegenerationCount(getRegenerationCount(boardData.regenerationDiminisher))
         setRegenerationDiminisher(boardData.regenerationDiminisher)
         setHasWon(false)
     }, [])

    const getRegenerationCount = (diminisher: number): number => {
        if (diminisher === 1) return 0
        if (diminisher === 0.9) return 1
        if (diminisher === 0.8) return 2
        return 3 // or more
    }

     const regenerateBoard = async () => {
         try {
             const boardData = await regenerateBoardAPI()
             setTiles(boardData.tiles)
             setRegenerationCount(getRegenerationCount(boardData.regenerationDiminisher))
             setRegenerationDiminisher(boardData.regenerationDiminisher)
             setHasWon(false)
         } catch (error) {
             console.error("Failed to regenerate board:", error)
             // Handle error, perhaps show toast
         }
     }



     const recordWin = useCallback(async () => {
         try {
             const disableAnnouncements = user?.settings?.gameplay?.disableWinAnnouncements || false
             const response = await fetch(`${getApiRoot()}/tiles/win`, {
                 method: 'POST',
                 credentials: "include",
                 headers: {"Content-Type": "application/json"},
                 body: JSON.stringify({
                     disableAnnouncement: disableAnnouncements
                 })
             })
             if (!response.ok) {
                 console.error("Failed to record win:", response.status)
             }
         } catch (error) {
             console.error("Failed to record win:", error)
         }
     }, [user?.settings?.gameplay?.disableWinAnnouncements])

    const toggleTile = (id: string) => {
        let mappedTile = tiles.filter((f) => f.id === id)[0];
        if (!mappedTile) return console.error("Unable to find tile matching id", id);
        if (mappedTile.id === "-1") return

        setTiles((prev) => {
            const newTiles = prev.map((tile) => (tile.id === id ? {...tile, marked: !tile.marked} : tile))

            // Check for bingo win
            const hasWin = checkBingoWin(newTiles)
            if (hasWin && !hasWon) {
                // Check if this is a valid win (all tiles in winning lines are confirmed)
                if (isValidWin(newTiles, confirmedTiles)) {
                    setHasWon(true)
                    // Record win on server for authenticated users
                    recordWin()
                    if (onWin) {
                        onWin()
                    }
                }
            } else if (!hasWin && hasWon) {
                // If they unmarked tiles and no longer have a win
                setHasWon(false)
            }

            return newTiles
        })
    }

    useEffect(() => {
        resetBoard(undefined).then(r => {
        })
        fetchConfirmed().then(r => {
        })

        // Listen for tile confirmed events from SSE
        const handleTileConfirmed = () => {
            fetchConfirmed()
        }

        window.addEventListener('tileConfirmed', handleTileConfirmed)

        return () => {
            window.removeEventListener('tileConfirmed', handleTileConfirmed)
        }
    }, [resetBoard, fetchConfirmed])

     // Re-check for wins when confirmed tiles change
     useEffect(() => {
         const hasWin = checkBingoWin(tiles)
         if (hasWin && !hasWon) {
             if (isValidWin(tiles, confirmedTiles)) {
                 setHasWon(true)
                 recordWin()
                 if (onWin) {
                     onWin()
                 }
             }
         } else if (!hasWin && hasWon) {
             setHasWon(false)
         }
     }, [confirmedTiles, tiles, hasWon, onWin, recordWin])

    return (
        <div className="flex h-full flex-col grow justify-start items-center">
                <Card
                    className="w-full h-full flex items-center justify-start p-3 sm:p-4 lg:p-6">
                    {/* Header with New Board button and stats */}
                     <BingoStatusBar ctx={ctx} hasWon={hasWon} resetBoard={resetBoard} regenerateBoard={regenerateBoard} regenerationCount={regenerationCount} tiles={tiles} confirmedTiles={confirmedTiles} regenerationDiminisher={regenerationDiminisher} showMaxScore={shouldShowMaxScore()} showMultiplier={shouldShowMultiplier()} showRegenerations={shouldShowRegenerations()}/>

                    {/* Board Area */}
                    <div className="grid grid-cols-5 aspect-square grid-rows-5 gap-1.5 sm:gap-2 lg:gap-3 w-full">
                          {tiles.map((tile) => (
                              <BingoTile
                                  key={tile.id}
                                  tile={tile}
                                  toggle={() => toggleTile(tile.id)}
                                  highlighted={shouldHighlightConfirmedTiles() && confirmedTiles.has(tile.id)}
                                  showScore={shouldShowTileScores()}
                                  textSize={getBoardTextSize()}
                              />
                          ))}
                    </div>
                </Card>
        </div>
    )
}
