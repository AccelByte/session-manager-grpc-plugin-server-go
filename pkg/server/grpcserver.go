// Copyright (c) 2024 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package server

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
	log.Println("session", request.GetSession())
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
	log.Println("old Session:", request.GetSessionOld())
	log.Println("new Session:", request.GetSessionNew())
	return &emptypb.Empty{}, nil
}

func (s *SessionManager) OnSessionDeleted(ctx context.Context, request *sessionmanager.SessionDeletedRequest) (*emptypb.Empty, error) {
	log.Println("got message from OnSessionDeleted")
	log.Println("session deleted", request.GetSession())
	return &emptypb.Empty{}, nil
}

func (s *SessionManager) OnPartyCreated(ctx context.Context, request *sessionmanager.PartyCreatedRequest) (*sessionmanager.PartyResponse, error) {
	log.Println("got message from OnSessionCreated")
	log.Println("session", request.GetSession())
	session := request.GetSession()
	if session.Session.Attributes == nil {
		session.Session.Attributes = &structpb.Struct{}
	}
	if session.Session.Attributes.Fields == nil {
		session.Session.Attributes.Fields = map[string]*structpb.Value{}
	}
	session.Session.Attributes.Fields["PARTY_SAMPLE"] = structpb.NewStringValue("party value from GRPC server")
	return &sessionmanager.PartyResponse{
		Session: session,
	}, nil
}

func (s *SessionManager) OnPartyUpdated(ctx context.Context, request *sessionmanager.PartyUpdatedRequest) (*emptypb.Empty, error) {
	log.Println("got message from OnPartyUpdated")
	log.Println("old session", request.GetSessionOld())
	log.Println("new session", request.GetSessionNew())
	return &emptypb.Empty{}, nil
}

func (s *SessionManager) OnPartyDeleted(ctx context.Context, request *sessionmanager.PartyDeletedRequest) (*emptypb.Empty, error) {
	log.Println("got message from OnPartyDeleted")
	log.Println("session deleted", request.GetSession())
	return &emptypb.Empty{}, nil
}
