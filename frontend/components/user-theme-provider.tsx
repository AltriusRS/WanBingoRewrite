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
            let fontFamily = 'var(--font-geist-sans), sans-serif' // default
            if (settings.appearance?.dyslexicFriendlyFont) {
                fontFamily = 'OpenDyslexic, Arial, sans-serif'
            } else if (font === 'serif') {
                fontFamily = 'serif'
            } else if (font === 'sans-serif') {
                fontFamily = 'sans-serif'
            } else if (font === 'roboto') {
                fontFamily = 'var(--font-roboto), sans-serif'
            } else if (font === 'lato') {
                fontFamily = 'var(--font-lato), sans-serif'
            } else if (font === 'open-sans') {
                fontFamily = 'var(--font-open-sans), sans-serif'
            } else if (font === 'montserrat') {
                fontFamily = 'var(--font-montserrat), sans-serif'
            } else if (font === 'atkinson-hyperlegible') {
                fontFamily = 'var(--font-atkinson-hyperlegible), sans-serif'
            } else if (font === 'lexend') {
                fontFamily = 'var(--font-lexend), sans-serif'
            } else if (font === 'open-dyslexic') {
                fontFamily = 'OpenDyslexic, Arial, sans-serif'
            }
            document.body.style.fontFamily = fontFamily
        }
    }, [user, setTheme])

    return <>{children}</>
}