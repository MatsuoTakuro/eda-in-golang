package rpc

import (
	"fmt"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type RpcConfig struct {
	Host     string   `default:"0.0.0.0"`
	Port     string   `default:":9000"`
	Services Services // for other service addresses
}

// Address returns your RPC server address in the format "host:port".
func (c RpcConfig) Address() string {
	return fmt.Sprintf("%s%s", c.Host, c.Port)
}

// Service returns the address of a specific service that your service depends on.
func (c RpcConfig) Service(service string) string {
	if address, exists := c.Services[service]; exists {
		return address
	}
	return c.Address()
}

type Services map[string]string // key: service name (e.g., "STORES"), value: address (e.g., "stores:9000")

var _ envconfig.Decoder = (*Services)(nil)

func (s *Services) Decode(v string) error {
	services := map[string]string{}
	pairs := strings.SplitSeq(v, ",")
	for p := range pairs {
		pair := strings.TrimSpace(p)
		if len(pair) == 0 {
			continue
		}
		kv := strings.Split(pair, "=")
		if len(kv) != 2 {
			return fmt.Errorf("invalid service pair: %q", pair)
		}
		services[strings.ToUpper(kv[0])] = kv[1]
	}

	*s = services
	return nil
}
