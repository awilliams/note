use std::slice::Iter;

use chrono::prelude::*;
use chrono::Duration;

pub struct Days {
    pub year: i32,
    pub week_num: u32,
    pub days: Vec<Date<Local>>,
}

impl Days {
    pub fn from_offset(week_offset: i64) -> Self {
        let monday = Self::monday(week_offset);

        let mut days = Vec::with_capacity(5);
        days.push(monday);
        for i in 1..=4 {
            days.push(monday + Duration::days(i));
        }

        let iso_week = monday.iso_week();
        Self {
            year: iso_week.year(),
            week_num: iso_week.week(),
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

impl<'a> IntoIterator for &'a Days {
    type Item = &'a Date<Local>;
    type IntoIter = Iter<'a, Date<Local>>;

    fn into_iter(self) -> Self::IntoIter {
        self.days.iter()
    }
}
