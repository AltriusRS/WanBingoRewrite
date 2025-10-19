"use client"

import {useEffect, useState} from "react"
import {Card} from "@/components/ui/card"
import {Button} from "@/components/ui/button"
import {Input} from "@/components/ui/input"
import {Label} from "@/components/ui/label"

import {ArrowLeft} from "lucide-react"
import Link from "next/link"
import {Switch} from "@/components/ui/switch"
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue} from "@/components/ui/select"
import {useAuth} from "@/components/auth"
import {getApiRoot} from "@/lib/auth";

export function AccountSettings() {
    const auth = useAuth()
    const user = auth.user
    const [displayName, setDisplayName] = useState("")
    const [chatColor, setChatColor] = useState("#FF6900")
    const [soundOnMention, setSoundOnMention] = useState(true)
    const [backgroundImageEnabled, setBackgroundImageEnabled] = useState(false)
    const [preferredTheme, setPreferredTheme] = useState("dark")
    const [autoYoutubePlayback, setAutoYoutubePlayback] = useState(false)
    const [highlightConfirmedTiles, setHighlightConfirmedTiles] = useState(true)
    const [confettiEnabled, setConfettiEnabled] = useState(true)
    const [disableWinAnnouncements, setDisableWinAnnouncements] = useState(false)
    const [showTileScores, setShowTileScores] = useState(true)
    const [showMaxScore, setShowMaxScore] = useState(true)
    const [showMultiplier, setShowMultiplier] = useState(true)
    const [boardTextSize, setBoardTextSize] = useState("medium")
    const [hostPanelTextSize, setHostPanelTextSize] = useState("medium")
    const [font, setFont] = useState("default")
    const [dyslexicFriendlyFont, setDyslexicFriendlyFont] = useState(false)
    const [saving, setSaving] = useState(false)

    const themeOptions = [
        { name: 'WAN Show', value: 'dark' },
        { name: 'Light', value: 'light' },
        { name: 'Winter', value: 'winter' },
        { name: 'Halloween', value: 'halloween' },
        { name: 'Easter', value: 'easter' },
        { name: 'Summer', value: 'summer' },
        { name: 'Pitch Black', value: 'pitch-black' },
    ]

    useEffect(() => {
        if (user) {
            setDisplayName(user.display_name || "")
            // Load settings
            if (user.settings) {
                const settings = user.settings as any
                // Chat settings
                if (settings.chat) {
                    setChatColor(settings.chat.color || "#FF6900")
                    setSoundOnMention(settings.chat.soundOnMention !== false) // Default to true
                } else {
                    // Fallback for old flat structure
                    setChatColor(settings.chatColor || "#FF6900")
                    setSoundOnMention(settings.soundOnMention !== false)
                }

                // Theme settings
                if (settings.themes) {
                    setPreferredTheme(settings.themes.preferred || "dark")
                } else {
                    // Fallback for old flat structure
                    setPreferredTheme(settings.preferredTheme || "dark")
                }

                // Video settings
                if (settings.video) {
                    setAutoYoutubePlayback(settings.video.autoYoutubePlayback || false)
                    setBackgroundImageEnabled(settings.video.backgroundImageEnabled || false)
                } else {
                    // Fallback for old flat structure
                    setAutoYoutubePlayback(settings.autoYoutubePlayback || false)
                    setBackgroundImageEnabled(settings.backgroundImageEnabled || false)
                }

                // Gameplay settings
                if (settings.gameplay) {
                    setHighlightConfirmedTiles(settings.gameplay.highlightConfirmedTiles !== false) // Default to true
                    setConfettiEnabled(settings.gameplay.confetti !== false) // Default to true
                    setDisableWinAnnouncements(settings.gameplay.disableWinAnnouncements || false) // Default to false
                    setShowTileScores(settings.gameplay.showTileScores !== false) // Default to true
                    setShowMaxScore(settings.gameplay.showMaxScore !== false) // Default to true
                    setShowMultiplier(settings.gameplay.showMultiplier !== false) // Default to true
                } else {
                    // Fallback for old flat structure
                    setHighlightConfirmedTiles(settings.highlightConfirmedTiles !== false)
                }

                // Appearance settings
                if (settings.appearance?.board?.textSize) {
                    setBoardTextSize(settings.appearance.board.textSize)
                }
                if (settings.appearance?.hostPanel?.textSize) {
                    setHostPanelTextSize(settings.appearance.hostPanel.textSize)
                }
                if (settings.appearance?.font) {
                    setFont(settings.appearance.font)
                }
                if (settings.appearance?.dyslexicFriendlyFont !== undefined) {
                    setDyslexicFriendlyFont(settings.appearance.dyslexicFriendlyFont)
                }
            }
        }
    }, [user])

    const handleSave = async () => {
        setSaving(true)

        try {
            const response = await fetch(`${getApiRoot()}/users/me`, {
                method: "PUT",
                credentials: "include",
                headers: {"Content-Type": "application/json"},
                body: JSON.stringify({
                    display_name: displayName,
                    settings: {
                        chat: {
                            color: chatColor,
                            soundOnMention,
                        },
                        themes: {
                            preferred: preferredTheme,
                        },
                        video: {
                            autoYoutubePlayback,
                            backgroundImageEnabled,
                        },
                          gameplay: {
                              highlightConfirmedTiles,
                              confetti: confettiEnabled,
                              disableWinAnnouncements,
                              showTileScores,
                              showMaxScore,
                              showMultiplier,
                          },
                          appearance: {
                              board: {
                                  textSize: boardTextSize,
                              },
                              hostPanel: {
                                  textSize: hostPanelTextSize,
                              },
                              font,
                              dyslexicFriendlyFont,
                          },
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
                            <ArrowLeft className="h-4 w-4"/>
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
                            <Input id="player-id" value={user.id} disabled/>
                            <p className="text-xs text-muted-foreground">Your unique player identifier</p>
                        </div>
                    </div>
                </Card>

                <Card className="p-6">
                    <h2 className="mb-4 text-lg font-semibold text-foreground">Chat Settings</h2>

                    <div className="space-y-6">
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
                                <Input value={chatColor} onChange={(e) => setChatColor(e.target.value)}
                                       placeholder="#FF6900"/>
                            </div>
                        </div>

                        <div className="flex items-center justify-between">
                            <div className="space-y-0.5">
                                <Label htmlFor="sound-mention">Sound on Mention</Label>
                                <p className="text-sm text-muted-foreground">Play a sound when someone mentions you in chat</p>
                            </div>
                            <Switch
                                id="sound-mention"
                                checked={soundOnMention}
                                onCheckedChange={setSoundOnMention}
                            />
                        </div>
                    </div>
                </Card>

                <Card className="p-6">
                    <h2 className="mb-4 text-lg font-semibold text-foreground">Themes</h2>

                    <div className="space-y-2">
                        <Label htmlFor="theme-select">Preferred Theme</Label>
                        <Select value={preferredTheme} onValueChange={setPreferredTheme}>
                            <SelectTrigger>
                                <SelectValue placeholder="Select a theme" />
                            </SelectTrigger>
                            <SelectContent>
                                {themeOptions.map((theme) => (
                                    <SelectItem key={theme.value} value={theme.value}>
                                        {theme.name}
                                    </SelectItem>
                                ))}
                            </SelectContent>
                        </Select>
                        <p className="text-sm text-muted-foreground">Choose your preferred theme for the application</p>
                    </div>
                 </Card>

                 <Card className="p-6">
                     <h2 className="mb-4 text-lg font-semibold text-foreground">Appearance</h2>

                     <div className="space-y-2">
                         <Label htmlFor="board-text-size">Bingo Board Text Size</Label>
                         <Select value={boardTextSize} onValueChange={setBoardTextSize}>
                             <SelectTrigger>
                                 <SelectValue placeholder="Select text size" />
                             </SelectTrigger>
                             <SelectContent>
                                 <SelectItem value="small">Small</SelectItem>
                                 <SelectItem value="medium">Medium</SelectItem>
                                 <SelectItem value="large">Large</SelectItem>
                             </SelectContent>
                         </Select>
                         <p className="text-sm text-muted-foreground">Choose the size of text on bingo tiles</p>
                     </div>

                     <div className="space-y-2">
                         <Label htmlFor="host-panel-text-size">Host Panel Text Size</Label>
                         <Select value={hostPanelTextSize} onValueChange={setHostPanelTextSize}>
                             <SelectTrigger>
                                 <SelectValue placeholder="Select text size" />
                             </SelectTrigger>
                             <SelectContent>
                                 <SelectItem value="small">Small</SelectItem>
                                 <SelectItem value="medium">Medium</SelectItem>
                                 <SelectItem value="large">Large</SelectItem>
                             </SelectContent>
                         </Select>
                         <p className="text-sm text-muted-foreground">Choose the size of text on host panel tiles</p>
                     </div>

                     <div className="space-y-2">
                         <Label htmlFor="font">Font</Label>
                         <Select value={font} onValueChange={setFont}>
                             <SelectTrigger>
                                 <SelectValue placeholder="Select font" />
                             </SelectTrigger>
                             <SelectContent>
                                 <SelectItem value="default">Default</SelectItem>
                                 <SelectItem value="serif">Serif</SelectItem>
                                 <SelectItem value="sans-serif">Sans Serif</SelectItem>
                             </SelectContent>
                         </Select>
                         <p className="text-sm text-muted-foreground">Choose the font for the website</p>
                     </div>

                     <div className="flex items-center justify-between">
                         <div className="space-y-0.5">
                             <Label htmlFor="dyslexic-friendly-font">Dyslexia Friendly Font</Label>
                             <p className="text-sm text-muted-foreground">Use a font designed for people with dyslexia</p>
                         </div>
                         <Switch
                             id="dyslexic-friendly-font"
                             checked={dyslexicFriendlyFont}
                             onCheckedChange={setDyslexicFriendlyFont}
                         />
                     </div>
                 </Card>

                 <Card className="p-6">
                     <h2 className="mb-4 text-lg font-semibold text-foreground">Video</h2>

                    <div className="space-y-6">
                        <div className="flex items-center justify-between">
                            <div className="space-y-0.5">
                                <Label htmlFor="youtube-autoplay">Auto YouTube Playback</Label>
                                <p className="text-sm text-muted-foreground">Automatically embed and play YouTube video when show is live</p>
                            </div>
                            <Switch
                                id="youtube-autoplay"
                                checked={autoYoutubePlayback}
                                onCheckedChange={setAutoYoutubePlayback}
                            />
                        </div>

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
                     </div>
                 </Card>

                 <Card className="p-6">
                     <h2 className="mb-4 text-lg font-semibold text-foreground">Gameplay</h2>

                     <div className="space-y-6">
                          <div className="flex items-center justify-between">
                              <div className="space-y-0.5">
                                  <Label htmlFor="highlight-confirmed">Highlight Confirmed Tiles</Label>
                                  <p className="text-sm text-muted-foreground">Add a subtle border highlight to tiles that have been confirmed</p>
                              </div>
                              <Switch
                                  id="highlight-confirmed"
                                  checked={highlightConfirmedTiles}
                                  onCheckedChange={setHighlightConfirmedTiles}
                              />
                          </div>

                          <div className="flex items-center justify-between">
                              <div className="space-y-0.5">
                                  <Label htmlFor="confetti-enabled">Confetti on Win</Label>
                                  <p className="text-sm text-muted-foreground">Show confetti animation when you get a valid bingo</p>
                              </div>
                              <Switch
                                  id="confetti-enabled"
                                  checked={confettiEnabled}
                                  onCheckedChange={setConfettiEnabled}
                              />
                          </div>

                          <div className="flex items-center justify-between">
                              <div className="space-y-0.5">
                                  <Label htmlFor="disable-win-announcements">Disable Win Announcements</Label>
                                  <p className="text-sm text-muted-foreground">Don&apos;t announce your wins in chat when you get bingo</p>
                              </div>
                              <Switch
                                  id="disable-win-announcements"
                                  checked={disableWinAnnouncements}
                                  onCheckedChange={setDisableWinAnnouncements}
                              />
                          </div>

                          <div className="flex items-center justify-between">
                              <div className="space-y-0.5">
                                  <Label htmlFor="show-tile-scores">Show Tile Scores</Label>
                                  <p className="text-sm text-muted-foreground">Display point values in the bottom right corner of tiles</p>
                              </div>
                              <Switch
                                  id="show-tile-scores"
                                  checked={showTileScores}
                                  onCheckedChange={setShowTileScores}
                              />
                          </div>

                          <div className="flex items-center justify-between">
                              <div className="space-y-0.5">
                                  <Label htmlFor="show-max-score">Show Max Score</Label>
                                  <p className="text-sm text-muted-foreground">Display the maximum possible score for the current board</p>
                              </div>
                              <Switch
                                  id="show-max-score"
                                  checked={showMaxScore}
                                  onCheckedChange={setShowMaxScore}
                              />
                          </div>

                          <div className="flex items-center justify-between">
                              <div className="space-y-0.5">
                                  <Label htmlFor="show-multiplier">Show Multiplier</Label>
                                  <p className="text-sm text-muted-foreground">Display the current score multiplier percentage</p>
                              </div>
                              <Switch
                                  id="show-multiplier"
                                  checked={showMultiplier}
                                  onCheckedChange={setShowMultiplier}
                              />
                          </div>
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
