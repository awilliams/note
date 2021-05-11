mod cli;
mod days;

use cli::CLI;
use days::Days;

use std::fs;
use std::process;
use std::{env, io::Write};

use home;

fn main() {
    let args = CLI::from_args(&mut env::args_os()).unwrap_or_else(|e| {
        bail(&e);
    });

    if args.print_md {
        bail("print is not implemented!");
    }

    let days = Days::from_offset(args.week_offset);

    let p = match home::home_dir() {
        Some(mut p) => {
            p.push("TEST");
            p.push(args.exe_name);
            p.push(days.year.to_string());
            p.push(format!("{}.md", days.week_num));
            p
        }
        None => {
            bail("unable to determine $HOME directory");
        }
    };
    let note_path = p.as_path();

    if args.print_path {
        println!("{}", note_path.display());
        return;
    }

    // Create note directory.
    let dir = note_path.parent().unwrap();
    fs::create_dir_all(dir).unwrap_or_else(|e| {
        bail(&format!("unable to create note directory: {}", e));
    });

    if !note_path.exists() {
        let mut fd = std::fs::File::create(note_path).unwrap_or_else(|e| bail(&e.to_string()));

        fd.write_all(format!("# Week {}, {}\n---\n\n", days.week_num, days.year).as_bytes())
            .expect("cannot write to file");

        for day in &days {
            fd.write_all(format!("## {}\n\n", day.format("%A, %d-%b-%Y")).as_bytes())
                .expect("cannot write to file");
        }
        drop(fd);
    }

    let result = process::Command::new(args.editor).arg(note_path).status();
    result.unwrap_or_else(|e| {
        bail(&e.to_string());
    });
}

fn bail(msg: &str) -> ! {
    eprintln!("error: {}", msg);
    process::exit(1);
}
