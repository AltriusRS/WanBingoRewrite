# Suggestions

## Checklist

- [x] Add hide/show chat toggle on host panel
- [x] Add pop out chat to separate window
- [x] Add text size setting for bingo board
- [x] Add text size setting for host panel
- [x] Add font selection for website
- [x] Add dyslexia friendly font toggle
- [x] Limit bingo board height
- [x] Persist bingo board state in local storage
- [x] Hide chat on mobile by default

Below are some suggestions for improvements to the project.

- Add a way to resize the video player / chat window on the host panel.
  Or otherwise make it able to be popped out into a separate tab or window.

> [!note]
> This is a good idea, the chat panel can be split out into a seperate window optionally,
> and if it is split out, we should re-add the late confirmation button onto
> the host panel. Additionally we should allow the tile set to take up the remaining
> space, but instead of just making the groups wider, we should limit them to 40dvw
> and have them displayed in a flex-row with flex-wrap enabled.
> If the window is wide enough, we should also swap to using the flex-row
> wrapping layout for the tile set.
>
> The current layout on a 1920x1080 monitor is fine, but when on a
> wider monitor, the layout is simply too restrictive. It is width
> limited to be 1536px wide, and this limits severely the number of
> tiles which can be displayed without scrolling. We should instead
> allow the page to expand out to the full width of the window (with
> a margin either side of approximately 100px), and then have the tiles
> wrap to the next line when they reach the end of the page.

- Add a way to choose the size of the text on the bingo board.

> [!note]
> Make it a setting under settings.appearance.board.textSize

- Add a way to choose the size of the text on the host panel.

> [!note]
> Make it a setting under settings.appearance.hostPanel.textSize

- Add a way to choose a different font for the website.

> [!note]
> This would allow us to have a different font for the website,

- Add a way to choose a dyslexia friendly font for the website.

> [!note]
> This would allow us to have a dyslexic friendly font for the website,
> available as a toggle which would override the font setting.
> Something like settings.appearance.dyslexicFriendlyFont

- Limit the size of the bingo board so that it does not expand off the height
  of the window when the window is resized.

> [!note]
> The maximum height of the board should be the height of the
> window minus approximately 100px.

- Add a way to persist the current bingo board state in the local storage so
  that it can be restored when the page is refreshed.

> [!note]
> Could even be persisted to the database too so that it can be restored across
> sessions for the same user. Allowing multiple devices to be used at the same time.

- Hide chat on mobile by default.

> [!note]
> This would be because it makes the available view very crowded and cramped, for a feature most people cannot use without being signed in, which we don't expect too many folks to be able to do.
