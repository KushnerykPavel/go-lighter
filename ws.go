package lighter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sonirico/vago/maps"
)

type logger interface {
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
}

type Subscription struct {
	ID      string
	Payload any
	Close   func()
}

type WebsocketClient struct {
	url           string
	conn          *websocket.Conn
	mu            sync.RWMutex
	writeMu       sync.Mutex
	done          chan struct{}
	reconnectWait time.Duration
	closeOnce     sync.Once
	debug         bool
	nextSubID     atomic.Int64
	logger        logger

	subscribers           map[string]*uniqSubscriber
	msgDispatcherRegistry map[string]MsgDispatcher
}

func NewWebsocketClient(baseURL string, opts ...WsOpt) (*WebsocketClient, error) {
	if baseURL == "" {
		baseURL = MainnetAPIURL
	}
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %v", err)
	}
	parsedURL.Scheme = "wss"
	parsedURL.Path = "/stream"
	wsUrl := parsedURL.String()

	c := &WebsocketClient{
		url:           wsUrl,
		done:          make(chan struct{}),
		reconnectWait: time.Second,
		subscribers:   make(map[string]*uniqSubscriber),
		msgDispatcherRegistry: map[string]MsgDispatcher{
			ChannelOrderBook: NewMsgDispatcher[*OrderBook](ChannelOrderBook),
			ChannelTrades:    NewMsgDispatcher[*Trades](ChannelTrades),
		},
	}
	for _, opt := range opts {
		opt.Apply(c)
	}

	return c, nil

}

func (w *WebsocketClient) Connect(ctx context.Context) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.conn != nil {
		return nil
	}

	dialer := websocket.Dialer{}

	conn, _, err := dialer.DialContext(ctx, w.url, nil)
	if err != nil {
		return fmt.Errorf("websocket dial: %w", err)
	}

	w.conn = conn

	go w.readPump(ctx)

	return w.resubscribeAll()
}

func (w *WebsocketClient) Close() error {
	var err error
	w.closeOnce.Do(func() {
		err = w.close()
	})
	return err
}

func (w *WebsocketClient) readPump(ctx context.Context) {
	defer func() {
		w.mu.Lock()
		if w.conn != nil {
			_ = w.conn.Close() // Ignore close error in defer
			w.conn = nil
		}
		w.mu.Unlock()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-w.done:
			return
		default:
			_, msg, err := w.conn.ReadMessage()
			if err != nil {
				if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					w.logErrf("websocket read error: %v", err)
				}
				return
			}

			if w.debug {
				w.logDebugf("[<] %s", string(msg))
			}
			fmt.Println(string(msg))
			var wsMsg wsMessage
			if err := json.Unmarshal(msg, &wsMsg); err != nil {
				w.logErrf("websocket message parse error: %v", err)
				continue
			}
			switch {
			case wsMsg.OrderBook != nil:
				wsMsg.MarketID = wsMsg.GetMarketID()
				wsMsg.Data = wsMsg.OrderBook
			}

			if err := w.dispatch(wsMsg); err != nil {
				w.logErrf("failed to dispatch websocket message: %v", err)
			}
		}
	}
}

func (w *WebsocketClient) dispatch(msg wsMessage) error {
	dispatcher, ok := w.msgDispatcherRegistry[msg.GetChannel()]
	if !ok {
		return fmt.Errorf("no dispatcher for channel: %s", msg.Channel)
	}

	w.mu.RLock()
	subscribers := maps.Values(w.subscribers)
	w.mu.RUnlock()

	return dispatcher.Dispatch(subscribers, msg)
}

func (w *WebsocketClient) subscribe(
	payload subscriptable,
	callback func(any),
) (*Subscription, error) {
	if callback == nil {
		return nil, fmt.Errorf("callback cannot be nil")
	}

	w.mu.Lock()

	pkey := payload.Key()
	subscriber, exists := w.subscribers[pkey]
	if !exists {
		subscriber = newUniqSubscriber(
			pkey,
			payload,
			// on subscribe
			func(p subscriptable) {
				if err := w.sendSubscribe(p); err != nil {
					w.logErrf("failed to subscribe: %v", err)
				}
			},
			// on unsubscribe
			func(p subscriptable) {
				w.mu.Lock()
				defer w.mu.Unlock()
				delete(w.subscribers, pkey)
				if err := w.sendUnsubscribe(p); err != nil {
					w.logErrf("failed to unsubscribe: %v", err)
				}
			},
		)

		w.subscribers[pkey] = subscriber
	}

	w.mu.Unlock()

	nextID := w.nextSubID.Add(1)
	subID := key(pkey, strconv.Itoa(int(nextID)))
	subscriber.subscribe(subID, callback)

	return &Subscription{
		ID: subID,
		Close: func() {
			subscriber.unsubscribe(subID)
		},
	}, nil
}

func (w *WebsocketClient) close() error {
	close(w.done)

	w.mu.Lock()
	defer w.mu.Unlock()

	if w.conn != nil {
		return w.conn.Close()
	}

	for _, subscriber := range w.subscribers {
		subscriber.clear()
	}
	return nil
}

func (w *WebsocketClient) resubscribeAll() error {
	for _, subscriber := range w.subscribers {
		if err := w.sendSubscribe(subscriber.subscriptionPayload); err != nil {
			return fmt.Errorf("resubscribe: %w", err)
		}
	}
	return nil
}

func (w *WebsocketClient) sendSubscribe(payload subscriptable) error {
	return w.writeJSON(payload)
}

func (w *WebsocketClient) sendUnsubscribe(payload subscriptable) error {
	return w.writeJSON(payload)
}

func (w *WebsocketClient) writeJSON(v any) error {
	w.writeMu.Lock()
	defer w.writeMu.Unlock()

	if w.conn == nil {
		return fmt.Errorf("connection closed")
	}

	if w.debug {
		bts, _ := json.Marshal(v)
		w.logDebugf("[>] %s", string(bts))
	}

	return w.conn.WriteJSON(v)
}

func (w *WebsocketClient) logDebugf(fmt string, args ...any) {
	if w.logger == nil {
		return
	}

	w.logger.Infof(fmt, args...)
}

func (w *WebsocketClient) logErrf(fmt string, args ...any) {
	if w.logger == nil {
		return
	}

	w.logger.Errorf(fmt, args...)
}
