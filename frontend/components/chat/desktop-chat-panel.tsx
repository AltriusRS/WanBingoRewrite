"use client"

import { Card } from "@/components/ui/card"
import Image from "next/image"
import { Radio, X, Users } from "lucide-react"
import { Button } from "@/components/ui/button"
import { ScrollArea } from "@/components/ui/scroll-area"
import { StandardMessage } from "@/components/chat/messages/standard-message"
import { SystemMessage } from "@/components/chat/messages/system-message"
import { Input } from "@/components/ui/input"
import { useEffect, useRef, useState } from "react"
import { useChat } from "@/components/chat/chat-context"
import { buildSubmitHandler, type ChatPanelProps, updateLiveTime } from "@/lib/chatUtils"
import { getCurrentUser, getSessionId, isChatBanned, isHost } from "@/lib/auth"
import { MemberList } from "./member-list"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"

export function DesktopChatPanel({ onClose }: ChatPanelProps) {
  const scrollRef = useRef<HTMLDivElement>(null)
  const chatContext = useChat()
  const [currentUserId, setCurrentUserId] = useState<string | null>(null)
  const [isCurrentUserHost, setIsCurrentUserHost] = useState(false)
  const [chatBanStatus, setChatBanStatus] = useState<{ banned: boolean; reason?: string; expiry?: Date }>({
    banned: false,
  })
  const [sessionId, setSessionId] = useState<string | null>(null)

  useEffect(() => {
    const loadUserData = async () => {
      const user = await getCurrentUser()
      const hostStatus = await isHost()
      const banStatus = await isChatBanned()
      const session = await getSessionId()

      setCurrentUserId(user?.id || null)
      setIsCurrentUserHost(hostStatus)
      setChatBanStatus(banStatus)
      setSessionId(session)
    }
    loadUserData()
  }, [])

  const handleSubmit = buildSubmitHandler(chatContext, sessionId || undefined)

  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollIntoView({ behavior: "smooth" })
    }
  }, [chatContext.messages])

  useEffect(() => {
    updateLiveTime(chatContext.episode, chatContext)
    const interval = setInterval(updateLiveTime, 60_000, chatContext.episode, chatContext)

    return () => clearInterval(interval)
  }, [chatContext.episode])

  return (
    <Card className="flex h-full max-h-full flex-col overflow-hidden">
      <div className="shrink-0 border-b border-border bg-card p-4">
        <div className="flex items-start justify-between gap-3">
          <div className="min-w-0 flex-1">
            {chatContext.episode.id && chatContext.episode.isLive ? (
              <div className="relative -mt-6 w-full overflow-hidden rounded-md pb-[56.25%]">
                <iframe
                  className="absolute left-0 top-0 h-full w-full rounded-md"
                  src={`https://www.youtube.com/embed/${chatContext.episode.id}?autoplay=1&mute=0&modestbranding=1&rel=0`}
                  title={`WAN Show Stream - ${chatContext.episode.title}`}
                  allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
                  allowFullScreen
                />
              </div>
            ) : (
              <Image
                className="-mt-6 w-full rounded-md"
                src={chatContext.episode.thumbnail ?? "https://cataas.com/cat?width=720&height=480"}
                alt={"Thumbnail for The WAN Show episode titled " + chatContext.episode.title}
                width={720}
                height={480}
              />
            )}
            <h3 className="truncate text-sm font-semibold text-foreground">{chatContext.episode.title}</h3>
            <div className="mt-2 flex items-center gap-3 text-xs text-muted-foreground">
              <div className="flex items-center gap-2">
                {chatContext.episode.isLive ? (
                  <>
                    <div className="relative flex h-2 w-2">
                      <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-primary/75 opacity-75"></span>
                      <span className="relative inline-flex h-2 w-2 rounded-full bg-primary"></span>
                    </div>
                    <span className="font-medium text-primary">LIVE</span>
                  </>
                ) : (
                  <>
                    <Radio className="h-4 w-4" />
                    <span>Scheduled</span>
                  </>
                )}
              </div>
              <div className="h-4 w-px bg-border" />
              <div className="truncate">
                {chatContext.episode.isLive ? `Live for ${chatContext.liveTime}` : `Starts ${chatContext.liveTime}`}
              </div>
            </div>
          </div>
        </div>
      </div>

      <Tabs defaultValue="chat" className="flex flex-1 flex-col overflow-hidden">
        <div className="-mt-2 flex shrink-0 items-center justify-between border-b border-border px-4 pb-2">
          <TabsList className="h-8">
            <TabsTrigger value="chat" className="text-xs">
              Chat
            </TabsTrigger>
            <TabsTrigger value="members" className="text-xs">
              <Users className="mr-1 h-3 w-3" />
              Members
            </TabsTrigger>
          </TabsList>
          <Button variant="ghost" size="icon" onClick={onClose} className="h-8 w-8 md:hidden">
            <X className="h-4 w-4" />
          </Button>
        </div>

        <TabsContent value="chat" className="mt-0 flex flex-1 flex-col overflow-hidden">
          <ScrollArea className="flex-1 p-4">
            <div className="space-y-3">
              {chatContext.messages.map((msg) => (
                <div key={msg.id}>
                  {msg.type === "user" ? (
                    <StandardMessage
                      msg={msg}
                      currentUserId={currentUserId || undefined}
                      isCurrentUserHost={isCurrentUserHost}
                    />
                  ) : (
                    <SystemMessage msg={msg} />
                  )}
                </div>
              ))}
              <div ref={scrollRef} />
            </div>
          </ScrollArea>

          {chatBanStatus.banned ? (
            <div className="border-t border-border bg-destructive/10 p-4">
              <p className="text-sm font-medium text-destructive">You are banned from chat</p>
              {chatBanStatus.reason && <p className="text-xs text-muted-foreground">{chatBanStatus.reason}</p>}
              {chatBanStatus.expiry && (
                <p className="text-xs text-muted-foreground">Ban expires: {chatBanStatus.expiry.toLocaleString()}</p>
              )}
            </div>
          ) : !currentUserId ? (
            <div className="border-t border-border p-4">
              <Button
                variant="outline"
                className="w-full bg-transparent"
                onClick={() => (window.location.href = "/api/auth/signin")}
              >
                Sign in to chat
              </Button>
            </div>
          ) : (
            <form onSubmit={handleSubmit} className="border-t border-border p-4">
              <div className="flex items-center gap-2">
                <Input
                  value={chatContext.text}
                  onChange={(e) => chatContext.setText(e.target.value)}
                  placeholder="Type a message"
                  disabled={chatContext.sending}
                />
                <Button type="submit" disabled={chatContext.sending || chatContext.text.trim().length === 0}>
                  Send
                </Button>
              </div>
            </form>
          )}
        </TabsContent>

        <TabsContent value="members" className="mt-0 flex-1 overflow-hidden">
          <MemberList />
        </TabsContent>
      </Tabs>
    </Card>
  )
}
