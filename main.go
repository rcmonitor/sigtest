package main

import (
	l "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"time"
)


type TData struct {
	Counter int
}

func init() {
	fInitLogger()
}



func main() {
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
			l.WithFields(l.Fields{
				"Signal": sigIncoming.String(),
			}).Warn("Interrupted")
			os.Exit(0)

		default:
			mainRoutine(d.Counter)
		}
	}
}

func mainRoutine(i int) {
		l.WithFields(l.Fields{
			"Step": i,
		}).Info("Cycle")

		time.Sleep(2 * time.Second)
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