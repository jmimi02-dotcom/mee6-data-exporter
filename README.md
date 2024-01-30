## Description
This CLI tool allows you to export user data from the [Mee6](https://mee6.xyz/en) Discord bot. The data includes XP and Levelling
information, allowing you to migrate data to your own bot. There is no official API documentation for this feature, however
the original code for this functionality can be found in [this repository](https://github.com/Colean128/mee6/blob/master/bot/plugins/levels.py), and the formula used to calculate the required XP can be found [here](https://github.com/Mee6/Mee6-documentation/blob/master/docs/levels_xp.md).

## To-Do

- [ ] Paginate the Mee6 API. As of writing, this program will only return page 0 containg the first 100 users in the guild.
- [ ] Write a fancy TUI using [Charm](https://charm.sh/) that guides users through the process.

## Installation

### Requirements
Before completing the installation process, ensure that you have the following dependencies installed and available at the system PATH:

- [Go 1.19](https://golang.org/dl/)
- [gcc](https://gcc.gnu.org/install/index.html) (Also available through [MSYS2](https://www.msys2.org/) on Windows)

1. Clone this repository and navigate to the project directory

```sh
git clone https://github.com/luisjones/mee6-xp-exporter.git && cd mee6-xp-exporter
```

2. Open the main.go file and change the serverId variable to represent the ID of the guild whose data you wish to export. This ID can be
found by enabling [Developer Mode](hhttps://support.discord.com/hc/en-us/articles/206346498-Where-can-I-find-my-User-Server-Message-ID) 
on Discord.

3. Ensure all the required libraries are installed locally by running
```sh
go mod download
```

4. Build the project into an executable
```sh
go build 
```

5. Run the compiled build from the previous step.
```sh
go run .
```

6. Ensure the database contains the requested information. This can be achieved online using a [SQLite Viewer](https://inloop.github.io/sqlite-viewer/)



