# Exploration


## choosing provider
```
$ merlion
You don't have any place to store your notes, please choose one:
- > sqlite database <
-   Obsidian vault
-   Cloud
```

If SQLite database is selected, the user will be prompted to select a path (optional).
-> `merlion new sqlite [<path>]` is also available

If Obsidian vault is selected, the user will be prompted to select a vault.
```
$ merlion  # Or `merlion new file [<path>]`
Select a vault, or create a new one:
Root Folder: ~/notes
```

If Cloud is selected, the user will be prompted to enter their credentials.
```
$ merlion  # Or `merlion new cloud`
Enter your email:
Enter your password:
-> To create an account, visit https://note.merlion.dev/login
-> To delete your credentials, run `merlion logout`
```

We will need to move where the credentials are stored `~/.merlion/credentials.json` instead of the .config folder.

the `~/.config/merlion/config.json` will be used to store the user's preferences.

