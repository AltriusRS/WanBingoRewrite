use rayon::prelude::*;
// use rand::prelude::*;
use indicatif::ProgressBar;
use rand::Rng;
use std::collections::HashSet;
use std::sync::Arc;
use std::sync::atomic::{AtomicBool, Ordering};
use std::time::Duration;
use std::{env, fs, thread};
use sysinfo::{Pid, ProcessesToUpdate, System};

pub(crate) mod definitions;
pub(crate) mod helpers;
pub(crate) mod plot;

use definitions::*;
use helpers::*;
use plot::*;

// ────────────────────── main ────────────────────────

/// The command‑line entry point – everything that follows is a single
/// function that can be dropped into your existing file.
fn main() {
    /* ------------------------------------------------------------------
    Parse command line arguments
    ------------------------------------------------------------------ */
    let mut args = env::args();
    args.next(); // skip program name

    //     let games_this_week: usize = match args.next() {
    //         Some(v) => v.parse().expect("<iterations> must be a positive integer"),
    //         None => panic!("usage: <games-per-week> <weeks> [tiles.json] [grow] [shrink]"),
    //     };

    let weeks: usize = match args.next() {
        Some(v) => v.parse().expect("<weeks> must be a positive integer"),
        None => panic!("<weeks> argument missing"),
    };

    let tiles_file = args.next().unwrap_or_else(|| "./tiles.sample.json".into());

    let grow_factor: f64 = match args.next() {
        Some(v) => v.parse().expect("<grow> must be a floating point value"),
        None => 1.05,
    };
    let shrink_factor: f64 = match args.next() {
        Some(v) => v.parse().expect("<shrink> must be a floating point value"),
        None => 0.95,
    };

    /* ------------------------------------------------------------------
    Read the tiles
    ------------------------------------------------------------------ */
    let json_str = fs::read_to_string(&tiles_file).expect("failed to read tiles.json");
    let defs: Vec<TileDef> =
        serde_json::from_str(&json_str).expect("JSON is not an array of objects with id, x, y");

    /* Convert the JSON into a `Vec<Tile>` – we keep the IDs as
    `'static str` by leaking them (this is fine for a short‑lived
    program).  If you prefer to avoid leaks, store the strings in a
    separate vector and take references from there. */
    let mut tiles: Vec<Tile> = defs
        .into_iter()
        .map(|d| Tile {
            id: Box::leak(d.id.into_boxed_str()),
            x: d.x,
            y: d.y,
        })
        .collect();

    /* ------------------------------------------------------------------
    Stats
    ------------------------------------------------------------------ */
    let mut final_scores = Vec::<f64>::new();
    let mut potential_maxes = Vec::<f64>::new();

    /* ------------------------------------------------------------------
    Random number generator
    ------------------------------------------------------------------ */
    let mut rng = rand::rng(); // the RNG defined in your file

    /* ------------------------------------------------------------------
    We keep a weight vector that holds the current weight for every
    tile – it is updated after each week.  Initially the weight
    equals the `y` value supplied by the master list.
    ------------------------------------------------------------------ */
    let mut weights: Vec<f64> = tiles.iter().map(|t| t.y).collect();

    let mut total_games: usize = 0;
    let mut total_wins: usize = 0;
    let mut total_losses: usize = 0;

    let mut weekly_means = Vec::<f64>::with_capacity(weeks);
    let mut weekly_raw: Vec<Vec<f64>> = Vec::new();

    println!("Starting simulation...");
    println!("Tiles: {}", tiles.len());
    println!("Weeks: {}", weeks);
    println!("Grow factor: {}", grow_factor);
    println!("Shrink factor: {}", shrink_factor);
    println!();

    /* ----- progress bar for weeks ----------------------------------- */
    let pb = Arc::new(ProgressBar::new(weeks as u64));
    pb.set_style(
        indicatif::ProgressStyle::default_bar()
            .template("[{elapsed_precise}] {bar:40.cyan/blue} {pos:>7}/{len:7} ({eta}) {msg}")
            .unwrap(),
    );

    pb.set_message("Initializing...");

    /* ----- memory printer thread ------------------------------------- */
    // let done_flag = Arc::new(AtomicBool::new(false));
    // let pb_clone = pb.clone();
    // let done_clone = done_flag.clone();
    // thread::spawn(move || {
    //     memory_bar_updater(pb_clone, done_clone);
    // });

    let start = std::time::Instant::now();

    /* ------------------------------------------------------------------
    Simulation
    ------------------------------------------------------------------ */
    for w in 0..weeks {
        simulate_week(
            &mut rng,
            &mut total_games,
            &mut total_wins,
            &mut total_losses,
            &mut tiles,
            &mut weights,
            &mut final_scores,
            &mut potential_maxes,
            &mut weekly_means,
            &mut weekly_raw,
            &grow_factor,
            &shrink_factor,
        );
        pb.inc(1);

        if w % 10000 == 0 {
            let elapsed = start.elapsed();
            let elapsed_sec = elapsed.as_secs() as f64 + elapsed.subsec_nanos() as f64 * 1e-9;
            let sims_per_sec = total_games as f64 / elapsed_sec;
            pb.set_message(format!("{:.2} games/sec | {:.2}% win rate", sims_per_sec, (total_wins as f64 / total_games as f64) * 100.0));
        }
    }

    // done_flag.store(true, Ordering::Relaxed);

    pb.finish_with_message("Done!");

    /* ------------------------------------------------------------------
    Global statistics
    ------------------------------------------------------------------ */
    print_global_stats(
        final_scores,
        potential_maxes,
        weeks,
        total_games,
        total_wins,
        total_losses,
        weekly_means,
        weekly_raw,
    );
}

