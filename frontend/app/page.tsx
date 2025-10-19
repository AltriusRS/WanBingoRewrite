"use client"

import {useState} from "react"
import confetti from "canvas-confetti"
import {Header} from "@/components/header"
import {BingoBoard} from "@/components/bingo/bingo-board"
import {ChatPanel} from "@/components/chat/chat-panel"
import {Button} from "@/components/ui/button"
import {MessageSquare} from "lucide-react"
import {useAuth} from "@/components/auth"

export default function Home() {
     const [isChatOpen, setIsChatOpen] = useState(true)
     const { user } = useAuth()

     const handleWin = () => {
         // Handle bingo win with confetti celebration (if enabled)
         const confettiEnabled = user?.settings?.gameplay?.confetti !== false
         if (confettiEnabled) {
             // Burst from left side - firing downwards into center
             confetti({
                 particleCount: 150,
                 angle: 90,
                 spread: 45,
                 origin: { x: 0.1, y: 0.2 },
                 decay: 0.92,
                 colors: ['#ff0000', '#00ff00', '#0000ff', '#ffff00', '#ff00ff', '#00ffff']
             })

             // Burst from right side - firing downwards into center
             confetti({
                 particleCount: 150,
                 angle: 90,
                 spread: 45,
                 origin: { x: 0.9, y: 0.2 },
                 decay: 0.92,
                 colors: ['#ff0000', '#00ff00', '#0000ff', '#ffff00', '#ff00ff', '#00ffff']
             })

             // Additional center burst for more celebration
             setTimeout(() => {
                 confetti({
                     particleCount: 100,
                     spread: 90,
                     origin: { y: 0.4 },
                     decay: 0.92,
                     colors: ['#ff0000', '#00ff00', '#0000ff', '#ffff00', '#ff00ff', '#00ffff']
                 })
             }, 200)
         }
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
