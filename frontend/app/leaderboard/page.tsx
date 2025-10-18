import {Card} from "@/components/ui/card"
import {Avatar, AvatarFallback, AvatarImage} from "@/components/ui/avatar"
import {Badge} from "@/components/ui/badge"
import {Button} from "@/components/ui/button"
import {ArrowLeft, Award, Medal, Trophy} from "lucide-react"
import Link from "next/link"
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow} from "@/components/ui/table"
import {getApiRoot} from "@/lib/auth";

export default async function LeaderboardPage() {
    // Fetch leaderboard data from API
    const leaderboardData = await fetch(`${getApiRoot()}/leaderboard`, {
        cache: "no-store",
    })
        .then((res) => res.json())
        .catch(() => [])

    const getRankIcon = (rank: number) => {
        switch (rank) {
            case 1:
                return <Trophy className="h-5 w-5 text-yellow-500"/>
            case 2:
                return <Medal className="h-5 w-5 text-gray-400"/>
            case 3:
                return <Award className="h-5 w-5 text-amber-600"/>
            default:
                return <span className="text-sm font-semibold text-muted-foreground">#{rank}</span>
        }
    }

    return (
        <div className="min-h-screen bg-background">
            <header className="border-b border-border bg-card">
                <div className="container mx-auto flex items-center gap-4 px-4 py-4">
                    <Link href="/">
                        <Button variant="ghost" size="icon">
                            <ArrowLeft className="h-4 w-4"/>
                        </Button>
                    </Link>
                    <div>
                        <h1 className="text-xl font-semibold text-foreground">Leaderboard</h1>
                        <p className="text-sm text-muted-foreground">Top WAN Show Bingo players</p>
                    </div>
                </div>
            </header>

            <div className="container mx-auto max-w-4xl space-y-6 p-4 py-8">
                {/* Top 3 Podium */}
                {leaderboardData.length >= 3 && (
                    <div className="grid gap-4 md:grid-cols-3">
                        {/* 2nd Place */}
                        <Card className="order-1 p-6 md:order-1">
                            <div className="flex flex-col items-center text-center">
                                <Medal className="mb-3 h-8 w-8 text-gray-400"/>
                                <Avatar className="mb-3 h-16 w-16">
                                    <AvatarImage
                                        src={leaderboardData[1]?.avatar || "/placeholder.svg"}
                                        alt={leaderboardData[1]?.username}
                                    />
                                    <AvatarFallback
                                        className="bg-primary/10">{leaderboardData[1]?.username?.[0]}</AvatarFallback>
                                </Avatar>
                                <h3 className="font-semibold text-foreground">{leaderboardData[1]?.username}</h3>
                                <Badge variant="secondary" className="mt-2">
                                    {leaderboardData[1]?.wins} wins
                                </Badge>
                                <p className="mt-1 text-sm text-muted-foreground">{leaderboardData[1]?.points} points</p>
                            </div>
                        </Card>

                        {/* 1st Place */}
                        <Card className="order-2 border-primary bg-primary/5 p-6 md:order-2">
                            <div className="flex flex-col items-center text-center">
                                <Trophy className="mb-3 h-10 w-10 text-yellow-500"/>
                                <Avatar className="mb-3 h-20 w-20 ring-2 ring-primary">
                                    <AvatarImage
                                        src={leaderboardData[0]?.avatar || "/placeholder.svg"}
                                        alt={leaderboardData[0]?.username}
                                    />
                                    <AvatarFallback
                                        className="bg-primary/10">{leaderboardData[0]?.username?.[0]}</AvatarFallback>
                                </Avatar>
                                <h3 className="text-lg font-bold text-foreground">{leaderboardData[0]?.username}</h3>
                                <Badge className="mt-2">{leaderboardData[0]?.wins} wins</Badge>
                                <p className="mt-1 text-sm text-muted-foreground">{leaderboardData[0]?.points} points</p>
                            </div>
                        </Card>

                        {/* 3rd Place */}
                        <Card className="order-3 p-6 md:order-3">
                            <div className="flex flex-col items-center text-center">
                                <Award className="mb-3 h-8 w-8 text-amber-600"/>
                                <Avatar className="mb-3 h-16 w-16">
                                    <AvatarImage
                                        src={leaderboardData[2]?.avatar || "/placeholder.svg"}
                                        alt={leaderboardData[2]?.username}
                                    />
                                    <AvatarFallback
                                        className="bg-primary/10">{leaderboardData[2]?.username?.[0]}</AvatarFallback>
                                </Avatar>
                                <h3 className="font-semibold text-foreground">{leaderboardData[2]?.username}</h3>
                                <Badge variant="secondary" className="mt-2">
                                    {leaderboardData[2]?.wins} wins
                                </Badge>
                                <p className="mt-1 text-sm text-muted-foreground">{leaderboardData[2]?.points} points</p>
                            </div>
                        </Card>
                    </div>
                )}

                {/* Full Leaderboard Table */}
                <Card>
                    <Table>
                        <TableHeader>
                            <TableRow>
                                <TableHead className="w-16">Rank</TableHead>
                                <TableHead>Player</TableHead>
                                <TableHead className="text-right">Wins</TableHead>
                                <TableHead className="text-right">Points</TableHead>
                                <TableHead className="text-right">Win Rate</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {leaderboardData.map((player: any, index: number) => (
                                <TableRow key={player.userId}>
                                    <TableCell className="font-medium">{getRankIcon(index + 1)}</TableCell>
                                    <TableCell>
                                        <div className="flex items-center gap-3">
                                            <Avatar className="h-8 w-8">
                                                <AvatarImage src={player.avatar || "/placeholder.svg"}
                                                             alt={player.username}/>
                                                <AvatarFallback
                                                    className="bg-primary/10 text-xs">{player.username?.[0]}</AvatarFallback>
                                            </Avatar>
                                            <span className="font-medium">{player.username}</span>
                                        </div>
                                    </TableCell>
                                    <TableCell className="text-right">{player.wins}</TableCell>
                                    <TableCell className="text-right">{player.points}</TableCell>
                                    <TableCell
                                        className="text-right">{((player.wins / player.gamesPlayed) * 100).toFixed(1)}%</TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                </Card>

                {leaderboardData.length === 0 && (
                    <Card className="p-8 text-center">
                        <p className="text-muted-foreground">
                            No leaderboard data available yet. Start playing to get on the board!
                        </p>
                    </Card>
                )}
            </div>
        </div>
    )
}
