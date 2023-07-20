package xray

import (
	"context"
	"fmt"
	"github.com/Yuzuki616/V2bX/limiter"

	"github.com/Yuzuki616/V2bX/api/panel"
	"github.com/Yuzuki616/V2bX/common/format"
	vCore "github.com/Yuzuki616/V2bX/core"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/proxy"
)

func (c *Core) GetUserManager(tag string) (proxy.UserManager, error) {
	handler, err := c.ihm.GetHandler(context.Background(), tag)
	if err != nil {
		return nil, fmt.Errorf("no such inbound tag: %s", err)
	}
	inboundInstance, ok := handler.(proxy.GetInbound)
	if !ok {
		return nil, fmt.Errorf("handler %s is not implement proxy.GetInbound", tag)
	}
	userManager, ok := inboundInstance.GetInbound().(proxy.UserManager)
	if !ok {
		return nil, fmt.Errorf("handler %s is not implement proxy.UserManager", tag)
	}
	return userManager, nil
}

func (c *Core) DelUsers(users []panel.UserInfo, tag string) error {
	userManager, err := c.GetUserManager(tag)
	if err != nil {
		return fmt.Errorf("get user manager error: %s", err)
	}
	var up, down, user string
	for i := range users {
		user = format.UserTag(tag, users[i].Uuid)
		err = userManager.RemoveUser(context.Background(), user)
		if err != nil {
			return err
		}
		up = "user>>>" + user + ">>>traffic>>>uplink"
		down = "user>>>" + user + ">>>traffic>>>downlink"
		c.shm.UnregisterCounter(up)
		c.shm.UnregisterCounter(down)
	}
	return nil
}

func (c *Core) GetUserTraffic(tag, uuid string, reset bool) (up int64, down int64) {
	upName := "user>>>" + format.UserTag(tag, uuid) + ">>>traffic>>>uplink"
	downName := "user>>>" + format.UserTag(tag, uuid) + ">>>traffic>>>downlink"
	upCounter := c.shm.GetCounter(upName)
	downCounter := c.shm.GetCounter(downName)
	if reset {
		if upCounter != nil {
			up = upCounter.Set(0)
		}
		if downCounter != nil {
			down = downCounter.Set(0)
		}
	} else {
		if upCounter != nil {
			up = upCounter.Value()
		}
		if downCounter != nil {
			down = downCounter.Value()
		}
	}
	return up, down
}

func (c *Core) AddUsers(p *vCore.AddUsersParams) (added int, err error) {
	users := make([]*protocol.User, 0, len(p.UserInfo))
	switch p.NodeInfo.Type {
	case "v2ray":
		if p.NodeInfo.ExtraConfig.EnableVless {
			users = buildVlessUsers(p.Tag, p.UserInfo, p.NodeInfo.ExtraConfig.VlessFlow)
		} else {
			users = buildVmessUsers(p.Tag, p.UserInfo)
		}
	case "trojan":
		users = buildTrojanUsers(p.Tag, p.UserInfo)
	case "shadowsocks":
		users = buildSSUsers(p.Tag,
			p.UserInfo,
			p.NodeInfo.Cipher,
			p.NodeInfo.ServerKey)
	default:
		return 0, fmt.Errorf("unsupported node type: %s", p.NodeInfo.Type)
	}
	man, err := c.GetUserManager(p.Tag)
	if err != nil {
		return 0, fmt.Errorf("get user manager error: %s", err)
	}
	for _, u := range users {
		mUser, err := u.ToMemoryUser()
		if err != nil {
			return 0, err
		}
		err = man.AddUser(context.Background(), mUser)
		if err != nil {
			return 0, err
		}
	}
	return len(users), nil
}

func (c *Core) AddUserSpeedLimit(l *limiter.Limiter, tag string, user *panel.UserInfo, speedLimit int, expire int64) error {
	return l.AddDynamicSpeedLimit(tag, user, speedLimit, expire)
}
