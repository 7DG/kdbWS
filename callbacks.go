package main

import (
	"os"

	kdb "github.com/sv/kdbgo"
)

func sendInitCallback(h *kdb.KDBConn) {
	if *masterConfig.onInitCallback == "" {
		return
	}
	masterLog.Println("INFO: Sending Init callback to kdb+ process...")
	err := h.WriteMessage(kdb.ASYNC, kdb.NewList(&kdb.K{kdb.KC, kdb.NONE, *masterConfig.onInitCallback}, &kdb.K{-kdb.KB, kdb.NONE, true}))
	if err != nil {
		masterLog.Printf("FATAL: Error sending Init callback to kdb+ process: %v", err)
		onFinish()
		os.Exit(10)
	}
	masterLog.Println("INFO: Successfully sent Init callback to kdb+")
}

func sendMsgCallback(h *kdb.KDBConn, data []byte) {
	if *masterConfig.onMsgCallback == "" {
		return
	}
	masterLog.Println("INFO: Sending Msg callback to kdb+ process...")
	err := h.WriteMessage(kdb.ASYNC, kdb.NewList(&kdb.K{kdb.KC, kdb.NONE, *masterConfig.onMsgCallback}, &kdb.K{kdb.KC, kdb.NONE, string(data)}))
	if err != nil {
		masterLog.Printf("FATAL: Error sending Msg callback to kdb+ process: %v", err)
		onFinish()
		os.Exit(11)
	}
	masterLog.Println("INFO: Successfully sent Msg callback to kdb+")
}

func sendAckCallback(h *kdb.KDBConn) {
	if *masterConfig.onAckCallback == "" {
		return
	}
	masterLog.Println("INFO: Sending Ack callback to kdb+ process...")
	err := h.WriteMessage(kdb.ASYNC, kdb.NewList(&kdb.K{kdb.KC, kdb.NONE, *masterConfig.onAckCallback}, &kdb.K{-kdb.KB, kdb.NONE, true}))
	if err != nil {
		masterLog.Printf("FATAL: Error sending Ack callback to kdb+ process: %v", err)
		onFinish()
		os.Exit(12)
	}
	masterLog.Println("INFO: Successfully sent Ack callback to kdb+")
}

func sendCloseCallback(h *kdb.KDBConn) {
	if *masterConfig.onCloseCallback == "" {
		return
	}
	masterLog.Println("INFO: Sending Close callback to kdb+ process...")
	err := h.WriteMessage(kdb.ASYNC, kdb.NewList(&kdb.K{kdb.KC, kdb.NONE, *masterConfig.onCloseCallback}, &kdb.K{-kdb.KB, kdb.NONE, true}))
	if err != nil {
		masterLog.Printf("FATAL: Error sending Close callback to kdb+ process: %v", err)
		onFinish()
		os.Exit(13)
	}
	masterLog.Println("INFO: Successfully sent Close callback to kdb+")
}

func sendErrorCallback(h *kdb.KDBConn, err error) {
	if *masterConfig.onErrorCallback == "" {
		return
	}
	masterLog.Println("INFO: Sending Error callback to kdb+ process...")
	serr := h.WriteMessage(kdb.ASYNC, kdb.NewList(&kdb.K{kdb.KC, kdb.NONE, *masterConfig.onErrorCallback}, kdb.Error(err)))
	if serr != nil {
		masterLog.Printf("FATAL: Error sending Close callback to kdb+ process: %v", serr)
		onFinish()
		os.Exit(14)
	}
	masterLog.Println("INFO: Successfully sent Error callback to kdb+")
}
