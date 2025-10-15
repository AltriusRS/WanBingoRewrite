import { Card } from "@/components/ui/card"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { ArrowLeft, ExternalLink } from "lucide-react"
import Link from "next/link"

export default function AboutPage() {
  const partners = [
    {
      name: "WhenPlane.com",
      description: 'Provides a "Lateness/Currently Live" timer at the bottom of the main page.',
      url: "https://whenplane.com",
    },
    {
      name: "TheWANDB.com",
      description:
        "Provides live Messages on the Website and Floatplane chat for verified tiles during the live show. WANDB also tracks all tiles.",
      url: "https://thewandb.com",
    },
  ]

  const team = [
    {
      name: "Brock Sexton",
      role: "Creator & Developer",
      description:
        "Brock handles hosting and development of WanShow.bingo, and created the site in March 2023 shortly after www.WanShowBingo.com went offline.",
      avatar: "/placeholder.svg",
    },
    {
      name: "Dax Anderson",
      role: "Chief Vision Officer - Floatplane OG",
      description: "Dax typically spends his time during WAN verifying tiles, and interacting with Floatplane Chat.",
      avatar: "/placeholder.svg",
    },
    {
      name: "ajgeiss0702",
      role: "Developer",
      description: "WhenPlane go brrr",
      avatar: "/placeholder.svg",
    },
    {
      name: "Arthur Amos",
      role: "Developer & Host",
      description: "TheWANDB go brrr - helps with hosting & google sheets",
      avatar: "/placeholder.svg",
    },
  ]

  const contributors = [
    {
      name: "Woofer21",
      description: "Pushed multiple fixes to the GitHub which made it to live production (Bug Fix!) - Thank You!",
    },
    {
      name: "Emperor Numerius",
      description:
        "Initial Pull request for Snapshot Functionality, QC Testing, and help with Bingo Suggestions Form through gforms",
    },
    {
      name: "Carl Ayres",
      description:
        "Carl (OSTycoon) was the original creator of WAN Show Bingo, unfortunately the GitHub has gone dormant and the website has gone offline. (Update Sept 2023, the website is backup)",
    },
    {
      name: "Skylar Ittner",
      description:
        "Skylar (skylarmt) was a contributor to the original repo for WAN Show Bingo, unfortunately the GitHub has gone dormant and the website has gone offline. (Update Sept 2023, the website is backup)",
    },
    {
      name: "James Anderson",
      description:
        "James Anderson was a contributor to the original repo for WAN Show Bingo, unfortunately the GitHub has gone dormant and the website has gone offline. (Update Sept 2023, the website is backup)",
    },
  ]

  return (
    <div className="min-h-screen bg-background">
      <header className="border-b border-border bg-card">
        <div className="container mx-auto flex items-center gap-4 px-4 py-4">
          <Link href="/">
            <Button variant="ghost" size="icon">
              <ArrowLeft className="h-4 w-4" />
            </Button>
          </Link>
          <div>
            <h1 className="text-xl font-semibold text-foreground">About WAN Show Bingo</h1>
            <p className="text-sm text-muted-foreground">Our team and partners</p>
          </div>
        </div>
      </header>

      <div className="container mx-auto max-w-4xl space-y-8 p-4 py-8">
        <section>
          <h2 className="mb-4 text-2xl font-bold text-foreground">Our Partner Websites</h2>
          <div className="grid gap-4 md:grid-cols-2">
            {partners.map((partner) => (
              <Card key={partner.name} className="p-6">
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <h3 className="mb-2 font-semibold text-foreground">{partner.name}</h3>
                    <p className="text-sm text-muted-foreground">{partner.description}</p>
                  </div>
                  <a href={partner.url} target="_blank" rel="noopener noreferrer">
                    <Button variant="ghost" size="icon">
                      <ExternalLink className="h-4 w-4" />
                    </Button>
                  </a>
                </div>
              </Card>
            ))}
          </div>
        </section>

        <section>
          <h2 className="mb-4 text-2xl font-bold text-foreground">Team</h2>
          <div className="grid gap-4 md:grid-cols-2">
            {team.map((member) => (
              <Card key={member.name} className="p-6">
                <div className="flex gap-4">
                  <Avatar className="h-16 w-16">
                    <AvatarImage src={member.avatar || "/placeholder.svg"} alt={member.name} />
                    <AvatarFallback className="bg-primary/10">
                      {member.name
                        .split(" ")
                        .map((n) => n[0])
                        .join("")}
                    </AvatarFallback>
                  </Avatar>
                  <div className="flex-1">
                    <h3 className="font-semibold text-foreground">{member.name}</h3>
                    <Badge variant="secondary" className="mb-2 text-xs">
                      {member.role}
                    </Badge>
                    <p className="text-sm text-muted-foreground">{member.description}</p>
                  </div>
                </div>
              </Card>
            ))}
          </div>
        </section>

        <section>
          <h2 className="mb-4 text-2xl font-bold text-foreground">Contributors</h2>
          <div className="space-y-4">
            {contributors.map((contributor) => (
              <Card key={contributor.name} className="p-4">
                <h3 className="mb-1 font-semibold text-foreground">{contributor.name}</h3>
                <p className="text-sm text-muted-foreground">{contributor.description}</p>
              </Card>
            ))}
          </div>
        </section>

        <Card className="bg-muted/50 p-6">
          <p className="text-center text-sm text-muted-foreground">
            WAN Show Bingo is a community project and is not affiliated with Linus Media Group.
          </p>
        </Card>
      </div>
    </div>
  )
}
