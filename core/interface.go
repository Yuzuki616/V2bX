package core

import (
	"github.com/Yuzuki616/V2bX/api/panel"
	"github.com/Yuzuki616/V2bX/conf"
	"github.com/Yuzuki616/V2bX/limiter"
)

type AddUsersParams struct {
	Tag      string
	Config   *conf.ControllerConfig
	Limiter  *limiter.Limiter
	UserInfo []panel.UserInfo
	NodeInfo *panel.NodeInfo
}
type Core interface {
	Start() error
	Close() error
	AddNode(tag string, info *panel.NodeInfo, config *conf.ControllerConfig) error
	DelNode(tag string) error
	AddUsers(p *AddUsersParams) (added int, err error)
	AddUserSpeedLimit(Limiter *limiter.Limiter, tag string, user *panel.UserInfo, speedLimit int, expire int64) error
	GetUserTraffic(tag, uuid string, reset bool) (up int64, down int64)
	DelUsers(users []panel.UserInfo, tag string) error
	Protocols() []string
}
