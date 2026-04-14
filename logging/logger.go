package logging

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/mrbagir/appfr/version"
)

// PrettyPrint defines an interface for objects that can render
// themselves in a human-readable format to the provided writer.
type PrettyPrint interface {
	PrettyPrint(writer io.Writer)
}

// Logger defines the interface for structured logging across
// different log levels such as Debug, Info, Warn, Error, and Fatal.
// It allows formatted logs and log level control.
type Logger interface {
	Debug(args ...any)
	Debugf(format string, args ...any)
	Log(args ...any)
	Logf(format string, args ...any)
	Info(args ...any)
	Infof(format string, args ...any)
	Notice(args ...any)
	Noticef(format string, args ...any)
	Warn(args ...any)
	Warnf(format string, args ...any)
	Error(args ...any)
	Errorf(format string, args ...any)
	Fatal(args ...any)
	Fatalf(format string, args ...any)
	ChangeLevel(level Level)
}

type logger struct {
	level      Level
	output     io.Writer
	isTerminal bool
	lock       sync.Mutex
}

type logEntry struct {
	Level       Level     `json:"level"`
	Time        time.Time `json:"time"`
	Message     any       `json:"message"`
	TraceID     string    `json:"trace_id,omitempty"`
	CoreVersion string    `json:"coreVersion"`
}

func (l *logger) logf(level Level, format string, args ...any) {
	if level < l.level {
		return
	}

	entry := logEntry{
		Level:       level,
		Time:        time.Now(),
		CoreVersion: version.Appcore,
	}

	traceID, filteredArgs := extractTraceIDAndFilterArgs(args)
	entry.TraceID = traceID

	switch {
	case len(filteredArgs) == 1 && format == "":
		entry.Message = filteredArgs[0]
	case len(filteredArgs) != 1 && format == "":
		entry.Message = filteredArgs
	case format != "":
		entry.Message = fmt.Sprintf(format, filteredArgs...)
	}

	if l.isTerminal {
		l.prettyPrint(&entry, l.output)
	} else {
		_ = json.NewEncoder(l.output).Encode(entry)
	}
}

func (l *logger) Debug(args ...any) {
	l.logf(DEBUG, "", args...)
}

func (l *logger) Debugf(format string, args ...any) {
	l.logf(DEBUG, format, args...)
}

func (l *logger) Log(args ...any) {
	l.logf(INFO, "", args...)
}

func (l *logger) Logf(format string, args ...any) {
	l.logf(INFO, format, args...)
}

func (l *logger) Info(args ...any) {
	l.logf(INFO, "", args...)
}

func (l *logger) Infof(format string, args ...any) {
	l.logf(INFO, format, args...)
}

func (l *logger) Notice(args ...any) {
	l.logf(NOTICE, "", args...)
}

func (l *logger) Noticef(format string, args ...any) {
	l.logf(NOTICE, format, args...)
}

func (l *logger) Warn(args ...any) {
	l.logf(WARN, "", args...)
}

func (l *logger) Warnf(format string, args ...any) {
	l.logf(WARN, format, args...)
}

func (l *logger) Error(args ...any) {
	l.logf(ERROR, "", args...)
}

func (l *logger) Errorf(format string, args ...any) {
	l.logf(ERROR, format, args...)
}

func (l *logger) Fatal(args ...any) {
	l.logf(FATAL, "", args...)
}

func (l *logger) Fatalf(format string, args ...any) {
	l.logf(FATAL, format, args...)
}

func (l *logger) ChangeLevel(level Level) {
	l.level = level
}

func (l *logger) prettyPrint(e *logEntry, out io.Writer) {
	l.lock.Lock()
	defer l.lock.Unlock()

	fmt.Fprintf(out, "\u001B[38;5;%dm%s\u001B[0m [%s]", e.Level.color(), e.Level.String()[0:4], e.Time.Format(time.TimeOnly))

	if e.TraceID != "" {
		fmt.Fprintf(out, " \u001B[38;5;8m%s\u001B[0m", e.TraceID)
	}

	fmt.Fprint(out, " ")

	// Print the message
	if fn, ok := e.Message.(PrettyPrint); ok {
		fn.PrettyPrint(out)
	} else {
		fmt.Fprintf(out, "%v\n", e.Message)
	}
}

func NewLogger(level Level) Logger {
	return &logger{
		output:     os.Stdout,
		level:      level,
		isTerminal: true,
	}
}

// extractTraceIDAndFilterArgs checks if any of the arguments contain a trace ID
// under the key "__trace_id__" and returns the extracted trace ID along with
// the remaining arguments excluding the trace metadata.
func extractTraceIDAndFilterArgs(args []any) (traceID string, filtered []any) {
	filtered = make([]any, 0, len(args))

	for _, arg := range args {
		if m, ok := arg.(map[string]any); ok {
			if tid, exists := m["__trace_id__"].(string); exists && traceID == "" {
				traceID = tid

				continue
			}
		}

		filtered = append(filtered, arg)
	}

	return traceID, filtered
}
