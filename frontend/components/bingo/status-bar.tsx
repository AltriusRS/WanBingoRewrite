import {Button} from "@/components/ui/button";
import {RotateCcw, Settings} from "lucide-react";
import {ChatContextValue} from "@/components/chat/chat-context";
import {BingoTile, checkBingoWin, isValidWin} from "@/lib/bingoUtils";
import {AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle, AlertDialogTrigger} from "@/components/ui/alert-dialog";
import {GameplaySettingsPopover} from "@/components/gameplay-settings-popover";
import {useState} from "react";
import {useAuth} from "@/components/auth";


interface BingoStatusBarProps {
     ctx: ChatContextValue,
     hasWon: boolean,
     resetBoard: (e: unknown) => void
     regenerateBoard: () => void
     regenerationCount: number
     tiles: BingoTile[]
     confirmedTiles: Set<string>
     regenerationDiminisher?: number
     showMaxScore?: boolean
     showMultiplier?: boolean
     showRegenerations?: boolean
 }

export function BingoStatusBar({tiles, hasWon, resetBoard, regenerateBoard, regenerationCount, ctx, confirmedTiles, regenerationDiminisher = 1, showMaxScore = true, showMultiplier = true, showRegenerations = true}: BingoStatusBarProps) {
     const { user } = useAuth()
     const [showRegenerateDialog, setShowRegenerateDialog] = useState(false)

     const canRegenerate = regenerationCount < 3

     const penaltyPercent = (regenerationCount + 1) * 10

     // Calculate max possible score
     const maxScore = tiles.reduce((sum, tile) => sum + (tile.score || 5), 0)

     // Calculate multiplier as percentage
     const multiplierPercent = Math.round(regenerationDiminisher * 100)

    // Check if there's a potential win (5 in a row) but not all tiles are confirmed
    const hasBingo = checkBingoWin(tiles)
    const isValidBingo = hasBingo && isValidWin(tiles, confirmedTiles)
    const hasPotentialWin = hasBingo && !hasWon && !isValidBingo

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
                     {showMaxScore && (
                         <div>
                             Max Score: {maxScore}
                         </div>
                     )}
                     {showMultiplier && (
                         <div>
                             Multiplier: {multiplierPercent}%
                         </div>
                     )}
                     {showRegenerations && (
                         <div>
                             Regenerations: {regenerationCount}/3
                         </div>
                     )}
                     {hasWon && (
                         <>
                             <div className="h-4 w-px bg-border"/>
                             <div className="font-medium text-primary">🎉 BINGO!</div>
                         </>
                     )}
                     {hasPotentialWin && (
                         <>
                             <div className="h-4 w-px bg-border"/>
                             <div className="font-medium text-orange-500">⚠️ Potential Win (awaiting confirmation)</div>
                         </>
                     )}
                 </div>
             </div>

             <div className="flex items-center gap-2">
                 <GameplaySettingsPopover />
                 {showRegenerations && (
                    <Button
                        onClick={() => {
                            if (user) {
                                setShowRegenerateDialog(true);
                            } else {
                                regenerateBoard();
                            }
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
                )}

                {user && showRegenerations && (
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
                 )}
            </div>
        </div>
    )
}