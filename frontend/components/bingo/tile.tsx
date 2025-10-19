import {BingoTile as IBingoTile} from "@/lib/bingoUtils";

interface BingoTileProps {
     tile: IBingoTile,
     toggle: () => void,
     highlighted?: boolean,
     showScore?: boolean,
     textSize?: string
 }

export function BingoTile(props: BingoTileProps) {

    const getTextSizeClass = (size?: string) => {
        switch (size) {
            case 'small': return 'text-[8px] sm:text-[10px] lg:text-xs'
            case 'large': return 'text-xs sm:text-sm lg:text-base'
            default: return 'text-[10px] sm:text-xs lg:text-sm'
        }
    }

    let baseClass = `relative w-full h-full aspect-square rounded-lg border-2 p-1.5 text-center font-medium transition-all hover:scale-105 active:scale-95 sm:p-2 ${getTextSizeClass(props.textSize)}`;

    if (props.tile.marked) {
        baseClass += "cursor-default border-primary bg-primary/20 text-primary-foreground"
    } else {
        baseClass += "cursor-pointer border-border bg-card text-card-foreground hover:border-primary/50 hover:bg-secondary"
    }

    // Add highlight border for confirmed tiles
    if (props.highlighted) {
        baseClass += " ring-2 ring-primary/20"
    }


    return (
         <button
             key={props.tile.id}
             onClick={props.toggle}
             className={baseClass}
         >
             <div className="flex h-full items-center justify-center">
                 <span className="text-balance leading-tight">{props.tile.title}</span>
             </div>
             {/* Score display in bottom right corner */}
             {props.showScore && (
                 <div className="absolute bottom-1 right-1 text-[8px] font-medium text-muted-foreground/60 sm:text-[9px]">
                     {props.tile.score || 5}
                 </div>
             )}
             {props.tile.marked && props.tile.title !== "FREE" && (
                 <div className="absolute inset-0 flex items-center justify-center">
                     <div className="h-1 w-full rotate-45 bg-primary/40"/>
                     <div className="absolute h-1 w-full -rotate-45 bg-primary/40"/>
                 </div>
             )}
         </button>
     )
}