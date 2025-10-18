"use client"

import React, {createContext, useContext, useEffect, useMemo, useState} from "react"
import {ChatMessage, handleSocketProtocol, Show, SSEMessage, Player} from "@/lib/chatUtils";

export interface ChatContextValue {
    messages: ChatMessage[]
    setMessages: React.Dispatch<React.SetStateAction<ChatMessage[]>>>
    memberCount: number
    setMemberCount: React.Dispatch<React.SetStateAction<number>>
    memberList: Player[]
    setMemberList: React.Dispatch<React.SetStateAction<Player[]>>
    memberListLoading: boolean
    setMemberListLoading: React.Dispatch<React.SetStateAction<boolean>>
    memberListError: string | null
    setMemberListError: React.Dispatch<React.SetStateAction<string | null>>
    liveTime: string
    setLiveTime: React.Dispatch<React.SetStateAction<string>>
    text: string
    setText: React.Dispatch<React.SetStateAction<string>>
    sending: boolean
    setSending: React.Dispatch<React.SetStateAction<boolean>>
    episode: Show
    setEpisode: React.Dispatch<React.SetStateAction<Show>>
}

const ChatContext = createContext<ChatContextValue | null>(null)


export function ChatProvider({children}: { children: React.ReactNode }) {
    const [messages, setMessages] = useState<ChatMessage[]>([])
    const [memberCount, setMemberCount] = useState<number>(0)
    const [memberList, setMemberList] = useState<Player[]>([])
    const [memberListLoading, setMemberListLoading] = useState<boolean>(true)
    const [memberListError, setMemberListError] = useState<string | null>(null)
    const [liveTime, setLiveTime] = useState<string>("")
    const [text, setText] = useState<string>("")
    const [sending, setSending] = useState<boolean>(false)
    const [episode, setEpisode] = useState<Show>({
        "id": "Y2kz75uBC8",
        "youtube_id": "YVHXYqMPyzc",
        "scheduled_time": "2025-10-11T00:30:00Z",
        "actual_start_time": "2025-10-11T00:05:06Z",
        "thumbnail": "https://pbs.floatplane.com/stream_thumbnails/5c13f3c006f1be15e08e05c0/733054221374526_1760139634263.jpeg",
        "metadata": {
            "fp_vod": "w3A5fKcfTi",
            "title": "Piracy Is Dangerous And Harmful"
        },
        "created_at": "2025-10-10T23:46:40Z",
        "updated_at": "2025-10-16T19:12:58.692094Z",
    })

    const value = useMemo(() => ({
        messages, setMessages
        , memberCount, setMemberCount
        , memberList, setMemberList
        , memberListLoading, setMemberListLoading
        , memberListError, setMemberListError
        , liveTime, setLiveTime
        , text, setText
        , sending, setSending
        , episode, setEpisode
    }), [messages, memberCount, memberList, memberListLoading, memberListError, liveTime, text, sending, episode]) as ChatContextValue;

    useEffect(() => {
        const apiRoot = process.env.NEXT_PUBLIC_API_ROOT || "http://localhost:8000"
        const es = new EventSource(`${apiRoot}/chat/stream`)

        es.onmessage = (ev) => {
            try {
                if (!ev.data) return
                // Ignore keep-alives like "{}"
                if (ev.data.trim() === "{}") return
                const parsed = JSON.parse(ev.data) as SSEMessage;
                handleSocketProtocol(parsed, value)
            } catch (e) {
                // ignore malformed events
                console.warn("Failed to parse SSE chat event", e)
            }
        }

        es.onerror = (e) => {
            // Keep the connection; browser will attempt to reconnect
            console.error("SSE error", e)
        }

        return () => {
            es.close()
        }
    }, [])

    return <ChatContext.Provider value={value}>{children}</ChatContext.Provider>
}

export function useChat() {
    const ctx = useContext(ChatContext)
    if (!ctx) throw new Error("useChat must be used within a ChatProvider")
    return ctx
}
