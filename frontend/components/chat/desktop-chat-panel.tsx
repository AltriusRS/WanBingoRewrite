import {Card} from "@/components/ui/card";
import Image from "next/image";
import {Radio, X} from "lucide-react";
import {Button} from "@/components/ui/button";
import {ScrollArea} from "@/components/ui/scroll-area";
import {StandardMessage} from "@/components/chat/messages/standard-message";
import {SystemMessage} from "@/components/chat/messages/system-message";
import {Input} from "@/components/ui/input";
import {useEffect, useRef} from "react";
import {useChat} from "@/components/chat/chat-context";
import {buildSubmitHandler, ChatPanelProps, updateLiveTime} from "@/lib/chatUtils";

export function DesktopChatPanel({onClose}: ChatPanelProps) {
    const isMobile = false;
    const scrollRef = useRef<HTMLDivElement>(null)
    const chatContext = useChat();

    const handleSubmit = buildSubmitHandler(chatContext);

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
                        {/* --- YouTube Mini Player --- */}
                        {chatContext.episode.id && chatContext.episode.isLive ? (
                            <div className="relative w-full pb-[56.25%] overflow-hidden rounded-md -mt-6">
                                <iframe
                                    className="absolute top-0 left-0 h-full w-full rounded-md"
                                    src={`https://www.youtube.com/embed/${chatContext.episode.id}?autoplay=1&mute=0&modestbranding=1&rel=0`}
                                    title={`WAN Show Stream - ${chatContext.episode.title}`}
                                    allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
                                    allowFullScreen
                                />
                            </div>
                        ) : (
                            <Image
                                className="w-full rounded-md -mt-6"
                                src={chatContext.episode.thumbnail ?? "https://cataas.com/cat?width=720&height=480"}
                                alt={"Thumbnail for The WAN Show episode titled " + chatContext.episode.title}
                                width={720}
                                height={480}
                            />
                        )}
                        {/* --- End Mini Player --- */}
                        <h3 className="truncate text-sm font-semibold text-foreground">{chatContext.episode.title}</h3>
                        <div className="mt-2 flex items-center gap-3 text-xs text-muted-foreground">
                            <div className="flex items-center gap-2">
                                {chatContext.episode.isLive ? (
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
                            <div
                                className="truncate">{chatContext.episode.isLive ? `Live for ${chatContext.liveTime}` : `Starts ${chatContext.liveTime}`}</div>
                        </div>
                    </div>
                </div>
            </div>

            {/* Header */}
            <div className="flex shrink-0 items-center justify-between border-b border-border -mt-2 px-4 pb-4">
                <div>
                    <h3 className="font-semibold text-foreground">WAN Show Chat</h3>
                    <p className="text-xs text-muted-foreground">{chatContext.memberCount} players</p>
                </div>
                <Button variant="ghost" size="icon" onClick={onClose} className="h-8 w-8 md:hidden">
                    <X className="h-4 w-4"/>
                </Button>
            </div>

            <ScrollArea className="flex-1 p-4 max-h-[calc(100vh-36rem)]">
                <div className="space-y-3">
                    {chatContext.messages.map((msg) => (
                        <div key={msg.id}>
                            {msg.type === "user" ? (
                                <StandardMessage msg={msg}/>
                            ) : (
                                <SystemMessage msg={msg}/>
                            )
                            }
                        </div>
                    ))}
                    <div ref={scrollRef}/>
                </div>
            </ScrollArea>

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
        </Card>
    )
}