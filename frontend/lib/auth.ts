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

export async function getCurrentUser(): Promise<DiscordUser | null> {
    try {
        const response = await fetch(`${getApiRoot()}/auth/discord/user`, {
            credentials: "include", // Include cookies
        })

        if (response.ok) {
            const data = await response.json()
            return data as DiscordUser
        } else if (response.status === 401) {
            return null
        } else {
            throw new Error("Failed to fetch user")
        }
    } catch (err) {
        console.error("Auth error:", err)
        return null
    }
}

export function getApiRoot(): string {
    return process.env.NEXT_PUBLIC_API_ROOT || "http://localhost:8000"
}