fn print_global_stats(
    final_scores: Vec<f64>,
    potential_maxes: Vec<f64>,
    weeks: usize,
    total_games: usize,
    total_wins: usize,
    total_losses: usize,
    weekly_means: Vec<f64>,
    weekly_raw: Vec<Vec<f64>>,
) {
    let mean_final: f64 = final_scores.iter().copied().sum::<f64>() / total_games as f64;
    let min_final = final_scores.iter().cloned().fold(f64::INFINITY, f64::min);
    let max_final = final_scores
        .iter()
        .cloned()
        .fold(f64::NEG_INFINITY, f64::max);

    println!("\n=== Global statistics over all {} games ===", total_games);
    println!("Mean final score : {:.2}", mean_final);
    println!("Min  final score : {:.2}", min_final);
    println!("Max  final score : {:.2}\n", max_final);
    println!("Total wins       : {}", total_wins);
    println!("Total losses     : {}", total_losses);
    println!("Weeks            : {}", weeks);
    println!(
        "Win rate         : {:.2}%",
        (total_wins as f64 / total_games as f64) * 100.0
    );

    /* ------------------------------------------------------------------
    Histograms
    ------------------------------------------------------------------ */
    if let Err(e) = save_histogram("final_scores.png", "Final Scores", &final_scores, 5.0) {
        eprintln!("Failed to save final_scores.png: {}", e);
    }
    if let Err(e) = save_histogram(
        "potential_maxes.png",
        "Potential Max Scores",
        &potential_maxes,
        5.0,
    ) {
        eprintln!("Failed to save potential_maxes.png: {}", e);
    }

    /* ------------------------------------------------------------------
    Weekly means
    ------------------------------------------------------------------ */

    if let Err(e) = draw_weekly_scores("weekly_scores.png", weeks, &weekly_means) {
        eprintln!("Failed to create weekly_scores.png: {}", e);
    }
    if let Err(e) = draw_weekly_boxplot("weekly_boxplot.png", weeks, &weekly_raw) {
        eprintln!("Failed to create weekly_boxplot.png: {}", e);
    }
}

