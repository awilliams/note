# note

A tool for managing markdown note files.

Notes are organized by year and week.
Editing is done using the configured `$EDITOR`.
New note files are created with a basic markdown template containing entries for Monday through Friday.

## Motivation

The author wanted an easy way to keep weekly markdown notes of things happening at work.

## Features

- ❌ Search
- ❌ Editor
- ❌ Backups
- ❌ Web Frontend
- ❌ GUI
- ✅ CLI
- ✅ Note file management

Most functionality, e.g. search, is outside the scope of this tool.
There are other tools, e.g. grep, that are better suited.

## Default Note Layout

Note files are stored under `$HOME/note` by default (this can be changed).
Notes are then stored under their [ISO 8601](https://golang.org/pkg/time/#Time.ISOWeek) year, then week number.

Example layout:
```
$HOME/note
├── 2020
│   ├── 25.md
│   └── 26.md
└── 2021
    ├── 01.md
    ├── 02.md
    └── 03.md
```

## Usage/Examples

* Open this week's note file:
  ```shell
  $ note
  $ cat ~/note/2021/16.md
  ### Monday, 19-April-2021

  ### Tuesday, 20-April-2021

  ### Wednesday, 21-April-2021

  ### Thursday, 22-April-2021

  ### Friday, 23-April-2021
  ```

* Open note file from two weeks ago:
  ```shell
  $ note -2
  ```

* Open next week's note file:
  ```shell
  $ note +1
  ```

* Search last week's note for TODOs:
  ```shell
  $ grep "TODO" $(note -1 -n)
  ```

* Change directory where the notes are stored (default is `$HOME/note`).
  This can be changed in two ways:
  * Using the `-d` flag:
      ```shell
      $ # Notes created under ~/other
      $ note -d ~/other
      ```
  * Changing the name of the executable:
      ```shell
      $ ln -s $(which note) ~/bin/work-notes
      $ # Notes created under $HOME/work-notes
      $ work-notes
      ```

## License

[MIT](/LICENSE)
