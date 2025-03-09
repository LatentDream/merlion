# 🌊 Merlion

**What is it?**
> Notion - For the terminal

**Current Progress?**
> MVP - I can use it everyday comfortably, but it's short in features

**Why?**
> Why not?
> - I wanted to check out Go
> - I'm using it every day, no more `.md` everywhere
> - And work across multiple computers (**Primary motivation**)

**Want to try it?**
> It's free, *for now at least*
> 1. Create an account [Merlion.dev](https://merlion.dev)
> 2. Build the project `go build -o merlion ./cmd/merlion`
> 3. Run the project `./merlion`

**What does it look like?**
> You'll have to try it to see it
> - It has Gruvbox & NeoTokyo Theme if that can convince you
> - `ctrl+t` to switch between them :)
> - I'll attach some pictures to the readme soon enough

**Where's the backend code?**
> Private at the moment

---

## Tmux-i-fy it
I usually run the app with `<tmux-leader>m` which simply open a window with Merlion & kill the window on exit

**To do so:**
1. Create a simple script 
    ```sh
    #!/bin/bash
    # Execute (Should be in your path)
    merlion
    # After exits, kill the window
    tmux send-keys "exit" C-m
    ```

2. Bind it in your tmux config
    ```conf
    bind-key m new-window -n "Merlion" -c "#{pane_current_path}" "~/.config/scripts/tmux-merlion.sh"
    bind-key C-m new-window -n "Merlion" -c "#{pane_current_path}" "~/.config/scripts/tmux-merlion.sh"
    ```


