import {ChatContextValue} from "@/components/chat/chat-context";
import {fromZonedTime} from "date-fns-tz";
import {ReactNode} from "react";


/**
 * Post a message to the chat server.
 * @param message {string} - The message content to send
 */
export async function postMessage(message: string): Promise<boolean> {
    try {
        const apiRoot = process.env.NEXT_PUBLIC_API_ROOT || "http://localhost:8000"
        await fetch(`${apiRoot}/chat`, {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            credentials: "include", // Include cookies for authentication
            body: JSON.stringify({contents:message}),
        })
        return true;
    } catch (err) {
        console.warn("Failed to send chat message", err)
        return false;
    }
}

export interface ChatPanelProps {
    onClose: () => void
    isMobile?: boolean
}

export interface Player {
    id: string;
    did: string;
    display_name: string;
    avatar?: string;
    settings?: { [key: string]: any };
    score: number;
    permissions: number;
    created_at: string;
    updated_at: string;
    deleted_at?: string;
}

export interface Show {
    id: string;
    state: string;
    youtube_id?: string;
    scheduled_time?: string;
    actual_start_time?: string;
    thumbnail?: string;
    metadata?: { [key: string]: any };
    created_at: string;
    updated_at: string;
    deleted_at?: string;
}

export interface Tile {
    id: string;
    title: string;
    category?: string;
    last_drawn?: string;
    created_by?: string;
    weight: number;
    score: number;
    settings: { [key: string]: any };
    created_at: string;
    updated_at: string;
    deleted_at?: string;
}

export interface ShowTile {
    show_id: string;
    tile_id: string;
    weight: number;
    score: number;
    created_at: string;
    updated_at: string;
    deleted_at?: string;
}

export interface Board {
    id: string;
    player_id: string;
    show_id: string;
    tiles: string[];
    winner: boolean;
    total_score: number;
    potential_score: number;
    regeneration_diminisher: number;
    created_at: string;
    updated_at: string;
    deleted_at?: string;
}

export interface TileConfirmation {
    id: string;
    show_id: string;
    tile_id: string;
    confirmed_by?: string;
    context?: string;
    confirmation_time: string;
    created_at: string;
    updated_at: string;
    deleted_at?: string;
}


export interface ChatMessage {
    id: string;
    show_id: string;
    player_id: string;
    contents: string;
    system: boolean;
    replying?: string;
    created_at: string;
    updated_at: string;
    deleted_at?: string;
    player?: Partial<Player>
    html: ReactNode
}

export interface MessageRequest {
    contents: string;
}


export function updateLiveTime(episodeInfo: Show, chatContext: ChatContextValue) {
    const now = new Date()

    const diff = Math.abs(now.getTime() - new Date(episodeInfo.actual_start_time ?? episodeInfo.scheduled_time ?? episodeInfo.created_at).getTime())
    const minutes = Math.floor(diff / 60000)
    const hours = Math.floor(minutes / 60)
    const days = Math.floor(hours / 24)
    const mins = minutes % 60
    const hrs = hours % 24

    if (episodeInfo.actual_start_time) {
        chatContext.setLiveTime(hours > 0 ? `${hours}h ${mins}m` : `${mins}m`)
    } else {
        if (episodeInfo.id) {
            chatContext.setLiveTime("imminently")
        } else if (new Date(episodeInfo.actual_start_time ?? episodeInfo.scheduled_time ?? episodeInfo.created_at) > now) {
            chatContext.setLiveTime(days > 0 ? `in ${days}d ${hrs}h ${mins}m` : hours > 0 ? `in ${hours}h ${mins}m` : `in ${mins}m`)
        }
    }
}

/**
 * Build a submit handler for the chat form.
 * @param chatContext {ChatContextValue} - The chat context to use
 * @returns {(e: React.FormEvent) => Promise<void>} - The submit handler function
 */
export function buildSubmitHandler(chatContext: ChatContextValue): (e: React.FormEvent) => Promise<void> {
    return async (e: React.FormEvent) => {
        e.preventDefault()
        const message = chatContext.text.trim()
        if (!message) return
        chatContext.setSending(true);

        const didSend = await postMessage(message);

        if (didSend) chatContext.setText("");
        chatContext.setSending(false)
    }
}

export interface SSEMessage {
    opcode: string
    data: any
}


export function handleSocketProtocol(protoMessage: SSEMessage, ctx: ChatContextValue) {
    switch (protoMessage.opcode) {
        case "chat.members.count":
            ctx.setMemberCount(() => protoMessage.data.count);
            break;

        case "chat.players":
            ctx.setMemberList(() => protoMessage.data.players as Player[]);
            ctx.setMemberListLoading(false);
            break;

        case 'chat.message':
            return handleChatMessage(protoMessage.data as ChatMessage, ctx);

        case 'whenplane.aggregate':
            const aggregate = protoMessage.data as Show;
            ctx.setEpisode((_) => aggregate)
            break;
        default:
            // Unknown opcode - silently ignore
            break;
    }
}


async function handleChatMessage(msg: ChatMessage, ctx: ChatContextValue) {
    // Truncate the list to 100 messages

    if (ctx.messages.filter((s) => s.id === msg.id).length > 0) {
        // Duplicate message - skip
        return;
    }
    ctx.setMessages((prev) => {
        if (prev.length > 99) {
            return [...prev.slice(1, 100), msg]
        } else {
            return [...prev, msg]
        }
    });
}


const timeZone = "America/Vancouver";

export function getNextWan(): Date {
    const now = new Date();

    // Step 1: get Vancouver local time equivalent of now
    const vancouverNow = new Date(now.toLocaleString("en-US", {timeZone}));

    // Step 2: figure out days to next Friday (5)
    const currentDay = vancouverNow.getDay(); // 0 = Sunday, 5 = Friday
    let daysUntilFriday = (5 - currentDay + 7) % 7;

    // If today is Friday but past 4:30pm, push to next Friday
    const isFriday = daysUntilFriday === 0;
    const targetTime = new Date(
        vancouverNow.getFullYear(),
        vancouverNow.getMonth(),
        vancouverNow.getDate() + daysUntilFriday,
        16,
        30,
        0,
        0
    );

    if (isFriday && vancouverNow > targetTime) {
        // after 4:30pm today, jump a week ahead
        targetTime.setDate(targetTime.getDate() + 7);
    }

    // Step 3: convert to UTC so it’s the same everywhere
    const utcTarget = fromZonedTime(targetTime, timeZone);

    return utcTarget;
}
