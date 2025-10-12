"use client"

import {ChatPanelProps} from "@/lib/chatUtils";
import {MobileChatPanel} from "@/components/chat/mobile-chat-panel";
import {DesktopChatPanel} from "@/components/chat/desktop-chat-panel";


export function ChatPanel(props: ChatPanelProps) {


    if (props.isMobile) {
        return (
            <MobileChatPanel {...props} />
        )
    } else {
        return (
            <DesktopChatPanel {...props} />
        )
    }
}
