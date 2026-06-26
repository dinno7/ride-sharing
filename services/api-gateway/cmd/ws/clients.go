package ws

import (
	wshub "github.com/dinno7/ride-sharing/shared/ws"
	"github.com/gorilla/websocket"
)

type driverClient struct {
	conn *websocket.Conn
	id   string
}

func newDriverClient(conn *websocket.Conn, id string) wshub.Client {
	return &driverClient{conn: conn, id: id}
}

func (dc *driverClient) ID() string {
	return dc.id
}

func (dc *driverClient) GetType() string {
	return "driver"
}

func (dc *driverClient) IsEqual(currentConn *websocket.Conn) bool {
	return dc.conn == currentConn
}

func (dc *driverClient) SendMessage(msg string) error {
	return dc.conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

func (dc *driverClient) SendJSON(msg any) error {
	return dc.conn.WriteJSON(msg)
}

func (dc *driverClient) ReadMessage() ([]byte, error) {
	_, msg, err := dc.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// -------------- Rider

type riderClient struct {
	id   string
	conn *websocket.Conn
}

func newRiderClient(conn *websocket.Conn, id string) wshub.Client {
	return &riderClient{conn: conn, id: id}
}

func (dc *riderClient) ID() string {
	return dc.id
}

func (dc *riderClient) GetType() string {
	return "driver"
}

func (dc *riderClient) IsEqual(currentConn *websocket.Conn) bool {
	return dc.conn == currentConn
}

func (dc *riderClient) SendMessage(msg string) error {
	return dc.conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

func (dc *riderClient) SendJSON(msg any) error {
	return dc.conn.WriteJSON(msg)
}

func (dc *riderClient) ReadMessage() ([]byte, error) {
	_, msg, err := dc.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	return msg, nil
}
