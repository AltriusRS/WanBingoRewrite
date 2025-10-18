"use client"

import {useEffect, useState} from "react"
import {Card} from "@/components/ui/card"
import {BingoBoardProps, BingoTile as IBingoTile, fetchBoardFromAPI} from "@/lib/bingoUtils";
import {useChat} from "@/components/chat/chat-context";
import {BingoTile} from "./tile"

import {BingoStatusBar} from "@/components/bingo/status-bar";

export function BingoBoard({onWin}: BingoBoardProps) {
    const ctx = useChat();
    const [tiles, setTiles] = useState<IBingoTile[]>([])
    const [hasWon, setHasWon] = useState(false)

    const resetBoard = async (e: unknown) => {
        let newTiles = await fetchBoardFromAPI()

        setTiles(newTiles)
        setHasWon(false)
    }

    const toggleTile = (id: number) => {
        let mappedTile = tiles.filter((f) => f.id === id)[0];
        if (!mappedTile) return console.error("Unable to find tile matching id", id);
        if (mappedTile.id === -1) return

        setTiles((prev) => prev.map((tile) => (tile.id === id ? {...tile, marked: !tile.marked} : tile)))
    }

    useEffect(() => {
        resetBoard(undefined).then(r => {
        })
    }, [])

    return (
        <div className="flex h-full flex-col grow justify-start items-center">
                <Card
                    className="w-full h-full flex items-center justify-start p-3 sm:p-4 lg:p-6">
                    {/* Header with New Board button and stats */}
                    <BingoStatusBar ctx={ctx} hasWon={hasWon} resetBoard={resetBoard} tiles={tiles}/>

                    {/* Board Area */}
                    <div className="grid grid-cols-5 aspect-square grid-rows-5 gap-1.5 sm:gap-2 lg:gap-3 w-full">
                        {tiles.map((tile) => (
                            <BingoTile key={tile.id} tile={tile} toggle={() => toggleTile(tile.id)}/>
                        ))}
                    </div>
                </Card>
        </div>
    )
}
