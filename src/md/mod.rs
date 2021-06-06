use std::fmt::Write;
use std::io;

use ansi_term::{Colour, Style};
use pulldown_cmark::{CowStr, Event, Options, Parser, Tag};

const LINE_BREAK: pulldown_cmark::CowStr = pulldown_cmark::CowStr::Borrowed("\n");

pub fn print_anscii_md(input: &String, output: &mut String) -> io::Result<()> {
    MD::new(input, output).write()
}

struct MD<'a> {
    iter: Parser<'a>,
    output: &'a mut String,
    cur_list: Option<Tag<'a>>,
    cur_link: Option<Tag<'a>>,
    style: Style,
    indentation: u8,
}

impl<'a> MD<'a> {
    fn new(input: &'a String, output: &'a mut String) -> Self {
        let mut options = Options::empty();
        options.insert(Options::ENABLE_STRIKETHROUGH);
        options.insert(Options::ENABLE_TASKLISTS);
        let parser = Parser::new_ext(input, options);
        MD {
            iter: parser,
            output: output,
            cur_list: None,
            cur_link: None,
            style: Style::new(),
            indentation: 0,
        }
    }

    fn write(&mut self) -> io::Result<()> {
        while let Some(event) = self.iter.next() {
            match event {
                Event::Start(tag) => {
                    self.start_tag(&tag);
                }
                Event::End(tag) => {
                    self.end_tag(&tag);
                }
                Event::Text(text) => {
                    if let Some(Tag::Link(_, dest, _)) = &self.cur_link {
                        let url = dest.clone().into_string();
                        let prev_style = self.style;
                        self.style = Style::new().fg(Colour::Blue);
                        self.text(&text);
                        self.style = prev_style;
                        if !url.starts_with("#") {
                            self.text(" (");
                            self.style.is_underline = true;
                            self.text(format!("{}", url).as_str());
                            self.style.is_underline = false;
                            self.text(")");
                        }
                    } else {
                        self.text(&text);
                    }
                }
                Event::Code(text) => {
                    self.style = Style::new().on(Colour::Fixed(239)).fg(Colour::White);
                    self.text(&text);
                    self.style = Style::new();
                }
                Event::Html(html) => {
                    writeln!(self.output, "?html: {}", html);
                }
                Event::SoftBreak => {
                    writeln!(self.output, "");
                }
                Event::FootnoteReference(name) => {
                    writeln!(self.output, "footnote: {}", name);
                }
                Event::HardBreak => {
                    self.newline(2);
                }
                Event::Rule => {
                    self.text("⎼⎼⎼");
                    self.newline(1);
                }
                Event::TaskListMarker(true) => {
                    self.text("[✓] ");
                }
                Event::TaskListMarker(false) => {
                    self.text("[ ] ");
                }
            }
        }
        Ok(())
    }

    fn start_tag(&mut self, tag: &Tag<'a>) {
        match tag {
            Tag::Heading(level) => {
                self.newline(1);
                let mut hdr = "#".repeat(*level as usize);
                hdr += " ";
                self.text(&hdr);
                self.style = Style::new().on(Colour::White).fg(Colour::Black).bold();
            }
            Tag::Item => match self.cur_list {
                Some(Tag::List(Some(i))) => {
                    self.text(&(format!("{}", i)));
                    self.cur_list = Some(Tag::List(Some(i + 1)));
                }
                Some(Tag::List(None)) => {
                    self.text("•");
                }
                _ => {}
            },
            Tag::List(_) => {
                self.cur_list = Some(tag.clone());
                self.indentation += 2;
                self.newline(1);
            }
            Tag::Paragraph => {}
            Tag::Link(_, dest, _) => {
                self.cur_link = Some(tag.clone());
            }
            Tag::BlockQuote => {
                self.indentation += 3;
                self.style = Style::new().fg(Colour::Fixed(244));
            }
            Tag::CodeBlock(_) => {
                self.indentation += 4;
                self.style = Style::new().fg(Colour::Green);
            }
            Tag::Strong => {
                self.style = self.style.bold();
            }
            Tag::Emphasis => {
                self.style = self.style.underline();
            }
            Tag::Strikethrough => {
                self.style = self.style.strikethrough();
            }
            _ => {}
        }
    }

    fn end_tag(&mut self, tag: &Tag) {
        match tag {
            Tag::Heading(level) => {
                self.style = Style::new();
                self.newline(2);
            }
            Tag::Item => {
                self.newline(1);
            }
            Tag::List(_) => {
                self.cur_list = None;
                self.indentation -= 2;
                self.newline(1);
            }
            Tag::Paragraph => {
                self.newline(2);
            }
            Tag::Link(_, _, _) => {
                self.cur_link = None;
            }
            Tag::BlockQuote => {
                self.indentation -= 3;
                self.style = Style::new();
            }
            Tag::Strong => {
                self.style.is_bold = false;
            }
            Tag::Emphasis => {
                self.style.is_underline = false;
            }
            Tag::Strikethrough => {
                self.style.is_strikethrough = false;
            }
            Tag::CodeBlock(_) => {
                self.indentation -= 4;
                self.style = Style::new();
            }
            _ => {}
        }
    }

    fn newline(&mut self, count: u32) {
        for i in 0..count {
            writeln!(self.output);
        }
    }

    fn text(&mut self, text: &str) {
        match self.indentation {
            i if i > 0 => {
                let prefix = " ".repeat(i.into());
                text.split_inclusive("\n").for_each(|l| {
                    write!(self.output, "{}", prefix);
                    write!(self.output, "{}", self.style.paint(l.to_string()));
                })
            }
            _ => {
                write!(self.output, "{}", self.style.paint(text.to_string()));
            }
        };
    }
}

#[test]
fn test_md() {
    let content = String::from(
        "# Week 19, 2021
---

## Monday, 10-May-2021
- Here is some text from a list
- Second item
  - Indented list
  - Another indented item

  ```
  code block here
  second line of code blcok
  ```

OK _what_ is **that**
Here is a [link](example.com)

## Tuesday, 11-May-2021

## Wednesday, 12-May-2021

## Thursday, 13-May-2021

## Friday, 14-May-2021


TODO:
- [ ] Do this
- [x] Done!
",
    );

    let mut out = String::new();
    let mut md = MD::new(&content);
    md.write(&mut out);
    println!("{}", out);
}
