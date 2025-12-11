package mocks

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type TestEntryRecord struct {
	Level   logrus.Level
	Message string
	Fields  logrus.Fields
	Time    time.Time
}

// TestHook implements logrus.Hook and stores fired entries in memory
type TestHook struct {
	mu      sync.Mutex
	Entries []TestEntryRecord
	levels  []logrus.Level
}

// NewTestHook returns a TestHook that listens to all levels
func NewTestHook() *TestHook {
	return &TestHook{
		Entries: make([]TestEntryRecord, 0),
		levels:  logrus.AllLevels,
	}
}

// Levels return the log levels this hook handles
func (h *TestHook) Levels() []logrus.Level {
	return h.levels
}

// Fire is called by logrus when an entry is logged
func (h *TestHook) Fire(entry *logrus.Entry) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// copy fields to avoid later mutation issues
	fieldsCopy := logrus.Fields{}
	for k, v := range entry.Data {
		fieldsCopy[k] = v
	}

	h.Entries = append(h.Entries, TestEntryRecord{
		Level:   entry.Level,
		Message: entry.Message,
		Fields:  fieldsCopy,
		Time:    entry.Time,
	})
	return nil
}

// Last returns the last recorded entry (or nil if none)
func (h *TestHook) Last() *TestEntryRecord {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(h.Entries) == 0 {
		return nil
	}
	e := h.Entries[len(h.Entries)-1]
	return &e
}

// All returns a copy of all recorded entries
func (h *TestHook) All() []TestEntryRecord {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make([]TestEntryRecord, len(h.Entries))
	copy(out, h.Entries)
	return out
}

// Reset clears recorded entries
func (h *TestHook) Reset() {
	h.mu.Lock()
	h.Entries = h.Entries[:0]
	h.mu.Unlock()
}

// NewTestLoggerEntry creates a *logrus.Entry backed by a logger configured with TestHook.
// Return values:
//   - *logrus.Entry: can be passed into your services (implements FieldLogger)
//   - *TestHook: inspect Entries in tests
func NewTestLoggerEntry() (*logrus.Entry, *TestHook) {
	base := logrus.New()
	// optional: keep output silent in tests
	base.SetOutput(nil)
	base.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true, DisableColors: true})
	base.SetLevel(logrus.TraceLevel)

	hook := NewTestHook()
	base.AddHook(hook)

	// create an entry with no extra fields
	entry := base.WithFields(logrus.Fields{})
	return entry, hook
}
