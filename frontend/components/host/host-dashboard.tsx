"use client"
import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { TileConfirmationPanel } from "./tile-confirmation-panel"
import { TimerPanel } from "./timer-panel"
import { HostChatPanel } from "./host-chat-panel"
import { TestMessagePanel } from "./test-message-panel"
import { TileManagementPanel } from "./tile-management-panel"
import { SuggestionManagementPanel } from "./suggestion-management-panel"
import { HostProvider } from "./host-context"
import { LogOut, ArrowLeft } from "lucide-react"
import { useAuth } from "@/components/auth"
import Link from "next/link"

export function HostDashboard() {
  const { user, logout } = useAuth()
  const [showChat, setShowChat] = useState(true)

  const handleSignOut = async () => {
    await logout()
    window.location.href = "/"
  }

  if (!user) {
    return <div>Loading...</div>
  }

  return (
    <HostProvider>
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
            <div className="flex items-center gap-2">
              <Button variant="outline" size="sm" onClick={() => setShowChat(!showChat)} className="gap-2 bg-transparent">
                {showChat ? "Hide Chat" : "Show Chat"}
              </Button>
              <Button variant="outline" size="sm" asChild className="gap-2 bg-transparent">
                <Link href="/">
                  <ArrowLeft className="h-4 w-4" />
                  Back to Player View
                </Link>
              </Button>
              <Button variant="outline" size="sm" onClick={handleSignOut} className="gap-2 bg-transparent">
                <LogOut className="h-4 w-4" />
                Sign Out
              </Button>
            </div>
          </div>
        </header>

        {/* Main Content */}
        <div className="w-full flex flex-1 gap-4 overflow-hidden px-24 py-4">
          <Tabs defaultValue="tiles" className="flex flex-1 flex-col">
            <TabsList className="w-full justify-start">
              <TabsTrigger value="tiles">Tile Confirmation</TabsTrigger>
              <TabsTrigger value="timers">Timers</TabsTrigger>
              <TabsTrigger value="manage-tiles">Manage Tiles</TabsTrigger>
              <TabsTrigger value="suggestions">Suggestions</TabsTrigger>
              <TabsTrigger value="test">Test Messages</TabsTrigger>
            </TabsList>

             <TabsContent value="tiles" className="flex-1 overflow-hidden">
              {showChat ? (
                <div className="grid h-full gap-4 lg:grid-cols-[1fr_600px]">
                  <TileConfirmationPanel />
                  <HostChatPanel />
                </div>
              ) : (
                <TileConfirmationPanel showLateButton={true} />
              )}
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
    </HostProvider>
  )
}
