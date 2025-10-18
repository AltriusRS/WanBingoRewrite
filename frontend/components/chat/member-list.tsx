"use client"

import {useChat} from "./chat-context"
import {ScrollArea} from "@/components/ui/scroll-area"
import {Avatar, AvatarFallback, AvatarImage} from "@/components/ui/avatar"
import {Badge} from "@/components/ui/badge"
import {HoverCard, HoverCardContent, HoverCardTrigger} from "@/components/ui/hover-card"
import {Shield, Users} from "lucide-react"
import {Card} from "@/components/ui/card"
import type {Player} from "@/lib/chatUtils"

export function MemberList() {
    const {memberList, memberCount, memberListLoading} = useChat()

    const getInitials = (name?: string) => {
        if (!name) return "?"
        return name
            .split(" ")
            .map((n) => n[0])
            .join("")
            .toUpperCase()
            .slice(0, 2)
    }

    const canModerate = (player: Player) => {
        return (player.permissions || 0) & 1024 // PermCanModerate = 1024
    }

    const getAccountCreated = (createdAt?: string) => {
        if (!createdAt) return "Unknown"
        return new Date(createdAt).toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'long',
            day: 'numeric'
        })
    }

    if (memberListLoading) {
        return (
            <Card className="mx-2 p-4">
                <p className="text-sm text-muted-foreground">Loading members...</p>
            </Card>
        )
    }

    return (
        <Card className="flex h-full flex-col mx-2">
            <div className="border-b border-border p-4">
                <div className="flex items-center gap-2">
                    <Users className="h-4 w-4 text-muted-foreground"/>
                    <h3 className="font-semibold text-foreground">Members ({memberCount})</h3>
                </div>
            </div>

            <ScrollArea className="flex-1 p-4">
                <div className="space-y-2">
                    {memberList.map((player) => (
                        <HoverCard key={player.id}>
                            <HoverCardTrigger asChild>
                                <div className="flex items-center gap-3 rounded-lg p-2 hover:bg-muted/50 cursor-pointer">
                                    <Avatar className="h-8 w-8">
                                        <AvatarImage src={player.avatar || "/placeholder.svg"} alt={player.display_name}/>
                                        <AvatarFallback className="bg-primary/10 text-xs">{getInitials(player.display_name)}</AvatarFallback>
                                    </Avatar>
                                    <div className="flex-1 truncate">
                                        <div className="flex items-center gap-2">
                                            <span
                                                className="truncate text-sm font-medium text-foreground"
                                                style={{color: player.settings?.chat_color || "#3b82f6"}}
                                            >
                                                {player.display_name}
                                            </span>
                                            {canModerate(player) && <Shield className="h-3 w-3 text-primary"/>}
                                        </div>
                                    </div>
                                </div>
                            </HoverCardTrigger>
                            <HoverCardContent className="w-64">
                                <div className="flex gap-3">
                                    <Avatar className="h-12 w-12">
                                        <AvatarImage src={player.avatar || "/placeholder.svg"} alt={player.display_name}/>
                                        <AvatarFallback className="bg-primary/10">{getInitials(player.display_name)}</AvatarFallback>
                                    </Avatar>
                                    <div className="flex-1 space-y-1">
                                        <div className="flex items-center gap-2">
                                            <h4 className="text-sm font-semibold">{player.display_name}</h4>
                                            {canModerate(player) && (
                                                <Badge variant="default" className="gap-1 text-xs">
                                                    <Shield className="h-3 w-3"/>
                                                    Mod
                                                </Badge>
                                            )}
                                        </div>
                                        <p className="text-xs text-muted-foreground">
                                            Chat Color: <span style={{color: player.settings?.chat_color || "#3b82f6"}}>{player.settings?.chat_color || "#3b82f6"}</span>
                                        </p>
                                        {player.settings?.bio && (
                                            <p className="text-xs text-muted-foreground">{player.settings.bio}</p>
                                        )}
                                        <p className="text-xs text-muted-foreground">Account Created {getAccountCreated(player.created_at)}</p>
                                    </div>
                                </div>
                            </HoverCardContent>
                        </HoverCard>
                    ))}

                    {memberList.length === 0 && (
                        <p className="py-8 text-center text-sm text-muted-foreground">No members online</p>
                    )}
                </div>
            </ScrollArea>
        </Card>
    )
}
