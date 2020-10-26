package log4libwrapper

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"testing"
)

func createTempDirOrFailTest(t *testing.T) (dir string, clean func()) {
	var err error = nil
	dir, err = ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}

	clean = func() {
		_ = os.RemoveAll(dir)
	}
	return
}

type StringBufferSink struct {
	builder *strings.Builder
}

func (s *StringBufferSink) Write(p []byte) (n int, err error) {
	return s.builder.Write(p)
}

func (s *StringBufferSink) Sync() error {
	return nil
}

func (s *StringBufferSink) Close() error {
	s.builder.Reset()
	return nil
}

// this is a trick to redirect zap logs to a string builder
func newTestStringBufferSink(builder *strings.Builder) func(*url.URL) (zap.Sink, error) {
	return func(*url.URL) (zap.Sink, error) {
		return &StringBufferSink{builder: builder}, nil
	}
}

func TestWrapZapLogger_MustLog(t *testing.T) {
	logReceiver := &strings.Builder{}

	err := zap.RegisterSink("testoutputMustLog", newTestStringBufferSink(logReceiver))
	if err != nil {
		t.Fatal(err)
	}

	zapLogger, err := zap.Config{
		Level:    zap.NewAtomicLevelAt(zap.DebugLevel),
		Encoding: "console",
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:    "level",
			EncodeLevel: zapcore.LowercaseLevelEncoder,
			MessageKey:  "msg",
		},
		OutputPaths:      []string{"stdout", "testoutputMustLog:///" + "nomatterwhat"},
		ErrorOutputPaths: []string{"stdout", "testoutputMustLog:///" + "nomatterwhat"},
	}.Build()
	if err != nil {
		t.Fatal(err)
	}

	logger := WrapZapLogger(zapLogger)
	logger.Debug("coucou")
	logger.Info("coucou")
	logger.Warn("coucou")
	logger.Error("coucou")
	logger.Debug("coucou", "joe", "la bidouille")
	logger.Info("coucou", "joe", "la bidouille")
	logger.Warn("coucou", "joe", "la bidouille")
	logger.Error("coucou", "joe", "la bidouille")
	logger.Debugf("coucou %s :)", "joe")
	logger.Infof("coucou %s :)", "joe")
	logger.Warnf("coucou %s :)", "joe")
	logger.Errorf("coucou %s :)", "joe")
	logger.Debugf("coucou %s %d :)", "joe", 12)
	logger.Infof("coucou %s %d :)", "joe", 12)
	logger.Warnf("coucou %s %d :)", "joe", 12)
	logger.Errorf("coucou %s %d :)", "joe", 12)

	expected := []string{
		"debug\tcoucou",
		"info\tcoucou",
		"warn\tcoucou",
		"error\tcoucou",
		"debug\tcoucoujoela bidouille",
		"info\tcoucoujoela bidouille",
		"warn\tcoucoujoela bidouille",
		"error\tcoucoujoela bidouille",
		"debug\tcoucou joe :)",
		"info\tcoucou joe :)",
		"warn\tcoucou joe :)",
		"error\tcoucou joe :)",
		"debug\tcoucou joe 12 :)",
		"info\tcoucou joe 12 :)",
		"warn\tcoucou joe 12 :)",
		"error\tcoucou joe 12 :)",
		"",
	}

	_ = zapLogger.Sync()

	if !strings.EqualFold(strings.Join(expected, "\n"), logReceiver.String()) {
		t.Fatalf("log output is not correct, expected:\n%s\nactual:\n%s", strings.Join(expected, "\n"), logReceiver.String())
	}
}

func TestWrapZapLogger_CallerMustBeUnwrapped(t *testing.T) {
	// when the caller is added to the zap output we don't want to see the Info, Warn, ... wrapping method from zapLoggerWrapper, we want the real caller
	logReceiver := &strings.Builder{}

	err := zap.RegisterSink("testoutputCallerMustBeUnwrapped", newTestStringBufferSink(logReceiver))
	if err != nil {
		t.Fatal(err)
	}

	zapLogger, err := zap.Config{
		Level:         zap.NewAtomicLevelAt(zap.InfoLevel),
		DisableCaller: false,
		Encoding:      "console",
		EncoderConfig: zapcore.EncoderConfig{
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout", "testoutputCallerMustBeUnwrapped:///" + "nomatterwhat"},
		ErrorOutputPaths: []string{"stdout", "testoutputCallerMustBeUnwrapped:///" + "nomatterwhat"},
	}.Build()
	if err != nil {
		t.Fatal(err)
	}

	WrapZapLogger(zapLogger).Info("coucou")

	_ = zapLogger.Sync()

	if !strings.Contains(logReceiver.String(), "go-log4libwrapper-zap/log4libwrapper_test.go") {
		t.Fatal("caller is now correct, the wrapper should be ignored and the caller must be the real calling function (here go-log4libwrapper-zap/log4libwrapper_test.go)")
	}
}
