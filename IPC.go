package main

import (
	"fmt"

	"github.com/gorilla/websocket"
	kdb "github.com/sv/kdbgo"
)

func kdbListener(handle *kdb.KDBConn, echan chan error, channel chan *kdb.K) {
	for {
		res, _, err := handle.ReadMessage()
		if err != nil {
			echan <- err
			return
		} else {
			channel <- res
		}
	}
}

func wsListener(handle *websocket.Conn, echan chan error, channel chan []byte) {
	for {
		_, res, err := handle.ReadMessage()
		if err != nil {
			echan <- err
			return
		} else {
			channel <- res
		}
	}
}

func handleKDBMessage(h *websocket.Conn, data *kdb.K) error {
	if data.Type != kdb.K0 {
		masterLog.Printf("ERROR: Received unsupported (object-type) message from kdb+: %v\n", data.Data)
		return fmt.Errorf("object-type")
	}
	if data.Len() != 2 {
		masterLog.Printf("ERROR: Received unsupported (object-length) message from kdb+: %v\n", data.Data)
		return fmt.Errorf("object-length")
	}

	innerdata := data.Data.([]*kdb.K)
	methodK := innerdata[0]
	msgK := innerdata[1]
	if methodK.Type != -kdb.KS {
		masterLog.Printf("ERROR: Received unsupported (method-type) message from kdb+: %v\n", methodK.Data)
		return fmt.Errorf("method-type")
	}
	if msgK.Type != kdb.KC {
		masterLog.Printf("ERROR: Received unsupported (message-type) message from kdb+: %v\n", msgK.Data)
		return fmt.Errorf("message-type")
	}

	method := methodK.Data.(string)
	msg := msgK.Data.(string)

	switch {
	case method == "message":
		masterLog.Printf("INFO: \"message\" method called from kdb+, sending data to Websocket target...\n")
		err := h.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			masterLog.Printf("ERROR: Error sending data to Websocket target: %v\n", err)
			return err
		}
		masterLog.Printf("INFO: Successfully sent data to Websocket target\n")
		return nil
	default:
		masterLog.Printf("ERROR: Received unsupported (method) message from kdb+: %v\n", method)
		return fmt.Errorf("unsupported method")
	}

}
