package constants

const LoggedInUserKeyPrefix = "agonyl:logged_in_user:"

const LoginFailedMsg = "Login failed."

const AccountAlreadyLoggedInMsg = "Account is already logged in."

const DuplicateCharacterMsg = "Character already exists."

const MaxCharactersPerAccountExceededMsg = "Max characters per account exceeded."

const MaxCharactersPerAccount = 5

const CharacterNotFoundMsg = "Character not found."

const (
	AccountStatusActive              = "active"
	AccountStatusInactive            = "inactive"
	AccountStatusBanned              = "banned"
	AccountStatusSuspended           = "suspended"
	AccountStatusPendingVerification = "pending_verification"
	AccountStatusDeleted             = "deleted"
)

const (
	CharacterStatusActive  = "active"
	CharacterStatusLocked  = "locked"
	CharacterStatusDeleted = "deleted"
)
