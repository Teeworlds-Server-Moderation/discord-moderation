package processors

import (
	"fmt"
	"strings"

	"github.com/Teeworlds-Server-Moderation/common/events"
	"github.com/Teeworlds-Server-Moderation/discord-log/markdown"
	"github.com/diamondburned/arikawa/v2/bot"
	"github.com/diamondburned/arikawa/v2/discord"
	a "github.com/streadway/amqp"
)

var (
	DiscordLogSkipJoinLeaveMessages = false
	DiscordLogSkipWhisperMessages   = false
)

func DiscordLog(ctx *bot.Context, channelID discord.ChannelID, eventType string, msg a.Delivery) error {
	switch eventType {
	case events.TypePlayerJoined, events.TypePlayerLeft:
		if DiscordLogSkipJoinLeaveMessages {
			return nil
		}
	case events.TypeChatWhisper:
		if DiscordLogSkipWhisperMessages {
			return nil
		}
	}

	_, err := ctx.SendMessage(channelID, fmtEvent(eventType, msg.Body), nil)
	return err
}

func fmtEvent(eventType string, data []byte) string {
	str := string(data)
	prefix := ""

	if strings.Contains(eventType, ":") {
		prefix = "[" + strings.ToLower(strings.Split(eventType, ":")[1]) + "]"
	}

	var err error
	switch eventType {
	case events.TypePlayerJoined:
		event := events.NewPlayerJoinedEvent()
		err = event.Unmarshal(str)
		if err != nil {
			break
		} else {
			str = fmt.Sprintf(
				"%s %s %s",
				markdown.Flag(event.Player.Country),
				markdown.WrapInInlineCodeBlock(event.Player.Name),
				markdown.WrapInInlineCodeBlock(event.Player.Clan),
			)
		}
	case events.TypePlayerLeft:
		event := events.NewPlayerLeftEvent()
		err = event.Unmarshal(str)
		if err != nil {
			break
		} else {
			str = fmt.Sprintf(
				"%s %s %s",
				markdown.Flag(event.Player.Country),
				markdown.WrapInInlineCodeBlock(event.Player.Name),
				markdown.WrapInInlineCodeBlock(event.Player.Clan),
			)
		}
	case events.TypeChat:
		event := events.NewChatEvent()
		err = event.Unmarshal(str)
		if err != nil {
			break
		} else {
			str = fmt.Sprintf(
				"%s (%d): %s",
				markdown.WrapInInlineCodeBlock(event.Source.Name),
				event.Source.ID,
				markdown.Escape(event.Text),
			)
		}
	case events.TypeChatTeam:
		event := events.NewChatTeamEvent()
		err = event.Unmarshal(str)
		if err != nil {
			break
		} else {
			str = fmt.Sprintf(
				"%s (%d): %s",
				markdown.WrapInInlineCodeBlock(event.Source.Name),
				event.Source.ID,
				markdown.Escape(event.Text),
			)
		}
	case events.TypeChatWhisper:
		event := events.NewChatWhisperEvent()
		err = event.Unmarshal(str)
		if err != nil {
			break
		} else {
			str = fmt.Sprintf(
				"%s (%d) -> %s (%d): %s",
				markdown.WrapInInlineCodeBlock(event.Source.Name),
				event.Source.ID,
				markdown.WrapInInlineCodeBlock(event.Target.Name),
				event.Target.ID,
				markdown.Escape(event.Text),
			)
		}
	case events.TypeMapChanged:
		event := events.MapChangedEvent{}
		err = event.Unmarshal(str)
		if err != nil {
			break
		} else {
			str = fmt.Sprintf(
				"from %s to %s",
				markdown.WrapInInlineCodeBlock(event.OldMap),
				markdown.WrapInInlineCodeBlock(event.NewMap),
			)
		}
	case events.TypeVoteKickStarted:
		event := events.VoteKickStartedEvent{}
		err = event.Unmarshal(str)

		if err != nil {
			break
		} else {
			str = fmt.Sprintf(
				"%s (%d) kickvotes %s (%d) with reason %s",
				markdown.WrapInInlineCodeBlock(event.Source.Name),
				event.Source.ID,
				markdown.WrapInInlineCodeBlock(event.Target.Name),
				event.Target.ID,
				markdown.WrapInInlineCodeBlock(event.Reason),
			)
		}
	case events.TypeVoteSpecStarted:
		event := events.VoteSpecStartedEvent{}
		err = event.Unmarshal(str)

		if err != nil {
			break
		} else {
			str = fmt.Sprintf(
				"%s (%d) specvotes %s (%d) with reason %s",
				markdown.WrapInInlineCodeBlock(event.Source.Name),
				event.Source.ID,
				markdown.WrapInInlineCodeBlock(event.Target.Name),
				event.Target.ID,
				markdown.WrapInInlineCodeBlock(event.Reason),
			)
		}
	case events.TypeVoteOptionStarted:
		event := events.VoteOptionStartedEvent{}
		err = event.Unmarshal(str)

		if err != nil {
			break
		} else {
			str = fmt.Sprintf(
				"%s (%d) voted option %s with reason %s",
				markdown.WrapInInlineCodeBlock(event.Source.Name),
				event.Source.ID,
				markdown.WrapInInlineCodeBlock(event.Option),
				markdown.WrapInInlineCodeBlock(event.Reason),
			)
		}
	}
	if err != nil {
		return fmt.Sprintf("[ERROR]: %v", err)
	}

	return fmt.Sprintf("%s %s", prefix, str)
}
