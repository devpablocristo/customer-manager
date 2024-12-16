package pkggomicro

import (
	"os"
	"sync"

	"github.com/go-micro/plugins/v4/registry/consul"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/registry"

	defs "github.com/devpablocristo/customer-manager/pkg/microservices/go-micro/v4/micro-service/defs"
)

var (
	instance  defs.Service
	once      sync.Once
	initError error
)

type service struct {
	s micro.Service
}

func newService(config defs.Config) (defs.Service, error) {
	once.Do(func() {
		setupLogger()

		instance = &service{
			s: setupService(config),
		}
	})

	if initError != nil {
		return nil, initError
	}

	return instance, nil
}

func setupService(config defs.Config) micro.Service {
	service := micro.NewService(
		micro.Server(config.GetServer()),
		micro.Client(config.GetClient()),
		micro.Registry(setupRegistry(config)),
	)

	service.Init()

	return service
}

func (s *service) Run() error {
	return s.s.Run()
}

func (s *service) GetService() micro.Service {
	return s.s
}

func setupRegistry(config defs.Config) registry.Registry {
	consulReg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{config.GetConsulAddress()}
	})
	return consulReg
}

func setupLogger() {
	logger.DefaultLogger = logger.NewLogger(
		logger.WithLevel(logger.InfoLevel),
		logger.WithOutput(os.Stdout),
	)
}