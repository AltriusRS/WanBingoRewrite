import {Button} from "@/components/ui/button";
import {RotateCcw} from "lucide-react";
import {ChatContextValue} from "@/components/chat/chat-context";
import {BingoTile} from "@/lib/bingoUtils";
import {AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle, AlertDialogTrigger} from "@/components/ui/alert-dialog";
import {useState} from "react";


interface BingoStatusBarProps {
    ctx: ChatContextValue,
    hasWon: boolean,
    resetBoard: (e: unknown) => void
    regenerateBoard: () => void
    regenerationCount: number
    tiles: BingoTile[]
}

export function BingoStatusBar({tiles, hasWon, resetBoard, regenerateBoard, regenerationCount, ctx}: BingoStatusBarProps) {
    const [showRegenerateDialog, setShowRegenerateDialog] = useState(false)

    const canRegenerate = regenerationCount < 3

    const penaltyPercent = (regenerationCount + 1) * 10

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
                    <div>
                        Regenerations: {regenerationCount}/3
                    </div>
                    {hasWon && (
                        <>
                            <div className="h-4 w-px bg-border"/>
                            <div className="font-medium text-primary">🎉 BINGO!</div>
                        </>
                    )}
                </div>

                <Button
                    onClick={() => {
                        setShowRegenerateDialog(true);
                    }}
                    variant="outline"
                    size="sm"
                    className="gap-2 bg-transparent"
                    disabled={!canRegenerate}
                >
                    <RotateCcw className="h-4 w-4"/>
                    <span className="hidden sm:inline">New Board</span>
                    <span className="sm:hidden">New</span>
                </Button>

                <AlertDialog open={showRegenerateDialog} onOpenChange={setShowRegenerateDialog}>
                    <AlertDialogContent>
                        <AlertDialogHeader>
                            <AlertDialogTitle>Regenerate Board?</AlertDialogTitle>
                            <AlertDialogDescription>
                                Regenerating your board will incur a -{penaltyPercent}% penalty on your final score. Are you sure you want to proceed?
                            </AlertDialogDescription>
                        </AlertDialogHeader>
                        <AlertDialogFooter>
                            <AlertDialogCancel>Cancel</AlertDialogCancel>
                            <AlertDialogAction onClick={() => {
                                regenerateBoard();
                                setShowRegenerateDialog(false);
                            }}>
                                Regenerate
                            </AlertDialogAction>
                        </AlertDialogFooter>
                    </AlertDialogContent>
                </AlertDialog>
            </div>
        </div>
    )
}