"use client"

import {Card} from "@/components/ui/card"
import {ChevronDown, Radio, X, Settings} from "lucide-react"
import {Button} from "@/components/ui/button"
import Link from "next/link"
import {ScrollArea} from "@/components/ui/scroll-area"
import {Input} from "@/components/ui/input"
import {useEffect, useRef} from "react"
import {useChat} from "@/components/chat/chat-context"
import {buildSubmitHandler, type ChatPanelProps, updateLiveTime} from "@/lib/chatUtils"
import {StandardMessage} from "@/components/chat/messages/standard-message"
import {SystemMessage} from "@/components/chat/messages/system-message"
import {useAuth} from "@/components/auth"

export function MobileChatPanel({onClose}: ChatPanelProps) {
    const scrollRef = useRef<HTMLDivElement>(null)
    const chatContext = useChat()
    const {user, login} = useAuth()

    const handleSubmit = buildSubmitHandler(chatContext)

    useEffect(() => {
        if (scrollRef.current) {
            scrollRef.current.scrollIntoView({behavior: "smooth"})
        }
    }, [chatContext.messages])

    useEffect(() => {
        updateLiveTime(chatContext.episode, chatContext)
        const interval = setInterval(updateLiveTime, 60_000, chatContext.episode, chatContext)

        return () => clearInterval(interval)
    }, [chatContext.episode, chatContext])

    return (
        <Card className="flex h-[60vh] max-h-[60vh] flex-col rounded-b-none border-b-0">
            {/* Header */}
            <div className="flex flex-col shrink-0 border-b border-border p-3">
                <h3 className="text-sm font-semibold text-foreground">WAN Show Bingo</h3>
                <div className="flex items-center justify-between mt-2">
                    <ChevronDown className="h-4 w-4 text-muted-foreground"/>
                    <div className="flex items-center gap-2">
                        <Link href="/account">
                            <Button variant="ghost" size="icon" className="h-8 w-8">
                                <Settings className="h-4 w-4"/>
                            </Button>
                        </Link>
                        <Button variant="ghost" size="icon" onClick={onClose} className="h-8 w-8">
                            <X className="h-4 w-4"/>
                        </Button>
                    </div>
                </div>
            </div>
            <div className="shrink-0 border-b border-border bg-card p-3">
                <div className="flex items-start justify-between gap-3">
                    <div className="min-w-0 flex-1">
                        <div className="flex flex-col">
                            <h3 className="truncate text-md font-semibold text-foreground">{chatContext.episode.metadata?.title}</h3>
                        </div>
                        <div className="mt-1.5 flex items-center gap-2 text-xs text-muted-foreground">
                            <div className="flex items-center gap-2">
                                {chatContext.episode.actual_start_time ? (
                                    <>
                                        <div className="relative flex h-2 w-2">
                                            <span
                                                className="absolute inline-flex h-full w-full animate-ping rounded-full bg-primary/75 opacity-75"></span>
                                            <span
                                                className="relative inline-flex h-2 w-2 rounded-full bg-primary"></span>
                                        </div>
                                        <span className="font-medium text-primary">LIVE</span>
                                    </>
                                ) : (
                                    <>
                                        <Radio className="h-4 w-4"/>
                                        <span>Scheduled</span>
                                    </>
                                )}
                            </div>
                            <div className="h-4 w-px bg-border"/>
                            <div className="truncate">
                                {chatContext.episode.actual_start_time ? `Live for ${chatContext.liveTime}` : `Starts ${chatContext.liveTime}`}
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <ScrollArea className="flex-1 p-3">
                <div className="space-y-3">
                    {chatContext.messages.map((msg) => (
                        <div key={msg.id}>
                            {msg.system ? (
                                <SystemMessage msg={msg}/>
                            ) : (
                                <StandardMessage
                                    msg={msg}
                                    currentUserId={user?.id}
                                    isCurrentUserHost={false}
                                />
                            )}
                        </div>
                    ))}
                    <div ref={scrollRef}/>
                </div>
            </ScrollArea>

            {!user ? (
                <div
                    className="border-t border-border bg-background/60 p-3 backdrop-blur supports-[backdrop-filter]:bg-background/40">
                    <Button
                        variant="outline"
                        className="w-full bg-transparent"
                        onClick={login}
                    >
                        Sign in with Discord to chat
                    </Button>
                </div>
            ) : (
                <form
                    onSubmit={handleSubmit}
                    className="border-t border-border bg-background/60 p-3 backdrop-blur supports-[backdrop-filter]:bg-background/40"
                >
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
        </Card>
    )
}
