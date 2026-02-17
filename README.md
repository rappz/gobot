# Chrisbot

A Discord bot built for Ratlantis. It's build to **Punish Chris** in particular, but anyone can be punished or absolved.

## Features

- **Punish** – Mark a user as punished (bad kitten).
- **Absolve** – Clear a user’s punishment (forgive the kitten).
- **Crissy** – Toggle “Crissy mode” (server-specific). When on, Chris gets a 1-minute timeout if he @’s someone. Toggle again to turn it off.

Bots can’t be punished—only God can do that.

### Crissy mode requirements

- **Server** – `/crissy` only works in one server (Ratlantis). In any other server, the bot responds with *"You are not allowed to use this command in this server"*.
- **Behavior** – When Crissy mode is on, Chris is timed out for 1 minute whenever he sends a message that @mentions someone.

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
| `/crissy` | Toggle Crissy mode (server-specific; times out Chris for 1 min when he @’s someone). See [Crissy mode requirements](#crissy-mode-requirements) below. |

---

Made for the server. Chris knows what he did.
