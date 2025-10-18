'use client'

import { useTheme } from 'next-themes'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Palette } from 'lucide-react'
import { useAuth } from '@/components/auth'
import { getApiRoot } from '@/lib/auth'

const themes = [
  { name: 'WAN Show', value: 'dark' },
  { name: 'Light', value: 'light' },
  { name: 'Winter', value: 'winter' },
  { name: 'Halloween', value: 'halloween' },
  { name: 'Easter', value: 'easter' },
  { name: 'Summer', value: 'summer' },
  { name: 'Pitch Black', value: 'pitch-black' },
]

export function ThemeSelector() {
  const { setTheme, theme } = useTheme()
  const { user, refetch } = useAuth()

  const handleThemeChange = async (newTheme: string) => {
    setTheme(newTheme)

    // Save to user settings if logged in
    if (user) {
      try {
        const currentSettings = user.settings as any || {}
        const updatedSettings = {
          ...currentSettings,
          themes: {
            ...currentSettings.themes,
            preferred: newTheme,
          },
        }

        await fetch(`${getApiRoot()}/users/me`, {
          method: "PUT",
          credentials: "include",
          headers: {"Content-Type": "application/json"},
          body: JSON.stringify({
            display_name: user.display_name,
            avatar: user.avatar,
            settings: updatedSettings,
          }),
        })

        // Refetch user data to update the theme provider
        await refetch()
      } catch (error) {
        console.error("Failed to save theme preference:", error)
      }
    }
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" size="sm" className="gap-2 bg-transparent">
          <Palette className="h-4 w-4" />
          <span className="hidden sm:inline">Theme</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        {themes.map((t) => (
          <DropdownMenuItem
            key={t.value}
            onClick={() => setTheme(t.value)}
            className={theme === t.value ? 'bg-accent' : ''}
          >
            {t.name}
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}