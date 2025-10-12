"use client"

import React, {createContext, useContext, useEffect, useMemo, useState} from "react"
import {ChatMessage, EpisodeInfo, getNextWan, handleSocketProtocol, SSEMessage} from "@/lib/chatUtils";
import {fromZonedTime} from "date-fns-tz";

export interface ChatContextValue {
    messages: ChatMessage[]
    setMessages: React.Dispatch<React.SetStateAction<ChatMessage[]>>

    memberCount: number
    setMemberCount: React.Dispatch<React.SetStateAction<number>>

    memberList: string[]
    setMemberList: React.Dispatch<React.SetStateAction<string[]>>

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

    episode: EpisodeInfo
    setEpisode: React.Dispatch<React.SetStateAction<EpisodeInfo>>
}

const ChatContext = createContext<ChatContextValue | null>(null)


export function ChatProvider({children}: { children: React.ReactNode }) {
    const [messages, setMessages] = useState<ChatMessage[]>([])
    const [memberCount, setMemberCount] = useState<number>(0)
    const [memberList, setMemberList] = useState<string[]>([])
    const [memberListLoading, setMemberListLoading] = useState<boolean>(true)
    const [memberListError, setMemberListError] = useState<string | null>(null)
    const [liveTime, setLiveTime] = useState<string>("")
    const [text, setText] = useState<string>("")
    const [sending, setSending] = useState<boolean>(false)
    const [episode, setEpisode] = useState<EpisodeInfo>({
        title: "Unknown",
        date: "",
        thumbnail: null,
        isLive: false,
        startTime: getNextWan()
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
        const es = new EventSource("http://localhost:8080/api/chat/stream")

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
