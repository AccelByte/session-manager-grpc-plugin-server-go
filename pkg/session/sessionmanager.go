package session

import (
	"context"
	"log"

	sessionmanager "accelbyte.net/session-manager-grpc-plugin-server-go/pkg/pb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
)

type SessionManager struct {
	sessionmanager.UnimplementedSessionManagerServer
}

func (s *SessionManager) OnSessionCreated(ctx context.Context, request *sessionmanager.SessionCreatedRequest) (*sessionmanager.SessionResponse, error) {
	log.Println("got message from OnSessionCreated")
	log.Println(request.GetSession())
	session := request.GetSession()
	if session.Session.Attributes == nil {
		session.Session.Attributes = &structpb.Struct{}
	}
	if session.Session.Attributes.Fields == nil {
		session.Session.Attributes.Fields = map[string]*structpb.Value{}
	}
	session.Session.Attributes.Fields["SAMPLE"] = structpb.NewStringValue("value from GRPC server")
	return &sessionmanager.SessionResponse{
		Session: session,
	}, nil
}

func (s *SessionManager) OnSessionUpdated(ctx context.Context, request *sessionmanager.SessionUpdatedRequest) (*emptypb.Empty, error) {
	log.Println("got message from OnSessionUpdated")
	log.Println(request.GetSessionOld())
	log.Println(request.GetSessionNew())
	return &emptypb.Empty{}, nil
}

func (s *SessionManager) OnSessionDeleted(ctx context.Context, request *sessionmanager.SessionDeletedRequest) (*emptypb.Empty, error) {
	log.Println("got message from OnSessionDeleted")
	log.Println(request.GetSession())
	return &emptypb.Empty{}, nil
}
