"use client"

import {useEffect, useState} from "react"
import {Card} from "@/components/ui/card"
import {Button} from "@/components/ui/button"
import {Input} from "@/components/ui/input"
import {Label} from "@/components/ui/label"
import {Switch} from "@/components/ui/switch"
import {ScrollArea} from "@/components/ui/scroll-area"
import {useAuth} from "@/components/auth"
import {getApiRoot} from "@/lib/auth"

export function ChatSettings() {
    const {user, refetch} = useAuth()
    const [chatColor, setChatColor] = useState("#FF6900")
    const [soundOnMention, setSoundOnMention] = useState(true)
    const [saving, setSaving] = useState(false)

    useEffect(() => {
        if (user?.settings) {
            const settings = user.settings as any
            // Load chat settings from nested structure
            if (settings.chat) {
                setChatColor(settings.chat.color || "#FF6900")
                setSoundOnMention(settings.chat.soundOnMention !== false)
            } else {
                // Fallback for old flat structure
                setChatColor(settings.chatColor || "#FF6900")
                setSoundOnMention(settings.soundOnMention !== false)
            }
        }
    }, [user])

    const handleSave = async () => {
        if (!user) return

        setSaving(true)
        const startTime = Date.now()

        try {
            // Get current settings and merge chat settings
            const currentSettings = user.settings as any || {}
            const updatedSettings = {
                ...currentSettings,
                chat: {
                    ...currentSettings.chat,
                    color: chatColor,
                    soundOnMention,
                },
            }

            const response = await fetch(`${getApiRoot()}/users/me`, {
                method: "PUT",
                credentials: "include",
                headers: {"Content-Type": "application/json"},
                body: JSON.stringify({
                    display_name: user.display_name,
                    avatar: user.avatar,
                    settings: updatedSettings,
                }),
            })

            if (!response.ok) {
                throw new Error("Failed to update chat settings")
            }

            // Refetch user data
            await refetch()
        } catch (error) {
            console.error("Failed to save chat settings:", error)
        } finally {
            // Ensure minimum loading time of 200ms
            const elapsed = Date.now() - startTime
            const remaining = Math.max(0, 200 - elapsed)
            setTimeout(() => {
                setSaving(false)
            }, remaining)
        }
    }

    if (!user) {
        return (
            <div className="flex h-full items-center justify-center p-4">
                <p className="text-sm text-muted-foreground">Sign in to access chat settings</p>
            </div>
        )
    }

    return (
        <ScrollArea className="h-full">
            <div className="space-y-6 p-4">
                <Card className="p-4">
                    <h3 className="mb-4 text-sm font-semibold text-foreground">Chat Appearance</h3>

                    <div className="space-y-4">
                        <div className="space-y-2">
                            <Label htmlFor="chat-color">Chat Name Color</Label>
                            <div className="flex gap-2">
                                <Input
                                    id="chat-color"
                                    type="color"
                                    value={chatColor}
                                    onChange={(e) => setChatColor(e.target.value)}
                                    className="h-10 w-20"
                                />
                                <Input
                                    value={chatColor}
                                    onChange={(e) => setChatColor(e.target.value)}
                                    placeholder="#FF6900"
                                />
                            </div>
                        </div>
                    </div>
                </Card>

                <Card className="p-4">
                    <h3 className="mb-4 text-sm font-semibold text-foreground">Notifications</h3>

                    <div className="flex items-center justify-between">
                        <div className="space-y-0.5">
                            <Label htmlFor="sound-mention">Sound on Mention</Label>
                            <p className="text-xs text-muted-foreground">Play a sound when someone mentions you</p>
                        </div>
                        <Switch
                            id="sound-mention"
                            checked={soundOnMention}
                            onCheckedChange={setSoundOnMention}
                        />
                    </div>
                </Card>

                <div className="flex justify-end">
                    <Button onClick={handleSave} disabled={saving} size="sm">
                        {saving ? "Saving..." : "Save Settings"}
                    </Button>
                </div>
            </div>
        </ScrollArea>
    )
}