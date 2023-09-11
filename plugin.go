package fresh

import pb "github.com/freshcloud-io/protos/go/freshcloud"

type FreshApplication interface {
	// EntryPoint defines the "main" function
	EntryPoint() error
	// Shutdown defines the place to tear down things where necessary
	Shutdown() error
}

type FreshCron interface {
	Create() (*Cron, error)
	Execute(args *pb.CronExecuteRequest) *pb.CronExecuteResponse
}
