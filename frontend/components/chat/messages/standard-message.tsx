"use client"
import type {ChatMessage} from "@/lib/chatUtils"
import {Button} from "@/components/ui/button"
import {Badge} from "@/components/ui/badge"
import {Avatar, AvatarFallback, AvatarImage} from "@/components/ui/avatar"
import {HoverCard, HoverCardContent, HoverCardTrigger} from "@/components/ui/hover-card"
import {Trash2, Shield} from "lucide-react"
import {MemoizedMarkdown} from "@/components/ui/markdown";
import {Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle} from "@/components/ui/dialog";
import React, {useState} from "react";

interface StandardMessageProps {
    msg: ChatMessage
    currentUserId?: string
    isCurrentUserHost?: boolean
}

export function StandardMessage({msg, currentUserId, isCurrentUserHost}: StandardMessageProps) {
    const [open, setOpen] = useState(false);
    const [clickedUrl, setClickedUrl] = useState("");

    const isOwnMessage = msg.player_id === currentUserId
    const canModerate = (msg.player?.permissions || 0) & 1024 // PermCanModerate = 1024
    const canDelete = isOwnMessage || isCurrentUserHost || canModerate

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

    const getInitials = (name?: string) => {
        if (!name) return "?"
        return name
            .split(" ")
            .map((n) => n[0])
            .join("")
            .toUpperCase()
            .slice(0, 2)
    }

    const chatColor = msg.player?.settings?.chat_color || "#3b82f6"
    const accountCreated = msg.player?.created_at
        ? new Date(msg.player.created_at).toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'long',
            day: 'numeric'
        })
        : "Unknown"


    return (
        <div className="group flex gap-3">
            <HoverCard>
                <HoverCardTrigger asChild>
                    <Avatar className="h-8 w-8 cursor-pointer">
                        <AvatarImage src={msg.player?.avatar || "/placeholder.svg"} alt={msg.player?.display_name}/>
                        <AvatarFallback className="bg-primary/10 text-xs">{getInitials(msg.player?.display_name)}</AvatarFallback>
                    </Avatar>
                </HoverCardTrigger>
                <HoverCardContent className="w-64">
                    <div className="flex gap-3">
                        <Avatar className="h-12 w-12">
                            <AvatarImage src={msg.player?.avatar || "/placeholder.svg"} alt={msg.player?.display_name}/>
                            <AvatarFallback className="bg-primary/10">{getInitials(msg.player?.display_name)}</AvatarFallback>
                        </Avatar>
                        <div className="flex-1 space-y-1">
                            <div className="flex items-center gap-2">
                                <h4 className="text-sm font-semibold">{msg.player?.display_name || msg.player_id}</h4>
                                {canModerate && (
                                    <Badge variant="default" className="gap-1 text-xs">
                                        <Shield className="h-3 w-3"/>
                                        Mod
                                    </Badge>
                                )}
                            </div>
                            <p className="text-xs text-muted-foreground">
                                Chat Color: <span style={{color: chatColor}}>{chatColor}</span>
                            </p>
                            {msg.player?.settings?.bio && (
                                <p className="text-xs text-muted-foreground">{msg.player.settings.bio}</p>
                            )}
                            <p className="text-xs text-muted-foreground">Account Created {accountCreated}</p>
                        </div>
                    </div>
                </HoverCardContent>
            </HoverCard>

            <div className="flex-1 space-y-1">
                <div className="flex items-baseline gap-2">
                    <HoverCard>
                        <HoverCardTrigger asChild>
                            <span
                                className="text-sm font-medium cursor-pointer"
                                style={{color: chatColor}}
                            >
                                {msg.player?.display_name || msg.player_id}
                            </span>
                        </HoverCardTrigger>
                        <HoverCardContent className="w-64">
                            <div className="flex gap-3">
                                <Avatar className="h-12 w-12">
                                    <AvatarImage src={msg.player?.avatar || "/placeholder.svg"} alt={msg.player?.display_name}/>
                                    <AvatarFallback className="bg-primary/10">{getInitials(msg.player?.display_name)}</AvatarFallback>
                                </Avatar>
                                <div className="flex-1 space-y-1">
                                    <div className="flex items-center gap-2">
                                        <h4 className="text-sm font-semibold">{msg.player?.display_name || msg.player_id}</h4>
                                        {canModerate && (
                                            <Badge variant="default" className="gap-1 text-xs">
                                                <Shield className="h-3 w-3"/>
                                                Mod
                                            </Badge>
                                        )}
                                    </div>
                                    <p className="text-xs text-muted-foreground">
                                        Chat Color: <span style={{color: chatColor}}>{chatColor}</span>
                                    </p>
                                    {msg.player?.settings?.bio && (
                                        <p className="text-xs text-muted-foreground">{msg.player.settings.bio}</p>
                                    )}
                                    <p className="text-xs text-muted-foreground">Account Created {accountCreated}</p>
                                </div>
                            </div>
                        </HoverCardContent>
                    </HoverCard>
                    {canModerate && <Shield className="h-3 w-3 text-primary"/>}
                    <span
                        className="text-xs text-muted-foreground">{new Date(msg.created_at).toLocaleTimeString()}</span>
                </div>
                <div className="flex items-start justify-between gap-2">
                    <MemoizedMarkdown
                        key={`${msg.id}-text`}
                        id={msg.id}
                        content={msg.contents}
                        onLinkClick={(href) => {
                            setClickedUrl(href);
                            setOpen(true);
                        }}
                    />
                    {canDelete && (
                        <Button
                            variant="ghost"
                            size="icon"
                            className="h-6 w-6 opacity-0 transition-opacity group-hover:opacity-100"
                            onClick={handleDelete}
                        >
                            <Trash2 className="h-3 w-3"/>
                        </Button>
                    )}
                    <Dialog open={open} onOpenChange={setOpen}>
                        <DialogContent>
                            <DialogHeader>
                                <DialogTitle>External Link</DialogTitle>
                            </DialogHeader>
                            <div className="py-2">You're about to open: {clickedUrl}</div>
                            <DialogFooter className="flex gap-2 justify-end">
                                <Button
                                    onClick={() => {
                                        window.open(clickedUrl, "_blank");
                                        setOpen(false);
                                    }}
                                >
                                    Proceed
                                </Button>
                                <Button variant="outline" onClick={() => setOpen(false)}>
                                    Cancel
                                </Button>
                            </DialogFooter>
                        </DialogContent>
                    </Dialog>

                </div>
            </div>
        </div>
    )
}
