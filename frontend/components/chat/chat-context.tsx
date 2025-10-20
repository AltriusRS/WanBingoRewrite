"use client"

import React, {createContext, useContext, useEffect, useMemo, useState} from "react"
import {ChatMessage, handleSocketProtocol, Player, Show, SSEMessage} from "@/lib/chatUtils";
import {useAuth} from "@/components/auth";

export interface ChatContextValue {
    messages: ChatMessage[]
    setMessages: React.Dispatch<React.SetStateAction<ChatMessage[]>>
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
    const {user} = useAuth()
    const [messages, setMessages] = useState<ChatMessage[]>([])
    const [memberCount, setMemberCount] = useState<number>(0)
    const [memberList, setMemberList] = useState<Player[]>([])
    const [memberListLoading, setMemberListLoading] = useState<boolean>(true)
    const [memberListError, setMemberListError] = useState<string | null>(null)
    const [liveTime, setLiveTime] = useState<string>("")
    const [text, setText] = useState<string>("")
    const [sending, setSending] = useState<boolean>(false)
    const [episode, setEpisode] = useState<Show | null>(null)

    const value = useMemo(() => ({
        messages, setMessages,
        memberCount, setMemberCount,
        memberList, setMemberList,
        memberListLoading, setMemberListLoading,
        memberListError, setMemberListError,
        liveTime, setLiveTime,
        text, setText,
        sending, setSending,
        episode, setEpisode
    }), [messages, memberCount, memberList, memberListLoading, memberListError, liveTime, text, sending, episode]) as ChatContextValue;

    useEffect(() => {
        const fetchEpisode = async () => {
            try {
                const apiRoot = process.env.NEXT_PUBLIC_API_ROOT || "http://localhost:8000"
                const response = await fetch(`${apiRoot}/shows/latest`)
                if (response.ok) {
                    const data = await response.json()
                    setEpisode(data)
                }
            } catch (error) {
                console.error("Failed to fetch episode:", error)
            }
        }

        fetchEpisode()

        const apiRoot = process.env.NEXT_PUBLIC_API_ROOT || "http://localhost:8000"
        const es = new EventSource(`${apiRoot}/chat/stream`, {
            withCredentials: true,
        })

        es.onmessage = (ev) => {
            try {
                if (!ev.data) return
                // Ignore keep-alives like "{}"
                if (ev.data.trim() === "{}") return
                const parsed = JSON.parse(ev.data) as SSEMessage;
                handleSocketProtocol(parsed, value, user)
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
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [user])

    return <ChatContext.Provider value={value}>{children}</ChatContext.Provider>
}

export function useChat() {
    const ctx = useContext(ChatContext)
    if (!ctx) throw new Error("useChat must be used within a ChatProvider")
    return ctx
}
