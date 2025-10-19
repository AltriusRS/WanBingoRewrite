import {getCurrentUserServer} from "@/lib/auth-server"
import {redirect} from "next/navigation"
import {getApiRoot} from "@/lib/auth";
import {DesktopChatPanel} from "@/components/chat/desktop-chat-panel"

export default async function ChatPage() {
    const user = await getCurrentUserServer()

    if (!user) {
        redirect(`${getApiRoot()}/auth/discord/login`)
    }

    return (
        <div className="h-screen bg-background">
            <DesktopChatPanel />
        </div>
    )
}