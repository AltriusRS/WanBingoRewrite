"use client"

import {useChat} from "./chat-context"
import {ScrollArea} from "@/components/ui/scroll-area"
import {Avatar, AvatarFallback, AvatarImage} from "@/components/ui/avatar"
import {Badge} from "@/components/ui/badge"
import {Shield, Users} from "lucide-react"
import {Card} from "@/components/ui/card"

export function MemberList() {
    const {memberList, memberCount, memberListLoading} = useChat()

    const getInitials = (username: string) => {
        return username
            .split(" ")
            .map((n) => n[0])
            .join("")
            .toUpperCase()
            .slice(0, 2)
    }

    const isHost = (username: string) => {
        return (
            username.toLowerCase().includes("linus") ||
            username.toLowerCase().includes("luke") ||
            username.toLowerCase().includes("dan")
        )
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
                    <h3 className="font-semibold text-foreground">Members</h3>
                    <Badge variant="secondary">{memberCount}</Badge>
                </div>
            </div>

            <ScrollArea className="flex-1 p-4">
                <div className="space-y-2">
                    {memberList.map((member, index) => (
                        <div key={index} className="flex items-center gap-3 rounded-lg p-2 hover:bg-muted/50">
                            <Avatar className="h-8 w-8">
                                <AvatarImage src={`https://api.dicebear.com/7.x/avataaars/svg?seed=${member}`}
                                             alt={member}/>
                                <AvatarFallback className="bg-primary/10 text-xs">{getInitials(member)}</AvatarFallback>
                            </Avatar>
                            <div className="flex-1 truncate">
                                <div className="flex items-center gap-2">
                                    <span className="truncate text-sm font-medium text-foreground">{member}</span>
                                    {isHost(member) && <Shield className="h-3 w-3 text-primary"/>}
                                </div>
                            </div>
                        </div>
                    ))}

                    {memberList.length === 0 && (
                        <p className="py-8 text-center text-sm text-muted-foreground">No members online</p>
                    )}
                </div>
            </ScrollArea>
        </Card>
    )
}
