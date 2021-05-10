use clap::{App, AppSettings, Arg};
use std::ffi::OsString;

#[derive(Debug)]
pub struct CLI {
    pub exe_name: String,
    pub week_offset: i64,
    pub print: bool,
}

impl CLI {
    pub fn from_args<I, T>(args: I) -> Result<Self, String>
    where
        I: IntoIterator<Item = T>,
        T: Into<OsString> + Clone,
    {
        let mut app = App::new("note")
            .version("0.1.0")
            .about("weekly note manager")
            .setting(AppSettings::AllowNegativeNumbers)
            .setting(AppSettings::ColorNever)
            .arg(
                Arg::with_name("OFFSET")
                    .help("Week offset")
                    .long_help("Week offset. Examples: '-2' for two weeks ago or '1' for next week")
                    .required(true)
                    .default_value("0")
                    .index(1),
            )
            .arg(
                Arg::with_name("print")
                    .short("p")
                    .long("print")
                    .takes_value(false)
                    .help("Print note to STDOUT"),
            );

        // Retain 'app' in order to get the binary name.
        let matches = app.get_matches_from_safe_borrow(args);

        let arg_matches = match matches {
            Err(e) => return Err(e.message),
            Ok(m) => m,
        };

        // Helper to construct an error message with
        // usage information appended.
        let err_msg = |msg: &str| -> String {
            let mut buf = msg.to_string();
            buf.push_str("\n");
            buf.push_str(arg_matches.usage());
            buf
        };

        let cli = CLI {
            exe_name: app.get_bin_name().unwrap_or("note").to_string(),
            print: arg_matches.is_present("print"),
            week_offset: match arg_matches.value_of("OFFSET") {
                Some(v) => match v.parse::<i64>() {
                    Ok(i) => i,
                    Err(e) => {
                        return Result::Err(err_msg(
                            format!("Unable to parse OFFSET: {}", e).as_str(),
                        ))
                    }
                },
                None => return Result::Err(err_msg("OFFSET must be set")),
            },
        };

        Result::Ok(cli)
    }
}
