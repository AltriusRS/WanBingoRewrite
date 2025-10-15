import {ChatContextValue} from "@/components/chat/chat-context";
import {fromZonedTime} from "date-fns-tz";


/**
 * Post a message to the chat server.
 * @param message {string} - The message content to send
 */
export async function postMessage(message: string): Promise<boolean> {
    try {
        await fetch("https://api.bingo.local/chat/message", {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify({message, username: "tester"}),
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

export interface EpisodeInfo {
    id?: string
    title: string
    date: string
    thumbnail: string | null
    isLive: boolean
    startTime: Date
}


export type ChatMessage = {
    id: number
    type: "user" | "system"
    username?: string
    message: string
    timestamp: string
}


export function updateLiveTime(episodeInfo: EpisodeInfo, chatContext: ChatContextValue) {
    const now = new Date()
    const diff = Math.abs(now.getTime() - episodeInfo.startTime.getTime())
    const minutes = Math.floor(diff / 60000)
    const hours = Math.floor(minutes / 60)
    const days = Math.floor(hours / 24)
    const mins = minutes % 60
    const hrs = hours % 24

    if (episodeInfo.isLive) {
        chatContext.setLiveTime(hours > 0 ? `${hours}h ${mins}m` : `${mins}m`)
    } else {
        if (episodeInfo.id) {
            chatContext.setLiveTime("imminently")
        } else if (episodeInfo.startTime > now) {
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
    console.log("Received SSE message:", protoMessage);
    switch (protoMessage.opcode) {
        case "chat.members.count":
            console.log("Received member count update:", protoMessage.data.count);
            ctx.setMemberCount(() => protoMessage.data.count);
            break;

        case 'chat.message':
            return handleChatMessage(protoMessage.data as ChatMessage, ctx);

        case 'whenplane.aggregate':
            console.log("Received aggregate update:", protoMessage.data);
            const aggregate = protoMessage.data as any;

            ctx.setEpisode(() => {
                let videoId = aggregate.youtube.videoId || undefined;

                let title = (aggregate.floatplane.title || "Unknown").split(" - ")[0] || "Unknown";

                let date = (new Date((aggregate.floatplane.title || "Unknown").split(" - ")[1]) || new Date()).toISOString();

                let thumbnail = videoId ? "https://i.ytimg.com/vi/" + videoId + "/maxresdefault_live.jpg" : (aggregate.floatplane.thumbnail || null)
                return {
                    id: videoId,
                    title,
                    date,
                    thumbnail,
                    isLive: aggregate.youtube.isLive || false,
                    startTime: getNextWan()
                } as EpisodeInfo;
            })

            break;
        default:
            console.warn("Unknown SSE opcode", protoMessage.opcode)
    }
}


async function handleChatMessage(msg: ChatMessage, ctx: ChatContextValue) {
    console.log("Received chat message:", msg);
    // Truncate the list to 100 messages
    ctx.setMessages((prev) => {
        if (prev.length > 99) {
            return [...prev.slice(0, 99), msg]
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
