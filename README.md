# 🧩 WAN Show Community Bingo

> **Fan-made replacement for [wanshow.bingo](https://wanshow.bingo)** — built for *The WAN Show* community to play along
> in real time.

**WAN Show Community Bingo** is a lightweight web app that lets fans mark off common moments and jokes during the WAN
Show.  
It’s designed to be built in a weekend — fast, fun, and extendable.

> ⚠️ **Disclaimer**  
> This is an **unofficial community project**.  
> It is **not affiliated with Linus Media Group, Linus Tech Tips, or The WAN Show**.  
> All references are made purely for fan enjoyment and parody.

---

## 🎯 Project Overview

The app provides:

- A **public bingo board** for fans to play along.
- A **collaborative host panel** for confirming tiles live.
- A **real-time chat feed** powered by Server-Sent Events (SSE).
- Integration with [`whenplane.com`](https://whenplane.com) to pull weekly show metadata.

Stretch goals include leaderboards, historical tracking, and playback integration — but this MVP focuses on the
essentials.

---

## 🧱 Core MVP Features

### 🎲 Public Bingo Board

- 5×5 grid generated from a predefined list of WAN Show tiles.
- “New Board” button reshuffles a fresh randomized grid.
- Clicking a tile toggles its “marked” state.
- Tiles automatically highlight when confirmed by hosts (via SSE events).
- Simple dark theme UI with Tailwind CSS and smooth Framer Motion animations.

---

### 💬 Real-Time Chat (Read-Only MVP)

- Displays a real-time chat feed using **Server-Sent Events (SSE)**.
- Shows both user messages and **system messages** (like tile confirmations).
- Collapsible panel (hidden by default on mobile).
- Currently **read-only** for public users.

---

### 🧑‍💻 Host Panel (Collaborative)

Replaces the Google Sheet workflow used by hosts to confirm tiles live.

#### Key Features

- Displays all current show tiles, grouped by category (Linus, Luke, Dan, Sponsors, etc.).
- Hosts can:
    - Click a tile → open **confirmation overlay**.
    - Add optional **context text** (e.g., `Tile Name | during sponsor segment`).
    - Confirm → sends system message to chat (`Tile confirmed: <name> | <context>`).
- Real-time **lock system**:
    - Prevents multiple hosts from confirming the same tile at once.
    - Lock automatically expires (e.g., after 5 seconds).
- “**Send Test Message**” button to broadcast a test system event.
- Integrates with [`whenplane.com`](https://whenplane.com) API to display current episode metadata.
- Access restricted to **staff accounts** (simple password or env variable check).

---

### 🧩 Basic Tile Weighting (Stubbed)

- Tile data stored in a static JSON file:
  ```json
  {
    "text": "Linus Drops Something",
    "category": "Linus",
    "weight": 1.0,
    "weeks_since_drawn": 12
  }
  ```

* Weighted random selection algorithm.
* Tiles not drawn in 20+ shows receive a reduced spawn chance (0.1× multiplier).
* No database required yet — logic runs in memory.

---

## 🧠 Technical Overview

### 🖥️ Frontend

* **Framework:** Next.js (React)
* **Styling:** Tailwind CSS (dark mode)
* **Animation:** Framer Motion
* **Communication:** Server-Sent Events (SSE) for chat + tile confirmations

**Pages**

| Path           | Purpose                                                          |
|----------------|------------------------------------------------------------------|
| `/`            | Public bingo board + chat panel                                  |
| `/host`        | Host dashboard (auth-protected)                                  |
| `/leaderboard` | The public leaderboard of all players and their validated scores |

**Key Components**

* `BingoGrid` – renders randomized board
* `Tile` – togglable visual component
* `ChatPanel` – SSE-based message feed
* `TileConfirmModal` – host confirmation dialog
* `FloatingChatToggle` – mobile toggle button

---

### ⚙️ Backend (Go API)

**Stack:** Go + Fiber or Echo + in-memory state

**Endpoints**

| Method | Endpoint             | Description                                                    |
|--------|----------------------|----------------------------------------------------------------|
| `GET`  | `/api/tiles`         | Returns a random weighted 25-tile set                          |
| `GET`  | `/api/chat/stream`   | SSE endpoint for live updates                                  |
| `POST` | `/api/chat/system`   | Auth-only endpoint to broadcast a system message               |
| `POST` | `/api/tiles/confirm` | Host confirms a tile (includes optional context)               |
| `GET`  | `/api/show`          | Fetches metadata from [`whenplane.com`](https://whenplane.com) |

**In-Memory Data**

* Tile list
* Current show metadata
* Host locks (`map[tileID]timestamp`)
* Active chat event queue

---

## 🚀 Quickstart

### Requirements

* Go 1.23+
* Node.js 20+
* pnpm or npm

### Setup

```bash
git clone https://github.com/yourname/wanshow-bingo
cd wanshow-bingo

# Install frontend deps
pnpm install

# Run backend
go run ./server/main.go

# Run frontend
pnpm dev
```

### Environment Variables

```bash
HOST_PASSWORD=supersecret
WHENPLANE_API_URL=https://whenplane.com/api/current
```

---

## 🧩 Architecture

```
/frontend
 ├── pages/
 │   ├── index.tsx        # Bingo board + chat
 │   └── host.tsx         # Host control panel
 ├── components/
 │   ├── BingoGrid.tsx
 │   ├── Tile.tsx
 │   ├── ChatPanel.tsx
 │   └── TileConfirmModal.tsx
 └── utils/
     └── sse.ts

/backend
 ├── main.go
 ├── handlers/
 │   ├── tiles.go
 │   ├── chat.go
 │   └── show.go
 └── data/
     └── tiles.json
```

---

## ✅ MVP Acceptance Criteria

* [ ] Public bingo board with interactive tiles and “New Board” button
* [ ] SSE-powered chat panel (read-only)
* [ ] Host dashboard for confirming tiles with context overlay
* [ ] Lock system to prevent duplicate confirmations
* [ ] “Send Test Message” button for broadcast testing
* [ ] Basic `whenplane.com` API integration
* [ ] Lightweight auth for hosts (env variable password)
* [ ] Simple, dark-themed responsive UI

---

## 🌟 Stretch Goals (Future Phases)

| Feature                     | Description                                 |
|-----------------------------|---------------------------------------------|
| 🏆 **Leaderboards**         | Track users’ validated wins and stats       |
| ⏱️ **Historical Tracking**  | Store and browse past tile confirmations    |
| 🎥 **Playback Mode**        | Replay confirmations synced to VOD timeline |
| 🧮 **Database Integration** | Switch from in-memory → SQLite/PostgreSQL   |
| 🧑‍💼 **Admin Dashboard**   | Spreadsheet-style tile editor               |
| 💬 **Authenticated Chat**   | Let users log in and send messages          |

---

## 🧭 Weekend Build Plan

**Day 1**

* ✅ Set up Next.js + Go project scaffolding
* ✅ Implement `/api/tiles` + `/api/chat/stream`
* ✅ Build public BingoGrid and ChatPanel

**Day 2**

* ✅ Create `/host` page with basic auth
* ✅ Add tile confirmation overlay + locking
* ✅ Integrate whenplane.com episode data
* ✅ Polish UI and deploy

---

## 🛠️ License

MIT License © 2025 Community Developers
Not affiliated with Linus Media Group or Linus Tech Tips.


---

## 🔌 Testing SSE with curl

You can use curl to verify the Server-Sent Events (SSE) chat stream and to send a test system message.

Prerequisites:
- Start the server
  - Example:
    - cd server
    - HOST_PASSWORD=changeme PORT=8080 go run ./main.go

### 1) Connect to the SSE stream
- Minimal (disables buffering so you see events as they arrive):
  - curl -N http://localhost:8080/api/chat/stream
- With an explicit Accept header (optional but clear):
  - curl -N -H "Accept: text/event-stream" http://localhost:8080/api/chat/stream

What you’ll see:
- An initial comment to confirm connection:
  - : connected
- Periodic keep-alives every ~30s (empty JSON payload):
  - data: {}

### 2) Send a test system message (POST)
Open a second terminal and run:
- export HOST_PASSWORD=changeme   # must match what the server was started with
- curl -sS -X POST \
    -H "Authorization: Bearer $HOST_PASSWORD" \
    -H "Content-Type: application/json" \
    -d '{"message":"Hello from curl","username":"tester"}' \
    http://localhost:8080/api/chat/system

If authorized, the SSE terminal will immediately show something like:
- data: {"type":"system","username":"tester","message":"Hello from curl","timestamp":"14:28"}

Tip: Pretty-print the JSON part of each event while streaming:
- curl -N http://localhost:8080/api/chat/stream \
  | sed -u 's/^data: //;/^:/d' \
  | jq -r .
  - Removes the leading "data: ", drops comment lines beginning with ":", and pipes JSON to jq.

### 3) Common troubleshooting
- 401 unauthorized when POSTing to /api/chat/system:
  - Ensure you pass: -H "Authorization: Bearer $HOST_PASSWORD"
  - Ensure the variable matches the one used to start the server.
- No events are appearing:
  - Make sure you used curl -N so output isn’t buffered.
  - Confirm the server logs show: "Server running on http://localhost:8080" and there are no runtime errors.
- Proxies or reverse proxies:
  - SSE works best without buffering. The server sets X-Accel-Buffering: no; ensure any reverse proxy honors it.
