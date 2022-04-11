package main

import (
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

// join時
func handleJoin(s *discordgo.Session, member *discordgo.GuildMemberAdd) {
	guestRoleID := os.Getenv("GUEST_ROLE_ID")

	if member.User.ID == s.State.User.ID {
		return
	}

	// 参加したのがBotなら無視
	if member.User.Bot {
		return
	}

	log.Println("join: ", member.Member.User.Username)

	// 体験入部期間
	trialPeriodMonth := 2 /* month */

	location := time.FixedZone("Asia/Tokyo", 9*60*60)
	rawLimitDay := time.Now().In(location).AddDate(0, trialPeriodMonth, 0)
	limitDay := rawLimitDay.Format("2006/01/02")

	guild, err := s.Guild(member.GuildID)
	if err != nil {
		log.Println(err)
		return
	}
	guildJoinChannelID := guild.SystemChannelID

	// Guestロールを追加
	s.GuildMemberRoleAdd(member.GuildID, member.User.ID, guestRoleID)

	// yyyy/mm/dd という名前の権限を追加
	// 同名のロールがある場合、そのロールを付与(同日に複数人が入ってくるとロールが重複するため)
	if !existSameRole(guild, limitDay) {
		role, err := s.GuildRoleCreate(member.GuildID)
		if err != nil {
			log.Println(err)
			return
		}
		s.GuildRoleEdit(member.GuildID, role.ID, limitDay, 0, false, 0, true)
		s.GuildMemberRoleAdd(member.GuildID, member.User.ID, role.ID)
	} else {
		// 既にロールが存在するため、当該ユーザにそのロールを付与
		roleID := getRoleID(guild, limitDay)
		s.GuildMemberRoleAdd(member.GuildID, member.User.ID, roleID)
	}

	// joinチャンネルにようこそメッセージと、体験入部期間を通知
	limitDayContent := rawLimitDay.Format("2006年1月2日")
	mention := member.User.Mention()
	content := mention + " さん、HUITにようこそ！\n体験入部期間は " + limitDayContent + " までとなります。\n #ようこそ！ を読んで、活動に楽しくご参加ください！"
	s.ChannelMessageSend(guildJoinChannelID, content)
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
