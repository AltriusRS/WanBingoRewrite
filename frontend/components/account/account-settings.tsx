"use client"

import { useState, useEffect } from "react"
import { Card } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { ArrowLeft, Upload } from "lucide-react"
import Link from "next/link"
import { Switch } from "@/components/ui/switch"
import { useAuth } from "@/components/auth"
import { Player } from "@/lib/auth"

export function AccountSettings() {
   const auth = useAuth()
   const user = auth.user
   const [displayName, setDisplayName] = useState("")
   const [avatarUrl, setAvatarUrl] = useState("")
   const [chatColor, setChatColor] = useState("#FF6900")
   const [backgroundImageEnabled, setBackgroundImageEnabled] = useState(false)
   const [saving, setSaving] = useState(false)

   useEffect(() => {
     if (user) {
       setDisplayName(user.display_name || "")
       setAvatarUrl(user.avatar || "")
       // Load settings
       if (user.settings) {
         const settings = user.settings as any
         setChatColor(settings.chatColor || "#FF6900")
         setBackgroundImageEnabled(settings.backgroundImageEnabled || false)
       }
     }
   }, [user])

   const handleSave = async () => {
     setSaving(true)

     try {
       const response = await fetch("https://api.bingo.local/users/me", {
         method: "PUT",
         credentials: "include",
         headers: { "Content-Type": "application/json" },
         body: JSON.stringify({
           display_name: displayName,
           avatar: avatarUrl,
           settings: {
             chatColor,
             backgroundImageEnabled,
           },
         }),
       })
       if (!response.ok) {
         throw new Error("Failed to update")
       }
       // Refetch user data
       await auth.refetch()
     } catch (error) {
       console.error("Failed to save settings:", error)
     } finally {
       setSaving(false)
     }
   }

  if (!user) {
    return <div>Loading...</div>
  }

  return (
    <div className="min-h-screen bg-background">
      <header className="border-b border-border bg-card">
        <div className="container mx-auto flex items-center gap-4 px-4 py-4">
          <Link href="/">
            <Button variant="ghost" size="icon">
              <ArrowLeft className="h-4 w-4" />
            </Button>
          </Link>
          <div>
            <h1 className="text-xl font-semibold text-foreground">Account Settings</h1>
            <p className="text-sm text-muted-foreground">Manage your profile and preferences</p>
          </div>
        </div>
      </header>

      <div className="container mx-auto max-w-2xl space-y-6 p-4 py-8">
        <Card className="p-6">
          <h2 className="mb-4 text-lg font-semibold text-foreground">Profile</h2>

          <div className="space-y-6">
             <div className="flex items-center gap-4">
               <Avatar className="h-20 w-20">
                 <AvatarImage src={avatarUrl || "/placeholder-user.jpg"} alt={displayName} />
                 <AvatarFallback className="bg-primary/10 text-lg">
                   {displayName
                     .split(" ")
                     .map((n) => n[0])
                     .join("")
                     .toUpperCase()}
                 </AvatarFallback>
               </Avatar>
               <div className="flex-1">
                 <Label htmlFor="avatar-url">Avatar URL</Label>
                 <Input
                   id="avatar-url"
                   value={avatarUrl}
                   onChange={(e) => setAvatarUrl(e.target.value)}
                   placeholder="https://example.com/avatar.jpg"
                 />
               </div>
             </div>

            <div className="space-y-2">
              <Label htmlFor="display-name">Display Name</Label>
              <Input
                id="display-name"
                value={displayName}
                onChange={(e) => setDisplayName(e.target.value)}
                placeholder="Your display name"
              />
            </div>

             <div className="space-y-2">
               <Label htmlFor="player-id">Player ID</Label>
               <Input id="player-id" value={user.id} disabled />
               <p className="text-xs text-muted-foreground">Your unique player identifier</p>
             </div>
          </div>
        </Card>

        <Card className="p-6">
          <h2 className="mb-4 text-lg font-semibold text-foreground">Chat Settings</h2>

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
                <Input value={chatColor} onChange={(e) => setChatColor(e.target.value)} placeholder="#FF6900" />
              </div>
            </div>
          </div>
        </Card>

        <Card className="p-6">
          <h2 className="mb-4 text-lg font-semibold text-foreground">Appearance</h2>

          <div className="flex items-center justify-between">
            <div className="space-y-0.5">
              <Label htmlFor="background-toggle">Background Image</Label>
              <p className="text-sm text-muted-foreground">Show episode thumbnail behind bingo tiles</p>
            </div>
            <Switch
              id="background-toggle"
              checked={backgroundImageEnabled}
              onCheckedChange={setBackgroundImageEnabled}
            />
          </div>
        </Card>

        <div className="flex justify-end gap-2">
          <Link href="/">
            <Button variant="outline">Cancel</Button>
          </Link>
          <Button onClick={handleSave} disabled={saving}>
            {saving ? "Saving..." : "Save Changes"}
          </Button>
        </div>
      </div>
    </div>
  )
}
