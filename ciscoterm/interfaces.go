package ciscoterm

import "context"

type Terminal interface {
	Connect(ciscoDev CiscoDevice) error
	Close() error
	EnableTerm(ctx context.Context, enablePasswd string) error
	DisablePagination(ctx context.Context) error
	ExecuteCommand(ctx context.Context, cmd string) ([]string, error)
}

func NewTerminal() Terminal {
	return &terminal{}
}
