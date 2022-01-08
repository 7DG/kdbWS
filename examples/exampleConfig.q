// USER CONFIG

// provide the path (absolute or relative) to the kdbWS binary
kdbWSbinary:"../kdbWS";

// provide the hostname of the machine (as it appears on the TLS certificate)
hostname:""

// provide the path (absolute or relative) of where to write the kdbWS process logs to
kdbWSlog:$[.z.o like "w*";first[system"echo %cd%"],"\\";first[system"pwd"],"/"],"exampleLogFile.log";

\c 100 1000
