import {getApiRoot} from "@/lib/auth";

export interface BingoTile {
    id: string
    title: string
    marked: boolean
    weight?: number
    score?: number,
    category?: string
    settings?: unknown
}

export interface BingoBoardProps {
    onWin?: () => void
}

export interface BoardData {
    tiles: BingoTile[]
    regenerationDiminisher: number
}

// Fetch confirmed tiles for the current show
export async function fetchConfirmedTiles(): Promise<Set<string>> {
    try {
        const response = await fetch(`${getApiRoot()}/tiles/confirmed`, {
            credentials: "include",
        })
        if (!response.ok) {
            console.error("Failed to fetch confirmed tiles:", response.status)
            return new Set()
        }
        const data = await response.json()
        return new Set(data)
    } catch (error) {
        console.error("Failed to fetch confirmed tiles:", error)
        return new Set()
    }
}

// Check if a bingo board has a winning line (5 in a row in any direction)
export function checkBingoWin(tiles: BingoTile[]): boolean {
    if (tiles.length !== 25) {
        return false // Not a valid 5x5 board
    }

    // Convert tiles array to 5x5 grid
    const grid: boolean[][] = []
    for (let i = 0; i < 5; i++) {
        grid[i] = []
        for (let j = 0; j < 5; j++) {
            grid[i][j] = tiles[i * 5 + j].marked
        }
    }

    // Check horizontal lines (rows)
    for (let row = 0; row < 5; row++) {
        if (grid[row].every(cell => cell)) {
            return true
        }
    }

    // Check vertical lines (columns)
    for (let col = 0; col < 5; col++) {
        let columnWin = true
        for (let row = 0; row < 5; row++) {
            if (!grid[row][col]) {
                columnWin = false
                break
            }
        }
        if (columnWin) {
            return true
        }
    }

    // Check main diagonal (top-left to bottom-right)
    let mainDiagonalWin = true
    for (let i = 0; i < 5; i++) {
        if (!grid[i][i]) {
            mainDiagonalWin = false
            break
        }
    }
    if (mainDiagonalWin) {
        return true
    }

    // Check anti-diagonal (top-right to bottom-left)
    let antiDiagonalWin = true
    for (let i = 0; i < 5; i++) {
        if (!grid[i][4 - i]) {
            antiDiagonalWin = false
            break
        }
    }
    if (antiDiagonalWin) {
        return true
    }

    return false
}

// Check if a bingo win is valid (all tiles in winning lines are confirmed)
export function isValidWin(tiles: BingoTile[], confirmedTiles: Set<string>): boolean {
    // Check horizontal lines (rows)
    for (let row = 0; row < 5; row++) {
        let allMarked = true
        let allConfirmed = true
        for (let col = 0; col < 5; col++) {
            const tile = tiles[row * 5 + col]
            if (!tile.marked) {
                allMarked = false
                break
            }
            if (!confirmedTiles.has(tile.id)) {
                allConfirmed = false
            }
        }
        if (allMarked && allConfirmed) return true
    }

    // Check vertical lines (columns)
    for (let col = 0; col < 5; col++) {
        let allMarked = true
        let allConfirmed = true
        for (let row = 0; row < 5; row++) {
            const tile = tiles[row * 5 + col]
            if (!tile.marked) {
                allMarked = false
                break
            }
            if (!confirmedTiles.has(tile.id)) {
                allConfirmed = false
            }
        }
        if (allMarked && allConfirmed) return true
    }

    // Check main diagonal (top-left to bottom-right)
    let allMarked = true
    let allConfirmed = true
    for (let i = 0; i < 5; i++) {
        const tile = tiles[i * 5 + i]
        if (!tile.marked) {
            allMarked = false
            break
        }
        if (!confirmedTiles.has(tile.id)) {
            allConfirmed = false
        }
    }
    if (allMarked && allConfirmed) return true

    // Check anti-diagonal (top-right to bottom-left)
    allMarked = true
    allConfirmed = true
    for (let i = 0; i < 5; i++) {
        const tile = tiles[i * 5 + (4 - i)]
        if (!tile.marked) {
            allMarked = false
            break
        }
        if (!confirmedTiles.has(tile.id)) {
            allConfirmed = false
        }
    }
    if (allMarked && allConfirmed) return true

    return false
}

export async function fetchBoardFromAPI(): Promise<BoardData> {
    try {
        let response = await fetch(`${getApiRoot()}/tiles/me`, {
            credentials: "include"
        })
        if (!response.ok) {
            if (response.status === 401) {
                // Try anonymous
                response = await fetch(`${getApiRoot()}/tiles/anonymous`, {
                    credentials: "include"
                })
            } else {
                throw new Error("Failed to fetch board")
            }
        }
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
        return {tiles: [], regenerationDiminisher: 1}
    }
}

export async function regenerateBoardAPI(): Promise<BoardData> {
    try {
        let response = await fetch(`${getApiRoot()}/tiles/me/regenerate`, {
            method: 'POST',
            credentials: "include"
        })
        if (!response.ok) {
            if (response.status === 401) {
                // Try anonymous regenerate
                response = await fetch(`${getApiRoot()}/tiles/anonymous/regenerate`, {
                    method: 'POST',
                    credentials: "include"
                })
            } else {
                const errorData = await response.json()
                throw new Error(errorData.message || "Failed to regenerate board")
            }
        }
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


