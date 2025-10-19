"use client"

import {useEffect, useState} from "react"
import {Card} from "@/components/ui/card"
import {Button} from "@/components/ui/button"
import {Input} from "@/components/ui/input"
import {Label} from "@/components/ui/label"
import {Pause, Play, Plus, RotateCcw, Trash2} from "lucide-react"
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
    const [currentTime, setCurrentTime] = useState(new Date())
    const [expiredTimers, setExpiredTimers] = useState<Set<string>>(new Set())

    useEffect(() => {
        fetchTimers()
        // Update current time every second for countdown
        const interval = setInterval(() => {
            setCurrentTime(new Date())
        }, 1000)
        return () => clearInterval(interval)
    }, [])

    // Check for expired timers and play audio cue
    useEffect(() => {
        timers.forEach(timer => {
            if (timer.is_active && timer.expires_at) {
                const expires = new Date(timer.expires_at)
                const timeLeft = expires.getTime() - currentTime.getTime()

                if (timeLeft <= 0 && !expiredTimers.has(timer.id)) {
                    // Timer just expired, play sound
                    setExpiredTimers(prev => new Set(prev).add(timer.id))
                    playExpirationSound()
                } else if (timeLeft > 1000 && expiredTimers.has(timer.id)) {
                    // Timer was reset, remove from expired set
                    setExpiredTimers(prev => {
                        const newSet = new Set(prev)
                        newSet.delete(timer.id)
                        return newSet
                    })
                }
            }
        })
    }, [currentTime, timers, expiredTimers])

    const playExpirationSound = () => {
        // Create a simple beep sound
        const audioContext = new (window.AudioContext || (window as any).webkitAudioContext)()
        const oscillator = audioContext.createOscillator()
        const gainNode = audioContext.createGain()

        oscillator.connect(gainNode)
        gainNode.connect(audioContext.destination)

        oscillator.frequency.setValueAtTime(800, audioContext.currentTime)
        oscillator.frequency.setValueAtTime(600, audioContext.currentTime + 0.1)

        gainNode.gain.setValueAtTime(0.3, audioContext.currentTime)
        gainNode.gain.exponentialRampToValueAtTime(0.01, audioContext.currentTime + 0.5)

        oscillator.start(audioContext.currentTime)
        oscillator.stop(audioContext.currentTime + 0.5)
    }

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

        const duration = Number.parseInt(newTimerDuration) * 60
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

    const resetTimer = async (id: string) => {
        try {
            await fetch(`${getApiRoot()}/timers/${id}/reset`, {
                method: "POST",
                credentials: "include",
            })
            fetchTimers()
        } catch (error) {
            console.error("Failed to reset timer:", error)
        }
    }

    const formatTime = (expiresAt?: string) => {
        if (!expiresAt) return "00:00"
        const expires = new Date(expiresAt)
        const diff = Math.max(0, Math.floor((expires.getTime() - currentTime.getTime()) / 1000))
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

            {timers.length > 0 && (
                <div>
                    <h3 className="mb-2 font-semibold text-foreground">Active Timers</h3>
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-2">
                {timers.map((timer) => (
                    <Button
                        key={timer.id}
                        variant="outline"
                        className="justify-start text-left h-auto py-2 relative"
                        onClick={() => toggleTimer(timer.id)}
                    >
                        <div className="flex-1 truncate text-sm">
                            <div className="font-medium">{timer.title}</div>
                            <div className="text-xs text-muted-foreground">
                                {formatTime(timer.expires_at)} / {timer.duration}s
                            </div>
                        </div>
                        <div className="flex gap-1 absolute top-1 right-1">
                            <Button
                                variant="ghost"
                                size="icon"
                                className="h-5 w-5"
                                onClick={(e) => {
                                    e.stopPropagation()
                                    resetTimer(timer.id)
                                }}
                                title="Reset timer"
                            >
                                <RotateCcw className="h-3 w-3"/>
                            </Button>
                            <Button
                                variant="ghost"
                                size="icon"
                                className="h-5 w-5"
                                onClick={(e) => {
                                    e.stopPropagation()
                                    deleteTimer(timer.id)
                                }}
                                title="Delete timer"
                            >
                                <Trash2 className="h-3 w-3"/>
                            </Button>
                        </div>
                        {timer.is_active ? (
                            <Pause className="ml-2 h-4 w-4 text-green-500 flex-shrink-0"/>
                        ) : (
                            <Play className="ml-2 h-4 w-4 text-muted-foreground flex-shrink-0"/>
                        )}
                    </Button>
                ))}
                    </div>
                </div>
            )}

            {timers.length === 0 && (
                <div className="text-center py-8">
                    <p className="text-muted-foreground">No timers created yet</p>
                </div>
            )}
        </div>
    )
}
