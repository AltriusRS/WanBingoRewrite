"use server"

export async function getCurrentUser() {
  try {
    const { user } = await getUser()
    return user
  } catch (error) {
    return null
  }
}

export async function getSessionId(): Promise<string | null> {
  try {
    const { sessionId } = await getUser()
    return sessionId || null
  } catch (error) {
    return null
  }
}

export async function isHost(): Promise<boolean> {
  try {
    const { user } = await getUser()
    if (!user) return false

    // Check if user email is from LMG domain or is in host list
    const hostEmails = (process.env.HOST_EMAIL || "").split(",").map((e) => e.trim())
    const isLMGDomain = user.email?.endsWith("@linusmediagroup.com")
    const isInHostList = hostEmails.includes(user.email || "")

    return isLMGDomain || isInHostList
  } catch (error) {
    return false
  }
}

export async function isChatBanned(): Promise<{ banned: boolean; reason?: string; expiry?: Date }> {
  try {
    const { user } = await getUser()
    if (!user) return { banned: false }

    // Check user metadata for chat ban
    const metadata = user.rawAttributes as any
    const chatBanned = metadata?.chatBanned === true
    const reason = metadata?.chatBanReason
    const expiry = metadata?.chatBanExpiry ? new Date(metadata.chatBanExpiry) : undefined

    // Check if ban has expired
    if (chatBanned && expiry && expiry < new Date()) {
      return { banned: false }
    }

    return { banned: chatBanned, reason, expiry }
  } catch (error) {
    return { banned: false }
  }
}



async function getUser(): Promise<UserProfile> {
  let response = await fetch(buildApiPath('/me'));

  let body = await response.json();

  console.log(body);
  return body;
}


/**
 * Build a correctly formatted api path based on the contextually known api 
 * base path and the provided route path
 * @param path {string} - The path to be built
 * @returns {string} - The constructed path
 */
export function buildApiPath(path: string): string {
  if(process.env.API_URL!.endsWith('/')) {
    if (path.startsWith('/')) path = path.substring(1,path.length-1)
  } else {
    if (!path.startsWith('/')) path = '/' + path;
  }

  return process.env.API_URL + path;
}