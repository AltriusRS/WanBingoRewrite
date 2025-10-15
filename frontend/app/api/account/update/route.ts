import { getUser } from "@workos-inc/authkit-nextjs"
import { NextResponse } from "next/server"

export async function POST(request: Request) {
  try {
    const { user } = await getUser()
    if (!user) {
      return NextResponse.json({ error: "Unauthorized" }, { status: 401 })
    }

    const body = await request.json()
    const { displayName, avatarUrl, chatColor, backgroundImageEnabled } = body

    // Update user metadata in WorkOS
    // Note: This requires WorkOS API integration
    // For now, we'll store in a separate database table via the Go backend
    await fetch("http://localhost:8080/api/user/preferences", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        userId: user.id,
        displayName,
        avatarUrl,
        chatColor,
        backgroundImageEnabled,
      }),
    })

    return NextResponse.json({ success: true })
  } catch (error) {
    console.error("Failed to update account:", error)
    return NextResponse.json({ error: "Failed to update account" }, { status: 500 })
  }
}
