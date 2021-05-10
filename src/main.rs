mod cli;
use cli::CLI;
mod days;
use days::Days;

use std::env;
use std::process;

fn main() {
    let args = CLI::from_args(&mut env::args_os()).unwrap_or_else(|e| {
        eprintln!("{}", e);
        process::exit(1);
    });

    println!("CLI: {:#?}", args);

    let days = Days::from_offset(args.week_offset);
    println!("week number: {}", days.week_num);
    for day in &days {
        println!("  {}", day);
    }
    for day in &days {
        println!("  {}", day);
    }
}
