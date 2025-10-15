"use client"

import { useState } from "react"
import { Card } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Textarea } from "@/components/ui/textarea"
import { Label } from "@/components/ui/label"
import { Send } from "lucide-react"

export function TestMessagePanel() {
  const [message, setMessage] = useState("")
  const [sending, setSending] = useState(false)

  const sendTestMessage = async () => {
    if (!message.trim()) return

    setSending(true)
    try {
      await fetch("http://localhost:8080/api/host/test-message", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ message }),
      })
      setMessage("")
    } catch (error) {
      console.error("Failed to send test message:", error)
    } finally {
      setSending(false)
    }
  }

  return (
    <Card className="p-4">
      <h3 className="mb-4 font-semibold text-foreground">Send Test System Message</h3>

      <div className="space-y-4">
        <div className="space-y-2">
          <Label htmlFor="test-message">Message</Label>
          <Textarea
            id="test-message"
            placeholder="Enter a test system message..."
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            rows={4}
          />
        </div>

        <Button onClick={sendTestMessage} disabled={sending || !message.trim()} className="gap-2">
          <Send className="h-4 w-4" />
          {sending ? "Sending..." : "Send Test Message"}
        </Button>
      </div>
    </Card>
  )
}
