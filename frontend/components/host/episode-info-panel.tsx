"use client"
import { useState, useEffect } from "react"
import { Card } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Clock } from "lucide-react"
import { useChat } from "@/components/chat/chat-context"
import { getApiRoot } from "@/lib/auth"
import { toast } from "sonner"
import type { BingoTile } from "@/lib/bingoUtils"

export function EpisodeInfoPanel() {
  const { episode, liveTime } = useChat()
  const [lateTile, setLateTile] = useState<BingoTile | null>(null)


  useEffect(() => {
    fetchLateTile()
  }, [])

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

      toast.success("Show Is Late has been confirmed successfully.")
    } catch (error) {
      console.error("Failed to confirm late tile:", error)
      toast.error(error instanceof Error ? error.message : "An error occurred")
    }
  }

  return (
    <Card className="p-4">
      <h3 className="mb-4 font-semibold text-foreground">Current Episode</h3>

      <div className="space-y-4">
        {episode.thumbnail && (
          <div className="aspect-video overflow-hidden rounded-lg">
            <img
              src={episode.thumbnail || "/placeholder.svg"}
              alt={episode.title}
              className="h-full w-full object-cover"
            />
          </div>
        )}

        <div>
          <h4 className="font-medium text-foreground">{episode.title}</h4>
          <p className="text-sm text-muted-foreground">{episode.date}</p>
        </div>

        <div className="flex items-center gap-2">
          {episode.isLive ? (
            <>
              <Badge variant="default" className="bg-primary">
                LIVE
              </Badge>
              <span className="text-sm text-muted-foreground">{liveTime}</span>
            </>
          ) : (
            <>
              <Badge variant="secondary">Scheduled</Badge>
              <span className="text-sm text-muted-foreground">{liveTime}</span>
            </>
          )}
        </div>

        {lateTile && (
          <Button
            variant="outline"
            size="sm"
            className="w-full gap-2"
            onClick={handleConfirmLate}
          >
            <Clock className="h-4 w-4" />
            Confirm Show Is Late
          </Button>
        )}
      </div>
    </Card>
  )
}
