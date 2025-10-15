"use client"
import { Card } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { useChat } from "@/components/chat/chat-context"

export function EpisodeInfoPanel() {
  const { episode, liveTime } = useChat()

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
      </div>
    </Card>
  )
}
