import {redirect} from "next/navigation"
import {AccountSettings} from "@/components/account/account-settings"
import {getApiRoot, getCurrentUser} from "@/lib/auth";

export default async function AccountPage() {
    const user = await getCurrentUser()

    if (!user) {
        redirect(`${getApiRoot()}/auth/discord/signin`)
    }

    return <AccountSettings user={user}/>
}
