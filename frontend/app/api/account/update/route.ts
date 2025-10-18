import { NextResponse } from "next/server"
import { buildApiPath } from "@/lib/auth"

export async function POST(request: Request) {
  try {
    // Get the Discord token from cookies
    const cookies = request.headers.get("cookie")
    if (!cookies || !cookies.includes("discord-token")) {
      return NextResponse.json({ error: "Unauthorized" }, { status: 401 })
    }

    const body = await request.json()
    const { displayName, avatarUrl, chatColor, backgroundImageEnabled } = body

    // Forward to Go backend API
    const response = await fetch(buildApiPath("/api/user/preferences"), {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Cookie: cookies, // Forward cookies to backend
      },
      body: JSON.stringify({
        displayName,
        avatarUrl,
        chatColor,
        backgroundImageEnabled,
      }),
    })

    if (!response.ok) {
      return NextResponse.json(
        { error: "Failed to update preferences" },
        { status: response.status }
      )
    }

    return NextResponse.json({ success: true })
  } catch (error) {
    console.error("Failed to update account:", error)
    return NextResponse.json({ error: "Failed to update account" }, { status: 500 })
  }
}
