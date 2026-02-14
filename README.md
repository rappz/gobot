# Chrisbot

A Discord bot built for Ratlantis. It's build to **Punish Chris** in particular, but anyone can be punished or absolved.

## Features

- **Punish** – Mark a user as punished (bad kitten).
- **Absolve** – Clear a user’s punishment (forgive the kitten).
- **Crissy** – Toggle “Crissy mode” (server-specific). When on, Chris gets a 1-minute timeout if he sends a message without @’ing someone. Toggle again to turn it off.

Bots can’t be punished—only God can do that.

## Requirements

- Go 1.24+
- A [Discord Bot](https://discord.com/developers/applications) token

## Environment Variables

| Variable | Description |
|---------|-------------|
| `TOKEN` | Your Discord bot token (required) |

Create a `.env` in the project root (or set `TOKEN` in your environment):

```env
TOKEN=your_bot_token_here
```

## Run locally

From the project root (so `.env` is found):

```bash
go run .
```

Optional: remove slash commands on exit:

```bash
go run . --rm-cmd
```

## Run with Docker

```bash
docker build -t gobot .
docker run -e TOKEN=your_bot_token_here gobot
```

## Commands (slash commands)

| Command   | Description                         |
|----------|-------------------------------------|
| `/punish @user` | Punish a user (bad kitten).        |
| `/absolve @user` | Absolve a user (forgive the kitten). |
| `/crissy` | Toggle Crissy mode (Chris timeout rule, times out Chris for 1 minute). |

---

Made for the server. Chris knows what he did.
