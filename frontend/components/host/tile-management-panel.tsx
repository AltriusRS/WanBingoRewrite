"use client"

import { useState, useEffect } from "react"
import { Card } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Badge } from "@/components/ui/badge"
import { Switch } from "@/components/ui/switch"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Edit2, Trash2, Plus } from "lucide-react"
import type { BingoTile } from "@/lib/bingoUtils"

interface TileStats {
  tileId: number
  winRate: number
  confirmRate: number
  timesConfirmed: number
  timesOnBoard: number
}

export function TileManagementPanel() {
  const [tiles, setTiles] = useState<BingoTile[]>([])
  const [stats, setStats] = useState<Map<number, TileStats>>(new Map())
  const [loading, setLoading] = useState(true)
  const [editingTile, setEditingTile] = useState<BingoTile | null>(null)
  const [showAddForm, setShowAddForm] = useState(false)

  useEffect(() => {
    fetchTiles()
    fetchStats()
  }, [])

  const fetchTiles = async () => {
    try {
      const response = await fetch("http://localhost:8080/tiles")
      const data = await response.json()
      setTiles(data)
    } catch (error) {
      console.error("Failed to fetch tiles:", error)
    } finally {
      setLoading(false)
    }
  }

  const fetchStats = async () => {
    try {
      const response = await fetch("http://localhost:8080/api/host/tile-stats")
      const data = await response.json()
      const statsMap = new Map(data.map((stat: TileStats) => [stat.tileId, stat]))
      setStats(statsMap)
    } catch (error) {
      console.error("Failed to fetch tile stats:", error)
    }
  }

  const toggleTileActive = async (tile: BingoTile) => {
    try {
      await fetch(`http://localhost:8080/api/host/tiles/${tile.id}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ active: !tile.marked }),
      })
      fetchTiles()
    } catch (error) {
      console.error("Failed to toggle tile:", error)
    }
  }

  const deleteTile = async (tileId: number) => {
    if (!confirm("Are you sure you want to delete this tile?")) return

    try {
      await fetch(`http://localhost:8080/api/host/tiles/${tileId}`, {
        method: "DELETE",
      })
      fetchTiles()
    } catch (error) {
      console.error("Failed to delete tile:", error)
    }
  }

  if (loading) {
    return (
      <Card className="flex items-center justify-center p-8">
        <p className="text-muted-foreground">Loading tiles...</p>
      </Card>
    )
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold text-foreground">Tile Management</h2>
        <Button onClick={() => setShowAddForm(true)} className="gap-2">
          <Plus className="h-4 w-4" />
          Add Tile
        </Button>
      </div>

      <Card>
        <ScrollArea className="h-[calc(100vh-16rem)]">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Tile Text</TableHead>
                <TableHead>Category</TableHead>
                <TableHead>Weight</TableHead>
                <TableHead>Win Rate</TableHead>
                <TableHead>Confirm Rate</TableHead>
                <TableHead>Active</TableHead>
                <TableHead>Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {tiles.map((tile) => {
                const tileStat = stats.get(tile.id)
                return (
                  <TableRow key={tile.id}>
                    <TableCell className="font-medium">{tile.text}</TableCell>
                    <TableCell>
                      <Badge variant="secondary">{tile.category || "General"}</Badge>
                    </TableCell>
                    <TableCell>{tile.weight || 1}</TableCell>
                    <TableCell>{tileStat ? `${(tileStat.winRate * 100).toFixed(1)}%` : "N/A"}</TableCell>
                    <TableCell>{tileStat ? `${(tileStat.confirmRate * 100).toFixed(1)}%` : "N/A"}</TableCell>
                    <TableCell>
                      <Switch checked={!tile.marked} onCheckedChange={() => toggleTileActive(tile)} />
                    </TableCell>
                    <TableCell>
                      <div className="flex gap-2">
                        <Button variant="ghost" size="icon" onClick={() => setEditingTile(tile)}>
                          <Edit2 className="h-4 w-4" />
                        </Button>
                        <Button variant="ghost" size="icon" onClick={() => deleteTile(tile.id)}>
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                )
              })}
            </TableBody>
          </Table>
        </ScrollArea>
      </Card>

      {showAddForm && <TileFormDialog onClose={() => setShowAddForm(false)} onSave={fetchTiles} />}
      {editingTile && <TileFormDialog tile={editingTile} onClose={() => setEditingTile(null)} onSave={fetchTiles} />}
    </div>
  )
}

function TileFormDialog({
  tile,
  onClose,
  onSave,
}: {
  tile?: BingoTile
  onClose: () => void
  onSave: () => void
}) {
  const [text, setText] = useState(tile?.text || "")
  const [category, setCategory] = useState(tile?.category || "General")
  const [weight, setWeight] = useState(tile?.weight?.toString() || "1")

  const handleSave = async () => {
    const data = {
      text,
      category,
      weight: Number.parseFloat(weight),
    }

    try {
      if (tile) {
        await fetch(`http://localhost:8080/api/host/tiles/${tile.id}`, {
          method: "PATCH",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(data),
        })
      } else {
        await fetch("http://localhost:8080/api/host/tiles", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(data),
        })
      }
      onSave()
      onClose()
    } catch (error) {
      console.error("Failed to save tile:", error)
    }
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-background/80">
      <Card className="w-full max-w-md p-6">
        <h3 className="mb-4 text-lg font-semibold">{tile ? "Edit Tile" : "Add New Tile"}</h3>

        <div className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="tile-text">Tile Text</Label>
            <Input id="tile-text" value={text} onChange={(e) => setText(e.target.value)} />
          </div>

          <div className="space-y-2">
            <Label htmlFor="tile-category">Category</Label>
            <Select value={category} onValueChange={setCategory}>
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="Linus">Linus</SelectItem>
                <SelectItem value="Luke">Luke</SelectItem>
                <SelectItem value="Dan">Dan</SelectItem>
                <SelectItem value="Technical">Technical</SelectItem>
                <SelectItem value="Sponsors">Sponsors</SelectItem>
                <SelectItem value="Tech News">Tech News</SelectItem>
                <SelectItem value="Show Flow">Show Flow</SelectItem>
                <SelectItem value="General">General</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <Label htmlFor="tile-weight">Weight (difficulty)</Label>
            <Input
              id="tile-weight"
              type="number"
              step="0.1"
              value={weight}
              onChange={(e) => setWeight(e.target.value)}
            />
            <p className="text-xs text-muted-foreground">Higher weight = more likely to appear</p>
          </div>
        </div>

        <div className="mt-6 flex justify-end gap-2">
          <Button variant="outline" onClick={onClose}>
            Cancel
          </Button>
          <Button onClick={handleSave}>Save</Button>
        </div>
      </Card>
    </div>
  )
}
