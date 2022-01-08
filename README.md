# kdbWS
WebSocket feedhandler for kdb+

# Config
kdbWS config is controlled by command line arguments:
| Flag | Required | Description |
|------|----------|-------------|
| kdbhost | Yes | Hostname (or IP) of kdb+ controller process
| kdbport | Yes | Port of kdb+ controller process
| kdbauth | No | Auth of kdb+ controller process (in Basic format e.g. 'user:pass')
| wshost | Yes | Host of target WebSocket (host:port or domain)
| wspath | No | Path of WebSocket endpoint on target host
| wsauth | Conditional: If wsauthtype is declared | Auth of WebSocket target (Basic format for Basic auth e.g. 'user:pass' or Bearer tokens as they are)
| wsauthtype | Conditional: If wsauth is declared | 'Basic' or 'Bearer' supported
| useTLS | No | Flag for whether or not to use TLS for WebSocket NOTE: Do not supply with any arguments, just the flag (i.e. "... -prevFlag prevFlagValue -useTLS -nextFlag nextFlagValue ...")
| tlskeyfile | Conditional: If useTLS is True | Location of the TLS key file
| tlscertfile | Conditional: If useTLS is True | Location of the TLS cert file
| proclogfile | No | Location to write output logs to; if not supplied no logs will be written
| onInitCallback | No* | kdb+ function to use for initialisation callback. Must be unary function with argument type -2
| onMsgCallback | No* | kdb+ function to use for initialisation callback. Must be unary function with argument type 10
| onAckCallback | No* | kdb+ function to use for initialisation callback. Must be unary function with argument type -2
| onErrorCallback | No | kdb+ function to use for initialisation callback. Must be unary function with argument type -128
| onCloseCallback | No | kdb+ function to use for initialisation callback. Must be unary function with argument type -2

*Note: At least one of the following callbacks must be defined: onInitCallback, onMsgCallback, onAckCallback
