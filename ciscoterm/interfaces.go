package ciscoterm

type Terminal interface {
	Connect(ciscoDev CiscoDevice) error
	Close() error
	EnableTerm(enablePasswd string) error
	DisablePagination() error
	ExecuteCommand(cmd string) ([]string, error)
}

func NewTerminal() Terminal {
	return &terminal{}
}
