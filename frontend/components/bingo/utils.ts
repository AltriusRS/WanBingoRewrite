import {BingoTile} from "@/lib/bingoUtils";
import confetti from "canvas-confetti";

export function checkWin(currentTiles: BingoTile[]): boolean {
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

    return diagonal1.every((tile) => tile.marked) || diagonal2.every((tile) => tile.marked);
}

export function triggerConfetti(): void {
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