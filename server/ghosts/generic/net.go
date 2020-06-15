package generic

import corepb "github.com/maxlandon/wiregost/proto/v1/gen/go/ghost/core"

// IfConfig - Get target network interfaces
func (g *Ghost) IfConfig() (net *corepb.IfConfig) {
	return
}

// Netstat - Get active connections from/to target
func (g *Ghost) Netstat() (net *corepb.Netstat) {
	return
}