"use client"

import {useEffect, useState} from "react"
import {Card} from "@/components/ui/card"
import {Button} from "@/components/ui/button"
import {Input} from "@/components/ui/input"
import {Label} from "@/components/ui/label"
import {Pause, Play, Plus, Trash2} from "lucide-react"
import {getApiRoot} from "@/lib/auth";

interface Timer {
    id: string
    title: string
    duration: number
    created_by?: string
    show_id?: string
    starts_at?: string
    expires_at?: string
    is_active: boolean
    settings: any
    created_at: string
    updated_at: string
}

export function TimerPanel() {
    const [timers, setTimers] = useState<Timer[]>([])
    const [newTimerTitle, setNewTimerTitle] = useState("")
    const [newTimerDuration, setNewTimerDuration] = useState("")

    useEffect(() => {
        fetchTimers()
    }, [])

    const fetchTimers = async () => {
        try {
            const response = await fetch(`${getApiRoot()}/timers?is_active=true`)
            const data = await response.json()
            setTimers(data.timers)
        } catch (error) {
            console.error("Failed to fetch timers:", error)
        }
    }

    const addTimer = async () => {
        if (!newTimerTitle || !newTimerDuration) return

        const duration = Number.parseInt(newTimerDuration)
        try {
            const response = await fetch(`${getApiRoot()}/timers`, {
                method: "POST",
                headers: {"Content-Type": "application/json"},
                credentials: "include",
                body: JSON.stringify({
                    title: newTimerTitle,
                    duration: duration,
                }),
            })
            if (response.ok) {
                setNewTimerTitle("")
                setNewTimerDuration("")
                fetchTimers()
            } else {
                console.error("Failed to create timer")
            }
        } catch (error) {
            console.error("Failed to create timer:", error)
        }
    }

    const toggleTimer = async (id: string) => {
        const timer = timers.find(t => t.id === id)
        if (!timer) return

        const action = timer.is_active ? "stop" : "start"
        try {
            await fetch(`${getApiRoot()}/timers/${id}/${action}`, {
                method: "POST",
                credentials: "include",
            })
            fetchTimers()
        } catch (error) {
            console.error("Failed to toggle timer:", error)
        }
    }

    const deleteTimer = async (id: string) => {
        try {
            await fetch(`${getApiRoot()}/timers/${id}`, {
                method: "DELETE",
                credentials: "include",
            })
            fetchTimers()
        } catch (error) {
            console.error("Failed to delete timer:", error)
        }
    }

    const formatTime = (expiresAt?: string) => {
        if (!expiresAt) return "00:00"
        const now = new Date()
        const expires = new Date(expiresAt)
        const diff = Math.max(0, Math.floor((expires.getTime() - now.getTime()) / 1000))
        const mins = Math.floor(diff / 60)
        const secs = diff % 60
        return `${mins}:${secs.toString().padStart(2, "0")}`
    }

    return (
        <div className="space-y-4">
            <Card className="p-4">
                <h3 className="mb-4 font-semibold text-foreground">Create Timer</h3>
                <div className="grid gap-4 sm:grid-cols-[1fr_auto_auto]">
                    <div className="space-y-2">
                        <Label htmlFor="timer-name">Timer Name</Label>
                        <Input
                            id="timer-title"
                            placeholder="Sponsor break"
                            value={newTimerTitle}
                            onChange={(e) => setNewTimerTitle(e.target.value)}
                        />
                    </div>
                    <div className="space-y-2">
                        <Label htmlFor="timer-duration">Duration (minutes)</Label>
                        <Input
                            id="timer-duration"
                            type="number"
                            placeholder="300"
                            value={newTimerDuration}
                            onChange={(e) => setNewTimerDuration(e.target.value)}
                        />
                    </div>
                    <div className="flex items-end">
                        <Button onClick={addTimer} className="gap-2">
                            <Plus className="h-4 w-4"/>
                            Add
                        </Button>
                    </div>
                </div>
            </Card>

            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                {timers.map((timer) => (
                    <Card key={timer.id} className="p-4">
                        <div className="mb-3 flex items-start justify-between">
                            <h4 className="font-medium text-foreground">{timer.title}</h4>
                            <Button variant="ghost" size="icon" onClick={() => deleteTimer(timer.id)}>
                                <Trash2 className="h-4 w-4"/>
                            </Button>
                        </div>

                        <div className="mb-4 text-center">
                            <p className="text-3xl font-bold text-foreground">{formatTime(timer.expires_at)}</p>
                            <p className="text-sm text-muted-foreground">Duration: {timer.duration}s</p>
                        </div>

                        <Button variant="outline" className="w-full gap-2 bg-transparent"
                                onClick={() => toggleTimer(timer.id)}>
                            {timer.is_active ? (
                                <>
                                    <Pause className="h-4 w-4"/>
                                    Stop
                                </>
                            ) : (
                                <>
                                    <Play className="h-4 w-4"/>
                                    Start
                                </>
                            )}
                        </Button>
                    </Card>
                ))}
            </div>

            {timers.length === 0 && (
                <Card className="p-8 text-center">
                    <p className="text-muted-foreground">No timers created yet</p>
                </Card>
            )}
        </div>
    )
}
