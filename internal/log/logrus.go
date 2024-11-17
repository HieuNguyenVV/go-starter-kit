package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"go-starter-kit/internal/server/config"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type logrusLogger struct {
	prefix string
	*logrus.Entry
}

func newLorusLogger(cfg *config.Config) (Logger, error) {
	logger := logrus.New()
	level, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		return nil, err
	}
	logger.SetLevel(level)

	logFormat := strings.ToLower(cfg.Log.Format)
	switch logFormat {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		})
	default:
		logger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: time.RFC3339Nano,
			FullTimestamp:   true,
		})
	}

	output := strings.ToLower(cfg.Log.Output)
	switch output {
	case "stdout":
		logger.SetOutput(os.Stdout)
	case "discard":
		logger.SetOutput(ioutil.Discard)
	default:
		logger.SetOutput(os.Stderr)

	}

	return &logrusLogger{Entry: logrus.NewEntry(logger)}, nil
}

func (l *logrusLogger) WithFields(fields map[string]interface{}) Logger {
	return &logrusLogger{
		prefix: l.prefix,
		Entry:  logrus.NewEntry(l.Entry.Logger).WithFields(l.Entry.Data).WithFields(fields),
	}
}

func (l *logrusLogger) WithPrefix(prefix string) Logger {
	if l.prefix != "" {
		prefix = fmt.Sprintf("%s/%s", l.prefix, prefix)
	}
	return &logrusLogger{
		prefix: prefix,
		Entry:  l.Entry,
	}
}

func (l *logrusLogger) Debug(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Debug(prefixHelper(l.prefix, args)...)
	}
	l.Entry.Debug(args...)
}

func (l *logrusLogger) Info(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Info(prefixHelper(l.prefix, args)...)
	}
	l.Entry.Info(args...)
}

func (l *logrusLogger) Print(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Print(prefixHelper(l.prefix, args)...)
	}
	l.Entry.Print(args...)
}

func (l *logrusLogger) Warn(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Warn(prefixHelper(l.prefix, args)...)
	}
	l.Entry.Warn(args...)
}

func (l *logrusLogger) Error(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Error(prefixHelper(l.prefix, args)...)
	}
	l.Entry.Error(args...)
}

func (l *logrusLogger) Panic(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Panic(prefixHelper(l.prefix, args)...)
	}
	l.Entry.Panic(args...)
}

func (l *logrusLogger) Fatal(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Fatal(prefixHelper(l.prefix, args)...)
	}
	l.Entry.Fatal(args...)
}

func (l *logrusLogger) Debugf(format string, args ...interface{}) {
	if l.prefix != "" {
		format = fmt.Sprintf("%s: %s", l.prefix, format)
	}
	l.Entry.Debugf(format, args...)
}

func (l *logrusLogger) Infof(format string, args ...interface{}) {
	if l.prefix != "" {
		format = fmt.Sprintf("%s: %s", l.prefix, format)
	}
	l.Entry.Infof(format, args...)
}

func (l *logrusLogger) Printf(format string, args ...interface{}) {
	if l.prefix != "" {
		format = fmt.Sprintf("%s: %s", l.prefix, format)
	}
	l.Entry.Printf(format, args...)
}

func (l *logrusLogger) Warnf(format string, args ...interface{}) {
	if l.prefix != "" {
		format = fmt.Sprintf("%s: %s", l.prefix, format)
	}
	l.Entry.Warnf(format, args...)
}

func (l *logrusLogger) Errorf(format string, args ...interface{}) {
	if l.prefix != "" {
		format = fmt.Sprintf("%s: %s", l.prefix, format)
	}
	l.Entry.Errorf(format, args...)
}

func (l *logrusLogger) Fatalf(format string, args ...interface{}) {
	if l.prefix != "" {
		format = fmt.Sprintf("%s: %s", l.prefix, format)
	}
	l.Entry.Fatalf(format, args...)
}

func (l *logrusLogger) Panicf(format string, args ...interface{}) {
	if l.prefix != "" {
		format = fmt.Sprintf("%s: %s", l.prefix, format)
	}
	l.Entry.Panicf(format, args...)
}

func (l *logrusLogger) Debugln(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Debugln(prefixHelper(l.prefix, args)...)
	}
	l.Entry.Debugln(args...)
}

func (l *logrusLogger) Infoln(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Infoln(prefixHelper(l.prefix, args)...)
	}
	l.Entry.Infoln(args...)
}

func (l *logrusLogger) Println(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Println(prefixHelper(l.prefix, args)...)
	}
	l.Entry.Println(args...)
}

func (l *logrusLogger) Warnln(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Warnln(prefixHelper(l.prefix, args)...)
	}
	l.Entry.Warnln(args...)
}

func (l *logrusLogger) Errorln(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Errorln(prefixHelper(l.prefix, args)...)
	}
	l.Entry.Errorln(args...)
}

func (l *logrusLogger) Panicln(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Panicln(prefixHelper(l.prefix, args)...)
	}
	l.Entry.Panicln(args...)
}

func (l *logrusLogger) Fatalln(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Fatalln(prefixHelper(l.prefix, args)...)
	}
	l.Entry.Fatalln(args...)
}
