package main

import (
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

func batch(dg *discordgo.Session) {
	for {
		if isBatchTime() {
			trialMemberBatch(dg)
			log.Println("info: end batch")
			time.Sleep(1 * time.Hour)
		}
	}
}

func isBatchTime() bool {
	location := time.FixedZone("Asia/Tokyo", 9*60*60)
	rawNowTime := time.Now().In(location)
	nowHour := rawNowTime.Format("15")

	return nowHour == "14"
}

func trialMemberBatch(dg *discordgo.Session) {
	guildID := os.Getenv("GUILD_ID")
	announceChannelID := os.Getenv("ANNOUNCE_CHANNEL_ID")
	moderatorsChannelID := os.Getenv("MODERATORS_CHANNEL_ID")

	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	nowTime := time.Now().In(jst)

	// 日数ロールを探す
	guild, err := dg.Guild(guildID)
	if err != nil {
		log.Println(err)
		return
	}

	layout := "2006/1/2"
	for _, role := range guild.Roles {
		trialTimeRole, err := time.ParseInLocation(layout, role.Name, jst)

		// パース出来ないロールは関係ないのでスキップ
		if err != nil {
			continue
		}

		// パースできて、かつ体験期間が今日より前の人はkickして、ロールを消す
		if trialTimeRole.Before(nowTime) {
			log.Println("kick: ", trialTimeRole.Format(layout))
			// kick
			members, err := searchRoleMembers(dg, guild.ID, role.ID)
			if err != nil {
				log.Println(err)
				return
			}

			// userごとにkick
			for _, mem := range members {
				roleUserID := mem.User.ID
				userName := mem.User.Username
				byeMessage := "体験入部期間が終了したため"
				dg.GuildMemberDeleteWithReason(guildID, roleUserID, byeMessage)
				content := userName + " さんの体験入部期間が終了しました。"
				dg.ChannelMessageSend(moderatorsChannelID, content)
			}
			// del role
			dg.GuildRoleDelete(guildID, role.ID)
			continue
		}

		// パースできて、かつ体験期間終了が1週間後の人は連絡
		if nowTime.AddDate(0, 0, 7).Format(layout) == trialTimeRole.Format(layout) {
			log.Println("kick after week: ", trialTimeRole)
			members, err := searchRoleMembers(dg, guild.ID, role.ID)
			if err != nil {
				log.Println(err)
				return
			}

			for _, mem := range members {
				mention := mem.Mention()
				content := mention + " さんの体験入部期間はあと1週間で終了します。\n今後も活動を続けたい場合は、ぜひ入部をお願いします。"
				dg.ChannelMessageSend(announceChannelID, content)
				continue
			}
			continue
		}

		// パースできて、かつ体験期間終了が明日の人がいる場合、運営に確認用の連絡
		if nowTime.AddDate(0, 0, 1).Format(layout) == trialTimeRole.Format(layout) {
			log.Println("kick tommorow: ", trialTimeRole)
			members, err := searchRoleMembers(dg, guild.ID, role.ID)
			if err != nil {
				log.Println(err)
				return
			}

			for _, mem := range members {
				mention := mem.Mention()
				content := "自動通知: " + mention + " さんの体験入部期間が明日で終了します。\n部費の支払いが終わっている場合、" + mention + " さんの体験入部期間ロールを解除してください"
				dg.ChannelMessageSend(moderatorsChannelID, content)
				continue
			}
			continue
		}
		// パースできて、かつ体験入部期間が直近に迫っていない人はスキップ
	}
}

func searchRoleMembers(dg *discordgo.Session, guildID, roleID string) (members []*discordgo.Member, err error) {
	// 引数で受け取ったロールIDを持つ メンバー(部員) をmembersスライスにappend
	// 体験入部期間のリミット権限は最大で1人1つのため、そのメンバーにリミット権限が既にあればそれ以上調べる必要はない
	mems, err := dg.GuildMembers(guildID, "", 1000)
	if err != nil {
		return members, err
	}

	for _, mem := range mems {
		log.Println(mem)
		for _, memRoleID := range mem.Roles {
			if memRoleID == roleID {
				members = append(members, mem)
				break
			}
		}
	}
	return members, nil
}
