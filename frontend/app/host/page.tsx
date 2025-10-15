import { getCurrentUser, isHost } from "@/lib/auth"
import { redirect } from "next/navigation"
import { HostDashboard } from "@/components/host/host-dashboard"

export default async function HostPage() {
  const user = await getCurrentUser()
  const isUserHost = await isHost()

  if (!user) {
    redirect("/api/auth/signin")
  }

  if (!isUserHost) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-background">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-foreground">Access Denied</h1>
          <p className="mt-2 text-muted-foreground">You do not have permission to access the host dashboard.</p>
        </div>
      </div>
    )
  }

  return <HostDashboard user={user} />
}
