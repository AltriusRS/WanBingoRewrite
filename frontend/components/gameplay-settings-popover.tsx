"use client"

import {useEffect, useState} from "react"
import {Button} from "@/components/ui/button"
import {Popover, PopoverContent, PopoverTrigger} from "@/components/ui/popover"
import {Label} from "@/components/ui/label"
import {Switch} from "@/components/ui/switch"
import {Settings} from "lucide-react"
import {useAuth} from "@/components/auth"
import {getApiRoot} from "@/lib/auth"

export function GameplaySettingsPopover() {
    const auth = useAuth()
    const user = auth.user
    const [highlightConfirmedTiles, setHighlightConfirmedTiles] = useState(true)
    const [confettiEnabled, setConfettiEnabled] = useState(true)
    const [disableWinAnnouncements, setDisableWinAnnouncements] = useState(false)
    const [showTileScores, setShowTileScores] = useState(true)
    const [showMaxScore, setShowMaxScore] = useState(true)
    const [showMultiplier, setShowMultiplier] = useState(true)
    const [saving, setSaving] = useState(false)
    const [open, setOpen] = useState(false)

    useEffect(() => {
        if (user?.settings?.gameplay) {
            const gameplay = user.settings.gameplay as any
            setHighlightConfirmedTiles(gameplay.highlightConfirmedTiles !== false)
            setConfettiEnabled(gameplay.confetti !== false)
            setDisableWinAnnouncements(gameplay.disableWinAnnouncements || false)
            setShowTileScores(gameplay.showTileScores !== false)
            setShowMaxScore(gameplay.showMaxScore !== false)
            setShowMultiplier(gameplay.showMultiplier !== false)
        }
    }, [user])

    const handleSave = async () => {
        if (!user) return

        setSaving(true)
        try {
            const response = await fetch(`${getApiRoot()}/users/me`, {
                method: "PUT",
                credentials: "include",
                headers: {"Content-Type": "application/json"},
                body: JSON.stringify({
                    settings: {
                        ...user.settings,
                        gameplay: {
                            highlightConfirmedTiles,
                            confetti: confettiEnabled,
                            disableWinAnnouncements,
                            showTileScores,
                            showMaxScore,
                            showMultiplier,
                        },
                    },
                }),
            })
            if (!response.ok) {
                throw new Error("Failed to update")
            }
            await auth.refetch()
            setOpen(false)
        } catch (error) {
            console.error("Failed to save settings:", error)
        } finally {
            setSaving(false)
        }
    }

    if (!user) return null

    return (
        <Popover open={open} onOpenChange={setOpen}>
            <PopoverTrigger asChild>
                <Button variant="ghost" size="sm" className="gap-2 bg-transparent">
                    <Settings className="h-4 w-4"/>
                    <span className="hidden sm:inline">Gameplay</span>
                </Button>
            </PopoverTrigger>
            <PopoverContent className="w-80" align="end">
                <div className="space-y-4">
                    <div className="font-medium text-sm">Gameplay Settings</div>

                    <div className="space-y-4">
                        <div className="flex items-center justify-between">
                            <div className="space-y-0.5">
                                <Label htmlFor="popover-highlight-confirmed" className="text-sm">Highlight Confirmed Tiles</Label>
                                <p className="text-xs text-muted-foreground">Add a subtle border highlight to tiles that have been confirmed</p>
                            </div>
                            <Switch
                                id="popover-highlight-confirmed"
                                checked={highlightConfirmedTiles}
                                onCheckedChange={setHighlightConfirmedTiles}
                            />
                        </div>

                        <div className="flex items-center justify-between">
                            <div className="space-y-0.5">
                                <Label htmlFor="popover-confetti-enabled" className="text-sm">Confetti on Win</Label>
                                <p className="text-xs text-muted-foreground">Show confetti animation when you get a valid bingo</p>
                            </div>
                            <Switch
                                id="popover-confetti-enabled"
                                checked={confettiEnabled}
                                onCheckedChange={setConfettiEnabled}
                            />
                        </div>

                        <div className="flex items-center justify-between">
                            <div className="space-y-0.5">
                                <Label htmlFor="popover-disable-win-announcements" className="text-sm">Disable Win Announcements</Label>
                                <p className="text-xs text-muted-foreground">Don&apos;t announce your wins in chat when you get bingo</p>
                            </div>
                            <Switch
                                id="popover-disable-win-announcements"
                                checked={disableWinAnnouncements}
                                onCheckedChange={setDisableWinAnnouncements}
                            />
                        </div>

                        <div className="flex items-center justify-between">
                            <div className="space-y-0.5">
                                <Label htmlFor="popover-show-tile-scores" className="text-sm">Show Tile Scores</Label>
                                <p className="text-xs text-muted-foreground">Display point values in the bottom right corner of tiles</p>
                            </div>
                            <Switch
                                id="popover-show-tile-scores"
                                checked={showTileScores}
                                onCheckedChange={setShowTileScores}
                            />
                        </div>

                        <div className="flex items-center justify-between">
                            <div className="space-y-0.5">
                                <Label htmlFor="popover-show-max-score" className="text-sm">Show Max Score</Label>
                                <p className="text-xs text-muted-foreground">Display the maximum possible score for the current board</p>
                            </div>
                            <Switch
                                id="popover-show-max-score"
                                checked={showMaxScore}
                                onCheckedChange={setShowMaxScore}
                            />
                        </div>

                        <div className="flex items-center justify-between">
                            <div className="space-y-0.5">
                                <Label htmlFor="popover-show-multiplier" className="text-sm">Show Multiplier</Label>
                                <p className="text-xs text-muted-foreground">Display the current score multiplier percentage</p>
                            </div>
                            <Switch
                                id="popover-show-multiplier"
                                checked={showMultiplier}
                                onCheckedChange={setShowMultiplier}
                            />
                        </div>
                    </div>

                    <div className="flex justify-end gap-2 pt-2 border-t">
                        <Button variant="outline" size="sm" onClick={() => setOpen(false)}>
                            Cancel
                        </Button>
                        <Button size="sm" onClick={handleSave} disabled={saving}>
                            {saving ? "Saving..." : "Save"}
                        </Button>
                    </div>
                </div>
            </PopoverContent>
        </Popover>
    )
}