package fresh

import pb "github.com/freshcloud-io/protos/go/freshcloud"

type FreshApplication interface {
	// EntryPoint defines the "main" function
	EntryPoint() error
	// Shutdown defines the place to tear down things where necessary
	Shutdown() error
	// GetApplication returns the actual application containing all the information
	GetApplication() (*Application, error)
}

type FreshCron interface {
	GetCron() (*Cron, error)
	Execute(args *pb.CronExecuteRequest) *pb.CronExecuteResponse
}
