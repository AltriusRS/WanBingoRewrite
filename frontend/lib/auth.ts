// Player type
export interface Player {
    id: string
    did: string
    display_name: string
    avatar?: string
    settings?: any
    score: number
    permissions: number
    created_at: string
    updated_at: string
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

export async function getCurrentUser(): Promise<Player | null> {
    try {
        const response = await fetch(`${getApiRoot()}/users/me`, {
            credentials: "include", // Include cookies
        })

        if (response.ok) {
            const data = await response.json()
            return data.user as Player
        } else {
            return null
        }
    } catch (err) {
        console.error("Auth error:", err)
        return null
    }
}

export function getApiRoot(): string {
    return process.env.NEXT_PUBLIC_API_ROOT || "http://localhost:8000"
}

export async function isHost(): Promise<boolean> {
    try {
        const response = await fetch(`${getApiRoot()}/users/me`, {
            credentials: "include",
        })

        if (response.ok) {
            const data = await response.json()
            const user = data.user
            // Check if user has host permission (PermCanHost = 512)
            return (user.permissions & 512) !== 0
        }
        return false
    } catch (err) {
        console.error("Auth error:", err)
        return false
    }
}

