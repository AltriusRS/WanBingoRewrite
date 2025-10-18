use serde::Deserialize;

#[derive(Debug, Clone, Deserialize)]
pub struct TileDef {
    pub id: String,
    pub x: f64,
    pub y: f64,
}

/// In-memory tile that will be mutated (weights change).
#[derive(Debug, Clone)]
pub struct Tile {
    pub id: &'static str,
    pub x: f64,
    pub y: f64,
}

// -----------------------------------------------------------------
// The regeneration diminisher – first regen = 1.00
// -----------------------------------------------------------------
// pub struct Diminish;
// impl Diminish {
//     pub const FIRST: f64 = 1.00;
// }

// The weight‑diminishing factor for the first round.
// pub const D1: f64 = 0.92;

// The weight‑diminishing factor for later rounds.
// pub const D2: f64 = 0.81;