fn simulate_week(
    rng: &mut impl Rng, // ── we keep the original RNG for the week counter
    total_games: &mut usize,
    total_wins: &mut usize,
    total_losses: &mut usize,
    tiles: &mut Vec<Tile>,
    weights: &mut Vec<f64>,
    final_scores: &mut Vec<f64>, // ── will be extended after the parallel part
    potential_maxes: &mut Vec<f64>, // ── same here
    weekly_means: &mut Vec<f64>,
    weekly_raw: &mut Vec<Vec<f64>>,
    grow_factor: &f64,
    shrink_factor: &f64,
) {
    /* ------------------------------------------------------------------
       Pick a random number of games to play – this part stays single‑threaded
    ------------------------------------------------------------------ */

    let games_this_week: usize = rng.random_range(5..=50);
    *total_games += games_this_week;

    /* ------------------------------------------------------------------
       Pick a random number of confirmed tiles (20‑50) ----------
    ------------------------------------------------------------------ */
    let confirmed_count = rng.random_range(20..=50);
    let confirmed_set: HashSet<usize> = random_sample(rng, confirmed_count, tiles.len());

    /* --------------------------------------------------------------------
       Parallel part – the heavy lifting happens here
    -------------------------------------------------------------------- */
    // we create a local copy of the tiles/weights that is read‑only
    let tiles_local: Vec<Tile> = tiles.iter().cloned().collect();
    let weights_local: Vec<f64> = weights.clone();

    // The `into_par_iter()` turns the range into a Rayon iterator.
    // Every iteration gets its own RNG instance – we use the same
    // helper that you already have (`rand::rng()`) so that the
    // behaviour stays identical to the single‑thread version.
    let game_results: Vec<(f64, f64)> = (0..games_this_week)
        .into_par_iter()
        .map(|_| {
            let mut rng_local = rand::rng();
            simulate_game(&mut rng_local, &tiles_local, &confirmed_set, &weights_local)
        })
        .collect();

    /* ------------------------------------------------------------------
       Collect the per‑week statistics from the parallel results
    ------------------------------------------------------------------ */
    // local vectors to avoid contention while collecting
    let mut week_final_scores: Vec<f64> = Vec::with_capacity(games_this_week);
    let mut week_potential_maxes: Vec<f64> = Vec::with_capacity(games_this_week);

    for (final_score, max_potential) in &game_results {
        week_final_scores.push(*final_score);
        week_potential_maxes.push(*max_potential);
    }

    // extend the global vectors – this is the *only* mutable operation
    final_scores.extend(&week_final_scores);
    potential_maxes.extend(&week_potential_maxes);

    /* ---- Print weekly mean -------------------------------------- */
    let week_mean_final: f64 =
        week_final_scores.iter().copied().sum::<f64>() / games_this_week as f64;

    // wins are simply the number of non‑zero final scores
    let week_wins: usize = week_final_scores.iter().filter(|&&s| s > 0.0).count();
    let week_losses: usize = games_this_week - week_wins;
    *total_wins += week_wins;
    *total_losses += week_losses;

    // println!(
    //     "Week {:03}: G: {:02} | WIN: {:02} | LOS: {:02} | MFS: {:06.2}",
    //     w + 1,
    //     games_this_week,
    //     week_wins,
    //     week_losses,
    //     week_mean_final,
    // );

    /* ---- Update tile weights for the master list (after week) ----- */
    for (idx, t) in tiles.iter_mut().enumerate() {
        if confirmed_set.contains(&idx) {
            // grow
            let new_wgt = t.y * grow_factor;
            weights[idx] = new_wgt.min(1.0);
        } else {
            // shrink
            let new_wgt = t.y * shrink_factor;
            weights[idx] = new_wgt.max(0.0);
        }
    }

    weekly_means.push(week_mean_final);
    weekly_raw.push(week_final_scores.clone()); // keep a copy for the box‑plot
}

