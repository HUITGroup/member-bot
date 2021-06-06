# member-bot

## .env の例

Discord のトークンが必須です。

https://discord.com/developers/applications/ から Bot のトークンを取得してください。

## 使用方法

1. .env を設定する (sample.env を.env にリネーム)

   ```.env
   DISCORD_TOKEN=<Discord Bot Token>
   GUILD_ID=<管理するギルド(Discord的にはサーバー)のID>
   ANNOUNCE_CHANNEL_ID=<HUITの入部期限が近いことを告知するチャンネルのID>
   MODERATORS_CHANNEL_ID=<管理者用チャンネルのID(kick前日やkickした日に連絡が行く)>
   ```

1. docker-compose で起動

   ```bash
   $ docker-compose up -d
   ```
