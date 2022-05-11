package migrate

import "github.com/rs/zerolog"

type Option []OptionFn
type OptionFn func(m *Migrator)

// WithClean clean database
func WithClean(db ...string) OptionFn {
	return func(m *Migrator) {
		m.cleanDb = db
	}
}

// WithLogger implement logger
func WithLogger(logger zerolog.Logger) OptionFn {
	return func(m *Migrator) {
		m.logger = logger
	}
}
