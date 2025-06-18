package loginserver

import (
	"os"

	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/network"
	"github.com/rs/zerolog"
)

type LoginServer struct {
	network.TCPServer
	loggedInAccounts *shared.SafeSet[string]
}

func NewLoginServer(addr string) *LoginServer {
	return &LoginServer{
		TCPServer: network.TCPServer{
			Addr:         addr,
			Name:         "login-server",
			NewSession:   newLoginServerSession,
			UidGenerator: shared.NewUidGenerator(0),
			Logger:       shared.NewZerologLogger(zerolog.New(os.Stdout), "login-server", zerolog.InfoLevel),
		},
		loggedInAccounts: shared.NewSafeSet[string](),
	}
}

func (s *LoginServer) AddLoggedInAccount(account string) {
	s.loggedInAccounts.Add(account)
}

func (s *LoginServer) RemoveLoggedInAccount(account string) {
	s.loggedInAccounts.Remove(account)
}

func (s *LoginServer) IsLoggedIn(account string) bool {
	return s.loggedInAccounts.Contains(account)
}
