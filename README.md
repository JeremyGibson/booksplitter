# Booksplitter

This is mostly a learning project, so that I can get to know Go, but it was driven by a desire
to have a tool that could split `.m4b` files faster than some of the existing projects out there.

The goal is to have a command that can be passed either a file path or a directory
path, and process any given or found `.m4b` files into chapter files.

### Requirements

* Posix OS
* FFMpeg >= 5.1.2

### Usage

```shell
$> ./booksplitter -p /some/path/to/an/audiobook.m4b
```


### Current State

The single file extraction works, but with a few magic assumptions that
need to be parameterized or set up with a config

1. All extracted books go to a set path `/data/Audiobooks/extracted`
2. All extracted files are converted to `.m4a` with a copy of the underlying audio bitstream, which in most cases is `aac`
    * This is fast
    * This is not flexible

### TODO
* Build out config or use environment variables to remove magic values
* Enable recursive discovery and extraction
* Add error handling for potentially broken metadata specifically chapters and title information.