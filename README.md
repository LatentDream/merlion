<h1 align="center">🌊 Merlion</h1>
<p align="center">
  Obsidian-inspired terminal note-taking app<br>
  <a href="https://merlion.dev">merlion.dev</a>
</p>

---

## What is it?

Merlion is a terminal-first, Markdown-based note-taking application, inspired by Obsidian but built from the ground up for command-line workflows.

---

## Current Status

Merlion is a work in progress, but already usable for daily note-taking.

It was born out of a personal need to:

- Avoid scattered `.md` files
- Work seamlessly across multiple machines
- Explore Go in a real-world project

---

## Sync and Accounts

Merlion works fully offline by default.

To sync notes across devices, Merlion uses a cloud backend.  
An account is required to enable sync — [merlion.dev](https://merlion.dev)

---

## Getting Started

1. Clone the repository  
2. Build the project:

```sh
just build
# Or:
go build -o merlion ./cmd/merlion


Run it:
```sh
./merlion
```

#### Features
- Keyboard (only) navigation
- Local-first note storage
- Optional cloud storage to sync notes across devices
- Built-in themes: Gruvbox and NeoTokyo
  - **Feel free to submit a PR to add more themes**.
  - Toggle themes with ctrl+t


#### Screenshots
<p align="center"> <img src="./screenshots/Screenshot_1.png" width="45%"> <img src="./screenshots/Screenshot_2.png" width="45%"> <img src="./screenshots/Screenshot_3.png" width="45%"> </p> 

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

#### Cloud Code

The cloud backend is currently private.
Roadmap

#### Planned features:
- [ ] Full-text search
- [ ] End-to-end encrypted
- [ ] Sync your local notes to the cloud
- [ ] Lightweight web UI (read-only)
