package handleEvent

import (
	"eolian/munc/tcp_server/types"
)

func AddTarget(room *types.Room, target *types.Target) {
	room.AddTarget(target)
}

func RemoveTarget(room *types.Room, target *types.Target) error {
	err := room.RemoveTarget(target)
	if err != nil {
		return err
	}
	return nil
}
