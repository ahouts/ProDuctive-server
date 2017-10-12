package tunnel

import "fmt"

type Endpoint struct {
	// Server host address
	Host string
	// Server port
	Port int
}

func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}
