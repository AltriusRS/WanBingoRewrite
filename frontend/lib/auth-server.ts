import 'server-only'
import {cookies} from 'next/headers'
import {Player} from './auth'

export async function getCurrentUserServer(): Promise<Player | null> {
    try {
        const cookieStore = await cookies()
        const sessionCookie = cookieStore.get('session_id')?.value

        if (!sessionCookie) {
            return null
        }

        const response = await fetch(`http://api:8000/users/me`, {
            headers: {
                'Cookie': `session_id=${sessionCookie}`,
            },
        })

        if (response.ok) {
            const data = await response.json()
            return data.user as Player
        } else {
            return null
        }
    } catch (err) {
        console.error("Server auth error:", err)
        return null
    }
}