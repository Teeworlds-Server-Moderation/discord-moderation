package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/Teeworlds-Server-Moderation/common/dto"
	"github.com/Teeworlds-Server-Moderation/common/events"
	"github.com/Teeworlds-Server-Moderation/common/topics"
)

// serverTopic may be either the server's ip:port address or the broadcast topic
const requestorID = "detect-vpn"

func (dvc *detectVPNConfig) RequestBan(player dto.Player, banReason, sourceServerAddr string) error {
	event := events.NewRequestCommandExecEvent()
	event.Timestamp = time.Now().Format("2006-01-02 15:04:05")
	event.Requestor = requestorID
	event.EventSource = requestorID

	// construct command and replace
	replacer := strings.NewReplacer(
		"{IP}", player.IP,
		"{ID}", fmt.Sprintf("%d", player.ID),
		"{DURATION:MINUTES}", fmt.Sprintf("%d", int64(dvc.BanDuration()/time.Minute)),
		"{DURATION:SECONDS}", fmt.Sprintf("%d", int64(dvc.BanDuration()/time.Second)),
		"{REASON}", banReason,
	)

	broadcastFeasible := true
	if strings.Contains(dvc.BanCommand(), "{ID}") {
		broadcastFeasible = false
	}

	banCommand := replacer.Replace(dvc.BanCommand())
	event.Command = banCommand

	if dvc.BroadcastBans() && broadcastFeasible {
		// ban on all servers
		// if broadcasting makes sense
		// if the ban command contains an ID,
		// it makes no sense to broadcast it
		Broker().Publisher().Publish(topics.Broadcast, "", event.Marshal())
	} else {
		// only ban on the server where the player joined
		// do not publish to exchange, but directly to the queue
		Broker().Publisher().Publish("", sourceServerAddr, event.Marshal())
	}
	return nil
}
