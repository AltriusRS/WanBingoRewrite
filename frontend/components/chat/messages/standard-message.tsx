"use client"
import type { ChatMessage } from "@/lib/chatUtils"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Button } from "@/components/ui/button"
import { HoverCard, HoverCardContent, HoverCardTrigger } from "@/components/ui/hover-card"
import { Badge } from "@/components/ui/badge"
import { Trash2, Shield } from "lucide-react"

interface StandardMessageProps {
  msg: ChatMessage
  currentUserId?: string
  isCurrentUserHost?: boolean
}

export function StandardMessage({ msg, currentUserId, isCurrentUserHost }: StandardMessageProps) {
  const isOwnMessage = msg.userId === currentUserId
  const canDelete = isOwnMessage || isCurrentUserHost

  const handleDelete = async () => {
    if (!confirm("Delete this message?")) return

    try {
      await fetch(`http://localhost:8080/api/chat/message/${msg.id}`, {
        method: "DELETE",
      })
    } catch (error) {
      console.error("Failed to delete message:", error)
    }
  }

  const getInitials = (username?: string) => {
    if (!username) return "?"
    return username
      .split(" ")
      .map((n) => n[0])
      .join("")
      .toUpperCase()
      .slice(0, 2)
  }

  // Check if user is a host based on username or metadata
  const isMessageFromHost =
    msg.username?.toLowerCase().includes("linus") ||
    msg.username?.toLowerCase().includes("luke") ||
    msg.username?.toLowerCase().includes("dan")

  return (
    <div className="group flex gap-3">
      <HoverCard>
        <HoverCardTrigger asChild>
          <Avatar className="h-8 w-8 cursor-pointer">
            <AvatarImage src={msg.avatar || "/placeholder.svg"} alt={msg.username} />
            <AvatarFallback className="bg-primary/10 text-xs">{getInitials(msg.username)}</AvatarFallback>
          </Avatar>
        </HoverCardTrigger>
        <HoverCardContent className="w-64">
          <div className="flex gap-3">
            <Avatar className="h-12 w-12">
              <AvatarImage src={msg.avatar || "/placeholder.svg"} alt={msg.username} />
              <AvatarFallback className="bg-primary/10">{getInitials(msg.username)}</AvatarFallback>
            </Avatar>
            <div className="flex-1 space-y-1">
              <div className="flex items-center gap-2">
                <h4 className="text-sm font-semibold">{msg.username}</h4>
                {isMessageFromHost && (
                  <Badge variant="default" className="gap-1 text-xs">
                    <Shield className="h-3 w-3" />
                    Host
                  </Badge>
                )}
              </div>
              <p className="text-xs text-muted-foreground">Member since {new Date().toLocaleDateString()}</p>
            </div>
          </div>
        </HoverCardContent>
      </HoverCard>

      <div className="flex-1 space-y-1">
        <div className="flex items-baseline gap-2">
          <span className="text-sm font-medium text-primary">{msg.username}</span>
          {isMessageFromHost && <Shield className="h-3 w-3 text-primary" />}
          <span className="text-xs text-muted-foreground">{msg.timestamp}</span>
        </div>
        <div className="flex items-start justify-between gap-2">
          <p className="text-sm text-foreground">{msg.message}</p>
          {canDelete && (
            <Button
              variant="ghost"
              size="icon"
              className="h-6 w-6 opacity-0 transition-opacity group-hover:opacity-100"
              onClick={handleDelete}
            >
              <Trash2 className="h-3 w-3" />
            </Button>
          )}
        </div>
      </div>
    </div>
  )
}
