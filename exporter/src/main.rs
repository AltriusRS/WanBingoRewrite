use anyhow::Result;
use anyhow::Result;
use google_sheets4::{
    api::ValueRange,
    hyper, hyper_rustls,
    oauth2::{read_service_account_key, ServiceAccountAuthenticator, ServiceAccountKey},
    Sheets,
};
use google_sheets4::{
    api::ValueRange,
    hyper, hyper_rustls,
    oauth2::{read_service_account_key, ServiceAccountAuthenticator, ServiceAccountKey},
    Sheets,
};
use std::{fs, path::Path};
use std::{fs, path::Path};
use tokio;
use tokio;

#[tokio::main]
async fn main() -> Result<()> {
    // --- CONFIG ---
    let creds_path = "service_account.json";
    let spreadsheet_id = "YOUR_SPREADSHEET_ID_HERE";
    let output_dir = "exports";

    fs::create_dir_all(output_dir)?;

    // --- AUTHENTICATION ---
    let sa_key: ServiceAccountKey = read_service_account_key(creds_path).await?;
    let auth = ServiceAccountAuthenticator::builder(sa_key).build().await?;

    // --- HTTPS client ---
    let connector = hyper_rustls::HttpsConnectorBuilder::new()
        .with_native_roots()
        .https_or_http()
        .enable_http1()
        .build();

    let client = hyper::Client::builder().build(connector);
    let hub = Sheets::new(client, auth);

    // --- FETCH SHEET METADATA ---
    let meta = hub.spreadsheets().get(spreadsheet_id).doit().await?;
    let spreadsheet = meta.1;
    let title = spreadsheet
        .properties
        .as_ref()
        .and_then(|p| p.title.clone())
        .unwrap_or_else(|| "Untitled".to_string());
    println!("Exporting spreadsheet: {}", title);

    // --- ITERATE SHEETS ---
    if let Some(sheets) = spreadsheet.sheets {
        for sheet in sheets {
            if let Some(props) = sheet.properties {
                let name = props.title.unwrap_or_else(|| "Untitled".to_string());
                println!("→ Exporting tab: {}", name);

                // Fetch all values in this sheet
                let range = format!("'{}'!A:ZZ", name);
                let (_, data): (_, ValueRange) = hub
                    .spreadsheets()
                    .values_get(spreadsheet_id, &range)
                    .doit()
                    .await?;

                if let Some(values) = data.values {
                    let out_path = Path::new(output_dir).join(format!("{}.csv", sanitize(&name)));
                    let mut wtr = csv::Writer::from_path(&out_path)?;

                    for row in values {
                        // Convert serde_json::Value → String
                        let row_strs: Vec<String> =
                            row.into_iter().map(|v| v.to_string()).collect();
                        wtr.write_record(&row_strs)?;
                    }

                    wtr.flush()?;
                    println!("   ✅ Saved: {}", out_path.display());
                } else {
                    println!("   (Empty sheet)");
                }
            }
        }
    }

    println!("\nAll sheets exported successfully to '{}'.", output_dir);
    Ok(())
}

/// Replace problematic characters for filenames
fn sanitize(name: &str) -> String {
    name.chars()
        .map(|c| {
            if c.is_ascii_alphanumeric() || c == '-' || c == '_' {
                c
            } else {
                '_'
            }
        })
        .collect()
}

