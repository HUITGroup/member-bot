package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("info: no .env file")
	}

	// Discord と接続
	token := os.Getenv("DISCORD_TOKEN")
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Println("Error creacting Discord session,", err)
	}

	// 誰かがjoinした時の処理
	// うまくいかない時はDiscord側の設定をしていないかも
	// https://github.com/bwmarrin/discordgo/issues/793#issuecomment-659021953 を参照

	dg.AddHandler(handleJoin)

	// ↑のURLとは違うが、こちらの方が適切(必要以上のデータにアクセスしないため)
	dg.Identify.Intents = discordgo.IntentsGuildMembers

	err = dg.Open()
	if err != nil {
		log.Fatalln("Error opening connection,", err)
		return
	}

	// バッチ処理
	go batch(dg)

	log.Println("Bot is now running. Press Ctrl-C to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc

	dg.Close()
	log.Println("shudown success bye!")
}
