'use client'

import { useEffect } from 'react'
import { useTheme } from 'next-themes'
import { useAuth } from '@/components/auth'

export function UserThemeProvider({ children }: { children: React.ReactNode }) {
    const { setTheme } = useTheme()
    const { user } = useAuth()

    useEffect(() => {
        if (user?.settings) {
            const settings = user.settings as any
            const preferredTheme = settings.preferredTheme
            if (preferredTheme) {
                setTheme(preferredTheme)
            }
        }
    }, [user, setTheme])

    return <>{children}</>
}