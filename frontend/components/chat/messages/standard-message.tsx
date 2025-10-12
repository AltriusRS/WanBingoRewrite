import {ChatMessage} from "@/lib/chatUtils";

interface StandardMessageProps {
    msg: ChatMessage
}


export function StandardMessage({msg}: StandardMessageProps) {
    return (
        <div className="space-y-1">
            <div className="flex items-baseline gap-2">
                <span className="text-sm font-medium text-primary">{msg.username}</span>
                <span className="text-xs text-muted-foreground">{msg.timestamp}</span>
            </div>
            <p className="text-sm text-foreground">{msg.message}</p>
        </div>
    )
}