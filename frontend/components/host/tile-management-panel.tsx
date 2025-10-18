"use client"

import { useState, useEffect } from "react"
import { Card } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Badge } from "@/components/ui/badge"
import { Switch } from "@/components/ui/switch"
import { Check, ChevronsUpDown } from "lucide-react"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Edit2, Trash2, Plus } from "lucide-react"
import type { BingoTile } from "@/lib/bingoUtils"
import { getApiRoot } from "@/lib/auth"

interface TileStats {
  tileId: string
  winRate: number
  confirmRate: number
  timesConfirmed: number
  timesOnBoard: number
}

interface TileSettings {
  requiresTimer?: boolean
  requiresContext?: boolean
  description?: string
  confirmationRules?: string
  timer?: {
    duration: number
    name: string
    description: string
  }
  contextExample?: string
}

export function TileManagementPanel() {
  const [tiles, setTiles] = useState<BingoTile[]>([])
  const [stats, setStats] = useState<Map<string, TileStats>>(new Map())
  const [loading, setLoading] = useState(true)
  const [editingTile, setEditingTile] = useState<BingoTile | null>(null)
  const [showAddForm, setShowAddForm] = useState(false)
  const [sortBy, setSortBy] = useState<string>("title")
  const [sortDir, setSortDir] = useState<"asc" | "desc">("asc")

  useEffect(() => {
    fetchTiles()
    fetchStats()
  }, [])

  const fetchTiles = async () => {
    try {
      const response = await fetch(`${getApiRoot()}/api/host/tiles`, {
        credentials: "include",
      })
      if (!response.ok) {
        console.error("Failed to fetch tiles:", response.status)
        setTiles([])
        return
      }
      const data = await response.json()
      setTiles(Array.isArray(data) ? data : [])
    } catch (error) {
      console.error("Failed to fetch tiles:", error)
      setTiles([])
    } finally {
      setLoading(false)
    }
  }

  const fetchStats = async () => {
    try {
      const response = await fetch(`${getApiRoot()}/api/host/tile-stats`, {
        credentials: "include",
      })
      if (!response.ok) {
        console.error("Failed to fetch stats:", response.status)
        setStats(new Map())
        return
      }
      const data = await response.json()
      const statsMap = new Map(Array.isArray(data) ? data.map((stat: TileStats) => [stat.tileId, stat]) : [])
      setStats(statsMap)
    } catch (error) {
      console.error("Failed to fetch tile stats:", error)
      setStats(new Map())
    }
  }



  const deleteTile = async (tileId: string) => {
    if (!confirm("Are you sure you want to delete this tile?")) return

    try {
      await fetch(`${getApiRoot()}/api/host/tiles/${tileId}`, {
        method: "DELETE",
        credentials: "include",
      })
      fetchTiles()
    } catch (error) {
      console.error("Failed to delete tile:", error)
    }
  }

  const sortedTiles = [...tiles].sort((a, b) => {
    let aVal: any = a[sortBy as keyof BingoTile]
    let bVal: any = b[sortBy as keyof BingoTile]

    if (sortBy === "category") {
      aVal = aVal || "General"
      bVal = bVal || "General"
    }

    if (typeof aVal === "string") {
      aVal = aVal.toLowerCase()
      bVal = bVal.toLowerCase()
    }

    if (aVal < bVal) return sortDir === "asc" ? -1 : 1
    if (aVal > bVal) return sortDir === "asc" ? 1 : -1
    return 0
  })

  const handleSort = (column: string) => {
    if (sortBy === column) {
      setSortDir(sortDir === "asc" ? "desc" : "asc")
    } else {
      setSortBy(column)
      setSortDir("asc")
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
                 <TableHead className="cursor-pointer select-none" onClick={() => handleSort("title")}>
                   Tile Text {sortBy === "title" && (sortDir === "asc" ? "↑" : "↓")}
                 </TableHead>
                 <TableHead className="cursor-pointer select-none" onClick={() => handleSort("category")}>
                   Category {sortBy === "category" && (sortDir === "asc" ? "↑" : "↓")}
                 </TableHead>
                 <TableHead className="cursor-pointer select-none" onClick={() => handleSort("weight")}>
                   Weight {sortBy === "weight" && (sortDir === "asc" ? "↑" : "↓")}
                 </TableHead>
                 <TableHead className="cursor-pointer select-none" onClick={() => handleSort("score")}>
                   Score {sortBy === "score" && (sortDir === "asc" ? "↑" : "↓")}
                 </TableHead>
                 <TableHead>Win Rate</TableHead>
                 <TableHead>Confirm Rate</TableHead>
                 <TableHead>Actions</TableHead>
               </TableRow>
             </TableHeader>
             <TableBody>
               {sortedTiles.map((tile, index) => {
                 const tileStat = stats.get(tile.id)
                 return (
                   <TableRow key={tile.id} className={index % 2 === 1 ? "bg-muted/50" : ""}>
                     <TableCell className="font-medium">{tile.title}</TableCell>
                     <TableCell>
                       <Badge variant="secondary">{tile.category || "General"}</Badge>
                     </TableCell>
                     <TableCell>{tile.weight || 1}</TableCell>
                     <TableCell>{tile.score || 0}</TableCell>
                     <TableCell>{tileStat ? `${(tileStat.winRate * 100).toFixed(1)}%` : "N/A"}</TableCell>
                     <TableCell>{tileStat ? `${(tileStat.confirmRate * 100).toFixed(1)}%` : "N/A"}</TableCell>
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
  const settings: TileSettings = tile?.settings || {}
  const [text, setText] = useState(tile?.title || "")
  const [category, setCategory] = useState(tile?.category || "General")
  const [weight, setWeight] = useState(tile?.weight?.toString() || "1")
  const [score, setScore] = useState(tile?.score?.toString() || "0")
  const [description, setDescription] = useState(settings.description || "")
  const [confirmationRules, setConfirmationRules] = useState(settings.confirmationRules || "")
  const [requiresTimer, setRequiresTimer] = useState(settings.requiresTimer || false)
  const [requiresContext, setRequiresContext] = useState(settings.requiresContext || false)
  const [timerName, setTimerName] = useState(settings.timer?.name || "")
  const [timerDuration, setTimerDuration] = useState(settings.timer?.duration?.toString() || "60")
  const [timerDescription, setTimerDescription] = useState(settings.timer?.description || "")
  const [contextExample, setContextExample] = useState(settings.contextExample || "")

  const existingCategories = ["Linus", "Luke", "Dan", "Late", "Sponsors", "Topics", "Set/Production", "Events"]

  useEffect(() => {
    if (requiresTimer && !timerName) {
      setTimerName(text)
    }
  }, [text, requiresTimer])

  const handleSave = async () => {
    const newSettings: TileSettings = {
      description: description || undefined,
      confirmationRules: confirmationRules || undefined,
      requiresTimer,
      requiresContext,
      contextExample: requiresContext ? contextExample : undefined,
      timer: requiresTimer ? {
        name: timerName,
        duration: Number.parseInt(timerDuration),
        description: timerDescription,
      } : undefined,
    }

    const data = {
      text,
      category: category === "General" ? null : category,
      weight: Number.parseFloat(weight),
      score: Number.parseFloat(score),
      settings: newSettings,
    }

    try {
      if (tile) {
        await fetch(`${getApiRoot()}/api/host/tiles/${tile.id}`, {
          method: "PATCH",
          headers: { "Content-Type": "application/json" },
          credentials: "include",
          body: JSON.stringify(data),
        })
      } else {
        await fetch(`${getApiRoot()}/api/host/tiles`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          credentials: "include",
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
             <Label htmlFor="tile-text">Tile Title</Label>
             <Input id="tile-text" value={text} onChange={(e) => setText(e.target.value)} />
           </div>

           <div className="space-y-2">
             <Label htmlFor="tile-category">Category</Label>
             <Input
               id="tile-category"
               value={category}
               onChange={(e) => setCategory(e.target.value)}
               list="categories"
             />
             <datalist id="categories">
               {existingCategories.map((cat) => (
                 <option key={cat} value={cat} />
               ))}
             </datalist>
           </div>

           <div className="grid grid-cols-2 gap-4">
             <div className="space-y-2">
               <Label htmlFor="tile-weight">Weight (difficulty)</Label>
               <Input
                 id="tile-weight"
                 type="number"
                 step="0.1"
                 value={weight}
                 onChange={(e) => setWeight(e.target.value)}
               />
               <p className="text-xs text-muted-foreground">Higher weight = more likely</p>
             </div>
             <div className="space-y-2">
               <Label htmlFor="tile-score">Score</Label>
               <Input
                 id="tile-score"
                 type="number"
                 step="1"
                 value={score}
                 onChange={(e) => setScore(e.target.value)}
               />
             </div>
           </div>

           <div className="space-y-2">
             <Label htmlFor="tile-description">Description</Label>
             <Input id="tile-description" value={description} onChange={(e) => setDescription(e.target.value)} />
           </div>

           <div className="space-y-2">
             <Label htmlFor="confirmation-rules">Confirmation Rules</Label>
             <Input id="confirmation-rules" value={confirmationRules} onChange={(e) => setConfirmationRules(e.target.value)} />
           </div>

           <div className="flex gap-4">
             <div className="flex items-center space-x-2">
               <input type="checkbox" id="requires-timer" checked={requiresTimer} onChange={(e) => setRequiresTimer(e.target.checked)} />
               <Label htmlFor="requires-timer">Requires Timer</Label>
             </div>
             <div className="flex items-center space-x-2">
               <input type="checkbox" id="requires-context" checked={requiresContext} onChange={(e) => setRequiresContext(e.target.checked)} />
               <Label htmlFor="requires-context">Requires Context</Label>
             </div>
           </div>

           {requiresTimer && (
             <div className="space-y-4 p-4 border rounded">
               <h4 className="font-medium">Timer Configuration</h4>
               <div className="space-y-2">
                 <Label htmlFor="timer-name">Timer Name</Label>
                 <Input id="timer-name" value={timerName} onChange={(e) => setTimerName(e.target.value)} />
               </div>
               <div className="space-y-2">
                 <Label htmlFor="timer-duration">Duration (seconds)</Label>
                 <Input id="timer-duration" type="number" value={timerDuration} onChange={(e) => setTimerDuration(e.target.value)} />
               </div>
               <div className="space-y-2">
                 <Label htmlFor="timer-description">Timer Description</Label>
                 <Input id="timer-description" value={timerDescription} onChange={(e) => setTimerDescription(e.target.value)} />
               </div>
             </div>
           )}

           {requiresContext && (
             <div className="space-y-4 p-4 border rounded">
               <h4 className="font-medium">Context Configuration</h4>
               <div className="space-y-2">
                 <Label htmlFor="context-example">Context Example</Label>
                 <Input id="context-example" value={contextExample} onChange={(e) => setContextExample(e.target.value)} />
               </div>
             </div>
           )}
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
