export interface BingoTile {
    id: string
    title: string
    marked: boolean
    weight?: number
    score?: number,
    category?: string
}

export interface BingoBoardProps {
    onWin?: () => void
}

export interface BoardData {
    tiles: BingoTile[]
    regenerationDiminisher: number
}

export async function fetchBoardFromAPI(): Promise<BoardData> {
    try {
        const response = await fetch("https://api.bingo.local/tiles/me", {
            credentials: "include"
        })
        if (!response.ok) {
            throw new Error("Failed to fetch board")
        }
        const data = await response.json()
        const tiles = data.tiles.map((tile: any) => ({
            id: tile.id,
            title: tile.title,
            marked: false,
            weight: tile.weight || 1,
            score: tile.score || 5,
            category: tile.category || "General",
        }))

        return {
            tiles,
            regenerationDiminisher: data.regeneration_diminisher || 1
        }
    } catch (error) {
        console.error("Error fetching board:", error)
        // Return empty array on error - component will handle fallback
        return { tiles: [], regenerationDiminisher: 1 }
    }
}

export async function regenerateBoardAPI(): Promise<BoardData> {
    try {
        const response = await fetch("https://api.bingo.local/tiles/me/regenerate", {
            method: 'POST',
            credentials: "include"
        })
        if (!response.ok) {
            const errorData = await response.json()
            throw new Error(errorData.message || "Failed to regenerate board")
        }
        const data = await response.json()
        const tiles = data.tiles.map((tile: any) => ({
            id: tile.id,
            title: tile.title,
            marked: false,
            weight: tile.weight || 1,
            score: tile.score || 5,
            category: tile.category || "General",
        }))

        return {
            tiles,
            regenerationDiminisher: data.regeneration_diminisher || 1
        }
    } catch (error) {
        console.error("Error regenerating board:", error)
        throw error
    }
}


