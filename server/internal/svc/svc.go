package svc

import (
	"ipashare/internal/model"
	logpkg "ipashare/pkg/log"

	"go.uber.org/zap"
)

// Service 每个业务对象都需要内嵌该结构体
type Service struct {
	store *model.Store
	log   *zap.Logger
}

func (s *Service) New(store *model.Store, log *zap.Logger) *Service {
	if log == nil {
		log = logpkg.New("").L()
	}
	s.store = store
	s.log = log.Named("svc")
	return s
}

func (s *Service) named(name string) *Service {
	if s.log == nil {
		s.log = logpkg.New("").L()
	}
	s.log = s.log.Named(name)
	return s
}
