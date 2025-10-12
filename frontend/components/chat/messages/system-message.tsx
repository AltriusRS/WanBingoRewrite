import {ChatMessage} from "@/lib/chatUtils";


interface StandardMessageProps {
    msg: ChatMessage
}

export function SystemMessage({msg}: StandardMessageProps) {
    return (
        <div className="rounded-lg border border-primary/30 bg-primary/10 p-3">
            <div className="flex items-center gap-2">
                <div className="h-1.5 w-1.5 rounded-full bg-primary"/>
                <span className="text-xs font-medium text-primary">SYSTEM</span>
                <span className="text-xs text-muted-foreground">{msg.timestamp}</span>
            </div>
            <p className="mt-1 text-sm text-foreground">{msg.message}</p>
        </div>
    )
}