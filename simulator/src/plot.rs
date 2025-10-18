use plotters::prelude::*;

/* ------------------------------------------------------------------
Histogram drawing helper
------------------------------------------------------------------ */

pub fn save_histogram(
    filename: &str,
    title: &str,
    data: &[f64],
    bin_width: f64,
) -> Result<(), Box<dyn std::error::Error>> {
    let root = BitMapBackend::new(filename, (1920, 1080)).into_drawing_area();
    root.fill(&WHITE)?;

    let min = data.iter().cloned().fold(f64::INFINITY, f64::min).floor();
    let max = data
        .iter()
        .cloned()
        .fold(f64::NEG_INFINITY, f64::max)
        .ceil();

    let bins = ((max - min) / bin_width).ceil() as usize;
    let mut hist = vec![0usize; bins];
    for &v in data {
        let idx = ((v - min) / bin_width).floor() as usize;
        if idx < bins {
            hist[idx] += 1;
        }
    }
    let max_count = *hist.iter().max().unwrap_or(&1);

    let mut chart = ChartBuilder::on(&root)
        .caption(title, ("sans-serif", 24))
        .margin(20)
        .x_label_area_size(40)
        .y_label_area_size(40)
        .build_cartesian_2d(min..max, 0..max_count)?;

    chart
        .configure_mesh()
        .x_desc("Score")
        .y_desc("Count")
        .draw()?;

    chart.draw_series(hist.iter().enumerate().map(|(i, &count)| {
        let x0 = min + i as f64 * bin_width;
        let x1 = x0 + bin_width;
        Rectangle::new([(x0, 0), (x1, count)], BLUE.mix(0.6).filled())
    }))?;

    root.present()?;
    println!("Saved histogram to {}", filename);
    Ok(())
}

/* ------------------------------------------------------------------
Line chart of weekly means
------------------------------------------------------------------ */
pub fn draw_weekly_scores(
    filename: &str,
    weeks: usize,
    weekly_means: &[f64],
) -> Result<(), Box<dyn std::error::Error>> {
    let root = BitMapBackend::new(filename, (1920, 1080)).into_drawing_area();
    root.fill(&WHITE)?;

    let y_min = weekly_means.iter().cloned().fold(f64::INFINITY, f64::min);
    let y_max = weekly_means
        .iter()
        .cloned()
        .fold(f64::NEG_INFINITY, f64::max);

    let mut chart = ChartBuilder::on(&root)
        .margin(20)
        .caption("Weekly mean final score", ("sans-serif", 24))
        .x_label_area_size(40)
        .y_label_area_size(40)
        .build_cartesian_2d(1f64..(weeks as f64 + 1.0), y_min..y_max)?;

    chart.configure_mesh().draw()?;

    // FIX: f64 x-values
    let line: Vec<(f64, f64)> = weekly_means
        .iter()
        .enumerate()
        .map(|(i, v)| (i as f64 + 1.0, *v))
        .collect();

    chart
        .draw_series(LineSeries::new(line, &RED))?
        .label("Mean")
        .legend(|(x, y)| PathElement::new(vec![(x, y), (x + 20, y)], &RED));

    chart.configure_series_labels().draw()?;
    root.present()?;
    Ok(())
}

/* ------------------------------------------------------------------
Box plot per week
------------------------------------------------------------------ */
pub fn draw_weekly_boxplot(
    filename: &str,
    weeks: usize,
    weekly_data: &[Vec<f64>],
) -> Result<(), Box<dyn std::error::Error>> {
    let root = BitMapBackend::new(filename, (1920, 1080)).into_drawing_area();
    root.fill(&WHITE)?;

    let y_min = weekly_data
        .iter()
        .flatten()
        .cloned()
        .fold(f64::INFINITY, f64::min);
    let y_max = weekly_data
        .iter()
        .flatten()
        .cloned()
        .fold(f64::NEG_INFINITY, f64::max);

    // FIX: f64 for x-axis
    let mut chart = ChartBuilder::on(&root)
        .margin(20)
        .caption("Weekly box-plots of final scores", ("sans-serif", 24))
        .x_label_area_size(40)
        .y_label_area_size(40)
        .build_cartesian_2d(1f64..(weeks as f64 + 1.0), y_min..y_max)?;

    chart.configure_mesh().draw()?;

    let box_width = 0.3;
    for (i, data) in weekly_data.iter().enumerate() {
        if data.is_empty() {
            continue;
        }

        let week = i as f64 + 1.0;
        let mut sorted = data.clone();
        sorted.sort_by(|a, b| a.partial_cmp(b).unwrap());

        let pct = |v: &[f64], p: f64| {
            let idx = ((v.len() - 1) as f64 * p).floor() as usize;
            v[idx]
        };

        let min_val = *sorted.first().unwrap();
        let max_val = *sorted.last().unwrap();
        let q1 = pct(&sorted, 0.25);
        let median = pct(&sorted, 0.5);
        let q3 = pct(&sorted, 0.75);

        chart.draw_series(std::iter::once(Rectangle::new(
            [(week - box_width, q1), (week + box_width, q3)],
            RED.mix(0.5).filled(),
        )))?;

        chart.draw_series(std::iter::once(PathElement::new(
            vec![(week - box_width, median), (week + box_width, median)],
            &BLACK,
        )))?;

        chart.draw_series(std::iter::once(PathElement::new(
            vec![(week, min_val), (week, max_val)],
            &BLACK,
        )))?;
    }

    root.present()?;
    Ok(())
}
