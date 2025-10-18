"use client"

import {useState} from "react"
import {BingoBoard} from "@/components/bingo/bingo-board"
import {ChatPanel} from "@/components/chat/chat-panel"
import {Button} from "@/components/ui/button"
import {MessageSquare} from "lucide-react"
import {Header} from "@/components/header";

export default function Home() {
    const [isChatOpen, setIsChatOpen] = useState(true)

    const handleWin = () => {
        // Handle bingo win
    }

    return (
        <div className="flex min-h-screen min-w-screen flex-col bg-background pb-4">
            {/* Header */}
            <Header/>

            <div
                className="container flex flex-1 flex-col gap-4 overflow-hidden px-0 pt-2 mx-auto md:flex-row md:max-h-[calc(100vh-5rem)]">
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
                    <ChatPanel onClose={() => setIsChatOpen(false)} isMobile/>
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
        </div>
    )
}
