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
	"strconv"
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/maxlandon/wiregost/client/util"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
	"github.com/olekukonko/tablewriter"
)

func registerProfileCommands() {

	profiles := &Command{
		Name: "profiles",
		// SubCommands: []string{
		//         "delete",
		// },
		Handle: func(r *Request) error {
			switch length := len(r.Args); {
			case length == 0:
				fmt.Println()
				listProfiles(r.context.Server.RPC)
			case length > 1:
				switch r.Args[0] {
				case "delete":
					deleteProfile(r.Args[1], r.context.Server.RPC)
				}

			}
			return nil
		},
	}

	AddCommand("main", profiles)
	AddCommand("module", profiles)
}

func deleteProfile(name string, rpc RPCServer) {
	pReq, _ := proto.Marshal(&clientpb.Profile{
		Name: name,
	})

	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgDeleteProfile,
		Data: pReq,
	}, defaultTimeout)

	if resp.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return
	}

	pRes := &clientpb.Profile{}
	err := proto.Unmarshal(resp.Data, pRes)
	if err != nil {
		fmt.Printf("\n"+Error+"%s\n", err.Error())
		return
	}

	if pRes.Name == name {
		fmt.Printf("\n"+Success+"Deleted profile %s\n", name)
	} else {
		fmt.Printf("\n"+Error+"%s\n", pRes.Name)
	}
}

func listProfiles(rpc RPCServer) {
	resp := <-rpc(&ghostpb.Envelope{
		Type: clientpb.MsgProfiles,
	}, defaultTimeout)

	if resp.Err != "" {
		fmt.Printf(RPCError+"%s\n", resp.Err)
		return
	}

	pbProfiles := &clientpb.Profiles{}
	err := proto.Unmarshal(resp.Data, pbProfiles)
	if err != nil {
		fmt.Println()
		fmt.Printf(Error+"%s", err.Error())
		fmt.Println()
		return
	}

	profiles := &map[string]*clientpb.Profile{}
	for _, profile := range pbProfiles.List {
		(*profiles)[profile.Name] = profile
	}

	table := util.Table()
	table.SetHeader([]string{"Name", "Platform", "Format", "Command & Control", "Limitations", "Debug"})
	table.SetColWidth(40)
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
	)

	for k, p := range *profiles {
		platform := fmt.Sprintf("%s/%s", p.Config.GOOS, p.Config.GOARCH)
		c2s := []string{}
		for _, c := range p.Config.C2 {
			c2s = append(c2s, c.URL)
		}
		limits := getLimitsString(p.Config)
		table.Append([]string{k, platform, p.Config.Format.String(), strings.Join(c2s, ","), limits, strconv.FormatBool(p.Config.Debug)})
	}

	table.Render()
}

func getLimitsString(config *clientpb.GhostConfig) string {
	limits := []string{}
	if config.LimitDatetime != "" {
		limits = append(limits, fmt.Sprintf("datetime=%s", config.LimitDatetime))
	}
	if config.LimitDomainJoined {
		limits = append(limits, fmt.Sprintf("domainjoined=%v", config.LimitDomainJoined))
	}
	if config.LimitUsername != "" {
		limits = append(limits, fmt.Sprintf("username=%s", config.LimitUsername))
	}
	if config.LimitHostname != "" {
		limits = append(limits, fmt.Sprintf("hostname=%s", config.LimitHostname))
	}
	return strings.Join(limits, "; ")
}
