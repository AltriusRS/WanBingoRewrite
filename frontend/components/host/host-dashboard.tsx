"use client"
import { Button } from "@/components/ui/button"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { TileConfirmationPanel } from "./tile-confirmation-panel"
import { TimerPanel } from "./timer-panel"
import { EpisodeInfoPanel } from "./episode-info-panel"
import { TestMessagePanel } from "./test-message-panel"
import { TileManagementPanel } from "./tile-management-panel"
import { SuggestionManagementPanel } from "./suggestion-management-panel"
import { LogOut } from "lucide-react"
import { useAuth } from "@/components/auth"

export function HostDashboard() {
  const { user, logout } = useAuth()

  const handleSignOut = async () => {
    await logout()
    window.location.href = "/"
  }

  if (!user) {
    return <div>Loading...</div>
  }

  return (
    <div className="flex min-h-screen flex-col bg-background">
      {/* Header */}
      <header className="shrink-0 border-b border-border bg-card">
        <div className="container mx-auto flex items-center justify-between px-4 py-4">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary">
              <span className="font-mono text-lg font-bold text-primary-foreground">H</span>
            </div>
            <div>
              <h1 className="text-xl font-semibold text-foreground">Host Dashboard</h1>
               <p className="text-sm text-muted-foreground">Logged in as {user.display_name}</p>
            </div>
          </div>
          <Button variant="outline" size="sm" onClick={handleSignOut} className="gap-2 bg-transparent">
            <LogOut className="h-4 w-4" />
            Sign Out
          </Button>
        </div>
      </header>

      {/* Main Content */}
      <div className="container mx-auto flex flex-1 gap-4 overflow-hidden p-4">
        <Tabs defaultValue="tiles" className="flex flex-1 flex-col">
          <TabsList className="w-full justify-start">
            <TabsTrigger value="tiles">Tile Confirmation</TabsTrigger>
            <TabsTrigger value="timers">Timers</TabsTrigger>
            <TabsTrigger value="manage-tiles">Manage Tiles</TabsTrigger>
            <TabsTrigger value="suggestions">Suggestions</TabsTrigger>
            <TabsTrigger value="test">Test Messages</TabsTrigger>
          </TabsList>

          <TabsContent value="tiles" className="flex-1 overflow-hidden">
            <div className="grid h-full gap-4 lg:grid-cols-[1fr_300px]">
              <TileConfirmationPanel />
              <EpisodeInfoPanel />
            </div>
          </TabsContent>

          <TabsContent value="timers" className="flex-1">
            <TimerPanel />
          </TabsContent>

          <TabsContent value="manage-tiles" className="flex-1 overflow-hidden">
            <TileManagementPanel />
          </TabsContent>

          <TabsContent value="suggestions" className="flex-1 overflow-hidden">
            <SuggestionManagementPanel />
          </TabsContent>

          <TabsContent value="test" className="flex-1">
            <TestMessagePanel />
          </TabsContent>
        </Tabs>
      </div>
    </div>
  )
}
