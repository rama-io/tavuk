# Tavuk

**Tavuk** is a small Discord bot written in Go lang for Rama discord server.

---

## Commands

Both subcommands live under `/announce`. The command is restricted to members
with the **Administrator** permission.

`/announce set <app> <channel> <role>`
Maps an app to the channel where its announcements go and the role to mention.

Example: `/announce set app:"Txori" channel:#updates role:@members`

`/announce publish <app> <version> <changes>`
Posts the release announcement. The app name and version combine into the title,
so `app:"Txori" version:"10"` becomes **Txori v10**. The `app` option
autocompletes from the apps you already set up. `<changes>` is a feature list in
the form `"- feat 1 - feat 2 - feat 3"`.

Example: `/announce publish app:"Txori" version:"10" changes:"- Fix a no-matching-query bug - Add media file filters"`

The announcement looks like this:

```
**Txori v10** has been released!

**What's Changed?**

- Fix a no-matching-query bug
- Add media file filters

**Where to download?**
- It's already on GitHub releases. It will take some time for F-Droid update.

|| <@&ROLE ID> ||
```

---

## Setup

Create an app and bot at <https://discord.com/developers/applications>, then
copy the example environment file and fill it in:

```
cp .env.example .env
```

```dotenv
DISCORD_TOKEN=your_bot_token
# GUILD_ID=only needed for instant, guild-local commands during development
```

Invite the bot with the `bot` and `applications.commands` scopes and at least
the **View Channels**, **Send Messages** and **Mention Everyone** permissions:

```
https://discord.com/oauth2/authorize?client_id=YOUR_CLIENT_ID&permissions=134144&scope=bot%20applications.commands
```

Run it:

```
go run ./cmd/tavuk
```

---

## Data

The data mapping lives in `data/apps.json` locally. Keep the directory gitignored, so it never ends up in the repository. Do not forget to gitignore it if you change it via .env.

---

## License

**Tavuk** is Free Software. You are free to use, study, share, and improve it
under the terms of the **GNU General Public License v3** or later.
