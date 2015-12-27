// Copyright 2014 G&W. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package service

import (
	"encoding/json"
	"net/rpc"
	"time"
	"zim/dao"
	//"zim/sys"
)

type pushSvr struct{}

func NewPushSvr() *pushSvr {
	return &pushSvr{}
}

func (p *pushSvr) PushTip(query string, data []byte) (res string, code int) {
	//android
	res, islogin := LoginSvr.CheckLogin(query)
	if res == "" || !islogin {
		code = 3000
		return
	}
	c := NewConnectSvr()
	if err := json.Unmarshal([]byte(res), c); err != nil || c.User.Uid == "" {
		code = 4000
		return
	}
	//client := sys.RpcPool.GetClient(c.SvrIp, c.SvrPort)
	client, err := rpc.DialHTTP("tcp", c.SvrIp+":"+c.SvrPort)
	defer client.Close()
	if client != nil && err == nil {
		//defer sys.RpcPool.PutClient(c.SvrIp, c.SvrPort, client)
		r := make(map[string]string, 0)
		r["token"] = c.Token
		r["query"] = string(data)
		var reply []byte
		err := client.Call("RpcSvr.PushTip", r, &reply)
		if err == nil {
			res = string(reply)
			return
		}
	} else {
		code = 5002
	}
	//ios
	//coding...
	return
}

func (p *pushSvr) PushForce(token string, ud *dao.UserDao) {
	sd := dao.NewSendDao()
	sd.Stime, sd.Fuid, sd.Fname, sd.To, sd.Tuid = time.Now().Unix(), ud.Uid, ud.Username, "u/"+ud.Uid, ud.Uid
	pd := dao.NewPushDao()
	pd.Assert("force", sd)
	data, _ := json.Marshal(pd)
	p.PushTip("c/"+token, data)
}

//online offline
func (p *pushSvr) PushStatusToRoster(ud *dao.UserDao, status string) {
	rd := dao.NewRosterDao()
	if rm, err := rd.GetRoster(ud.Uid, ""); err == nil {
		for _, r := range rm {
			sd := dao.NewSendDao()
			sd.Stime, sd.Fuid, sd.Fname, sd.To, sd.Tuid = time.Now().Unix(), ud.Uid, ud.Username, "u/"+r.User.Uid, r.User.Uid
			pd := dao.NewPushDao()
			pd.Assert(status, sd)
			data, _ := json.Marshal(pd)
			p.PushTip("u/"+r.User.Uid, data)
		}
	}
}

func (p *pushSvr) PushStatusToGroup(ud *dao.UserDao, status string) {
	if c, err := ConnectHub.getConnectSvr("u/" + ud.Uid); err == nil {
		if u, err := c.getUser(); err == nil {
			if sm, err := u.GetGroupUser(ud.Uid, "normal"); err == nil {
				for _, gu := range sm {
					if gu.Uid == ud.Uid {
						continue
					}
					sd := dao.NewSendDao()
					sd.Stime, sd.Fuid, sd.Fname, sd.To, sd.Tuid = time.Now().Unix(), ud.Uid, ud.Username, "u/"+gu.Uid, gu.Uid
					pd := dao.NewPushDao()
					pd.Assert(status, sd)
					data, _ := json.Marshal(pd)
					p.PushTip("u/"+gu.Uid, data)
				}
			}
		}
	}
}
