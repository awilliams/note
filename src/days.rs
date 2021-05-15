use std::iter::Iterator;

use chrono::Duration;
use chrono::{prelude::*, IsoWeek};

pub struct DayRange {
    first: Date<Local>,
    iso_week: IsoWeek,
}

impl DayRange {
    pub fn from_monday(rel_week_offset: i64) -> Self {
        let monday = Self::monday(rel_week_offset);
        Self::from(monday)
    }

    pub fn from(first: Date<Local>) -> Self {
        Self {
            first: first,
            iso_week: first.iso_week(),
        }
    }

    fn monday(rel_week_offset: i64) -> Date<Local> {
        let today: Date<Local> = Local::today();
        let monday_offset = Duration::days(1 - today.weekday().num_days_from_sunday() as i64);
        let offset = monday_offset + Duration::weeks(rel_week_offset);
        today + offset
    }

    pub fn range(&self, num_days: i64) -> DayRangeIterator {
        DayRangeIterator::new(self.first, num_days)
    }

    pub fn week_num(&self) -> u32 {
        self.iso_week.week()
    }

    pub fn year(&self) -> i32 {
        return self.iso_week.year();
    }
}

pub struct DayRangeIterator {
    first: Date<Local>,
    days: i64,
}

impl DayRangeIterator {
    fn new(first: Date<Local>, days: i64) -> Self {
        DayRangeIterator {
            first: first,
            days: days,
        }
    }
}

impl Iterator for DayRangeIterator {
    type Item = Date<Local>;
    fn next(&mut self) -> Option<Date<Local>> {
        if self.days <= 0 {
            return None;
        }
        let ret = self.first;
        self.first = self.first + Duration::days(1);
        self.days -= 1;
        Some(ret)
    }
}

#[test]
fn test_day_range() {
    let n = 5;

    let r = DayRange::from(Local.ymd(2021, 10, 10));
    assert_eq!(r.week_num(), 40);
    assert_eq!(r.year(), 2021);

    let days: Vec<Date<Local>> = r.range(n).collect();
    assert_eq!(days.len(), n as usize);

    for i in 0..n {
        assert_eq!(days[i as usize], Local.ymd(2021, 10, 10 + i as u32));
    }
}
