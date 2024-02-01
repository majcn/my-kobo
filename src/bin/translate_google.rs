use my_kobo::common::*;

use std::collections::HashMap;

use ureq::serde_json::from_str as json_from_str;
use ureq::serde_json::Value;

const BASE_URL: &str = "https://translate.google.com";

fn get_raw_object_get_parameters<'a>() -> Result<Vec<(&'a str, String)>, TranslateError> {
    fn parse_value<'b>(result: &'b str, key: &str) -> &'b str {
        let key_index = result.find(key).unwrap();
        let start_index = key_index + key.len() + 3; // {key}":"
        let end_index = start_index + result[start_index..].find('\"').unwrap();

        &result[start_index..end_index]
    }

    let result = ureq::get(BASE_URL)
        .call()
        .map_err(|_| TranslateError::NETWORK)?
        .into_string()
        .map_err(|_| TranslateError::NETWORK)?;

    let query_pairs = vec![
        ("rpcids", String::from("MkEWBc")),
        ("f.sid", String::from(parse_value(&result, "FdrFJe"))),
        ("bl", String::from(parse_value(&result, "cfb2h"))),
        ("hl", String::from("en-US")),
        ("soc-app", String::from("1")),
        ("soc-platform", String::from("1")),
        ("soc-device", String::from("1")),
        ("_reqid", fastrand::u32(1000..10000).to_string()),
        ("rt", String::from("c")),
    ];

    Ok(query_pairs)
}

fn get_raw_object(
    source: &str,
    source_lang: &str,
    target_lang: &str,
) -> Result<String, TranslateError> {
    let query_pairs_string = get_raw_object_get_parameters()?;
    let query_pairs = query_pairs_string
        .iter()
        .map(|(a, b)| (*a, b.as_str()))
        .collect::<Vec<_>>();

    let mut url = String::from(BASE_URL);
    url.push_str("/_/TranslateWebserverUi/data/batchexecute");

    let body = format!(
        "[[[\"MkEWBc\",\"[[\\\"{}\\\",\\\"{}\\\",\\\"{}\\\",true],[null]]\",null,\"generic\"]]]",
        source, source_lang, target_lang
    );

    let raw_result = ureq::post(&url)
        .query_pairs(query_pairs)
        .send_form(&[("f.req", &body)])
        .map_err(|_| TranslateError::NETWORK)?
        .into_string()
        .map_err(|_| TranslateError::NETWORK)?;

    let length_length = raw_result[6..].find('\n').ok_or(TranslateError::JSON)?;
    let length = &raw_result[6..6 + length_length];
    let length = length.parse::<usize>().map_err(|_| TranslateError::JSON)?;

    let result = raw_result
        .chars()
        .skip(7 + length_length)
        .take(length - 2)
        .collect::<String>();

    Ok(result)
}

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

    let result =
        get_raw_object(source, SOURCE_LANG, TARGET_LANG).map_err(|_| TranslateError::NETWORK)?;

    let value = json_from_str::<Value>(&result).map_err(|_| TranslateError::JSON)?;
    let value = get_result(&value, &[0, 2])
        .and_then(Value::as_str)
        .unwrap_or_default();
    let value = json_from_str::<Value>(value).map_err(|_| TranslateError::JSON)?;

    let translation = String::from(
        get_result(&value, &[1, 0, 0, 5, 0, 0])
            .and_then(Value::as_str)
            .unwrap_or_default(),
    );

    let mut details: HashMap<String, Vec<String>> = HashMap::new();

    for x in get_result(&value, &[3, 5, 0])
        .and_then(Value::as_array)
        .unwrap_or(&vec![])
    {
        let translate_type = x.get(0).and_then(Value::as_str).unwrap_or_default();
        let entry = details.entry(String::from(translate_type)).or_default();

        if let Some(v) = x.get(1) {
            if let Some(va) = v.as_array() {
                let values_iter = va
                    .iter()
                    .filter_map(|x| x.get(0))
                    .filter_map(Value::as_str)
                    .map(String::from);
                entry.extend(values_iter);
            }
        }
    }

    Ok(TranslateResult::detail(translation, details))
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
