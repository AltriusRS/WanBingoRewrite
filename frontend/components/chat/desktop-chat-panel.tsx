"use client"

import {Card} from "@/components/ui/card"
import Image from "next/image"
import {MessagesSquare, Radio, Users, X} from "lucide-react"
import {Button} from "@/components/ui/button"
import {ScrollArea} from "@/components/ui/scroll-area"
import {StandardMessage} from "@/components/chat/messages/standard-message"
import {SystemMessage} from "@/components/chat/messages/system-message"
import {Input} from "@/components/ui/input"
import {useEffect, useRef} from "react"
import {useChat} from "@/components/chat/chat-context"
import {buildSubmitHandler, type ChatPanelProps, updateLiveTime} from "@/lib/chatUtils"
import {useAuth} from "@/components/auth"
import {MemberList} from "./member-list"
import {Tabs, TabsContent, TabsList, TabsTrigger} from "@/components/ui/tabs"

export function DesktopChatPanel({onClose}: ChatPanelProps) {
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
    }, [chatContext.episode])

    return (
        <Card className="flex h-full max-h-full flex-col overflow-hidden">
            <div className="shrink-0 border-b border-border bg-card p-4">
                <div className="flex items-start justify-between gap-3">
                    <div className="min-w-0 flex-1">
                        {chatContext.episode.youtube_id && chatContext.episode.actual_start_time ? (
                            <div className="relative -mt-6 w-full overflow-hidden rounded-md pb-[56.25%]">
                                <iframe
                                    className="absolute left-0 top-0 h-full w-full rounded-md"
                                    src={`https://www.youtube.com/embed/${chatContext.episode.youtube_id}?autoplay=1&mute=0&modestbranding=1&rel=0`}
                                    title={`WAN Show Stream - ${chatContext.episode.metadata?.title}`}
                                    allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
                                    allowFullScreen
                                />
                            </div>
                        ) : (
                            <Image
                                className="-mt-6 w-full rounded-md"
                                src={chatContext.episode.thumbnail ?? "https://cataas.com/cat?width=720&height=480"}
                                alt={"Thumbnail for The WAN Show episode titled " + chatContext.episode.metadata?.title}
                                width={720}
                                height={480}
                            />
                        )}
                        <h3 className="truncate text-sm font-semibold text-foreground">{chatContext.episode.metadata?.title}</h3>
                        <div className="mt-2 flex items-center gap-3 text-xs text-muted-foreground">
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

            <Tabs defaultValue="chat" className="flex flex-1 flex-col overflow-hidden -mt-4">
                <div className=" flex shrink-0 items-center justify-between border-b border-border px-4 pb-2">
                    <TabsList className="h-8">
                        <TabsTrigger value="chat" className="text-xs">
                            <MessagesSquare className="mr-1 h-3 w-3"/>
                            Chat
                        </TabsTrigger>
                        <TabsTrigger value="members" className="text-xs">
                            <Users className="mr-1 h-3 w-3"/>
                            Members
                        </TabsTrigger>
                    </TabsList>
                    <Button variant="ghost" size="icon" onClick={onClose} className="h-8 w-8 md:hidden">
                        <X className="h-4 w-4"/>
                    </Button>
                </div>

                <TabsContent value="chat" className="mt-0 flex flex-1 flex-col overflow-hidden">
                    <ScrollArea className="flex-1 p-4">
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
                        <div className="border-t border-border p-4">
                            <Button
                                variant="outline"
                                className="w-full bg-transparent"
                                onClick={login}
                            >
                                Sign in with Discord to chat
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
                                <Button type="submit"
                                        disabled={chatContext.sending || chatContext.text.trim().length === 0}>
                                    Send
                                </Button>
                            </div>
                        </form>
                    )}
                </TabsContent>

                <TabsContent value="members" className="mt-0 flex-1 overflow-hidden">
                    <MemberList/>
                </TabsContent>
            </Tabs>
        </Card>
    )
}
