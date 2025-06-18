package shared

import (
	"github.com/rs/zerolog"
)

type Field struct {
	Key   string
	Value any
}

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	With(fields ...Field) Logger
}

type zerologLogger struct {
	logger zerolog.Logger
}

func NewZerologLogger(l zerolog.Logger, serviceName string, level zerolog.Level) Logger {
	return &zerologLogger{
		logger: l.With().Str("service", serviceName).Logger().Level(level),
	}
}

func (z *zerologLogger) Debug(msg string, fields ...Field) {
	z.logger.Debug().Fields(toMap(fields)).Msg(msg)
}

func (z *zerologLogger) Info(msg string, fields ...Field) {
	z.logger.Info().Fields(toMap(fields)).Msg(msg)
}

func (z *zerologLogger) Warn(msg string, fields ...Field) {
	z.logger.Warn().Fields(toMap(fields)).Msg(msg)
}

func (z *zerologLogger) Error(msg string, fields ...Field) {
	z.logger.Error().Fields(toMap(fields)).Msg(msg)
}

func (z *zerologLogger) With(fields ...Field) Logger {
	return &zerologLogger{
		logger: z.logger.With().Fields(toMap(fields)).Logger(),
	}
}

func toMap(fields []Field) map[string]any {
	m := make(map[string]any, len(fields))
	for _, f := range fields {
		m[f.Key] = f.Value
	}

	return m
}
