use rand::prelude::*;
use std::collections::HashSet;

// 1. pre‑select a winning set (5 distinct indices)
// pub fn pick_winning_set(rng: &mut ThreadRng) -> HashSet<usize> {
//     let mut win = HashSet::new();
//     while win.len() < 5 {
//         win.insert(rng.random_range(0..25));
//     }
//     win
// }

// 2. verify that the set contains a valid line on the board
// pub fn has_valid_line(board: &[[usize; 5]; 5], win: &HashSet<usize>) -> bool {
//     // rows
//     for r in 0..5 {
//         if board[r].iter().copied().all(|i| win.contains(&i)) {
//             return true;
//         }
//     }
//     // columns
//     for c in 0..5 {
//         if (0..5).map(|r| board[r][c]).all(|i| win.contains(&i)) {
//             return true;
//         }
//     }
//     // main diagonals
//     let diag1: Vec<_> = (0..5).map(|i| board[i][i]).collect();
//     if diag1.iter().copied().all(|i| win.contains(&i)) {
//         return true;
//     }
//     let diag2: Vec<_> = (0..5).map(|i| board[i][4 - i]).collect();
//     if diag2.iter().copied().all(|i| win.contains(&i)) {
//         return true;
//     }
//
//     false
// }

/// Convert a tile identifier into the flat board index (0‑24).
// pub fn id_to_index(board: &Vec<Tile>, id: &str) -> usize {
//     board.iter().position(|t| t.id == id).expect("invalid id")
// }

pub fn random_sample(rng: &mut impl Rng, size: usize, max: usize) -> HashSet<usize> {
    let mut set = HashSet::new();
    while set.len() < size {
        set.insert(rng.random_range(0..max));
    }
    set
}
