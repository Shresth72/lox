use miette::{Diagnostic, SourceSpan};
use thiserror::Error;

#[derive(Diagnostic, Debug, Error)]
#[error("Unexpected token '{token}' in input")]
pub struct SingleTokenError {
    #[source_code]
    pub src: String,

    pub token: char,

    #[label = "this input character"]
    pub err_span: SourceSpan,
}

impl SingleTokenError {
    pub fn line(&self) -> usize {
        let until_unrecog = &self.src[..=self.err_span.offset()];
        until_unrecog.lines().count()
    }
}
