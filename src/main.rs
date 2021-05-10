mod cli;

use std::env;
use std::process;

use std::vec::IntoIter;

use chrono::prelude::*;
use chrono::Duration;

fn main() {
    let args = cli::CLI::from_args(&mut env::args_os()).unwrap_or_else(|e| {
        eprintln!("{}", e);
        process::exit(1);
    });

    println!("CLI: {:#?}", args);

    let days = Days::from_offset(args.week_offset);
    println!("week number: {}", days.week_num);
    for day in days {
        println!("  {}", day);
    }
}

struct Days {
    week_num: u32,
    days: Vec<Date<Local>>,
}

impl Days {
    fn from_offset(week_offset: i64) -> Self {
        let monday = Self::monday(week_offset);

        let mut days = Vec::with_capacity(5);
        days.push(monday);
        for i in 1..=4 {
            days.push(monday + Duration::days(i));
        }

        Self {
            week_num: monday.iso_week().week(),
            days: days,
        }
    }

    fn monday(week_offset: i64) -> Date<Local> {
        let today: Date<Local> = Local::today();
        let monday_offset = Duration::days(1 - today.weekday().num_days_from_sunday() as i64);
        let offset = monday_offset + Duration::weeks(week_offset);
        today + offset
    }
}

impl IntoIterator for Days {
    type Item = Date<Local>;
    type IntoIter = IntoIter<Self::Item>;

    fn into_iter(self) -> Self::IntoIter {
        self.days.into_iter()
    }
}
