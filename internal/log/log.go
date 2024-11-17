package log

import (
	"fmt"
	"go-starter-kit/internal/server/config"
)

type Logger interface {
	WithPrefix(prefix string) Logger
	WithFields(fields map[string]interface{}) Logger

	Debug(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Panic(args ...interface{})
	Fatal(args ...interface{})

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Println(args ...interface{})
	Warnln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})
}

func prefixHelper(prefix string, s []interface{}) []interface{} {
	if len(s) == 0 {
		return []interface{}{prefix}
	}

	s[0] = fmt.Sprintf("%s: %v", prefix, s[0])
	return s
}

func NewLogger(cfg *config.Config) (Logger, error) {
	switch cfg.Log.Core {
	case "logrus":
		return newLorusLogger(cfg)
	default:
		return nil, nil
	}
}
