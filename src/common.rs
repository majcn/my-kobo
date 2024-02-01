use std::collections::HashMap;

pub const TEXT_WRAP_SIZE: usize = 40;

#[derive(Debug)]
pub enum TranslateError {
    NETWORK,
    JSON,
}

impl ToString for TranslateError {
    fn to_string(&self) -> String {
        match self {
            TranslateError::NETWORK => {
                String::from("Something went wrong when accessing internet.")
            }
            TranslateError::JSON => String::from("Something went wrong while parsing content."),
        }
    }
}

#[derive(Debug)]
pub struct TranslateResult {
    translation: String,
    details: HashMap<String, Vec<String>>,
}

impl TranslateResult {
    pub fn simple(translation: String) -> Self {
        Self {
            translation,
            details: HashMap::new(),
        }
    }

    pub fn detail(translation: String, details: HashMap<String, Vec<String>>) -> Self {
        Self {
            translation,
            details,
        }
    }
}

impl ToString for TranslateResult {
    fn to_string(&self) -> String {
        let mut result = String::new();

        result.push_str(&self.translation);
        result.push('\n');
        result.push('\n');

        for (k, v) in &self.details {
            result.push_str(k);
            result.push('\n');
            for vv in v {
                result.push_str("  ");
                result.push_str(vv);
                result.push('\n');
            }
            result.push('\n');
        }

        result
    }
}
