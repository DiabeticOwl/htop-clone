# htop-clone

![appPreview](https://i.imgur.com/OgaTFqS.gif)

htop-clone is a simplified version of the [htop](https://htop.dev/) command in Linux.

A description of the optimization process done over this application can be found
[here](appOptimization.md).

## Code Guide

The ui of this application can be found in the [ui.go](ui.go) and [ui-tables.go](ui-tables.go)
files. The packages used is [bubble-table](https://github.com/Evertras/bubble-table).

### UI

* **ui.go:** This file describes the UI/UX that the user will interact with.
Here the structure of each table to be shown is instantiated and updated.

* **ui-tables.go:** This file stores the functions used by ui.go when creating
and populating each table.

### Data

The collected information about the system is found in the [stats.go](stats.go)
file. Here the [gopsutil](https://github.com/shirou/gopsutil) is used.
Due to differences in the logic of generating the data about the processes
running in the machine in each tested operating system (Linux, MacOS and Windows),
the files [processes_linux.go](processes_linux.go), [processes_darwin.go](processes_darwin.go)
and [processes_windows.go](processes_windows.go) were created.
