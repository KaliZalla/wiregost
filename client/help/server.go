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

package help

import (
	"fmt"

	"github.com/evilsocket/islazy/tui"
)

var (
	serverHelp = fmt.Sprintf(`%s%s Server Commands%s 

%s About:%s Server management

%s Commands:%s
    server          %sList saved servers in the ~/.wiregost-client/configs/ directory%s
    server connect  %sConnect to one of the available Wiregost servers (Example below) (Completed)%s

%s Examples:%s
    server connect localhost:8443 as wiregost   %s(The 'localhost:8443 as wiregost' is automatically completed)%s`,
		tui.BLUE, tui.BOLD, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.DIM, tui.RESET,
		tui.DIM, tui.RESET,
		tui.YELLOW, tui.RESET,
		tui.DIM, tui.RESET,
	)
)
