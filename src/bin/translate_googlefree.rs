use my_kobo::common::*;

use ureq::serde_json::Value;

const BASE_URL: &str = "https://translate.googleapis.com/translate_a/single";

fn get_result<'a>(json_value: &'a Value, path: &[usize]) -> Option<&'a Value> {
    let mut result = json_value;
    for p in path {
        result = result.get(p)?
    }
    Some(result)
}

pub fn translate(source: &str) -> Result<TranslateResult, TranslateError> {
    const SOURCE_LANG: &str = "en";
    const TARGET_LANG: &str = "sl";

    let query_pairs = vec![
        ("client", "gtx"),
        ("dt", "t"),
        ("sl", SOURCE_LANG),
        ("tl", TARGET_LANG),
        ("q", source),
    ];

    let response = ureq::get(BASE_URL)
        .query_pairs(query_pairs)
        .call()
        .map_err(|_| TranslateError::NETWORK)?;

    let json_value = response
        .into_json::<Value>()
        .map_err(|_| TranslateError::NETWORK)?;

    let translation = get_result(&json_value, &[0, 0, 0]).ok_or(TranslateError::JSON)?;
    let translation = translation.as_str().ok_or(TranslateError::JSON)?;
    let translation = String::from(translation);

    Ok(TranslateResult::simple(translation))
}

fn main() {
    if let Some(source) = std::env::args().nth(1) {
        let text = match translate(&source) {
            Ok(v) => v.to_string(),
            Err(e) => e.to_string(),
        };

        println!("{}", textwrap::fill(&text, TEXT_WRAP_SIZE));
    }
}
