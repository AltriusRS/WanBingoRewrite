# Chat

Chat is powered by Server Sent Events (SSE). Interactions which
trigger an SSE are described in this document, but are labeled
with the relevant HTTP call and payloads.

## Events

### Format

The basic format of all payloads sent over the SSE gateway is
as follows:

```json
{
  "id": "ofCsrk2uu4bM04KKr2UIn",
  "opcode": "some-event",
  "data": T
}
```

- The `id` field is used to identify the event. It is unique to
  every single event coming from the gateway.
- The `opcode` field is used to identify the type of event.
- The `data` field is used to send data to the client. This field
  can be any JSON serializable data. (\<T\>)

> **Note**: The gateway does not guarantee delivery order of events.

From this point onwards, the payloads for each opcode are considered
to be wrapped in this payload interface.

### hub.connected

This is sent when the client has connected to the SSE gateway.
The data provided by this event is simply "Connected to the chat hub."

### hub.authenticated

This is sent when the client has successfully been authenticated.
The data provided by this event is important, and should be used,
when enabling or disabling chat features. See the
[chat permissions](#chat-permissions) section for more information.

```json5
{
  // Only true when the user is granted chat permissions.
  "canChat": true,
  // Only true when the user is granted hosting permissions.
  "canHost": true,
  // Only true when the user is granted moderation permissions.
  "canModerate": true,
}
```

### hub.connections.count

Status event telling the client how many users are currently connected
to the chat hub

The data provided by this event is an integer.

### whenplane.aggregate

This event is sent when a new aggregate is received from the whenplane
server. It is immediately forwarded to the clients over the SSE gateway.

This data is used to display the thumbnail, title, and live status of
The WAN Show on the website.

## Chat Permissions

These are the permissions that can be granted to a user in the chat hub.
They are stored server side and can only be modified by a moderator.

```json
{
  "canSendMessages": true,
  "canSendWhispers": true,
  "canDeleteOwnMessages": true,
  "canDeleteMessages": true,
  "canBanUsers": true,
  "canKickUsers": true,
  "canMuteUsers": true,
  "canUnmuteUsers": true,
  "canHost": true,
  "canModerate": true,
  "canManageChat": true,
  "canManageChatPermissions": true,
  "canSuggestTiles": true,
  "canReviewTiles": true,
  "canApproveTiles": true,
  "canManageTiles": true,
  "canPromotePlayers": true,
  "canModifyShowData": true
}
```
