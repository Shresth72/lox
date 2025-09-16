use clap::Parser;
use clap::Subcommand;
use lox::error::SingleTokenError;
use lox::error::StringTerminationError;
use miette::{IntoDiagnostic, WrapErr};
use std::fs;
use std::path::PathBuf;

use lox::*;

#[derive(Parser, Debug)]
#[command(version, about, long_about = None)]
struct Args {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand, Debug)]
enum Commands {
    Tokenize { filename: PathBuf },
}

fn main() -> miette::Result<()> {
    let args = Args::parse();
    let mut any_cc_error = false;

    match args.command {
        Commands::Tokenize { filename } => {
            let file_contents = fs::read_to_string(&filename)
                .into_diagnostic()
                .wrap_err_with(|| format!("Reading '{}' failed", filename.display()))?;

            let lexer = Lexer::new(&file_contents);
            for token in lexer {
                let token = match token {
                    Ok(t) => t,
                    Err(e) => {
                        eprintln!("{e:?}");
                        if let Some(unrecognized) = e.downcast_ref::<SingleTokenError>() {
                            any_cc_error = true;
                            eprintln!(
                                "[line {}] Error: Unexpected character: {}",
                                unrecognized.line(),
                                unrecognized.token
                            );
                        } else if let Some(unterminated) =
                            e.downcast_ref::<StringTerminationError>()
                        {
                            any_cc_error = true;
                            eprintln!("[line {}] Error: Unterminated string", unterminated.line());
                        }
                        continue;
                    }
                };
                println!("{token}");
            }
            println!("EOF  null");
        }
    }

    if any_cc_error {
        std::process::exit(65);
    }

    Ok(())
}
