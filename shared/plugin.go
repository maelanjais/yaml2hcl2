package shared

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// ConverterRPCClient est le client RPC pour communiquer avec le plugin
type ConverterRPCClient struct {
	client *rpc.Client
}

func (g *ConverterRPCClient) Convert(input []byte) ([]byte, error) {
	var resp []byte
	err := g.client.Call("Plugin.Convert", input, &resp)
	return resp, err
}

// ConverterRPCServer est le serveur RPC qui encapsule l'implémentation du plugin
type ConverterRPCServer struct {
	Impl Converter
}

func (s *ConverterRPCServer) Convert(args []byte, resp *[]byte) error {
	res, err := s.Impl.Convert(args)
	*resp = res
	return err
}

// ConverterPlugin est l'implémentation du plugin HashiCorp
type ConverterPlugin struct {
	Impl Converter
}

func (p *ConverterPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &ConverterRPCServer{Impl: p.Impl}, nil
}

func (ConverterPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &ConverterRPCClient{client: c}, nil
}

// HandshakeConfig est utilisé pour la validation entre le client et le serveur
var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "YAML2HCL2_PLUGIN",
	MagicCookieValue: "hello",
}

// PluginMap est la map des plugins utilisables
var PluginMap = map[string]plugin.Plugin{
	"converter": &ConverterPlugin{},
}
