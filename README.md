# m text viewer

The `m` program is a simple text mode, text file viewer, developed for my own use,
so don't expect too much.

The `m` viewer is designated specifically for viewing large text files or input streams, as viewed file is not fully loaded into memory, but are stored in a temporary file instead and read into memory in blocks as needed, not exceeding of specified total occupied memory size.

The maximum size of occupied memory and the block size can be specified by `-total` and `-block` parameters respecively.


> Note:
The viewer is designated to be used in Unix and/or Mac OS terminals but should also work
in Windows command prompt environment.

## Usage

```console
$ m -h
Program m is designated to view and browse flat, text files.
Usage:
	m <options> [file]
where <options> are:
	-h	help, shows this text
	-b	remove backspaces
		default: false
	-block	single data block size limit (MB)
		default: 4
	-t	title to show
	-total	total data size limit (MB)
		default: 64
Press h when browsing, to see list of available shortcuts.
Configuration file: ~/Library/Application Support/m/config.yaml
Copyright (C) 2020 Bartek Rybak (licensed under the MIT license).
```

When browsing press `h` or `F1` to display list of available keys and related actions.

## Dependencies

The program is dependent on [github.com/rivo/tview](https://github.com/rivo/tview)
package and its dependencies.

## License

Copyright (c) 2020, Bartek Rybak.  All rights reserved.  
Copyrights licensed under the **MIT License**.  
See the accompanying LICENSE file for terms.  
