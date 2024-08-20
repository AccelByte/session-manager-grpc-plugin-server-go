package session

import (
	"context"
	"fmt"

	sessionmanager "accelbyte.net/session-manager-grpc-plugin-server-go/pkg/pb"
)

type SessionManager struct {
	sessionmanager.UnimplementedSessionManagerServer
}

func (s *SessionManager) OnSessionCreated(ctx context.Context, request *sessionmanager.SessionRequest) (*sessionmanager.SessionResponse, error) {
	fmt.Println("test Session created")
	fmt.Println(request.GetSession())
	return &sessionmanager.SessionResponse{
		Session: request.GetSession(),
	}, nil
}

func (s *SessionManager) OnSessionUpdated(ctx context.Context, request *sessionmanager.SessionUpdatedRequest) (*sessionmanager.SessionResponse, error) {
	fmt.Println("test Session updated")
	return &sessionmanager.SessionResponse{
		Session: request.GetSessionNew(),
	}, nil
}

func (s *SessionManager) OnSessionDeleted(ctx context.Context, request *sessionmanager.SessionRequest) (*sessionmanager.SessionResponse, error) {
	fmt.Println("test Session deleted")
	return &sessionmanager.SessionResponse{
		Session: request.GetSession(),
	}, nil
}
