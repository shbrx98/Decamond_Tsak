package logger

import (
    "os"

    "github.com/sirupsen/logrus"
)

type Logger interface {
    Info(args ...interface{})
    Infof(format string, args ...interface{})
    Warn(args ...interface{})
    Warnf(format string, args ...interface{})
    Error(args ...interface{})
    Errorf(format string, args ...interface{})
    Fatal(args ...interface{})
    Fatalf(format string, args ...interface{})
    WithField(key string, value interface{}) Logger
    WithFields(fields map[string]interface{}) Logger
    WithError(err error) Logger
}

type logrusLogger struct {
    entry *logrus.Entry
}

func New() Logger {
    l := logrus.New()
    l.SetOutput(os.Stdout)
    l.SetFormatter(&logrus.JSONFormatter{})
    l.SetLevel(logrus.InfoLevel)
    return &logrusLogger{entry: logrus.NewEntry(l)}
}

func (l *logrusLogger) Info(a ...interface{})                    { l.entry.Info(a...) }
func (l *logrusLogger) Infof(f string, a ...interface{})         { l.entry.Infof(f, a...) }
func (l *logrusLogger) Warn(a ...interface{})                    { l.entry.Warn(a...) }
func (l *logrusLogger) Warnf(f string, a ...interface{})         { l.entry.Warnf(f, a...) }
func (l *logrusLogger) Error(a ...interface{})                   { l.entry.Error(a...) }
func (l *logrusLogger) Errorf(f string, a ...interface{})        { l.entry.Errorf(f, a...) }
func (l *logrusLogger) Fatal(a ...interface{})                   { l.entry.Fatal(a...) }
func (l *logrusLogger) Fatalf(f string, a ...interface{})        { l.entry.Fatalf(f, a...) }
func (l *logrusLogger) WithField(k string, v interface{}) Logger { return &logrusLogger{entry: l.entry.WithField(k, v)} }
func (l *logrusLogger) WithFields(m map[string]interface{}) Logger { return &logrusLogger{entry: l.entry.WithFields(m)} }
func (l *logrusLogger) WithError(err error) Logger               { return &logrusLogger{entry: l.entry.WithError(err)} }