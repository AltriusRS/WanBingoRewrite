import {Tooltip, TooltipContent, TooltipTrigger} from "@/components/ui/tooltip";
import {Button} from "@/components/ui/button";
import {RotateCcw} from "lucide-react";
import {ChatContextValue} from "@/components/chat/chat-context";
import {BingoTile} from "@/lib/bingoUtils";


interface BingoStatusBarProps {
    ctx: ChatContextValue,
    hasWon: boolean,
    resetBoard: (e: unknown) => void
    tiles: BingoTile[]
}

export function BingoStatusBar({tiles, hasWon, resetBoard, ctx}: BingoStatusBarProps) {
    return (
        <div className="flex shrink-0 items-center justify-between w-full max-w-[min(90vw,90vh)] mb-2 sm:mb-3">
            <div className="flex items-center gap-3">
                <div className="flex gap-3 text-xs text-muted-foreground sm:text-sm">
                    <div>
            <span className="font-medium text-foreground">
              {tiles.filter((t) => t.marked).length}
            </span>{" "}
                        / 25
                    </div>
                    {hasWon && (
                        <>
                            <div className="h-4 w-px bg-border"/>
                            <div className="font-medium text-primary">🎉 BINGO!</div>
                        </>
                    )}
                </div>

                <Tooltip>
                    <TooltipTrigger asChild>
                        <Button
                            onClick={(e) => {
                                if (ctx.episode.isLive) return;
                                resetBoard(e);
                            }}
                            variant="outline"
                            size="sm"
                            className="gap-2 bg-transparent"
                        >
                            <RotateCcw className="h-4 w-4"/>
                            <span className="hidden sm:inline">New Board</span>
                            <span className="sm:hidden">New</span>
                        </Button>
                    </TooltipTrigger>
                    <TooltipContent>
                        {ctx.episode.isLive
                            ? "You can't reset the board while the show is live"
                            : "Reset the board to a new randomized set of tiles"}
                    </TooltipContent>
                </Tooltip>
            </div>
        </div>
    )
}