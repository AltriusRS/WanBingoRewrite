import {redirect} from "next/navigation"
import {AccountSettings} from "@/components/account/account-settings"
import {getApiRoot} from "@/lib/auth";
import {getCurrentUserServer} from "@/lib/auth-server";

export default async function AccountPage() {
    const user = await getCurrentUserServer()

    if (!user) {
        redirect(`${getApiRoot()}/auth/discord/login`)
    }

    return <AccountSettings/>
}
