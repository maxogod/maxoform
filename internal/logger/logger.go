package logger

import (
	"go.uber.org/zap"
)

var Log *zap.SugaredLogger

// TODO: add --quiet flag
func Init(dev bool) error {
	var cfg zap.Config

	if dev {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}

	l, err := cfg.Build()
	if err != nil {
		return err
	}

	Log = l.Sugar()
	return nil
}

func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
