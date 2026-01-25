<h1 align="center">🌊 Merlion</h1>
<p align="center">
  Obsidian-inspired terminal note-taking app<br>
  <a href="https://note.merlion.dev">merlion.dev</a>
</p>

---

## What is it?

Merlion is a TUI, Markdown-based note-taking application, inspired by Obsidian but built for command-line workflows.
- Compatible with [Obsidian](https://obsidian.md/) vaults
- Ability to use a single SQLite database for all notes stored locally
- Cloud storage support - work in progress, see below

_Merlion works fully offline by default, no account needed, all files are on your computer in a SqliteDB or directly as .md files._

<p align="center"> <img src="./screenshots/Screenshot_1.png" width="45%"> <img src="./screenshots/Screenshot_2.png" width="45%"> <img src="./screenshots/Screenshot_3.png" width="45%"> </p> 

---

#### Features
- Keyboard (only) navigation
- Local-first note storage
- Optional cloud storage to sync notes across devices
  - Lightweight web UI (this will be removed in favor of a sync feature)
- Built-in themes: Gruvbox and NeoTokyo
  - **Feel free to submit a PR to add more themes**.
  - Or to ask for a new theme to be added
  - Toggle themes with ctrl+t
- Naviguate between note base on note title
- Markdown support
- Use your `$EDITOR` as note editor

### Keymap

| Key(s) | Action | Description |
|--------|--------|-------------|
| `↑` or `k` | Up | Move selection up |
| `↓` or `j` | Down | Move selection down |
| `←` or `h` | Left | Go back to list view |
| `→` or `l` | Right | View selected note |
| `delete` | Delete | Delete selected item |
| `tab` | Next Tab | Switch to next tab |
| `shift+tab` | Previous Tab | Switch to previous tab |
| `pgup` or `ctrl+u` | Page Up | Scroll up one page |
| `pgdn` or `ctrl+d` | Page Down | Scroll down one page |
| `enter` | Select | Confirm selection |
| `e` | Edit | Edit the current note |
| `m` | Manage | Manage note information |
| `esc` | Clear Filter/Back | Clear current filter or go back |
| `q` or `ctrl+c` | Quit | Exit the application |
| `ctrl+t` | Toggle Theme | Switch between light/dark theme |
| `i` | Toggle Info | Show/hide note information panel |
| `ctrl+p` | Toggle Info Position | Change note info panel position |
| `ctrl+f` | Toggle Compact View | Toggle compact view (large screens only) |
| `c` | Create | Create a new note |
| `(` or `)` | Toggle Store | Toggle store view |

---

## Getting Started

1. Clone the repository  
2. Build the project:

```sh
just build

# Or:
go build -o merlion ./cmd/merlion
```

Run it:
```sh
./merlion
```

#### Tmux Integration

Add the following to your .tmux.conf to launch Merlion in a popup window:

```.tmux.conf
bind C-m display-popup \
  -d "#{pane_current_path}" \
  -w 90% \
  -h 90% \
  -E "merlion --compact"
```
Then launch it with `<tmux-leader> + m`.

#### Cloud Storage

Merlion supports cloud storage, you can create an account at [note.Merlion.dev](https://note.merlion.dev) to get your notes across devices.
- This is still a WIP and subject to change / removal in favor of a sync feature
- The notes are not encrypted on the server (yet), it's still a work in progress
    - **DON'T use this for sensitive data**

- To have online note for sharing note across computer, you can create an account at [note.Merlion.dev](https://note.merlion.dev)
    - Your notes will then be local, or online, switch between the two workplace with `(` or `)`
    - Sync capability between online & offline note will soon be added


---


<p align="center">
  <i>This software is provided "as is", the code ain't perfect, the app ain't perfect and I'm having fun </i>
</p>
