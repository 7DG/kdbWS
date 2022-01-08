package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
	kdb "github.com/sv/kdbgo"
)

var masterLog *log.Logger

func openWebsocket() (*websocket.Conn, error) {
	scheme := "ws"
	// Use specified TLS key/cert if useTLS is declared
	if *masterConfig.useTLS {
		websocket.DefaultDialer.TLSClientConfig = masterConfig.tlsConfig
		scheme = "wss"
	}
	u := url.URL{Scheme: scheme, Host: *masterConfig.wshost, Path: *masterConfig.wspath}
	var reqheader http.Header

	// Auth handling
	if *masterConfig.wsauthtype != "" {
		reqheader = http.Header{}
		reqheader.Set("Authorization", *masterConfig.wsauthtype+" "+base64.StdEncoding.EncodeToString([]byte(*masterConfig.wsauth)))
	}

	masterLog.Printf("INFO: Connecting to WebSocket target at %v\n", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), reqheader)
	return c, err
}

func main() {
	var err error
	// Create logger
	setup()

	masterLog.Printf("INFO: Opening connection to kdb+ process at %v:%v ...\n", *masterConfig.kdbhost, *masterConfig.kdbport)
	kdbHandle, err := kdb.DialKDB(*masterConfig.kdbhost, *masterConfig.kdbport, *masterConfig.kdbauth)
	if err != nil {
		masterLog.Printf("FATAL: Error connecting to kdb+ process: %v", err)
		onFinish()
		os.Exit(6)
	}
	masterLog.Printf("INFO: Successfully connected to kdb+ process\n")
	// Open WebSocket to target
	wsHandle, err := openWebsocket()
	if err != nil {
		masterLog.Printf("FATAL: Error connecting to Websocket target: %v", err)
		onFinish()
		os.Exit(7)
	}
	masterLog.Println("INFO: Successfully connected to WebSocket target")

	// Send initialisation success signal to KDB
	sendInitCallback(kdbHandle)

	// Create channels
	wsResChannel := make(chan []byte)
	wsErrorChannel := make(chan error)
	kdbResChannel := make(chan *kdb.K)
	kdbErrorChannel := make(chan error)

	// Beginner dual listener
	go kdbListener(kdbHandle, kdbErrorChannel, kdbResChannel)
	go wsListener(wsHandle, wsErrorChannel, wsResChannel)
	// Begin Listening
	for {
		select {

		case err = <-kdbErrorChannel:
			masterLog.Printf("FATAL: Error reading from kdb+ handle: %v", err)
			onFinish()
			os.Exit(9)

		case err = <-wsErrorChannel:
			masterLog.Printf("FATAL: Error reading from WebSocket handle: %v", err)
			sendCloseCallback(kdbHandle)
			onFinish()
			os.Exit(10)

		case res := <-kdbResChannel:
			err = handleKDBMessage(wsHandle, res)
			if err != nil {
				sendErrorCallback(kdbHandle, err)
			} else {
				sendAckCallback(kdbHandle)
			}

		case res := <-wsResChannel:
			sendMsgCallback(kdbHandle, res)
		}
	}
}
