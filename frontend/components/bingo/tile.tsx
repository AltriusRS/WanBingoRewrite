import {BingoTile as IBingoTile} from "@/lib/bingoUtils";

interface BingoTileProps {
    tile: IBingoTile,
    toggle: () => void
}

export function BingoTile(props: BingoTileProps) {

    let baseClass = "relative w-full h-full aspect-square rounded-lg border-2 p-1.5 text-center text-[10px] font-medium transition-all\n                hover:scale-105 active:scale-95 sm:p-2 sm:text-xs lg:text-sm";

    if (props.tile.title === "FREE") {
        baseClass += "border-primary bg-primary/20 text-foreground shadow-lg shadow-primary/20"
    } else {
        baseClass += "border-border bg-card text-card-foreground hover:border-primary/50 hover:bg-secondary"
    }

    if (props.tile.marked) {
        baseClass += "cursor-default border-primary bg-primary text-primary-foreground"
    } else {
        baseClass += "cursor-pointer"
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
            {props.tile.marked && props.tile.title !== "FREE" && (
                <div className="absolute inset-0 flex items-center justify-center">
                    <div className="h-1 w-full rotate-45 bg-primary/40"/>
                    <div className="absolute h-1 w-full -rotate-45 bg-primary/40"/>
                </div>
            )}
        </button>
    )
}