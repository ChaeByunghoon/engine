package adapter

import "github.com/it-chain/it-chain-Engine/p2p/api"

type CommandHandler struct {
	leaderApi api.LeaderApi
	peerApi   api.PeerApi
}
