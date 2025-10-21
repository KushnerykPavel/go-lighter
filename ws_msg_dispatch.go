package lighter

import (
	"encoding/json"
	"fmt"
)

type MsgDispatcher interface {
	Dispatch(subs []*uniqSubscriber, msg wsMessage) error
}
type msgDispatcherFunc[T any] func(subs []*uniqSubscriber, msg wsMessage) error

func (d msgDispatcherFunc[T]) Dispatch(subs []*uniqSubscriber, msg wsMessage) error {
	return d(subs, msg)
}

func NewMsgDispatcher[T remoteProcessor](channel string) MsgDispatcher {
	return msgDispatcherFunc[T](func(subs []*uniqSubscriber, msg wsMessage) error {
		if msg.GetChannel() != channel {
			return nil
		}
		var x T
		if err := json.Unmarshal(msg.Data, &x); err != nil {
			return fmt.Errorf("failed to unmarshal message: %v", err)
		}
		x.SetMarketID(msg.GetMarketID())

		for _, subscriber := range subs {
			if subscriber.id == x.Key() {
				subscriber.dispatch(x)
			}
		}

		return nil
	})
}
