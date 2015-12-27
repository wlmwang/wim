package service

import (
	"encoding/json"
	"zim/common"
	"zim/dao"
	"zim/sys"
)

type receiveSvr struct {
	req chan *dao.RequestDao
}

var ReceiveSvr *receiveSvr

func NewReceiveSvr() *receiveSvr {
	return &receiveSvr{
		req: make(chan *dao.RequestDao, 10000), //消息中间件中最多10000条数据同时处理
	}
}

func (r *receiveSvr) Run() {
	for {
		select {
		case req := <-r.req:
			r.handle(req)
		}
	}
}

/**
 * 获取消息处理器
 * websocket
 */
func (r *receiveSvr) handle(req *dao.RequestDao) (err error) {
	//defer common.HandleError()
	if req.Tuid != "" {
		if c, err := ConnectHub.getConnectSvr("u/" + req.Tuid); err == nil {
			rd := dao.NewReceiveDao()
			rd.Get(req)
			if len(rd.Message) > 0 {
				if message, err := json.Marshal(rd); err == nil {
					c.sendText(message)
				}
			} else {
				c.sendText([]byte{})
			}
		} else {
			common.LogSvr.Println("info:" + err.Error())
		}
	} else {
		common.LogSvr.Println("info: "+sys.LangConf.Get("4000").MustString(), err)
		return
	}
	return
}
