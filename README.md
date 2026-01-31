# Chatty - Terminal Chat Interface

![Docker Version](https://img.shields.io/docker/v/bwalheim1205/chatty?sort=semver)

```chatty``` is an open-source, Go-based TUI for working with LLMs from your terminal. Built with developer workflows in mind and inspired by Neovim, chatty aims to make AI a natural part of your command-line environment through a fast, customizable, and keyboard-first interface.

**Warning**: This is still in early stage development so expect frequent (and breaking) changes

# Features
- **Chat TUI**: A TUI interface for interacting with different LLM models
- **Multiple Model Support**: Supports ability to connect to verity of model providers. (Ollama, ChatGPT, etc.)

### Prerequisites
*   Go 1.21 or higher
*   Git

# Installation

### Clone Repository

First clone down the repository then follow instructions for either Go executable or docker build

```sh
git clone https://github.com/bwalheim1205/chatty.git
cd chatty
```

### Go

```sh
go build -o chatty ./cmd/chatty
```

### Docker

```sh
docker build . -t chatty
```

# Usage

Once you've completed either build you can run the chatty using executable or docker image.

### Go
```sh
./chatty
```

### Docker
```sh
docker run -v chatty
```

## Key Bindings

This chat TUI uses **Vim-inspired modal key bindings**. There are three primary modes:

* **Read Mode** (default): navigate messages and history
* **Chat Mode**: type and send chat messages
* **Command Mode**: enter `:` or `/` commands

---

### Read Mode (Navigation)

Used for scrolling and moving the cursor through chat history.

| Key       | Action                         |
| --------- | ------------------------------ |
| `j` / `↓` | Move cursor down one line      |
| `k` / `↑` | Move cursor up one line        |
| `h` / `←` | Move cursor left               |
| `l` / `→` | Move cursor right              |
| `gg`      | Jump to top of chat history    |
| `G`       | Jump to bottom of chat history |
| `:`       | Enter command mode             |
| `/`       | Enter search command mode      |
| `i`       | Enter chat mode                |
| `c`       | Enter chat mode                |

---

### Chat Mode (Typing Messages)

Used to compose and send chat messages.

| Key                         | Action                                 |
| --------------------------- | -------------------------------------- |
| `Enter`                     | Send chat message                      |
| `Backspace`                 | Delete last character                  |
| `Ctrl+C`                    | Cancel message and return to read mode |

Notes:

* If the input buffer becomes empty and you press `Backspace`, the app returns to Read Mode.

---

### Command Mode (`:` and `/`)

Used for commands and search (Vim-style).

| Key                         | Action                                 |
| --------------------------- | -------------------------------------- |
| `Enter`                     | Execute command                        |
| `Backspace`                 | Delete last character                  |
| `Esc` / `Ctrl+C`            | Cancel command and return to read mode |
| *(any printable character)* | Append to command                      |


Commands

* :m <model-name> / :model <model-name> - Swap models
* :p <provider-name> / :provider <provider-name> - Swap provider (chatgpt or ollama)
* :c / :clear - clear chat history
* :q



# Roadmap

If there's something you'd like to see implemented we'd love to her from you just open an issue. Here are some of the functionality next on the horizon:

- **Chat History**: To give a full featured chat interface we want to track previous chat histories allowing it to remember your previous conversations and resume them.
- **MCP Connections**: Tools! Tools! Tools! Similar to must GUI chat interfaces we want the ability to leverage tools to improve responses
