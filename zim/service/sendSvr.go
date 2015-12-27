package service

import (
	"encoding/json"
	"zim/dao"
)

type sendSrv struct {
	req chan *dao.RequestDao
}

var SendSrv *sendSrv

func NewSendSrv() *sendSrv {
	return &sendSrv{
		req: make(chan *dao.RequestDao, 10000), //消息中间件中最多10000条数据同时处理
	}
}

func (s *sendSrv) Run() {
	for {
		select {
		case req := <-s.req:
			s.handle(req)
		}
	}
}

/**
 * 发送消息处理器
 * 包含：单聊消息、群聊消息
 */
func (s *sendSrv) handle(req *dao.RequestDao) (err error) {
	//defer common.HandleError()
	if req.To == "" {
		return
	}
	switch req.To[0:2] {
	case "u/":
		sd := dao.NewSendDao()
		sd.SetTuid(string(req.To[2:]))
		sd.Assert(req)
		sd.Save()
		//push tip
		pd := dao.NewPushDao()
		pd.Assert("push", sd)
		p := NewPushSvr()
		data, _ := json.Marshal(pd)
		p.PushTip("u/"+req.To[2:], data)
		//end
	case "g/":
		u := dao.NewUserDao()
		if sm, err := u.GetGroupUser(req.To[2:], "normal"); err == nil {
			for _, gu := range sm {
				if gu.Uid == req.Fuid {
					continue
				}
				sd := dao.NewSendDao()
				sd.SetTuid(gu.Uid)
				sd.Assert(req)
				sd.Save()
				//push tip
				pd := dao.NewPushDao()
				pd.Assert("push", sd)
				p := NewPushSvr()
				data, _ := json.Marshal(pd)
				p.PushTip("u/"+gu.Uid, data)
				//end
			}
		}
	case "s/":
		u := dao.NewUserDao()
		if sm, err := u.GetGroupUser(req.To[2:], "session"); err == nil {
			for _, gu := range sm {
				if gu.Uid == req.Fuid {
					continue
				}
				sd := dao.NewSendDao()
				sd.SetTuid(gu.Uid)
				sd.Assert(req)
				sd.Save()
				//push tip
				pd := dao.NewPushDao()
				pd.Assert("push", sd)
				p := NewPushSvr()
				data, _ := json.Marshal(pd)
				p.PushTip("u/"+gu.Uid, data)
				//end
			}
		}
	case "t/":
		var sUid []string
		if req.To == "t/all" {
			u := dao.NewUserDao()
			sUid, err = u.GetAllUid()
		} else {
			td := dao.NewTagDao()
			sUid, err = td.GetUid(req.To[2:])
		}
		if err == nil {
			for _, v := range sUid {
				if v == req.Fuid {
					continue
				}
				sd := dao.NewSendDao()
				sd.SetTuid(v)
				sd.Assert(req)
				sd.Save()
				//push tip
				pd := dao.NewPushDao()
				pd.Assert("push", sd)
				p := NewPushSvr()
				data, _ := json.Marshal(pd)
				p.PushTip("u/"+v, data)
				//end
			}
		}
	case "b/":
		req.Cmd = "broadcast"
		if req.To == "b/online" {
			sd := dao.NewSendDao()
			sd.Assert(req)
			rd := dao.NewReceiveDao()
			rd.Cmd = req.Cmd
			rd.Message = append(rd.Message, *sd)
			if message, err := json.Marshal(rd); err == nil {
				ConnectHub.Broadcast(message)
			}
		}
	case "c/":
	}
	return
}

func (s *sendSrv) autoAsk(req *dao.RequestDao) (err error) {
	//defer common.HandleError()
	if connectSrv, err := ConnectHub.getConnectSvr("u/" + req.Fuid); err == nil {
		sd := dao.NewSendDao()
		sd.Assert(req)
		sd.SetTuid(string(req.Fuid))
		sd.Fuid = string(req.To[2:])
		sd.Message["content"] = "这是自动回应消息"

		rd := dao.NewReceiveDao()
		rd.Cmd = "message"
		rd.Message = append(rd.Message, *sd)
		if message, err := json.Marshal(rd); err == nil {
			connectSrv.sendText(message)
		}
	}
	return
}
