// Wiregost - Golang Exploitation Framework
// Copyright © 2020 Para
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package commands

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/evilsocket/islazy/tui"
	"github.com/gogo/protobuf/proto"
	"github.com/olekukonko/tablewriter"

	"github.com/maxlandon/wiregost/client/util"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

func registerSessionCommands() {

	sessions := &Command{
		Name: "sessions",
		SubCommands: []string{
			"interact",
			"kill",
			"kill-all",
		},
		Handle: func(r *Request) error {
			switch length := len(r.Args); {
			case length == 0:
				fmt.Println()
				listSessions(*r.context, r.context.Server.RPC)
			case length >= 1:
				switch r.Args[0] {
				case "interact":
					if len(r.Args) < 2 {
						fmt.Printf("\n" + Error + "Provide a ghost name or session number\n")
					}
					interactGhost(r.Args[1], *r.context, r.context.Server.RPC)
				case "kill":
					fmt.Println()
					killSession(r.Args[1:], *r.context, r.context.Server.RPC)
				case "kill-all":
					fmt.Println()
					killAllSessions(*r.context, r.context.Server.RPC)
				}
			}

			return nil
		},
	}

	AddCommand("main", sessions)
	AddCommand("module", sessions)

	interact := &Command{
		Name: "interact",
		Handle: func(r *Request) error {
			switch length := len(r.Args); {
			case length == 0:
				fmt.Println()
				fmt.Printf(Error + "Provide a session name to interact with\n")
			case length >= 1:
				interactGhost(r.Args[0], *r.context, r.context.Server.RPC)
			}
			return nil
		},
	}

	AddCommand("main", interact)
	AddCommand("module", interact)

	background := &Command{
		Name: "background",
		Handle: func(r *Request) error {
			fmt.Println()
			*r.context.CurrentAgent = clientpb.Ghost{}
			fmt.Printf(Info + " Background ...\n")
			return nil
		},
	}

	AddCommand("agent", background)
}

func listSessions(ctx ShellContext, rpc RPCServer) {

	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgSessions,
		Data: []byte{},
	}, defaultTimeout)
	if resp.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return
	}
	sessions := &clientpb.Sessions{}
	proto.Unmarshal(resp.Data, sessions)

	ghosts := map[uint32]*clientpb.Ghost{}
	for _, ghost := range sessions.Ghosts {
		ghosts[ghost.ID] = ghost
	}
	if 0 < len(ghosts) {
		printGhosts(ghosts, rpc)
	} else {
		fmt.Printf(Info + "No ghosts connected\n")
	}
}

func killSession(ghosts []string, ctx ShellContext, rpc RPCServer) {

	if len(ghosts) == 0 {
		fmt.Printf(Warn + "Provide a session name or ID\n")
		return
	}

	for _, g := range ghosts {
		ghost := getGhost(g, rpc)
		if ghost != nil {
			data, _ := proto.Marshal(&ghostpb.KillReq{
				GhostID: ghost.ID,
				Force:   true,
			})
			rpc(&ghostpb.Envelope{
				Type: ghostpb.MsgKill,
				Data: data,
			}, 5*time.Second)

			fmt.Printf(Info+"Killed agent %s (Session %d)\n", ghost.Name, ghost.ID)
		} else {
			fmt.Printf(Error+"Invalid ghost name: %s", g)
		}

	}
}

func killAllSessions(ctx ShellContext, rpc RPCServer) {

	sessions := GetGhosts(rpc)
	for _, session := range sessions.Ghosts {
		data, _ := proto.Marshal(&ghostpb.KillReq{
			GhostID: session.ID,
			Force:   true,
		})
		rpc(&ghostpb.Envelope{
			Type: ghostpb.MsgKill,
			Data: data,
		}, 5*time.Second)

		fmt.Printf(Info+"Killed %s (%d)\n", ctx.CurrentAgent.Name, ctx.CurrentAgent.ID)
	}
}

func printGhosts(sessions map[uint32]*clientpb.Ghost, rpc RPCServer) {
	table := util.Table()
	table.SetHeader([]string{"WsID", "ID", "Name", "Proto", "Remote Address", "user@host", "Platform", "Status"})
	table.SetColWidth(40)
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
	)

	// Sort the keys because maps have a randomized order
	var keys []int
	for _, sliver := range sessions {
		keys = append(keys, int(sliver.ID))
	}
	sort.Ints(keys) // Fucking Go can't sort int32's, so we convert to/from int's

	for _, key := range keys {
		ghost := sessions[uint32(key)]
		workspace := ""
		if ghost.WorkspaceID != 0 {
			workspace = strconv.Itoa(int(ghost.WorkspaceID))
		}
		os := fmt.Sprintf("%s/%s", ghost.OS, ghost.Arch)
		userHost := fmt.Sprintf("%s@%s", ghost.Username, ghost.Hostname)

		status := getSessionStatus(ghost, rpc)
		table.Append([]string{workspace, strconv.Itoa(int(ghost.ID)), ghost.Name, ghost.Transport,
			ghost.RemoteAddress, userHost, os, status})
	}

	table.Render()
}

