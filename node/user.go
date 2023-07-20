package node

import (
	"github.com/Yuzuki616/V2bX/api/panel"
	log "github.com/sirupsen/logrus"
	"runtime"
	"strconv"
	"time"
)

func (c *Controller) reportUserTrafficTask() (err error) {
	// Get User traffic
	userTraffic := make([]panel.UserTraffic, 0)
	for i := range c.userList {
		up, down := c.server.GetUserTraffic(c.Tag, c.userList[i].Uuid, true)
		if up > 0 || down > 0 {
			if c.LimitConfig.EnableDynamicSpeedLimit {
				c.userList[i].Traffic += up + down
			}
			userTraffic = append(userTraffic, panel.UserTraffic{
				UID:      (c.userList)[i].Id,
				Upload:   up,
				Download: down})
		}
	}
	if len(userTraffic) > 0 {
		err = c.apiClient.ReportUserTraffic(userTraffic)
		if err != nil {
			log.WithFields(log.Fields{
				"tag": c.Tag,
				"err": err,
			}).Info("Report user traffic failed")
		} else {
			log.WithField("tag", c.Tag).Infof("Report %d online users", len(userTraffic))
		}
	}
	userTraffic = nil
	runtime.GC()
	return nil
}

func compareUserList(old, new []panel.UserInfo) (deleted, added []panel.UserInfo) {
	tmp := map[string]struct{}{}
	tmp2 := map[string]struct{}{}
	for i := range old {
		tmp[old[i].Uuid+strconv.Itoa(old[i].SpeedLimit)] = struct{}{}
	}
	l := len(tmp)
	for i := range new {
		e := new[i].Uuid + strconv.Itoa(new[i].SpeedLimit)
		tmp[e] = struct{}{}
		tmp2[e] = struct{}{}
		if l != len(tmp) {
			added = append(added, new[i])
			l++
		}
	}
	tmp = nil
	l = len(tmp2)
	for i := range old {
		tmp2[old[i].Uuid+strconv.Itoa(old[i].SpeedLimit)] = struct{}{}
		if l != len(tmp2) {
			deleted = append(deleted, old[i])
			l++
		}
	}
	return deleted, added
}

func (c *Controller) dynamicSpeedLimit() error {
	if c.LimitConfig.EnableDynamicSpeedLimit {
		for i := range c.userList {
			up, down := c.server.GetUserTraffic(c.Tag, c.userList[i].Uuid, false)
			if c.userList[i].Traffic+down+up/1024/1024 > c.LimitConfig.DynamicSpeedLimitConfig.Traffic {
				err := c.server.AddUserSpeedLimit(c.limiter, c.Tag,
					&c.userList[i],
					c.LimitConfig.DynamicSpeedLimitConfig.SpeedLimit,
					time.Now().Add(time.Second*time.Duration(c.LimitConfig.DynamicSpeedLimitConfig.ExpireTime)).Unix())
				if err != nil {
					log.Print(err)
				}
			}
			c.userList[i].Traffic = 0
		}
	}
	return nil
}
