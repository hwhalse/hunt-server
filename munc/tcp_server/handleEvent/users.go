package handleEvent

import "eolian/munc/tcp_server/types"

func AddUser(room *types.Room, user *types.User) error {
	err := room.AddUser(user)
	if err != nil {
		return err
	}
	return nil
}