func interactGhost(name string, ctx ShellContext, rpc RPCServer) {

	ghost := &clientpb.Ghost{}

	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgSessions,
		Data: []byte{},
	}, defaultTimeout)
	if resp.Err != "" {
		fmt.Printf(Error+"Impossible to establish communication with session %s\n", name)
		return
	}

	sessions := &clientpb.Sessions{}
	proto.Unmarshal((resp).Data, sessions)

	for _, g := range sessions.Ghosts {
		if strconv.Itoa(int(g.ID)) == name || g.Name == name {
			ghost = g
		}
	}

	if ghost != nil {
		// Get cwd, and check that session is connected by the same way
		data, _ := proto.Marshal(&ghostpb.PwdReq{
			GhostID: ghost.ID,
		})
		resp := <-rpc(&ghostpb.Envelope{
			Type: ghostpb.MsgPwdReq,
			Data: data,
		}, defaultTimeout)
		if resp.Err != "" {
			fmt.Printf("\n"+Error+"Impossible to establish communication with session %s\n", name)
			return
		}

		pwd := &ghostpb.Pwd{}
		err := proto.Unmarshal(resp.Data, pwd)
		if err != nil {
			fmt.Printf(Warn+"Unmarshaling envelope error: %v\n", err)
			return
		}
		*ctx.CurrentAgent = *ghost
		*ctx.AgentPwd = pwd.Path

	} else {
		fmt.Printf(Error+"Invalid ghost name or session number: %s", name)
	}
}

// Get Ghost by session ID or name
func getGhost(arg string, rpc RPCServer) *clientpb.Ghost {
	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgSessions,
		Data: []byte{},
	}, defaultTimeout)
	sessions := &clientpb.Sessions{}
	proto.Unmarshal((resp).Data, sessions)

	for _, ghost := range sessions.Ghosts {
		if strconv.Itoa(int(ghost.ID)) == arg || ghost.Name == arg {
			return ghost
		}
	}
	return nil
}

// GetGhosts - Get all connected sessions
func GetGhosts(rpc RPCServer) *clientpb.Sessions {
	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgSessions,
		Data: []byte{},
	}, defaultTimeout)
	sessions := &clientpb.Sessions{}
	proto.Unmarshal((resp).Data, sessions)

	return sessions
}

// GhostSessionsByName - Get a session by name
func GhostSessionsByName(name string, rpc RPCServer) []*clientpb.Ghost {
	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgSessions,
		Data: []byte{},
	}, defaultTimeout)
	allSessions := &clientpb.Sessions{}
	proto.Unmarshal((resp).Data, allSessions)

	sessions := []*clientpb.Ghost{}
	for _, ghost := range allSessions.Ghosts {
		if ghost.Name == name {
			sessions = append(sessions, ghost)
		}
	}
	return sessions
}

func getSessionStatus(ghost *clientpb.Ghost, rpc RPCServer) string {

	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgListGhostBuilds,
	}, defaultTimeout)
	if resp.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return ""
	}

	config := &clientpb.GhostConfig{}

	builds := &clientpb.GhostBuilds{}
	proto.Unmarshal(resp.Data, builds)
	for _, b := range builds.Configs {
		if ghost.Name == b.Name {
			config = b
		}
	}

	dur, errDur := time.ParseDuration(fmt.Sprintf("%ds", config.ReconnectInterval))
	if errDur != nil {
		fmt.Println(errDur)
	}
	fmt.Println(dur)

	lastCheckin, err := time.Parse(time.RFC1123, ghost.LastCheckin)
	fmt.Println(lastCheckin)
	if err != nil {
		fmt.Println(err)
	}

	if lastCheckin.Add(dur).After(time.Now()) {
		return tui.Green("Alive")
	} else if lastCheckin.Add(dur * time.Duration(config.MaxConnectionErrors+1)).After(time.Now()) {
		return tui.Yellow("Delayed")
	} else {
		return tui.Red("Dead")
	}
}
