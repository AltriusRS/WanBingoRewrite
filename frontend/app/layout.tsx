import type {Metadata} from "next";
import {Geist, Geist_Mono} from "next/font/google";
import "./globals.css";
import {ChatProvider} from "@/components/chat/chat-context";
import {AuthProvider} from "@/components/auth";
import {ThemeProvider} from "@/components/theme-provider";
import {UserThemeProvider} from "@/components/user-theme-provider";
import {PostHogProvider} from "@/components/posthog-provider";

const geistSans = Geist({
    variable: "--font-geist-sans",
    subsets: ["latin"],
});

const geistMono = Geist_Mono({
    variable: "--font-geist-mono",
    subsets: ["latin"],
});

export const metadata: Metadata = {
    title: "WAN Show Bingo",
    description: "The classic tech news bingo game!",
};

export default function RootLayout({
                                       children,
                                   }: Readonly<{
    children: React.ReactNode;
}>) {
    return (
        <html lang="en" suppressHydrationWarning>
        <body
            className={`${geistSans.variable} ${geistMono.variable} antialiased`}
        >
        <PostHogProvider>
            <ThemeProvider
                attribute="class"
                defaultTheme="dark"
                enableSystem
                disableTransitionOnChange
            >
                <AuthProvider>
                    <UserThemeProvider>
                        <ChatProvider>
                            {children}
                        </ChatProvider>
                    </UserThemeProvider>
                </AuthProvider>
            </ThemeProvider>
        </PostHogProvider>
        </body>
        </html>
    );
}
