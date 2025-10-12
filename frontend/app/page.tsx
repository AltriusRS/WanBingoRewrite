"use client"

import {useState} from "react"
import {BingoBoard} from "@/components/bingo-board"
import {ChatPanel} from "@/components/chat/chat-panel"
import {SuggestTileModal} from "@/components/suggest-tile-modal"
import {Button} from "@/components/ui/button"
import {Lightbulb, Menu, MessageSquare} from "lucide-react"

export default function Home() {
    const [isChatOpen, setIsChatOpen] = useState(true)
    const [isSuggestModalOpen, setIsSuggestModalOpen] = useState(false)

    const handleWin = () => {
        // In a future iteration, this could POST a system message to the server.
        console.log("Bingo!")
    }

    const handleSuggestTile = (data: { name: string; tileName: string; reason: string }) => {
        // In a future iteration, this could POST a system message to the server.
        console.log("Suggested tile", data)
    }

    return (
        <div className="flex min-h-screen flex-col bg-background">
            {/* Header */}
            <header className="shrink-0 border-b border-border bg-card">
                <div className="container mx-auto flex items-center justify-between px-4 py-4">
                    <div className="flex items-center gap-3">
                        <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary">
                            <span className="font-mono text-lg font-bold text-primary-foreground">W</span>
                        </div>
                        <div>
                            <h1 className="text-xl font-semibold text-foreground">WAN Show Bingo</h1>
                            <p className="text-sm text-muted-foreground">Not affiliated with Linus Media Group</p>
                        </div>
                    </div>
                    <div className="flex items-center gap-2">
                        <Button
                            variant="outline"
                            size="sm"
                            className="gap-2 bg-transparent"
                            onClick={() => setIsSuggestModalOpen(true)}
                        >
                            <Lightbulb className="h-4 w-4"/>
                            <span className="hidden sm:inline">Suggest Tiles</span>
                        </Button>
                        <Button variant="ghost" size="sm" onClick={() => setIsChatOpen(!isChatOpen)}
                                className="gap-2 md:hidden">
                            <Menu className="h-4 w-4"/>
                        </Button>
                    </div>
                </div>
            </header>

            <div
                className="container mx-auto flex flex-1 flex-col gap-4 overflow-hidden p-4 md:flex-row md:max-h-[calc(100vh-5rem)]">
                {/* Bingo Board - Main Content */}
                <div className="flex flex-1 flex-col overflow-hidden md:min-w-0">
                    <BingoBoard onWin={handleWin}/>
                </div>

                {isChatOpen && (
                    <div className="hidden w-96 shrink-0 md:block md:max-h-full">
                        <ChatPanel onClose={() => setIsChatOpen(false)}/>
                    </div>
                )}
            </div>

            {isChatOpen && (
                <div className="fixed inset-x-0 bottom-0 z-50 h-[60vh] md:hidden">
                    <ChatPanel onClose={() => setIsChatOpen(false)}
                               isMobile/>
                </div>
            )}

            {!isChatOpen && (
                <Button
                    onClick={() => setIsChatOpen(true)}
                    className="fixed bottom-4 right-4 h-12 w-12 rounded-full shadow-lg md:hidden"
                    size="icon"
                >
                    <MessageSquare className="h-5 w-5"/>
                </Button>
            )}

            <SuggestTileModal open={isSuggestModalOpen} onOpenChange={setIsSuggestModalOpen}
                              onSubmit={handleSuggestTile}/>
        </div>
    )
}
