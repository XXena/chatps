package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/XXena/chatps/pkg/logger"

	"github.com/gorilla/websocket"
)

type WebsocketClient struct {
	wsConn   *websocket.Conn
	userID   string
	chatID   string
	sendChan chan []byte
	logger   *logger.Logger
	Hub      IHub
}

func NewWebsocketClient(hub IHub, conn *websocket.Conn, chatID string, logger *logger.Logger) Client {
	return &WebsocketClient{
		Hub:    hub,
		wsConn: conn,
		chatID: chatID,
		logger: logger,
	}
}

func (c WebsocketClient) GetChatID() ChatID {
	return ChatID(c.chatID)
}

func (c WebsocketClient) GetHub() IHub {
	return c.Hub
}

func (c WebsocketClient) GetSendChan() chan []byte {
	return c.sendChan
}

// SendMessage reads message from the websocket connection and promotes to chat
func (c WebsocketClient) SendMessage() error {
	defer func() {
		c.Hub.Unregister(c)
		err := c.wsConn.Close()
		if err != nil {
			c.logger.Fatal(fmt.Errorf("websocket close connection error: %w", err))
			return
		}
	}()

	c.wsConn.SetReadLimit(maxMessageSize)
	err := c.wsConn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		c.logger.Error(fmt.Errorf("websocket connection set read deadline error: %w", err))
		return err
	}
	c.wsConn.SetPongHandler(func(string) error {
		err := c.wsConn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			c.logger.Error(fmt.Errorf("websocket connection set read deadline error: %w", err))
			return err
		}
		return nil
	})

	for {
		_, msg, err2 := c.wsConn.ReadMessage()
		if err2 != nil {
			if websocket.IsUnexpectedCloseError(err2, websocket.CloseGoingAway) {
				c.logger.Error(fmt.Errorf("read message from ws connection error: %w", err2))
				return err2 // todo или убрать?
			}
			break
		}

		m := Message{
			ChatID: c.GetChatID(),
			Data:   msg,
		}

		c.Hub.Broadcast(m)
	}
	return nil
}

// PullMessage gets message from the send channel to the websocket connection
func (c WebsocketClient) PullMessage() error {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		err := c.wsConn.Close()
		if err != nil {
			c.logger.Fatal(fmt.Errorf("websocket close connection ticker error: %w", err)) // todo см. подобные места, возможно, это избыточно - тут логировать
			return
		}
	}()
	for {
		select {
		case message, ok := <-c.sendChan:
			if !ok {
				err := c.Write(websocket.CloseMessage, []byte{})
				if err != nil {
					c.logger.Error(fmt.Errorf("websocket message not received, then Write error: %w", err))
					return err
				}
				return errors.New("websocket message not received")
			}
			if err := c.Write(websocket.TextMessage, message); err != nil {
				c.logger.Error(fmt.Errorf("websocket message recieved but Write error: %w", err))
				return err
			}
			fmt.Printf("user %s connection wrote message %s to myself \n", c.userID, message) // todo remove

		case <-ticker.C:
			if err := c.Write(websocket.PingMessage, []byte{}); err != nil {
				c.logger.Fatal(fmt.Errorf("websocket ping message Write error:  %w", err))
				return err
			}
		}
	}
}

// Write writes a message with the given message type and payload
func (c WebsocketClient) Write(mt int, payload []byte) error {
	err := c.wsConn.SetWriteDeadline(time.Now().Add(writeWait))
	if err != nil {
		return err
	}
	return c.wsConn.WriteMessage(mt, payload)
}
