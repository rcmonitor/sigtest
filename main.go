package main

import (
	l "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)


type TData struct {
	Counter int
}

func init() {
	fLoadConfig()
	fInitLogger()

	strWD, err := os.Getwd()
	if nil != err {
		l.WithField("Working directory", strWD).Info("init")
	}else{
		l.Error(err)
	}
}

var GitCommit, BuildDate, Version string
var delay int

func main() {

	l.WithFields(l.Fields{
		"Commit": GitCommit,
		"Build Date": BuildDate,
		"Version": Version,
	}).Info()

	d := &TData{}

	cSignal := make(chan os.Signal, 1)
	signal.Notify(cSignal)

	mainLoop(cSignal, d)
}

func mainLoop(cSignal chan os.Signal, d *TData) {
	for true {
		d.Counter ++
		select {
		case sigIncoming := <-cSignal:
			if !fHandleSignal(sigIncoming) { os.Exit(0) }

		default:
			mainRoutine(d.Counter)
		}
	}
}

func mainRoutine(i int) {
		l.WithFields(l.Fields{
			"Step": i,
			"Delay": delay,
		}).Info("Cycle")

		time.Sleep(time.Duration(delay) * time.Second)
}

func fInitLogger() {
	fLog, err := os.OpenFile("/var/log/sigtest/sigtest.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		l.SetOutput(fLog)
	} else {
		l.Info("Failed to log to file, using default stderr")
	}
}

func fHandleSignal(sigIncoming os.Signal) (bContinue bool) {
	l.WithField("Signal", sigIncoming.String()).Info("Got incoming", )
	switch sigIncoming {
	case syscall.SIGTERM:
		l.Info("Got Termination signal, finalizing")
	case syscall.SIGHUP:
		l.Info("Got Reload signal")
		fLoadConfig()
		bContinue = true

	default:
		l.Info("Unknown signal, quit")
	}

	return
}

func fLoadConfig() {
	l.Info("Loading configuration")
	delay = fGetDelay()
}

func fGetDelay() (intDelay int) {


	intDelay = 2
	strDelay := os.Getenv("SIGTEST_DELAY")
	if "" == strDelay {
		l.Info("no '$SIGTEST_DELAY' provided")
		return
	}
	intTemp, err := strconv.ParseInt(strDelay, 10, 0)
	if nil != err {
		l.Error("please, provide $SIGTEST_DELAY in numeric format")
	}else{ intDelay = int(intTemp) }

	return
}