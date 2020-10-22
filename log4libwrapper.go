package log4libwrapper

import (
	"github.com/anthonyraymond/go-log4lib"
	"go.uber.org/zap"
)

type zapLoggerWrapper struct {
	delegate *zap.SugaredLogger
}

func (w *zapLoggerWrapper) Info(args ...interface{}) {
	w.delegate.Info(args)
}

func (w *zapLoggerWrapper) Warn(args ...interface{}) {
	w.delegate.Warn(args)
}

func (w *zapLoggerWrapper) Error(args ...interface{}) {
	w.delegate.Error(args)
}

func (w *zapLoggerWrapper) Panic(args ...interface{}) {
	w.delegate.Panic(args)
}

func (w *zapLoggerWrapper) Fatal(args ...interface{}) {
	w.delegate.Fatal(args)
}


func (w *zapLoggerWrapper) Infof(template string, args ...interface{}) {
	w.delegate.Infof(template, args)
}

func (w *zapLoggerWrapper) Warnf(template string, args ...interface{}) {
	w.delegate.Warnf(template, args)
}

func (w *zapLoggerWrapper) Errorf(template string, args ...interface{}) {
	w.delegate.Errorf(template, args)
}

func (w *zapLoggerWrapper) Panicf(template string, args ...interface{}) {
	w.delegate.Panicf(template, args)
}

func (w *zapLoggerWrapper) Fatalf(template string, args ...interface{}) {
	w.delegate.Fatalf(template, args)
}

func WrapZapLogger(pointerToLogger *zap.SugaredLogger) log4lib.LibLogger {
	return &zapLoggerWrapper{delegate: pointerToLogger}
}
