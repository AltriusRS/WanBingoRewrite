"use client";

import {ChatMessage} from "@/lib/chatUtils";


import React, {useState} from "react";
import {MemoizedMarkdown} from "@/components/ui/markdown";
import {Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle} from "@/components/ui/dialog";
import {Button} from "@/components/ui/button";


interface StandardMessageProps {
    msg: ChatMessage
}

export function SystemMessage({msg}: StandardMessageProps) {
    const [open, setOpen] = useState(false);
    const [clickedUrl, setClickedUrl] = useState("");

    return (
        <div className="rounded-lg border border-primary/30 bg-primary/10 p-3">
            <div className="flex items-center gap-2">
                <div className="h-1.5 w-1.5 rounded-full bg-primary"/>
                <span className="text-xs font-medium text-primary">SYSTEM</span>
                <span className="text-xs text-muted-foreground">{new Date(msg.created_at).toLocaleTimeString()}</span>
            </div>
            <div className="mt-1 font-medium text-sm text-foreground prose">
                <MemoizedMarkdown
                    key={`${msg.id}-text`}
                    id={msg.id}
                    content={msg.contents}
                    onLinkClick={(href) => {
                        setClickedUrl(href);
                        setOpen(true);
                    }}
                />
                <Dialog open={open} onOpenChange={setOpen}>
                    <DialogContent>
                        <DialogHeader>
                            <DialogTitle>External Link</DialogTitle>
                        </DialogHeader>
                        <div className="py-2">You&apos;re about to open: {clickedUrl}</div>
                        <DialogFooter className="flex gap-2 justify-end">
                            <Button
                                onClick={() => {
                                    window.open(clickedUrl, "_blank");
                                    setOpen(false);
                                }}
                            >
                                Proceed
                            </Button>
                            <Button variant="outline" onClick={() => setOpen(false)}>
                                Cancel
                            </Button>
                        </DialogFooter>
                    </DialogContent>
                </Dialog>
            </div>
        </div>
    )
}
