package logger

import (
	"context"
	"fmt"
	"os"
	"time"
)

type Logger struct {
	log chan<- string
}

func (l *Logger) Info(text string) {
	l.log <- fmt.Sprintf("Info: %s \n", text)
}

func (l *Logger) Warning(text string) {
	l.log <- fmt.Sprintf("Warning: %s \n", text)
}

func (l *Logger) Error(text string) {
	l.log <- fmt.Sprintf("Error: %s \n", text)
}

func SetupLogger(ctx context.Context, directory string) Logger {
	channel := make(chan string)
	logger := Logger{log: channel}
	go func() {
		filetime := fmt.Sprintf("%s/%d.txt", directory, time.Now().UnixNano())
		file, err := os.OpenFile(filetime, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			panic("fuck")
		}
		for {
			select {
			case <-ctx.Done():
				file.Close()
				return
			case str := <-channel:
				fmt.Fprintf(file, "%s", str)
			}
		}

	}()
	return logger
}
