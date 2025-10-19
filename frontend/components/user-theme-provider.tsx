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
            let preferredTheme = "dark"
            let font = "default"

            // Check nested structure first, then fallback to flat structure
            if (settings.themes?.preferred) {
                preferredTheme = settings.themes.preferred
            } else if (settings.preferredTheme) {
                preferredTheme = settings.preferredTheme
            }

            if (settings.appearance?.font) {
                font = settings.appearance.font
            }

            if (preferredTheme) {
                setTheme(preferredTheme)
            }

            // Apply font
            let fontFamily = 'inherit'
            if (settings.appearance?.dyslexicFriendlyFont) {
                fontFamily = 'Arial, sans-serif' // Simple dyslexia friendly font
            } else if (font === 'serif') {
                fontFamily = 'serif'
            } else if (font === 'sans-serif') {
                fontFamily = 'sans-serif'
            }
            document.body.style.fontFamily = fontFamily
        }
    }, [user, setTheme])

    return <>{children}</>
}