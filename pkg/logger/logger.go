package logger

import (
	"fmt"
	"log"
	"net"
	"os"
	"runtime/debug"
	"strings"
)

var (
	hostname, _ = os.Hostname()
	srcAddr, _  = getLocalIP()
	logPrefix   = fmt.Sprintf("%s |ziglu.com|news|%s|", hostname, os.Getenv("VERSION"))
	logSuffix   = fmt.Sprintf("src=%s spt=%d", srcAddr.IP, srcAddr.Port)
)

const logLine = " msg=%s"
const debugLine = " msg=%s, extra=%s"
const errLine = " msg=%s trace=%s"

func Info(action string, msg interface{}, severity int) {
	l := getInfoLogConfig(action, severity)
	// Get logging level from runtime
	level := strings.ToLower(os.Getenv("LOGGING_LEVEL"))

	if level == "error" {
		return
	}

	l.Printf(logLine, msg)
}

func Debug(action string, msg interface{}, extra interface{}, severity int) {
	l := getDebugLogConfig(action, severity)

	// Get logging level from runtime
	level := strings.ToLower(os.Getenv("LOGGING_LEVEL"))

	if level != "debug" {
		return
	}

	l.Printf(debugLine, msg, extra)
}

func Error(action string, msg error, severity int) {
	trace := fmt.Sprintf(errLine, msg.Error(), debug.Stack())
	l := getErrorLogConfig(action, severity)
	err := l.Output(2, trace)
	if err != nil {
		fmt.Print(err.Error())
	}
}

// getInfoLogConfig Info log format LEVEL YYYY/mm/dd message
func getInfoLogConfig(action string, severity int) *log.Logger {
	details := fmt.Sprintf("%s:%s|INFO|%d|%s", logPrefix, action, severity, logSuffix)
	return log.New(os.Stdout, details, log.Lmsgprefix|log.Ldate|log.Ltime|log.LUTC)
}

// getDebugLogConfig Info log format LEVEL YYYY/mm/dd message
func getDebugLogConfig(action string, severity int) *log.Logger {
	details := fmt.Sprintf("%s:%s|DEBUG|%d|%s", logPrefix, action, severity, logSuffix)
	return log.New(os.Stdout, details, log.Lmsgprefix|log.Ldate|log.Ltime|log.LUTC)
}

// getErrorLogConfig Info log format LEVEL YYYY/mm/dd message
func getErrorLogConfig(action string, severity int) *log.Logger {
	details := fmt.Sprintf("%s:%s|ERROR|%d|%s", logPrefix, action, severity, logSuffix)
	return log.New(os.Stderr, details, log.Lmsgprefix|log.Ldate|log.Ltime|log.Lshortfile)
}

// getLocalIP Get IP address of server. Uses UDP so it never connects or cares if it does
func getLocalIP() (net.UDPAddr, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return net.UDPAddr{}, err
	}

	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return *localAddr, nil

}
