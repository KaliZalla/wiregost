package context

import (
	"github.com/google/uuid"
	"github.com/lmorg/readline"
	"google.golang.org/grpc"

	clientpb "github.com/maxlandon/wiregost/proto/v1/gen/go/client"
	dbpb "github.com/maxlandon/wiregost/proto/v1/gen/go/db"
	ghostpb "github.com/maxlandon/wiregost/proto/v1/gen/go/ghost"
	modulepb "github.com/maxlandon/wiregost/proto/v1/gen/go/module"
)

var (
	// Context - The console context object
	Context = newContext()
)

// Menu Contexts
const (
	// MAIN_CONTEXT - Available only in main menu
	MAIN_CONTEXT = "main"
	// MODULE_CONTEXT - Available only when a module is loaded
	MODULE_CONTEXT = "module"
	// GHOST_CONTEXT - Available only when interacting with a ghost implant
	GHOST_CONTEXT = "ghost"
)

// ConsoleContext - Stores all variables needed for console context
type ConsoleContext struct {
	ClientID  uuid.UUID              // Unique user ID for module requests
	User      dbpb.User              // User information sent back after auth
	Shell     *readline.Instance     // Shell object
	Config    clientpb.ConsoleConfig // Shell configuration
	Menu      string                 // Current shell menu
	Workspace dbpb.Workspace         // Current workspace
	Module    modulepb.Module        // Current module
	Ghost     ghostpb.Ghost          // Current implant
	Jobs      int                    // Number of jobs
	Ghosts    int                    // Number of connected implants
}

func newContext() (ctx *ConsoleContext) {

	ctx = &ConsoleContext{}
	ctx.User = dbpb.User{}
	ctx.Config = clientpb.ConsoleConfig{}
	ctx.Workspace = dbpb.Workspace{}
	ctx.Module = modulepb.Module{}
	ctx.Ghost = ghostpb.Ghost{}

	return
}

// SetConsoleContext - Set the context used by commands & shell
func SetConsoleContext(cli clientpb.ConnectionRPCClient) (info *clientpb.ConnectionInfo) {

	// Info Request
	info, _ = cli.GetConnectionInfo(base, &clientpb.ConnectionInfoRequest{}, grpc.EmptyCallOption{})

	// Set fields (beware of nil fields in pb message)
	// Context.Workspace = (*info.Workspace)
	// *Context.Jobs = int(info.Jobs)
	// *Context.Ghosts = int(info.Ghosts)
	Context.Menu = MAIN_CONTEXT
	return
}

// GetVersion - Get client & server version information upon connection
func GetVersion(cli clientpb.ConnectionRPCClient) (info *clientpb.Version) {
	info, _ = cli.GetVersion(base, &clientpb.Empty{}, grpc.EmptyCallOption{})
	return
}