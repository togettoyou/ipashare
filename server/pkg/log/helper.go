package log

import (
	"strings"

	"go.uber.org/zap"
)

type module struct {
	name string
}

func New(name string) *module {
	return &module{
		name: "[Module]" + name,
	}
}

func (m module) Named(s string) module {
	if s == "" {
		return m
	}
	if m.name == "" {
		cm := New(s)
		m.name = cm.name
	} else {
		m.name = strings.Join([]string{m.name, s}, ".")
	}
	return m
}

func (m module) L() *zap.Logger {
	return zap.L().Named(m.name)
}

func (m module) S() *zap.SugaredLogger {
	return zap.S().Named(m.name)
}
