mod cli;
mod day_range;
mod md;

use cli::CLI;
use day_range::DayRange;
use md::print_anscii_md;

use std::env;
use std::fs;
use std::io::{self, Write};
use std::process;

use home;

fn main() {
    let args = CLI::from_args(&mut env::args_os()).unwrap_or_else(|e| {
        bail(&e);
    });

    let date_range = DayRange::from_monday(args.week_offset);

    let mut note_path =
        home::home_dir().unwrap_or_else(|| bail("unable to determine $HOME directory"));
    note_path.push("TEST");
    note_path.push(args.exe_name);
    note_path.push(date_range.year().to_string());
    note_path.push(format!("{}.md", date_range.week_num()));
    let note_path = note_path.as_path();

    if args.print_path {
        println!("{}", note_path.display());
        return;
    }

    if args.print_md {
        let note_contents = fs::read_to_string(note_path).unwrap_or_else(|e| bail(&e.to_string()));
        let mut rendered = String::new();
        let md =
            print_anscii_md(&note_contents, &mut rendered).unwrap_or_else(|e| bail(&e.to_string()));
        io::copy(&mut rendered.as_bytes(), &mut io::stdout());
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
