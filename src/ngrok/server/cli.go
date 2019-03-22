package server

import (
	"flag"
	"strconv"
	"strings"
)

//jimmy add port range config
type Range struct {
	min int
	max int
}

type Options struct {
	httpAddr   string
	httpsAddr  string
	tunnelAddr string
	domain     string
	tlsCrt     string
	tlsKey     string
	logto      string
	loglevel   string
	portRange  Range
	restPort   string
}

func parseArgs() *Options {
	httpAddr := flag.String("httpAddr", ":80", "Public address for HTTP connections, empty string to disable")
	httpsAddr := flag.String("httpsAddr", ":443", "Public address listening for HTTPS connections, emptry string to disable")
	tunnelAddr := flag.String("tunnelAddr", ":4443", "Public address listening for ngrok client")
	domain := flag.String("domain", "dbc.com", "Domain where the tunnels are hosted")
	tlsCrt := flag.String("tlsCrt", "", "Path to a TLS certificate file")
	tlsKey := flag.String("tlsKey", "", "Path to a TLS key file")
	logto := flag.String("log", "stdout", "Write log messages to this file. 'stdout' and 'none' have special meanings")
	loglevel := flag.String("log-level", "DEBUG", "The level of messages to log. One of: DEBUG, INFO, WARNING, ERROR")
	portRange := flag.String("portRange", "20000-20100", "The tunnel port range")
	restPort := flag.String("restPort", "8000", "The restful port")
	flag.Parse()

	port := strings.Split(*portRange, "-")
	min, max := 20000, 20100
	var err error

	if len(port) == 2 {

		if min, err = strconv.Atoi(port[0]); err != nil {
			min = 20000
		}

		if max, err = strconv.Atoi(port[1]); err != nil {
			max = min + 100
		}

		if min > max {
			max = min + 100
		}
	}

	return &Options{
		httpAddr:   *httpAddr,
		httpsAddr:  *httpsAddr,
		tunnelAddr: *tunnelAddr,
		domain:     *domain,
		tlsCrt:     *tlsCrt,
		tlsKey:     *tlsKey,
		logto:      *logto,
		loglevel:   *loglevel,
		portRange:  Range{min, max},
		restPort:   *restPort,
	}
}
