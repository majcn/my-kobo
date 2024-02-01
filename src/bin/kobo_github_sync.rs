use std::fs;
use std::io::Read;
use std::io::{self, Cursor};
use std::process::exit;

fn download(url: &str) -> Result<Vec<u8>, ureq::Error> {
    let resp = ureq::get(url).call()?;

    let mut reader = io::BufReader::new(resp.into_reader());
    let mut bytes = vec![];
    reader.read_to_end(&mut bytes)?;

    Ok(bytes)
}

pub fn unzip_to_kobo_sdcard(bytes: Vec<u8>) -> Result<(), zip::result::ZipError> {
    let mut archive = zip::ZipArchive::new(Cursor::new(bytes))?;

    for i in 0..archive.len() {
        let mut file = archive.by_index(i)?;
        let outpath = match file.enclosed_name() {
            Some(path) => path.to_owned(),
            None => continue,
        };

        let outpath = match outpath.strip_prefix("my-kobo-master/sdcard") {
            Ok(path) => path,
            Err(_) => continue,
        };

        if !outpath.starts_with("mnt/onboard/.adds") {
            continue;
        }

        if (*file.name()).ends_with('/') {
            println!("File {} extracted to \"{}\"", i, outpath.display());
            fs::create_dir_all(&outpath)?;
        } else {
            println!(
                "File {} extracted to \"{}\" ({} bytes)",
                i,
                outpath.display(),
                file.size()
            );
            if let Some(p) = outpath.parent() {
                if !p.exists() {
                    fs::create_dir_all(p)?;
                }
            }
            let mut outfile = fs::File::create(&outpath)?;
            io::copy(&mut file, &mut outfile)?;
        }

        // Get and Set permissions
        #[cfg(unix)]
        {
            use std::os::unix::fs::PermissionsExt;

            if let Some(mode) = file.unix_mode() {
                fs::set_permissions(outpath, fs::Permissions::from_mode(mode))?;
            }
        }
    }

    Ok(())
}

fn main() {
    const ARCHIVE_URL: &str = "https://github.com/majcn/my-kobo/archive/refs/heads/master.zip";

    let zip_bytes = match download(ARCHIVE_URL) {
        Ok(bytes) => bytes,
        Err(_) => {
            println!("Error while downloading from github.");
            exit(1);
        }
    };

    match unzip_to_kobo_sdcard(zip_bytes) {
        Ok(_) => {}
        Err(_) => {
            println!("Error while extracting files.");
            exit(1);
        }
    }
}
