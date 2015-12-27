package service

import (
	"encoding/json"
	"errors"
	"zim/dao"
)

type RpcSvr struct{}

func NewRpcSvr() *RpcSvr {
	return &RpcSvr{}
}

func (r *RpcSvr) CheckLogin(args string, res *[]byte) (err error) {
	LoginSvr.lock.RLock()
	defer LoginSvr.lock.RUnlock()
	if args[0:2] == "c/" {
		if ls, ok := LoginSvr.tToServer[args[2:]]; ok {
			*res, _ = json.Marshal(ls)
			return
		}
	} else if args[0:2] == "u/" {
		if token, ok := LoginSvr.uToToken[args[2:]]; ok {
			if ls, ok := LoginSvr.tToServer[token]; ok {
				*res, _ = json.Marshal(ls)
				return
			}
		}
	}
	err = errors.New("4011")
	return
}

func (r *RpcSvr) Logout(args string, res *string) (err error) {
	LoginSvr.lock.RLock()
	defer LoginSvr.lock.RUnlock()
	if args[0:2] == "c/" {
		if ls, ok := LoginSvr.tToServer[args[2:]]; ok {
			//离线通知push tip
			p := NewPushSvr()
			p.PushStatusToRoster(ls.User, "offline")
			//end
			delete(LoginSvr.tToServer, args[2:])
			delete(LoginSvr.uToToken, ls.User.Uid)
			LoginSvr.online--
			return
		}
	} else if args[0:2] == "u/" {
		if token, ok := LoginSvr.uToToken[args[2:]]; ok {
			if ls, ok := LoginSvr.tToServer[token]; ok {
				//离线通知push tip
				p := NewPushSvr()
				p.PushStatusToRoster(ls.User, "offline")
				//end
				delete(LoginSvr.tToServer, token)
				delete(LoginSvr.uToToken, ls.User.Uid)
				LoginSvr.online--
				return
			}
		}
	}
	err = errors.New("4011")
	return
}

func (r *RpcSvr) PushTip(args map[string]string, res *[]byte) (err error) {
	token := args["token"]
	query := args["query"]
	if token == "" || query == "" {
		err = errors.New("4008")
		return
	}
	pd := dao.NewPushDao()
	if err = json.Unmarshal([]byte(query), pd); err != nil {
		return
	}
	c, err := ConnectHub.getConnectSvr("c/" + token)
	if err != nil {
		return
	}
	if *res, err = json.Marshal(pd); err == nil {
		c.sendText(*res)
	}
	return
}
