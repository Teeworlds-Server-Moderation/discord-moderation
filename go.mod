module github.com/Teeworlds-Server-Moderation/discord-moderation

go 1.16

replace github.com/Teeworlds-Server-Moderation/discord-moderation => ./

require (
	github.com/Teeworlds-Server-Moderation/common v0.7.4
	github.com/diamondburned/arikawa/v2 v2.1.0
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/jxsl13/goripr v1.1.1
	github.com/jxsl13/simple-configo v1.18.0
	github.com/onsi/gomega v1.14.0 // indirect
	github.com/streadway/amqp v1.0.0
)
