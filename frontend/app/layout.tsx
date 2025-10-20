import type {Metadata} from "next";
import {Geist, Geist_Mono, Roboto, Lato, Open_Sans, Montserrat, Atkinson_Hyperlegible, Lexend} from "next/font/google";
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

const roboto = Roboto({
    weight: ['400', '700'],
    subsets: ['latin'],
    variable: '--font-roboto',
});

const lato = Lato({
    weight: ['400', '700'],
    subsets: ['latin'],
    variable: '--font-lato',
});

const openSans = Open_Sans({
    weight: ['400', '700'],
    subsets: ['latin'],
    variable: '--font-open-sans',
});

const montserrat = Montserrat({
    weight: ['400', '700'],
    subsets: ['latin'],
    variable: '--font-montserrat',
});

const atkinsonHyperlegible = Atkinson_Hyperlegible({
    weight: ['400', '700'],
    subsets: ['latin'],
    variable: '--font-atkinson-hyperlegible',
});

const lexend = Lexend({
    weight: ['400', '700'],
    subsets: ['latin'],
    variable: '--font-lexend',
});

export default function RootLayout({
                                       children,
                                   }: Readonly<{
    children: React.ReactNode;
}>) {
    return (
        <html lang="en" suppressHydrationWarning>
        <head>
            <link rel="preconnect" href="https://fonts.googleapis.com" />
            <link rel="preconnect" href="https://fonts.gstatic.com" crossOrigin="" />
            <link href="https://fonts.cdnfonts.com/css/open-dyslexic" rel="stylesheet" />
        </head>
        <body
            className={`${geistSans.variable} ${geistMono.variable} ${roboto.variable} ${lato.variable} ${openSans.variable} ${montserrat.variable} ${atkinsonHyperlegible.variable} ${lexend.variable} antialiased`}
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
                            <div id="app-container">
                                {children}
                            </div>
                        </ChatProvider>
                    </UserThemeProvider>
                </AuthProvider>
            </ThemeProvider>
        </PostHogProvider>
        </body>
        </html>
    );
}
