"use client"

import { useState, useEffect } from "react"
import { Card } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Check, X, Eye } from "lucide-react"
import { getApiRoot } from "@/lib/auth"
import { SuggestionAcceptanceModal } from "./suggestion-acceptance-modal"

interface Suggestion {
  id: string
  name: string
  tile_name: string
  reason: string
  status: string
  reviewed_by?: string
  reviewed_at?: string
  created_at: string
}

export function SuggestionsPanel() {
  const [suggestions, setSuggestions] = useState<Suggestion[]>([])
  const [loading, setLoading] = useState(true)
  const [acceptingSuggestion, setAcceptingSuggestion] = useState<Suggestion | null>(null)

  useEffect(() => {
    fetchSuggestions()
  }, [])

  const fetchSuggestions = async () => {
    try {
      const response = await fetch(`${getApiRoot()}/suggestions`, {
        credentials: "include",
      })
      if (!response.ok) {
        console.error("Failed to fetch suggestions:", response.status)
        setSuggestions([])
        return
      }
      const data = await response.json()
      setSuggestions(data.suggestions || [])
    } catch (error) {
      console.error("Failed to fetch suggestions:", error)
      setSuggestions([])
    } finally {
      setLoading(false)
    }
  }

  const updateSuggestionStatus = async (id: string, status: string) => {
    try {
      const response = await fetch(`${getApiRoot()}/suggestions/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ status }),
      })
      if (!response.ok) {
        console.error("Failed to update suggestion:", response.status)
        return
      }
      fetchSuggestions() // Refresh list
    } catch (error) {
      console.error("Failed to update suggestion:", error)
    }
  }

  const handleAccept = (suggestion: Suggestion) => {
    setAcceptingSuggestion(suggestion)
  }

  const handleReject = (id: string) => {
    updateSuggestionStatus(id, "rejected")
  }

  const getStatusBadge = (status: string) => {
    switch (status) {
      case "pending":
        return <Badge variant="secondary">Pending</Badge>
      case "accepted":
        return <Badge variant="default">Accepted</Badge>
      case "rejected":
        return <Badge variant="destructive">Rejected</Badge>
      default:
        return <Badge>{status}</Badge>
    }
  }

  if (loading) {
    return <div>Loading suggestions...</div>
  }

  return (
    <Card className="p-6">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Name</TableHead>
            <TableHead>Tile Name</TableHead>
            <TableHead>Reason</TableHead>
            <TableHead>Status</TableHead>
            <TableHead>Created</TableHead>
            <TableHead>Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {suggestions.map((suggestion) => (
            <TableRow key={suggestion.id}>
              <TableCell>{suggestion.name}</TableCell>
              <TableCell>{suggestion.tile_name}</TableCell>
              <TableCell className="max-w-xs truncate">{suggestion.reason}</TableCell>
              <TableCell>{getStatusBadge(suggestion.status)}</TableCell>
              <TableCell>{new Date(suggestion.created_at).toLocaleDateString()}</TableCell>
              <TableCell>
                <div className="flex gap-2">
                  {suggestion.status === "pending" && (
                    <>
                      <Button size="sm" onClick={() => handleAccept(suggestion)}>
                        <Check className="h-4 w-4" />
                      </Button>
                      <Button size="sm" variant="destructive" onClick={() => handleReject(suggestion.id)}>
                        <X className="h-4 w-4" />
                      </Button>
                    </>
                  )}
                  <Button size="sm" variant="outline">
                    <Eye className="h-4 w-4" />
                  </Button>
                </div>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>

      {acceptingSuggestion && (
        <SuggestionAcceptanceModal
          suggestion={acceptingSuggestion}
          onClose={() => setAcceptingSuggestion(null)}
          onAccept={() => {
            updateSuggestionStatus(acceptingSuggestion.id, "accepted")
            setAcceptingSuggestion(null)
          }}
        />
      )}
    </Card>
  )
}