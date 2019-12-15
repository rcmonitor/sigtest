package main

import (
	l "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
	"time"
)


type TData struct {
	Counter int
}

func init() {
	fInitLogger()

	fInitConfig()

	fGetConfig()

}

var GitCommit, BuildDate, Version string
var delay int
var workingDir string

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
	for {
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
	//	temporarily add SIGINT to handle "Ctrl + C" as reload command
	case syscall.SIGHUP:
		l.Info("Got Reload signal")
		fGetConfig()
		bContinue = true

	default:
		l.Info("Unknown signal, quit")
	}

	return
}

func fInitConfig() {
	l.Info("Loading configuration")

	var err error
	workingDir, err = os.Getwd()
	if nil != err {
		l.Error(err)
	}else{
		l.WithField("Working directory", workingDir).Info("init")
	}


	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath(workingDir)
	viper.SetConfigType("hcl")
	viper.SetDefault("delay", 2)


}

func fGetConfig() () {

	err := viper.ReadInConfig()

	if nil != err {
		l.Error(err)
	}

	delay = viper.GetInt("delay")
}