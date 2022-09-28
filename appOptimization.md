# Optimizing the htop-clone application

The interactive process viewer htop is an application used for monitoring and
controlling the processes on an operating system, found in terminal-based
environments such as Linux. Given the task of creating an htop “clone” by using
the Go programming language some performance issues made this application far
from ideal compared to htop. The process of solving these issues will be described.

The initial application is built with packages such as bubbletea, bubble-table
and gopsutil. Bubbletea and bubble-table are used for displaying the UI that the
user interacts with and gopsutil extracts the information to be displayed.

## Progress bars animation

The first issue found was the animation used in the progress bars. The usage of
the application with this feature is shown in the following images:

| ![htop-clone - Before First Optimization](https://i.imgur.com/EgEKhZ0.png?1) |
| :-: |
| *Running htop-clone with animated progress bars.* |

| ![time command - Before First Optimization](https://i.imgur.com/BqsdkHy.png) |
| :-: |
| *Result of running the `time` command.* |

| ![top 10 Profiler functions - Before First Optimization](https://i.imgur.com/od9ZoiL.png) |
| :-: |
| *Top 10 functions ran by the program and sorted by their usage of CPU.* |

**By removing the animation we got the following results:**

| ![time command - After First Optimization](https://i.imgur.com/CmLd9OJ.png) |
| :-: |
| *Result of running the `time` command.* |

| ![top 10 Profiler functions - After First Optimization](https://i.imgur.com/s4hiz9v.png) |
| :-: |
| *Top 10 functions ran by the program and sorted by their usage of CPU.* |

## The Syscall bottleneck

The profiler revealed a huge bottleneck, the syscall.Syscall6 function. The
origin of this phenomenon was attached to the usage of the gopsutil package as
it [demands access to the filesystem to the PC](https://stackoverflow.com/a/69301915/7987716),
adding a heavy toll to the CPU. This was tested on Linux and a similar pattern
was founded on Windows. MacOS didn't present this behavior.

### Solution

By replacing the gopsutil package with a less elegant alternative using the
terminal command `ps` from Linux and MacOs further improvements were achievable.

| ![time command - With Syscall bottleneck](https://i.imgur.com/OfE7Bd9.png) |
| :-: |
| *Result of running the `time` command with the shirou/gopsutil/v3/process package.* |

| ![time command - Without Syscall bottleneck](https://i.imgur.com/ijL4fNW.png) |
| :-: |
| *Result of running the `time` command without the shirou/gopsutil/v3/process package.* |

| ![top 10 Profiler functions - Without Syscall bottleneck](https://i.imgur.com/hbuOmdC.png) |
| :-: |
| *Top 10 functions ran by the program and sorted by their usage of CPU, without the shirou/gopsutil/v3/process package.* |
