# member-bot

# .envの例
Discord のトークンが必須です。

https://discord.com/developers/applications/ からBotのトークンを取得してください。


## 使用方法
1. .envを設定する (sample.envを.envにリネーム)
   
    ```.env
    DISCORD_TOKEN=<Discord Bot Token>
    GUILD_ID=<管理するギルド(Discord的にはサーバー)のID>
    ANNOUNCE_CHANNEL_ID=<HUITの入部期限が近いことを告知するチャンネルのID>
    MODERATORS_CHANNEL_ID=<管理者用チャンネルのID(kick前日やkickした日に連絡が行く)>
    ```

1. docker-compose で MYSQL と一緒に起動

    ```bash
    $ docker-compose up -d
    ```
