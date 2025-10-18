"use client"

import { useState, useEffect } from "react"
import { Card } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Badge } from "@/components/ui/badge"
import { Check, X, Clock } from "lucide-react"
import {getApiRoot} from "@/lib/auth";

interface TileSuggestion {
  id: number
  name: string
  tileName: string
  reason: string
  status: "pending" | "approved" | "denied"
  createdAt: string
}

export function SuggestionManagementPanel() {
  const [suggestions, setSuggestions] = useState<TileSuggestion[]>([])
  const [loading, setLoading] = useState(true)
  const [filter, setFilter] = useState<"all" | "pending" | "approved" | "denied">("pending")

  useEffect(() => {
    fetchSuggestions()
  }, [])

  const fetchSuggestions = async () => {
    try {
      const response = await fetch("${getApiRoot()}/host/suggestions")
      const data = await response.json()
      setSuggestions(data)
    } catch (error) {
      console.error("Failed to fetch suggestions:", error)
    } finally {
      setLoading(false)
    }
  }

  const handleApprove = async (id: number) => {
    try {
      await fetch(`${getApiRoot()}/host/suggestions/${id}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ status: "approved" }),
      })
      fetchSuggestions()
    } catch (error) {
      console.error("Failed to approve suggestion:", error)
    }
  }

  const handleDeny = async (id: number) => {
    try {
      await fetch(`${getApiRoot()}/host/suggestions/${id}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ status: "denied" }),
      })
      fetchSuggestions()
    } catch (error) {
      console.error("Failed to deny suggestion:", error)
    }
  }

  const filteredSuggestions = suggestions.filter((s) => filter === "all" || s.status === filter)

  if (loading) {
    return (
      <Card className="flex items-center justify-center p-8">
        <p className="text-muted-foreground">Loading suggestions...</p>
      </Card>
    )
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold text-foreground">Tile Suggestions</h2>
        <div className="flex gap-2">
          <Button variant={filter === "all" ? "default" : "outline"} size="sm" onClick={() => setFilter("all")}>
            All
          </Button>
          <Button variant={filter === "pending" ? "default" : "outline"} size="sm" onClick={() => setFilter("pending")}>
            Pending
          </Button>
          <Button
            variant={filter === "approved" ? "default" : "outline"}
            size="sm"
            onClick={() => setFilter("approved")}
          >
            Approved
          </Button>
          <Button variant={filter === "denied" ? "default" : "outline"} size="sm" onClick={() => setFilter("denied")}>
            Denied
          </Button>
        </div>
      </div>

      <ScrollArea className="h-[calc(100vh-14rem)]">
        <div className="space-y-4">
          {filteredSuggestions.map((suggestion) => (
            <Card key={suggestion.id} className="p-4">
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <div className="mb-2 flex items-center gap-2">
                    <h3 className="font-semibold text-foreground">{suggestion.tileName}</h3>
                    <Badge
                      variant={
                        suggestion.status === "approved"
                          ? "default"
                          : suggestion.status === "denied"
                            ? "destructive"
                            : "secondary"
                      }
                    >
                      {suggestion.status === "pending" && <Clock className="mr-1 h-3 w-3" />}
                      {suggestion.status}
                    </Badge>
                  </div>

                  <p className="mb-2 text-sm text-muted-foreground">Suggested by: {suggestion.name}</p>
                  <p className="text-sm text-foreground">{suggestion.reason}</p>
                  <p className="mt-2 text-xs text-muted-foreground">
                    {new Date(suggestion.createdAt).toLocaleDateString()}
                  </p>
                </div>

                {suggestion.status === "pending" && (
                  <div className="flex gap-2">
                    <Button variant="outline" size="icon" onClick={() => handleApprove(suggestion.id)}>
                      <Check className="h-4 w-4 text-green-500" />
                    </Button>
                    <Button variant="outline" size="icon" onClick={() => handleDeny(suggestion.id)}>
                      <X className="h-4 w-4 text-red-500" />
                    </Button>
                  </div>
                )}
              </div>
            </Card>
          ))}

          {filteredSuggestions.length === 0 && (
            <Card className="p-8 text-center">
              <p className="text-muted-foreground">No {filter !== "all" && filter} suggestions found</p>
            </Card>
          )}
        </div>
      </ScrollArea>
    </div>
  )
}
