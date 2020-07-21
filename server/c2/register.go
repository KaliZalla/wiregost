package c2

// HandleGhostRegistration - Of all the process starting from TCP handshake to
// complete registration and usage of the ghost implant by users, this function
// is the first one that all implants, independently from their transports, target OS,
// have in common.
// Generally, security details linked to the transport mechanism are already dealt with.
func HandleGhostRegistration() {

	// Register RPC services/handlers if the ghost reverse-calls us (we are the server)

	// If bind, then either we wait for registration message to come in, or we request it.

	// This should include the logging infrastructure

	// Populate new ghostpb object with all registration info, and register user/module interfaces
	// This means, at this point, that although all OS-specific commands are technically available,
	// much of the implant state/information is not disseminated in the ghost object that will be
	// further used by consoles/modules.
	// registrar := &ghostpb.Ghost{}
	// registered := ghosts.NewGhost(registrar)

	// Send all necessary information to Database

	// Register/check ghost owner & permissions

	// Register/instantiate/populate OS specific objects in the registered ghost.

	// Send registration notification to user consoles
}
