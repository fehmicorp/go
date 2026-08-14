package logger

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

// Level defines log severity
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
	LevelDebug Level = "DEBUG"
)

type LogEntry struct {
	Level     Level
	Category  string
	Message   string
	Metadata  map[string]interface{}
	CreatedAt time.Time
}

type Logger struct {
	db          *sql.DB
	buffer      chan LogEntry
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	batchSize   int
	flushPeriod time.Duration
}

var (
	globalLogger *Logger
	once         sync.Once
)

func Init(dbDir string, dbFileName string, bufferSize int, batchSize int, flushPeriod time.Duration) error {
	var err error
	once.Do(func() {
		if err = os.MkdirAll(dbDir, 0755); err != nil {
			return
		}

		dbPath := filepath.Join(dbDir, dbFileName)
		// Note: driver name is "sqlite" for modernc.org/sqlite
		db, e := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)")
		if e != nil {
			err = e
			return
		}

		db.SetMaxOpenConns(1)

		createTableQuery := `
		CREATE TABLE IF NOT EXISTS application_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			level TEXT NOT NULL,
			category TEXT NOT NULL,
			message TEXT NOT NULL,
			metadata TEXT,
			created_at DATETIME NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_logs_category ON application_logs(category);
		CREATE INDEX IF NOT EXISTS idx_logs_level ON application_logs(level);
		`
		if _, e := db.Exec(createTableQuery); e != nil {
			err = e
			return
		}

		ctx, cancel := context.WithCancel(context.Background())
		globalLogger = &Logger{
			db:          db,
			buffer:      make(chan LogEntry, bufferSize),
			ctx:         ctx,
			cancel:      cancel,
			batchSize:   batchSize,
			flushPeriod: flushPeriod,
		}

		globalLogger.wg.Add(1)
		go globalLogger.worker()
	})
	return err
}

func (l *Logger) Log(level Level, category, message string, metadata map[string]interface{}) {
	if l == nil {
		log.Printf("[%s] [%s] %s (Logger not initialized)", level, category, message)
		return
	}

	entry := LogEntry{
		Level:     level,
		Category:  category,
		Message:   message,
		Metadata:  metadata,
		CreatedAt: time.Now().UTC(),
	}

	select {
	case l.buffer <- entry:
	default:
		log.Printf("[LOGGER_DROPPED] [%s] [%s] %s", level, category, message)
	}
}

func Info(category, message string, metadata map[string]interface{}) {
	if globalLogger != nil {
		globalLogger.Log(LevelInfo, category, message, metadata)
	}
}

func Warn(category, message string, metadata map[string]interface{}) {
	if globalLogger != nil {
		globalLogger.Log(LevelWarn, category, message, metadata)
	}
}

func Error(category, message string, metadata map[string]interface{}) {
	if globalLogger != nil {
		globalLogger.Log(LevelError, category, message, metadata)
	}
}

func Debug(category, message string, metadata map[string]interface{}) {
	if globalLogger != nil {
		globalLogger.Log(LevelDebug, category, message, metadata)
	}
}

func (l *Logger) worker() {
	defer l.wg.Done()

	ticker := time.NewTicker(l.flushPeriod)
	defer ticker.Stop()

	batch := make([]LogEntry, 0, l.batchSize)

	flush := func() {
		if len(batch) == 0 {
			return
		}
		if err := l.writeBatchToDB(batch); err != nil {
			log.Printf("[LOGGER_ERROR] Failed to write log batch to SQLite: %v", err)
		}
		batch = batch[:0]
	}

	for {
		select {
		case <-l.ctx.Done():
			close(l.buffer)
			for entry := range l.buffer {
				batch = append(batch, entry)
				if len(batch) >= l.batchSize {
					flush()
				}
			}
			flush()
			return

		case entry := <-l.buffer:
			batch = append(batch, entry)
			if len(batch) >= l.batchSize {
				flush()
			}

		case <-ticker.C:
			flush()
		}
	}
}

func (l *Logger) writeBatchToDB(batch []LogEntry) error {
	tx, err := l.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO application_logs (level, category, message, metadata, created_at)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, entry := range batch {
		var metaJSON []byte
		if entry.Metadata != nil {
			metaJSON, _ = json.Marshal(entry.Metadata)
		} else {
			metaJSON = []byte("{}")
		}

		_, err := stmt.Exec(entry.Level, entry.Category, entry.Message, string(metaJSON), entry.CreatedAt)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func Close() {
	if globalLogger != nil {
		globalLogger.cancel()
		globalLogger.wg.Wait()
		globalLogger.db.Close()
	}
}
