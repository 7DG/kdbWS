package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type cmdArgs struct {
	kdbhost         *string
	kdbport         *int
	kdbauth         *string
	wshost          *string
	wspath          *string
	wsauthtype      *string
	wsauth          *string
	useTLS          *bool
	tlskeyfile      *string
	tlscertfile     *string
	proclogfile     *string
	onInitCallback  *string
	onMsgCallback   *string
	onAckCallback   *string
	onErrorCallback *string
	onCloseCallback *string
	tlsConfig       *tls.Config
}

var masterConfig = &cmdArgs{}
var onFinish = func() {}

// required flags: kdbhost, kdbport, wshost,
// conditional flags:
//   useSSL -> tlskeyfile, tlscertfile
// optional flags: kdbauth, wspath, proclogfile, onInitCallback, onMsgCallback, onAckCallback, onCloseCallback, onErrorCallback
// at least one of the following must be present: onInitCallback, onMsgCallback, onCloseCallback

func parseCmdArgs() error {
	masterConfig.kdbhost = flag.String("kdbhost", "", "Host address of KDB+ master process")
	masterConfig.kdbport = flag.Int("kdbport", 0, "Listening port of KDB+ master process")
	masterConfig.wshost = flag.String("wshost", "", "Host address of Websocket target")
	masterConfig.kdbauth = flag.String("kdbauth", "", "OPTIONAL: auth of KDB+ master process to connect to (user:pass)")
	masterConfig.wspath = flag.String("wspath", "", "OPTIONAL: Path of Websocket endpoint")
	masterConfig.wsauth = flag.String("wsauth", "", "OPTIONAL: Auth for the Websocket endpoint. REQUIREFLAGS: wsauthtype")
	masterConfig.wsauthtype = flag.String("wsauthtype", "", "OPTIONAL: Auth type (Basic or Bearer) for the Websocket endpoint. REQUIREFLAGS: wsauth")
	masterConfig.useTLS = flag.Bool("useTLS", false, "OPTIONAL: Flag to use TLS Websocket connection. REQUIREFLAGS: tlskeyfile, tlscertflags")
	masterConfig.tlskeyfile = flag.String("tlskeyfile", "", "OPTIONAL: TLS key to use for Secure WS connection")
	masterConfig.tlscertfile = flag.String("tlscertfile", "", "OPTIONAL: TLS cert to use for Secure WS connection")
	masterConfig.proclogfile = flag.String("proclogfile", "", "OPTIONAL: File to write process output/error logs to")
	masterConfig.onInitCallback = flag.String("onInitCallback", "", "OPTIONAL: Name of callback function in kdb+ to call after WebSocket is opened")
	masterConfig.onMsgCallback = flag.String("onMsgCallback", "", "OPTIONAL: Name of callback function in kdb+ to call when a message is received")
	masterConfig.onAckCallback = flag.String("onAckCallback", "", "OPTIONAL: Name of callback function in kdb+ to call to acknowledge a message was sent")
	masterConfig.onErrorCallback = flag.String("onErrorCallback", "", "OPTIONAL: Name of callback function in kdb+ to call after an error sending data to Websocket")
	masterConfig.onCloseCallback = flag.String("onCloseCallback", "", "OPTIONAL: Name of callback function in kdb+ to call after WebSocket is closed (before exitting)")
	flag.Parse()

	switch {
	case *masterConfig.kdbhost == "":
		return fmt.Errorf("kdbhost flag not defined")
	case *masterConfig.kdbport == 0:
		return fmt.Errorf("kdbport flag not defined or is 0")
	case *masterConfig.wshost == "":
		return fmt.Errorf("wshost flag not defined")
	case *masterConfig.useTLS && *masterConfig.tlskeyfile == "":
		return fmt.Errorf("useTLS flag is provided but tlskeyfile is not defined")
	case *masterConfig.useTLS && *masterConfig.tlscertfile == "":
		return fmt.Errorf("useTLS flag is provided but tlscertfile is not defined")
	case *masterConfig.onInitCallback == "" && *masterConfig.onMsgCallback == "" && *masterConfig.onCloseCallback == "":
		return fmt.Errorf("no kdb+ callbacks defined; at least one of the following must be defined: onInitCallback, onMsgCallback, onCloseCallback")

	case *masterConfig.wsauth == "" && *masterConfig.wsauthtype != "":
		return fmt.Errorf("wsauthtype is provided but no wsauth is defined")
	case *masterConfig.wsauth != "" && *masterConfig.wsauthtype == "":
		return fmt.Errorf("wsauth is provided but no wsauthtype is defined")
	case *masterConfig.wsauthtype != "" && *masterConfig.wsauthtype != "Bearer" && *masterConfig.wsauthtype != "Basic":
		return fmt.Errorf("wsauthtype unsupported (only Bearer or Basic auth supported)")

	default:
		_, fchk := os.Stat(*masterConfig.tlskeyfile)
		if os.IsNotExist(fchk) {
			return fmt.Errorf("tlskeyfile file does not exist: %v", *masterConfig.tlskeyfile)
		}
		_, fchk = os.Stat(*masterConfig.tlscertfile)
		if os.IsNotExist(fchk) {
			return fmt.Errorf("tlscertfile file does not exist: %v", *masterConfig.tlscertfile)
		}
		return nil
	}
}

