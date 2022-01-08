// Ensure this script is started with q example.q -E 1 -p XXXXX

// load config
\l exampleConfig.q

// scripts
wshandle:0i;
tcphandle:0i;

.cfg.kdbport:system"p";
if[.cfg.kdbport=0;
  0N!"NO PORT ASSIGNED, MUST START KDB+ WITH A LISTENING PORT";
  0N!"EXITTING...";
  exit 3;
  ];

.z.ws:{[x] 
  show "RECEIVED MESSAGE FROM WEBSOCKET:";
  show .Q.s1 x;
  :.j.j `time`ack!(.z.p;1b);
  };

.z.wo:{[h] wshandle::h;show "KDB'S WS LISTENER OPENED A CONNECTION";};
.z.po:{[h] tcphandle::h;show "KDB'S TCP LISTENER OPENED A CONNECTION";};
.z.pc:{[h] tcphandle::0i;show "KDB'S TCP LISTENER CLOSED A CONNECTION";};

.z.pw:{[u;p]
  if[not (`kdbWSuser;"kdbWSpass")~(u;p);:0b];
  :1b;
  };

.cfg.kdbhost:hostname;
.cfg.kdbport:string .cfg.kdbport;
.cfg.useTLS:"";
.cfg.tlskeyfile:getenv[`SSL_KEY_FILE];
.cfg.tlscertfile:getenv[`SSL_CERT_FILE];
.cfg.kdbauth:"kdbWSuser:kdbWSpass";
.cfg.wshost:.cfg.kdbhost,":",.cfg.kdbport;
.cfg.wsauthtype:"Basic";
.cfg.wsauth:.cfg.kdbauth;
.cfg.proclogfile:kdbWSlog;
.cfg.onInitCallback:"initCallback";
.cfg.onMsgCallback:"msgCallback";
.cfg.onAckCallback:"ackCallback";
.cfg.onErrorCallback:"errorCallback";
.cfg.onCloseCallback:"closeCallback";

initCallback:{[success]
  show "initCallback: RECEIVED WEBSOCKET INIT SIGNAL FROM KDBWS";
  show "initCallback: SENDING EXAMPLE MESSAGE TO KDBWS:";
  show "initCallback: (`message;.j.j `example`object!(\"message\";21))";
  sendmessagetokdbWS .j.j `example`object!("message";21);
  };

msgCallback:{[msg]
  show "msgCallback: RECEIVED WEBSOCKET MESSAGE FROM KDBWS:";
  show "msgCallback: ",msg;
  };

ackCallback:{[success]
  show "ackCallback: RECEIVED ACK SIGNAL FROM KDBWS";
  show "ackCallback: CLOSING KDB'S WEBSOCKET CONNECTION TO KDBWS, THIS WILL TRIGGER closeCallback:";
  hclose first key [.z.W] except .z.w
  };

errorCallback:{[err]
  show "errorCallback: RECEIVED WEBSOCKET ERROR SIGNAL FROM KDBWS:";
  show "errorCallback: ",err;
  };

closeCallback:{[success]
  show "closeCallback: RECEIVED WEBSOCKET CLOSE SIGNAL FROM KDBWS";
  };

startkdbWS:{[]
  buildstartcmd:{[x] "nohup ",kdbWSbinary," ",x," >/dev/null 2>&1 &"};
  flags:{[x]
    x:(string key x;value x);
    :" " sv "-",/:x[0],'" ",/:x[1];
    }[1_.cfg];
  0N!"STARTING KDBWS...";
  0N!startline:buildstartcmd flags;
  system startline;
  0N!"KDBWS STARTED";
  };

sendmessagetokdbWS:{[msg]
  if[tcphandle=0;'"NO KDBWS CONNECTION"];
  0N!"SENDING MESSAGE TO KDBWS: ",msg;
  neg[tcphandle](`message;msg);neg[tcphandle][];
  };

startkdbWS[];