fn simulate_game(
    rng: &mut impl Rng,
    tiles: &Vec<Tile>,
    confirmed_set: &HashSet<usize>,
    weights: &Vec<f64>,
) -> (f64, f64) {
    /* ---- Generate a fresh board of 25 random tiles from the
    master list ------------------------------------------ */
    let board_indices: HashSet<usize> = random_sample(rng, 25, tiles.len());

    // clone the selected tiles – this copies the current weights
    let board: Vec<Tile> = board_indices.iter().map(|&i| tiles[i].clone()).collect();

    /* ---- Build a boolean mask that tells whether each tile
    on the board is confirmed for *this* week ---------- */
    let mut confirm_mask = [false; 25];
    for &idx in confirmed_set {
        if idx == 12 {
            // Simulate a free tile right in the center of the board
            confirm_mask[idx] = true;
            continue;
        }
        if idx < board.len() {
            // only tiles actually used
            confirm_mask[idx] = true;
        }
    }

    let win_found = check_win(&confirm_mask);

    /* ---- Compute weighted sum ----------------------------------- */
    let mut final_score: f64 = 0.0;
    if win_found {
        for (idx, tile) in board.iter().enumerate() {
            if confirm_mask[idx] {
                // use the current weight of this tile
                let wgt = weights[tiles.iter().position(|t| t.id == tile.id).unwrap()];
                final_score += tile.x * wgt;
            }
        }
    }

    /* ---- Max potential (unchanged) ------------------------------ */
    let max_potential: f64 = board
        .iter()
        .enumerate()
        .map(|(_idx, t)| {
            let wgt = weights[tiles.iter().position(|x| x.id == t.id).unwrap()];
            t.x * wgt
        })
        .sum();

    return (final_score, max_potential);
}

fn check_win(confirm_mask: &[bool; 25]) -> bool {
    /* ---- Check whether any winning line (row/column/diagonal)
    consists entirely of confirmed tiles ------------------- */

    // rows
    for r in 0..5 {
        if (r * 5..r * 5 + 5).all(|i| confirm_mask[i]) {
            return true;
        }
    }

    // columns – only test further if still not found
    for c in 0..5 {
        let mut col_ok = true;
        for r in 0..5 {
            if !confirm_mask[r * 5 + c] {
                col_ok = false;
                break;
            }
        }
        if col_ok {
            return true;
        }
    }

    // Check diagonals for win conditions
    check_diagonals(confirm_mask)
}

fn check_diagonals(confirm_mask: &[bool; 25]) -> bool {
    // main diagonal
    let mut diag_ok = true;
    for i in &[0, 6, 12, 18, 24] {
        if !confirm_mask[*i] {
            diag_ok = false;
            break;
        }
    }
    if diag_ok {
        return true;
    }

    // anti‑diagonal
    let mut adiag_ok = true;
    for i in &[4, 8, 12, 16, 20] {
        if !confirm_mask[*i] {
            adiag_ok = false;
            break;
        }
    }
    if adiag_ok {
        return true;
    }

    false
}

// fn memory_bar_updater(pb: Arc<ProgressBar>, done: Arc<AtomicBool>) {
//     let mut sys = System::new_all(); // the lib we need
//     let pid = Pid::from_u32(std::process::id());
//     while !done.load(Ordering::Relaxed) {
//         sys.refresh_processes(ProcessesToUpdate::Some(&[pid]), true); // keep data current
//         if let Some(p) = sys.process(pid) {
//             let used_mb = p.virtual_memory() / 1024 / 1024; // KiB → MiB
//
//             if used_mb > 1000 {
//                 let used_gb = used_mb as f64 / 1024.0;
//                 pb.set_message(format!("{:>6.2} GiB", used_gb));
//                 continue;
//             }
//
//             pb.set_message(format!("{:>6} MiB", used_mb));
//         } else {
//             pb.set_message("??");
//         }
//         thread::sleep(Duration::from_secs(1)); // once per second
//     }
// }
