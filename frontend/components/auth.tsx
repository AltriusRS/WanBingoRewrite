"use client"

import {createContext, ReactNode, useContext, useEffect, useState} from "react"

// Discord user type
export interface DiscordUser {
    id: string
    username: string
    discriminator: string
    email: string
    avatar: string
    verified: boolean
}

// Auth context type
interface AuthContextType {
    user: DiscordUser | null
    loading: boolean
    error: string | null
    login: () => void
    logout: () => Promise<void>
    refetch: () => Promise<void>
}

const AuthContext = createContext<AuthContextType | null>(null)

// Get API root from environment
function getApiRoot(): string {
    return process.env.NEXT_PUBLIC_API_ROOT || "http://localhost:8000"
}

// Auth provider component
export function AuthProvider({children}: { children: ReactNode }) {
    const [user, setUser] = useState<DiscordUser | null>(null)
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState<string | null>(null)

    const fetchUser = async () => {
        try {
            setLoading(true)
            setError(null)

            const response = await fetch(`${getApiRoot()}/users/me`, {
                credentials: "include", // Include cookies
            })

            if (response.ok) {
                const data = await response.json()
                setUser(data.user)
            } else if (response.status === 401) {
                // Not authenticated
                setUser(null)
            } else {
                console.error(await response.json())
                throw new Error("Failed to fetch user")
            }
        } catch (err) {
            console.error("Auth error:", err)
            setError(err instanceof Error ? err.message : "Unknown error")
            setUser(null)
        } finally {
            setLoading(false)
        }
    }

    const login = () => {
        window.location.href = `${getApiRoot()}/auth/discord/login`
    }

    const logout = async () => {
        try {
            await fetch(`${getApiRoot()}/auth/discord/logout`, {
                method: "POST",
                credentials: "include",
            })
            setUser(null)
        } catch (err) {
            console.error("Logout error:", err)
        }
    }

    useEffect(() => {
        fetchUser()
    }, [])

    return (
        <AuthContext.Provider
            value={{
                user,
                loading,
                error,
                login,
                logout,
                refetch: fetchUser,
            }}
        >
            {children}
        </AuthContext.Provider>
    )
}

// Hook to use auth context
export function useAuth() {
    const context = useContext(AuthContext)
    if (!context) {
        throw new Error("useAuth must be used within an AuthProvider")
    }
    return context
}

// Build API path helper
export function buildApiPath(path: string): string {
    const apiRoot = getApiRoot()

    if (apiRoot.endsWith("/")) {
        if (path.startsWith("/")) path = path.substring(1)
    } else {
        if (!path.startsWith("/")) path = "/" + path
    }

    return apiRoot + path
}

// Export API root getter
export {getApiRoot}
