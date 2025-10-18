// Player represents a user account
interface Player {
    id: string;
    did: string;
    display_name: string;
    avatar?: string | null;
    settings?: string | null;
    score: number;
    created_at: string; // ISO date string
    updated_at: string; // ISO date string
    deleted_at?: string | null; // ISO date string
}

// Show represents a WAN show episode
interface Show {
    id: string;
    youtube_id?: string | null;
    scheduled_time?: string | null; // ISO date string
    actual_start_time?: string | null; // ISO date string
    thumbnail?: string | null;
    metadata?: string | null;
    created_at: string; // ISO date string
    updated_at: string; // ISO date string
    deleted_at?: string | null; // ISO date string
}

// Tile represents a bingo tile definition
interface Tile {
    id: string;
    title: string;
    category?: string | null;
    last_drawn?: string | null; // ISO date string
    created_by?: string | null;
    settings?: string | null;
    created_at: string; // ISO date string
    updated_at: string; // ISO date string
    deleted_at?: string | null; // ISO date string
}

// ShowTile represents the junction table linking tiles to shows
interface ShowTile {
    show_id: string;
    tile_id: string;
    weight: number;
    score: number;
    created_at: string; // ISO date string
    updated_at: string; // ISO date string
    deleted_at?: string | null; // ISO date string
}

// Board represents a player's bingo board for a specific show
interface Board {
    id: string;
    player_id: string;
    show_id: string;
    tiles: string[]; // Array of tile IDs
    winner: boolean;
    total_score: number;
    regeneration_diminisher: number;
    created_at: string; // ISO date string
    updated_at: string; // ISO date string
    deleted_at?: string | null; // ISO date string
}

// TileConfirmation records when tiles are confirmed during a show
interface TileConfirmation {
    id: string;
    show_id: string;
    tile_id: string;
    confirmed_by?: string | null;
    context?: string | null;
    confirmation_time: string; // ISO date string
    created_at: string; // ISO date string
    updated_at: string; // ISO date string
    deleted_at?: string | null; // ISO date string
}

// Message records chat messages during a show
interface Message {
    id: string;
    show_id: string;
    player_id: string;
    contents: string;
    system: boolean;
    replying?: string | null;
    created_at: string; // ISO date string
    updated_at: string; // ISO date string
    deleted_at?: string | null; // ISO date string
}
