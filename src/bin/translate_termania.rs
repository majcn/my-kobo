use my_kobo::common::*;

use std::collections::HashMap;

use scraper::Element;
use scraper::ElementRef;
use scraper::Html;
use scraper::Selector;

const BASE_URL: &str = "https://www.termania.net/iskanje";

fn get_details_key(element: &ElementRef) -> Option<String> {
    let selector = Selector::parse(".content span").unwrap();
    let content = element.select(&selector).next()?;

    let raw_result = content.text().next().unwrap();
    let result = raw_result
        .chars()
        .filter(|x| !['\n', '\t', '(', ')'].contains(x))
        .collect::<String>();

    Some(result)
}

fn get_details_value(element: &ElementRef) -> Option<Vec<String>> {
    let selector = Selector::parse(".lang_sl").unwrap();
    let content = element.select(&selector).next()?;

    let mut result = vec![];

    let mut element_option = content.next_sibling_element();
    while let Some(element) = element_option {
        if element.value().name() != "strong" {
            break;
        }

        result.push(element.inner_html());

        element_option = element.next_sibling_element();
    }

    Some(result)
}

fn get_details_part(source: &str, element: &ElementRef) -> Option<(String, Vec<String>)> {
    let element = element.parent_element()?;

    let source_from_page = element
        .select(&Selector::parse("h4").unwrap())
        .next()?
        .inner_html();
    if source != source_from_page {
        return None;
    }

    Some((get_details_key(&element)?, get_details_value(&element)?))
}

fn get_details(source: &str, document: &Html) -> HashMap<String, Vec<String>> {
    let mut result = HashMap::new();

    let selector = Selector::parse(".lang_en.headword").unwrap();
    for element in document.select(&selector) {
        if let Some((key, value)) = get_details_part(source, &element) {
            result.insert(key, value);
        }
    }

    result
}

pub fn translate(source: &str) -> Result<TranslateResult, TranslateError> {
    const SOURCE_LANG: &str = "2"; // en
    const TARGET_LANG: &str = "61"; // sl

    let query_pairs = vec![
        ("SearchIn", "Linked"),
        ("ld", "70"),
        ("sl", SOURCE_LANG),
        ("tl", TARGET_LANG),
        ("query", source),
    ];

    let response = ureq::get(BASE_URL)
        .query_pairs(query_pairs)
        .call()
        .map_err(|_| TranslateError::NETWORK)?;

    let raw_result = response
        .into_string()
        .map_err(|_| TranslateError::NETWORK)?;

    let document = Html::parse_document(&raw_result);

    let translation = String::from(
        document
            .select(&Selector::parse(".lang_sl").unwrap())
            .next()
            .and_then(|x| x.next_sibling_element())
            .and_then(|x| x.text().next())
            .unwrap(),
    );

    let details = get_details(source, &document);

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
