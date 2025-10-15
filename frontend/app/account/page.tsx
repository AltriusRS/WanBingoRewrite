import { getCurrentUser } from "@/lib/auth"
import { redirect } from "next/navigation"
import { AccountSettings } from "@/components/account/account-settings"

export default async function AccountPage() {
  const user = await getCurrentUser()

  if (!user) {
    redirect("/api/auth/signin")
  }

  return <AccountSettings user={user} />
}
