"use client"
import {useEffect, useState} from "react"
import {Button} from "@/components/ui/button"
import {Tabs, TabsContent, TabsList, TabsTrigger} from "@/components/ui/tabs"
import {TileConfirmationPanel} from "./tile-confirmation-panel"
import {TimerPanel} from "./timer-panel"
import {HostChatPanel} from "./host-chat-panel"
import {TestMessagePanel} from "./test-message-panel"
import {TileManagementPanel} from "./tile-management-panel"
import {SuggestionManagementPanel} from "./suggestion-management-panel"
import {HostProvider} from "./host-context"
import {ArrowLeft, ExternalLink, LogOut, Clock, Eye, EyeOff} from "lucide-react"
import {useAuth} from "@/components/auth"
import Link from "next/link"
import {Card} from "@/components/ui/card"
import {getApiRoot} from "@/lib/auth"
import {toast} from "sonner"

export function HostDashboard() {
    const {user, logout} = useAuth()
    const [showChat, setShowChat] = useState(true)
    const [showLateModal, setShowLateModal] = useState(false)
    const [lateReason, setLateReason] = useState("")
    const [confirmingLate, setConfirmingLate] = useState(false)
    const [hideTabs, setHideTabs] = useState(false)

    useEffect(() => {
        const isMobile = window.innerWidth < 768
        setShowChat(!isMobile)
    }, [])

    useEffect(() => {
        const handleKeyDown = (e: KeyboardEvent) => {
            if (e.ctrlKey && e.key === '\\') {
                setHideTabs(prev => !prev)
            }
        }
        window.addEventListener('keydown', handleKeyDown)
        return () => window.removeEventListener('keydown', handleKeyDown)
    }, [])

    const handleSignOut = async () => {
        await logout()
        window.location.href = "/"
    }

    const handleConfirmLate = async () => {
        if (!lateReason.trim()) {
            toast.error("Please provide a reason for the delay")
            return
        }

        setConfirmingLate(true)
        try {
            const response = await fetch(`${getApiRoot()}/tiles/confirmations`, {
                method: "POST",
                headers: {"Content-Type": "application/json"},
                credentials: "include",
                body: JSON.stringify({
                    tile_id: "late", // Special ID for late show
                    context: lateReason,
                }),
            })

            if (!response.ok) {
                const errorText = await response.text()
                throw new Error(`Failed to confirm late show: ${response.status} ${errorText}`)
            }

            toast.success("Show delay has been confirmed and announced to viewers.")
            setShowLateModal(false)
            setLateReason("")
        } catch (error) {
            console.error("Failed to confirm late show:", error)
            toast.error(error instanceof Error ? error.message : "Failed to confirm late show")
        } finally {
            setConfirmingLate(false)
        }
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
                            <Button
                                variant="destructive"
                                size="sm"
                                onClick={() => setShowLateModal(true)}
                                className="gap-2"
                            >
                                <Clock className="h-4 w-4"/>
                                Show Is Late
                            </Button>
                            <Button variant="outline" size="sm" onClick={() => setShowChat(!showChat)}
                                    className="gap-2 bg-transparent">
                                {showChat ? "Hide Chat" : "Show Chat"}
                            </Button>
                            <Button variant="outline" size="sm" onClick={() => {
                                window.open('/chat', '_blank', 'width=500,height=800');
                                setShowChat(false);
                            }} className="gap-2 bg-transparent">
                                <ExternalLink className="h-4 w-4"/>
                                Pop Out Chat
                            </Button>
                            <Button variant="outline" size="sm" asChild className="gap-2 bg-transparent">
                                <Link href="/">
                                    <ArrowLeft className="h-4 w-4"/>
                                    Back to Player View
                                </Link>
                            </Button>
                             <Button variant="outline" size="sm" onClick={() => setHideTabs(prev => !prev)}
                                     className="gap-2 bg-transparent">
                                 {hideTabs ? <Eye className="h-4 w-4"/> : <EyeOff className="h-4 w-4"/>}
                                 {hideTabs ? "Show Tabs" : "Hide Tabs"}
                             </Button>
                             <Button variant="outline" size="sm" onClick={handleSignOut}
                                     className="gap-2 bg-transparent">
                                 <LogOut className="h-4 w-4"/>
                                 Sign Out
                             </Button>
                        </div>
                    </div>
                </header>

                {/* Main Content */}
                <div className="w-full flex flex-1 gap-4 overflow-hidden px-24 py-4">
                     <Tabs defaultValue="tiles" className="flex flex-1 flex-col">
                         {!hideTabs && (
                             <TabsList className="w-full justify-start">
                                 <TabsTrigger value="tiles">Tile Confirmation</TabsTrigger>
                                 <TabsTrigger value="timers">Timers</TabsTrigger>
                                 <TabsTrigger value="manage-tiles">Manage Tiles</TabsTrigger>
                                 <TabsTrigger value="suggestions">Suggestions</TabsTrigger>
                                 <TabsTrigger value="test">Test Messages</TabsTrigger>
                             </TabsList>
                         )}

                        <TabsContent value="tiles" className="flex-1 overflow-hidden">
                            {showChat ? (
                                <div className="grid h-full gap-4 lg:grid-cols-[2fr_1fr]">
                                    <TileConfirmationPanel columns={2}/>
                                    <HostChatPanel/>
                                </div>
                             ) : (
                                 <TileConfirmationPanel columns={3}/>
                             )}
                        </TabsContent>

                        <TabsContent value="timers" className="flex-1">
                            <TimerPanel/>
                        </TabsContent>

                        <TabsContent value="manage-tiles" className="flex-1 overflow-hidden">
                            <TileManagementPanel/>
                        </TabsContent>

                        <TabsContent value="suggestions" className="flex-1 overflow-hidden">
                            <SuggestionManagementPanel/>
                        </TabsContent>

                        <TabsContent value="test" className="flex-1">
                            <TestMessagePanel/>
                        </TabsContent>
                    </Tabs>
                 </div>

                {/* Late Show Modal */}
                {showLateModal && (
                    <div className="fixed inset-0 z-50 flex items-center justify-center bg-background/80">
                        <Card className="w-full max-w-md p-6">
                            <div className="space-y-4">
                                <div className="flex items-center justify-between">
                                    <h3 className="text-lg font-semibold">Confirm Show Delay</h3>
                                    <Button variant="ghost" size="icon" onClick={() => setShowLateModal(false)}>
                                        ×
                                    </Button>
                                </div>

                                <div className="space-y-3">
                                    <div className="rounded-lg border border-border bg-muted p-4">
                                        <h4 className="font-medium mb-2">The WAN Show is running late</h4>
                                        <p className="text-sm text-muted-foreground">
                                            This will announce to all viewers that the show has been delayed.
                                        </p>
                                    </div>

                                    <div className="space-y-3">
                                        <div>
                                            <label className="text-sm font-medium">Reason for delay (optional)</label>
                                            <textarea
                                                className="w-full mt-1 p-2 border border-border rounded-md bg-background text-sm"
                                                placeholder="Brief explanation of the delay..."
                                                rows={3}
                                                value={lateReason}
                                                onChange={(e) => setLateReason(e.target.value)}
                                            />
                                        </div>
                                    </div>

                                    <div className="flex gap-2">
                                        <Button
                                            variant="outline"
                                            onClick={() => setShowLateModal(false)}
                                            className="flex-1"
                                            disabled={confirmingLate}
                                        >
                                            Cancel
                                        </Button>
                                        <Button
                                            variant="destructive"
                                            onClick={handleConfirmLate}
                                            className="flex-1"
                                            disabled={confirmingLate}
                                        >
                                            {confirmingLate ? "Confirming..." : "Confirm Delay"}
                                        </Button>
                                    </div>
                                </div>
                            </div>
                        </Card>
                    </div>
                )}
             </div>
         </HostProvider>
     )
}
