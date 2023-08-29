package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
)

var token, guestRoleID string

func init() {
	token = os.Getenv("DISCORD_TOKEN")
	guestRoleID = os.Getenv("GUEST_ROLE_ID")

	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	time.Local = jst
}

func main() {
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

	log.Println("Bot is now running. Press Ctrl-C to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc

	dg.Close()
	log.Println("shudown success bye!")
}

func handleJoin(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	if m.User.Bot || (m.User.ID == s.State.User.ID) {
		return
	}

	// 体験入部期間
	trialPeriodMonth := 2 /* month */

	rawLimitDay := time.Now().AddDate(0, trialPeriodMonth, 0)
	roleName := rawLimitDay.Format("2006/01/02")

	guild, err := s.Guild(m.GuildID)
	if err != nil {
		log.Println(err)
		return
	}

	// Guestロールを追加
	if err := s.GuildMemberRoleAdd(m.GuildID, m.User.ID, guestRoleID); err != nil {
		log.Println(errors.New("Guestロールの付与に失敗しました"))
	}

	// 体験入部期限ロールを付与 (ex. yyyy/mm/dd)
	// 同名のロールがある場合、そのロールを付与(同じ日に複数人が入るとロール名が重複する)
	if existSameRole(guild, roleName) {
		// 既にロールが存在するため、当該ユーザにそのロールを付与
		roleID := getRoleID(guild, roleName)
		if err := s.GuildMemberRoleAdd(m.GuildID, m.User.ID, roleID); err != nil {
			log.Println(errors.New("既存の体験入部期限ロールの付与に失敗しました"))
			return
		}
	} else {
		param := &discordgo.RoleParams{Name: roleName}
		role, err := s.GuildRoleCreate(m.GuildID, param)
		if err != nil {
			log.Println(errors.New("新規体験入部期限ロールの作成に失敗しました"))
			return
		}

		if _, err := s.GuildRoleEdit(m.GuildID, role.ID, param); err != nil {
			log.Println(errors.New("新規体験入部期限ロールの編集に失敗しました"))
		}

		if err := s.GuildMemberRoleAdd(m.GuildID, m.User.ID, role.ID); err != nil {
			log.Println(errors.New("新規体験入部期限ロールの付与に失敗しました"))
		}
	}

	// joinチャンネルにようこそメッセージと、体験入部期間を通知
	if _, err := s.ChannelMessageSend(guild.SystemChannelID, welcomeMessageContent(m.Mention(), roleName)); err != nil {
		log.Println(errors.New("joinチャンネルへのウェルカムメッセージ送信に失敗しました"))
	}
}

func existSameRole(guild *discordgo.Guild, roleName string) bool {
	roleID := getRoleID(guild, roleName)
	return roleID != ""
}

func getRoleID(guild *discordgo.Guild, roleName string) (roleID string) {
	for _, role := range guild.Roles {
		// もしギルド内に既に同名のロールがある場合、ロールIDを返す
		if role.Name == roleName {
			return role.ID
		}
	}

	// 同名のロールがない場合, 空文字列を返す
	return ""
}

func welcomeMessageContent(targetMember, guestRoleName string) string {
	return fmt.Sprintf("%v さん、HUITにようこそ！\n体験入部期間は %v までとなります。", targetMember, guestRoleName)
}
