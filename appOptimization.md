# Optimizing the htop-clone application.

The interactive process viewer htop is an application used for monitoring and controlling the processes on an operating system, found in terminal-based environments such as Linux. Given the task of creating an htop “clone” by using the Go programming language some performance issues made this application far from ideal compared to htop. The process of solving these issues will be described.

The initial application is built with packages such as bubbletea, bubble-table and gopsutil. Bubbletea and bubble-table are used for displaying the UI that the user interacts with and gopsutil extracts the information to be displayed.

The first issue found was the animation used in the progress bars. The usage of the application with this feature is shown in the following images:

![htop-clone | Before First Optimization](https://i.imgur.com/EgEKhZ0.png?1)
