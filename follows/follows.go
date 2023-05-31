package follows

import "github.com/hakonslie/twitchfriendmaker/session"

type Channel string

type FollowStorage map[session.SessionID][]Channel

func (f FollowStorage) GetFollows(id session.SessionID) []Channel {
	return f[id]
}

func (f FollowStorage) AddFollows(id session.SessionID, channels []Channel) {
	f[id] = channels
}
