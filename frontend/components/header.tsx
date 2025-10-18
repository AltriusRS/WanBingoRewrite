import Link from "next/link";
import {Button} from "@/components/ui/button";
import {Info, Lightbulb, Menu, Trophy, User} from "lucide-react";
import {SuggestTileModal} from "@/components/suggest-tile-modal";
import {useState} from "react";

export function Header() {
    const [isSuggestModalOpen, setIsSuggestModalOpen] = useState(false)


    const handleSuggestTile = (data: { name: string; tileName: string; reason: string }) => {
        console.log("Suggested tile", data)
    }


    return (
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
                    <Link href="/leaderboard">
                        <Button variant="ghost" size="sm" className="gap-2 bg-transparent">
                            <Trophy className="h-4 w-4"/>
                            <span className="hidden sm:inline">Leaderboard</span>
                        </Button>
                    </Link>
                    <Link href="/about">
                        <Button variant="ghost" size="sm" className="gap-2 bg-transparent">
                            <Info className="h-4 w-4"/>
                            <span className="hidden sm:inline">About</span>
                        </Button>
                    </Link>
                    <Link href="/account">
                        <Button variant="ghost" size="sm" className="gap-2 bg-transparent">
                            <User className="h-4 w-4"/>
                            <span className="hidden sm:inline">Account</span>
                        </Button>
                    </Link>
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

            <SuggestTileModal open={isSuggestModalOpen} onOpenChange={setIsSuggestModalOpen}
                              onSubmit={handleSuggestTile}/>
        </header>
    )
}