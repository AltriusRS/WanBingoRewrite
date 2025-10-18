"use client"

import { createContext, useContext, useEffect, useState, ReactNode } from "react"
import { getApiRoot } from "@/lib/auth"

interface HostContextType {
  confirmedTiles: Set<string>
  locks: Map<string, { tileId: string; lockedBy: string; expiresAt: number }>
  refreshConfirmedTiles: () => Promise<void>
}

const HostContext = createContext<HostContextType | null>(null)

export function HostProvider({ children }: { children: ReactNode }) {
  const [confirmedTiles, setConfirmedTiles] = useState<Set<string>>(new Set())
  const [locks, setLocks] = useState<Map<string, { tileId: string; lockedBy: string; expiresAt: number }>>(new Map())

  const refreshConfirmedTiles = async () => {
    try {
      console.log("[HostContext] Fetching confirmed tiles...")
      const response = await fetch(`${getApiRoot()}/api/host/confirmed-tiles`, {
        credentials: "include",
      })
      if (!response.ok) {
        console.error("Failed to fetch confirmed tiles:", response.status)
        return
      }
      const data = await response.json()
      console.log("[HostContext] Confirmed tiles loaded:", data)
      setConfirmedTiles(new Set(data))
    } catch (error) {
      console.error("Failed to fetch confirmed tiles:", error)
    }
  }

  useEffect(() => {
    // Initial fetch of confirmed tiles
    refreshConfirmedTiles()

    // Connect to host SSE for real-time updates
    console.log("[HostContext] Connecting to host SSE stream:", `${getApiRoot()}/host/stream`)
    const eventSource = new EventSource(`${getApiRoot()}/host/stream`)

    eventSource.onopen = () => {
      console.log("[HostContext] Host SSE connection opened")
    }

    eventSource.onerror = (error) => {
      console.error("[HostContext] Host SSE connection error:", error)
    }

    eventSource.onmessage = (event) => {
      try {
        console.log("[HostContext] Received SSE event:", event.data)
        const envelope = JSON.parse(event.data)
        const { opcode, data } = envelope

        switch (opcode) {
          case "tile.lock":
            console.log(`[HostContext] Processing tile.lock: ${data.tileId} by ${data.user}`)
            setLocks((prev) => {
              const newLocks = new Map(prev)
              newLocks.set(data.tileId, {
                tileId: data.tileId,
                lockedBy: data.user,
                expiresAt: Date.now() + 5000,
              })
              console.log("[HostContext] Updated locks:", Array.from(newLocks.entries()))
              return newLocks
            })
            break

          case "tile.unlock":
            console.log(`[HostContext] Processing tile.unlock: ${data.tileId}`)
            setLocks((prev) => {
              const updated = new Map(prev)
              updated.delete(data.tileId)
              console.log("[HostContext] Updated locks after unlock:", Array.from(updated.entries()))
              return updated
            })
            break

          case "tile.confirm":
            console.log(`[HostContext] Processing tile.confirm: ${data.tileId}`)
            setConfirmedTiles((prev) => {
              const newSet = new Set(prev)
              newSet.add(data.tileId)
              console.log("[HostContext] Updated confirmed tiles:", Array.from(newSet))
              return newSet
            })
            break

          case "tile.revoke":
            console.log(`[HostContext] Processing tile.revoke: ${data.tileId}`)
            setConfirmedTiles((prev) => {
              const newSet = new Set(prev)
              newSet.delete(data.tileId)
              console.log("[HostContext] Updated confirmed tiles after revoke:", Array.from(newSet))
              return newSet
            })
            break

          default:
            console.log("[HostContext] Unknown opcode:", opcode)
            break
        }
      } catch (error) {
        console.error("Failed to parse SSE event:", error)
      }
    }

    const cleanupExpiredLocks = () => {
      const now = Date.now()
      setLocks((prev) => {
        const updated = new Map(prev)
        let hasChanges = false
        for (const [tileId, lock] of updated.entries()) {
          if (lock.expiresAt < now) {
            updated.delete(tileId)
            hasChanges = true
          }
        }
        return hasChanges ? updated : prev
      })
    }

    const interval = setInterval(cleanupExpiredLocks, 1000)

    return () => {
      eventSource.close()
      clearInterval(interval)
    }
  }, [])

  return (
    <HostContext.Provider value={{
      confirmedTiles,
      locks,
      refreshConfirmedTiles,
    }}>
      {children}
    </HostContext.Provider>
  )
}

export function useHost() {
  const context = useContext(HostContext)
  if (!context) {
    throw new Error("useHost must be used within a HostProvider")
  }
  return context
}