import type {ChatContextValue} from "@/components/chat/chat-context"

export interface BingoTile {
    id: number
    title: string
    marked: boolean
    weight?: number
    score?: number,
    category?: string
}

export interface BingoBoardProps {
    onWin?: () => void
}

export async function fetchTilesFromAPI(): Promise<BingoTile[]> {
    try {
        const response = await fetch("https://api.bingo.local/tiles")
        if (!response.ok) {
            throw new Error("Failed to fetch tiles")
        }
        const tiles = await response.json()
        return tiles.map((tile: any) => ({
            id: tile.id,
            title: tile.title,
            marked: false,
            weight: tile.weight || 1,
            score: tile.score || 5,
            category: tile.category || "General",
        }))
    } catch (error) {
        console.error("Error fetching tiles:", error)
        // Return empty array on error - component will handle fallback
        return []
    }
}

function weightedRandomSelection(tiles: BingoTile[], count: number): BingoTile[] {
    const selected: BingoTile[] = []
    const available = [...tiles]

    while (selected.length < count && available.length > 0) {
        // Calculate total weight
        const totalWeight = available.reduce((sum, tile) => sum + (tile.weight || 1), 0)

        // Random selection based on weight
        let random = Math.random() * totalWeight
        let selectedIndex = 0

        for (let i = 0; i < available.length; i++) {
            random -= available[i].weight || 1
            if (random <= 0) {
                selectedIndex = i
                break
            }
        }

        selected.push(available[selectedIndex])
        available.splice(selectedIndex, 1)
    }

    return selected
}

export function randomizeTiles(
    currentTiles: BingoTile[],
    ctx: ChatContextValue,
    availableTiles: BingoTile[],
): BingoTile[] {
    if (ctx.episode.isLive) return currentTiles

    const selected = weightedRandomSelection(availableTiles, 24)

    const newTiles: BingoTile[] = selected

    // Add FREE space in the center
    newTiles.splice(12, 0, {
        id: -1,
        title: "FREE",
        marked: true,
    })

    return newTiles
}
