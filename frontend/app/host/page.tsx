import {getCurrentUserServer} from "@/lib/auth-server"
import {redirect} from "next/navigation"
import {HostDashboard} from "@/components/host/host-dashboard"
import {getApiRoot} from "@/lib/auth";

export default async function HostPage() {
    const user = await getCurrentUserServer()

    if (!user) {
        redirect(`${getApiRoot()}/auth/discord/login`)
    }

    const isUserHost = (user.permissions & 512) !== 0

    if (!isUserHost) {
        return (
            <div className="flex min-h-screen items-center justify-center bg-background">
                <div className="text-center">
                    <h1 className="text-2xl font-bold text-foreground">Access Denied</h1>
                    <p className="mt-2 text-muted-foreground">You do not have permission to access the host
                        dashboard.</p>
                </div>
            </div>
        )
    }

    return <HostDashboard/>
}
