// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"github.com/google/wire"
	"go-chat/config"
	"go-chat/internal/cache"
	"go-chat/internal/dao"
	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/provider"
	"go-chat/internal/service"
	"go-chat/internal/websocket/internal/handler"
	"go-chat/internal/websocket/internal/process"
	"go-chat/internal/websocket/internal/process/handle"
	"go-chat/internal/websocket/internal/router"
)

// Injectors from wire.go:

func Initialize(ctx context.Context, conf *config.Config) *Providers {
	client := provider.NewRedisClient(ctx, conf)
	sidServer := cache.NewSid(client)
	wsClientSession := cache.NewWsClientSession(client, conf, sidServer)
	clientService := service.NewClientService(wsClientSession)
	room := cache.NewRoom(client)
	db := provider.NewMySQLClient(conf)
	baseService := service.NewBaseService(db, client)
	baseDao := dao.NewBaseDao(db, client)
	relation := cache.NewRelation(client)
	groupMemberDao := dao.NewGroupMemberDao(baseDao, relation)
	groupMemberService := service.NewGroupMemberService(baseService, groupMemberDao)
	defaultWebSocket := handler.NewDefaultWebSocket(client, conf, clientService, room, groupMemberService)
	exampleWebsocket := handler.NewExampleWebsocket()
	handlerHandler := &handler.Handler{
		DefaultWebSocket: defaultWebSocket,
		ExampleWebsocket: exampleWebsocket,
	}
	session := cache.NewSession(client)
	engine := router.NewRouter(conf, handlerHandler, session)
	websocketServer := provider.NewWebsocketServer(conf, engine)
	health := process.NewHealth(conf, sidServer)
	talkVote := cache.NewTalkVote(client)
	talkRecordsVoteDao := dao.NewTalkRecordsVoteDao(baseDao, talkVote)
	filesystemFilesystem := filesystem.NewFilesystem(conf)
	talkRecordsDao := dao.NewTalkRecordsDao(baseDao)
	talkRecordsService := service.NewTalkRecordsService(baseService, talkVote, talkRecordsVoteDao, filesystemFilesystem, groupMemberDao, talkRecordsDao)
	usersFriendsDao := dao.NewUsersFriendsDao(baseDao, relation)
	contactService := service.NewContactService(baseService, usersFriendsDao)
	subscribeConsume := handle.NewSubscribeConsume(conf, wsClientSession, room, talkRecordsService, contactService)
	wsSubscribe := process.NewWsSubscribe(client, conf, subscribeConsume)
	coroutine := process.NewCoroutine(health, wsSubscribe)
	providers := &Providers{
		Config:    conf,
		WsServer:  websocketServer,
		Coroutine: coroutine,
	}
	return providers
}

// wire.go:

var providerSet = wire.NewSet(provider.NewMySQLClient, provider.NewRedisClient, provider.NewWebsocketServer, router.NewRouter, process.NewCoroutine, process.NewHealth, process.NewWsSubscribe, handle.NewSubscribeConsume, cache.NewSession, cache.NewSid, cache.NewRedisLock, cache.NewWsClientSession, cache.NewRoom, cache.NewTalkVote, cache.NewRelation, dao.NewBaseDao, dao.NewTalkRecordsDao, dao.NewTalkRecordsVoteDao, dao.NewGroupMemberDao, dao.NewUsersFriendsDao, filesystem.NewFilesystem, service.NewBaseService, service.NewTalkRecordsService, service.NewClientService, service.NewGroupMemberService, service.NewContactService, handler.NewDefaultWebSocket, handler.NewExampleWebsocket, wire.Struct(new(handler.Handler), "*"), wire.Struct(new(Providers), "*"))
