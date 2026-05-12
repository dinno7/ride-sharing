package main

import "github.com/gorilla/websocket"

type driverClient struct {
	conn *websocket.Conn
}

func newDriverClient(conn *websocket.Conn) *driverClient {
	return &driverClient{conn: conn}
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

func (dc *driverClient) ReadMessage() (string, error) {
	_, msg, err := dc.conn.ReadMessage()
	if err != nil {
		return "", err
	}
	return string(msg), nil
}

// -------------- Rider

type riderClient struct {
	conn *websocket.Conn
}

func newRiderClient(conn *websocket.Conn) *riderClient {
	return &riderClient{conn: conn}
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

func (dc *riderClient) ReadMessage() (string, error) {
	_, msg, err := dc.conn.ReadMessage()
	if err != nil {
		return "", err
	}
	return string(msg), nil
}
