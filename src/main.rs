mod cli;
mod days;

use cli::CLI;
use days::DayRange;

use std::env;
use std::fs;
use std::io::Write;
use std::process;

use home;

fn main() {
    let args = CLI::from_args(&mut env::args_os()).unwrap_or_else(|e| {
        bail(&e);
    });

    if args.print_md {
        bail("print not yet implemented!");
    }

    let date_range = DayRange::from_monday(args.week_offset);

    let p = match home::home_dir() {
        Some(mut p) => {
            p.push("TEST");
            p.push(args.exe_name);
            p.push(date_range.year().to_string());
            p.push(format!("{}.md", date_range.week_num()));
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
        // Note doesn't exist. Create it and write default template.

        let mut fd = std::fs::File::create(note_path).unwrap_or_else(|e| bail(&e.to_string()));

        write!(
            fd,
            "# Week {}, {}\n---\n\n",
            date_range.week_num(),
            date_range.year()
        )
        .unwrap_or_else(|e| bail(&e.to_string()));

        for day in date_range.range(5) {
            write!(fd, "## {}\n\n", day.format("%A, %d-%b-%Y"))
                .unwrap_or_else(|e| bail(&e.to_string()));
        }

        fd.flush().unwrap_or_else(|e| bail(&e.to_string()));
    }

    // Open note in $EDITOR.
    let result = process::Command::new(args.editor).arg(note_path).status();
    result.unwrap_or_else(|e| {
        bail(&e.to_string());
    });
}

// print error to STDERR and exit process.
fn bail(msg: &str) -> ! {
    eprintln!("error: {}", msg);
    process::exit(1);
}
