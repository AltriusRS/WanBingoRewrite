"use client"

import { useState, useEffect } from "react"
import { Card } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Plus, Play, Pause, Trash2 } from "lucide-react"

interface Timer {
  id: string
  name: string
  duration: number
  remaining: number
  isRunning: boolean
}

export function TimerPanel() {
  const [timers, setTimers] = useState<Timer[]>([])
  const [newTimerName, setNewTimerName] = useState("")
  const [newTimerDuration, setNewTimerDuration] = useState("")

  useEffect(() => {
    const interval = setInterval(() => {
      setTimers((prev) =>
        prev.map((timer) => {
          if (timer.isRunning && timer.remaining > 0) {
            return { ...timer, remaining: timer.remaining - 1 }
          }
          return timer
        }),
      )
    }, 1000)

    return () => clearInterval(interval)
  }, [])

  const addTimer = () => {
    if (!newTimerName || !newTimerDuration) return

    const duration = Number.parseInt(newTimerDuration) * 60
    const newTimer: Timer = {
      id: Date.now().toString(),
      name: newTimerName,
      duration,
      remaining: duration,
      isRunning: false,
    }

    setTimers((prev) => [...prev, newTimer])
    setNewTimerName("")
    setNewTimerDuration("")

    // Broadcast timer to all viewers
    fetch("http://localhost:8080/api/host/timer", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ action: "create", timer: newTimer }),
    })
  }

  const toggleTimer = (id: string) => {
    setTimers((prev) =>
      prev.map((timer) => {
        if (timer.id === id) {
          const updated = { ...timer, isRunning: !timer.isRunning }
          // Broadcast update
          fetch("http://localhost:8080/api/host/timer", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ action: "toggle", timer: updated }),
          })
          return updated
        }
        return timer
      }),
    )
  }

  const deleteTimer = (id: string) => {
    setTimers((prev) => prev.filter((t) => t.id !== id))
    // Broadcast deletion
    fetch("http://localhost:8080/api/host/timer", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ action: "delete", timerId: id }),
    })
  }

  const formatTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60)
    const secs = seconds % 60
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
              id="timer-name"
              placeholder="Sponsor break"
              value={newTimerName}
              onChange={(e) => setNewTimerName(e.target.value)}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="timer-duration">Duration (minutes)</Label>
            <Input
              id="timer-duration"
              type="number"
              placeholder="5"
              value={newTimerDuration}
              onChange={(e) => setNewTimerDuration(e.target.value)}
            />
          </div>
          <div className="flex items-end">
            <Button onClick={addTimer} className="gap-2">
              <Plus className="h-4 w-4" />
              Add
            </Button>
          </div>
        </div>
      </Card>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {timers.map((timer) => (
          <Card key={timer.id} className="p-4">
            <div className="mb-3 flex items-start justify-between">
              <h4 className="font-medium text-foreground">{timer.name}</h4>
              <Button variant="ghost" size="icon" onClick={() => deleteTimer(timer.id)}>
                <Trash2 className="h-4 w-4" />
              </Button>
            </div>

            <div className="mb-4 text-center">
              <p className="text-3xl font-bold text-foreground">{formatTime(timer.remaining)}</p>
              <p className="text-sm text-muted-foreground">of {formatTime(timer.duration)}</p>
            </div>

            <Button variant="outline" className="w-full gap-2 bg-transparent" onClick={() => toggleTimer(timer.id)}>
              {timer.isRunning ? (
                <>
                  <Pause className="h-4 w-4" />
                  Pause
                </>
              ) : (
                <>
                  <Play className="h-4 w-4" />
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
