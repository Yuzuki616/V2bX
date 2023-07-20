package hy

import (
	"encoding/base64"
	"errors"
	"github.com/Yuzuki616/V2bX/api/panel"
	"github.com/Yuzuki616/V2bX/core"
	"github.com/Yuzuki616/V2bX/limiter"
)

func (h *Hy) AddUsers(p *core.AddUsersParams) (int, error) {
	s, ok := h.servers.Load(p.Tag)
	if !ok {
		return 0, errors.New("the node not have")
	}
	u := &s.(*Server).users
	for i := range p.UserInfo {
		u.Store(p.UserInfo[i].Uuid, struct{}{})
	}
	return len(p.UserInfo), nil
}

func (h *Hy) GetUserTraffic(tag, uuid string, reset bool) (up int64, down int64) {
	v, _ := h.servers.Load(tag)
	s := v.(*Server)
	auth := base64.StdEncoding.EncodeToString([]byte(uuid))
	up = s.counter.getCounters(auth).UpCounter.Load()
	down = s.counter.getCounters(auth).DownCounter.Load()
	if reset {
		s.counter.Reset(auth)
	}
	return
}

func (h *Hy) DelUsers(users []panel.UserInfo, tag string) error {
	v, e := h.servers.Load(tag)
	if !e {
		return errors.New("the node is not have")
	}
	s := v.(*Server)
	for i := range users {
		s.users.Delete(users[i].Uuid)
		s.counter.Delete(base64.StdEncoding.EncodeToString([]byte(users[i].Uuid)))
	}
	return nil
}

func (h *Hy) AddUserSpeedLimit(l *limiter.Limiter, tag string, user *panel.UserInfo, speedLimit int, expire int64) error {
	return l.AddDynamicSpeedLimit(tag, user, speedLimit, expire)
}
