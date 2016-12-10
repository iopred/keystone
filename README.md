# Keystone
A Discord bot for tracking World of Warcraft mythic keystones

## Usage:

### Installation:

`go get github.com/iopred/keystone/cmd/keystone`

`go install github.com/iopred/keystone/cmd/keystone`

`cd $GOPATH/bin`

### Setup

`keystone -discordtoken "Bot <discord bot token>"`

It is suggested that you set `-discordapplicationclientid` if you are running a bot account.

It is suggested that you set `-discordowneruserid` as this prevents anyone from calling `playingplugin`.

To invite your bot to a server, visit: `https://discordapp.com/oauth2/authorize?client_id=<discord client id>&scope=bot`