func setup() {
	// Create null logger
	masterLog = log.New(ioutil.Discard, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	// Parse command line arguments
	err := parseCmdArgs()

	// See if a proclogfile is defined, if so redirect masterLog
	ferr := assignLogOutput()
	if ferr != nil {
		masterLog.Printf("Error creating handle to log file: %v", ferr.Error())
		onFinish()
		os.Exit(3)
	}

	// Then print error logs from parseCmdArgs call if there were any
	if err != nil {
		masterLog.Printf("FATAL: Error parsing command-line arguments: %v", err.Error())
		onFinish()
		os.Exit(4)
	}

	// Print config to log file
	printMasterConfigToLog()

	// If TLS is enabled, load cert and key
	if *masterConfig.useTLS {
		err = loadTlsConfig()
		if err != nil {
			masterLog.Printf("FATAL: Error loading TLS Certificate and Key: %v", err.Error())
			onFinish()
			os.Exit(5)
		}
		masterLog.Println("INFO: Successfully loaded TLS Certificate and Key")
	}
}

func assignLogOutput() error {
	if *masterConfig.proclogfile == "" {
		return nil
	}
	var logFile *os.File
	var ferr error

	if _, fchk := os.Stat(*masterConfig.proclogfile); os.IsNotExist(fchk) {
		logFile, ferr = os.Create(*masterConfig.proclogfile)
	} else {
		logFile, ferr = os.OpenFile(*masterConfig.proclogfile, os.O_RDWR, 0666)
	}
	if ferr != nil {
		return ferr
	}
	masterLog.SetOutput(logFile)
	onFinish = func() { logFile.Close() }
	masterLog.Printf("Beginning process logging to file %v ...\n", *masterConfig.proclogfile)
	return nil
}

func printMasterConfigToLog() {
	masterLog.Printf("INFO: Parsed command line arguments\n")
	// kdb master args
	masterLog.Printf("INFO: kdbhost: %v\n", *masterConfig.kdbhost)
	masterLog.Printf("INFO: kdbport: %v\n", *masterConfig.kdbport)
	if *masterConfig.kdbauth == "" {
		masterLog.Println("INFO: kdbauth: undefined")
	} else {
		masterLog.Println("INFO: kdbauth: defined")
	}
	// WS target args
	masterLog.Printf("INFO: wshost: %v\n", *masterConfig.wshost)
	if *masterConfig.wspath != "" {
		masterLog.Printf("INFO: wspath: %v\n", *masterConfig.wspath)
	}
	// TLS args
	masterLog.Printf("INFO: useTLS: %v\n", *masterConfig.useTLS)
	if *masterConfig.useTLS {
		masterLog.Printf("INFO: tlskeyfile: %v\n", *masterConfig.tlskeyfile)
		masterLog.Printf("INFO: tlscertfile: %v\n", *masterConfig.tlscertfile)
	}
	// callback args
	if *masterConfig.onInitCallback != "" {
		masterLog.Printf("INFO: onInitCallback: %v\n", *masterConfig.onInitCallback)
	}
	if *masterConfig.onMsgCallback != "" {
		masterLog.Printf("INFO: onMsgCallback: %v\n", *masterConfig.onMsgCallback)
	}
	if *masterConfig.onAckCallback != "" {
		masterLog.Printf("INFO: onAckCallback: %v\n", *masterConfig.onAckCallback)
	}
	if *masterConfig.onCloseCallback != "" {
		masterLog.Printf("INFO: onCloseCallback: %v\n", *masterConfig.onCloseCallback)
	}
}

func loadTlsConfig() error {
	masterLog.Println("INFO: Loading TLS Certificate and Key...")
	tlscert, err := tls.LoadX509KeyPair(*masterConfig.tlscertfile, *masterConfig.tlskeyfile)
	if err != nil {
		return err
	}
	tlscfg := &tls.Config{Certificates: []tls.Certificate{tlscert}}
	masterConfig.tlsConfig = tlscfg
	return nil
}
