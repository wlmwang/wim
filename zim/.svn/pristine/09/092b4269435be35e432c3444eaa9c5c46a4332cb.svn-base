package service

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/rpc"
	"strconv"
	"sync"
	"time"
	"zim/common"
	"zim/dao"
	"zim/sys"
)

type loginSvr struct {
	lock      *sync.RWMutex
	online    int
	tToServer map[string]connection //token conn
	uToToken  map[string]string     //uid token
}

var LoginSvr *loginSvr

func NewLoginSvr() *loginSvr {
	return &loginSvr{
		lock:      new(sync.RWMutex),
		tToServer: make(map[string]connection),
		uToToken:  make(map[string]string),
	}
}

func (l *loginSvr) LoginSvrHandle(w http.ResponseWriter, r *http.Request) (res string, code int) {
	//defer common.HandleError()
	l.lock.RLock()
	defer l.lock.RUnlock()
	if l.online >= sys.BaseConf.Get("online").MustInt() {
		code = 4015
		return
	}
	switch r.FormValue("act") {
	case "login":
		identity := r.FormValue("identity")
		appid := r.FormValue("appid")
		username := r.FormValue("username")
		password := r.FormValue("password")
		force := r.FormValue("force")
		if appid == "" || username == "" || password == "" {
			code = 4008
		}
		res, code = l.Login(username, password, appid, force, identity, r.RemoteAddr, common.GetDevice(r))
	case "logout":
		token := r.FormValue("token")
		if token != "" {
			l.Logout("c/" + token)
		}
	default:
		code = 4008
	}
	return
}
func (l *loginSvr) Login(username, password, appid, force, identity, addr, device string) (res string, code int) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	//rand token start
	t, i := common.RandomStr(), 0
	for _, ok := l.tToServer[t]; ok; i++ {
		t = common.RandomStr()
		_, ok = l.tToServer[t]
		if i > 50 {
			code = 5008
			break
		}
	}
	//end
	userDao := dao.NewUserDao()
	if userDao.Login(username, password, appid); userDao.Uid == "" {
		code = 4014
		return
	}
	oldInfo, islogin := l.CheckLogin("u/" + userDao.Uid)
	if oldInfo != "" && islogin && force == "" {
		code = 4010
		return
	} else if oldInfo != "" && islogin && force != "" {
		//强登通知push tip
		c := NewConnectSvr()
		if err := json.Unmarshal([]byte(oldInfo), c); err == nil && c.Identity != identity {
			p := NewPushSvr()
			p.PushForce(c.Token, userDao)
		}
		//end
		l.Logout("u/" + userDao.Uid)
	}
	//rand server start
	allNs, _ := sys.BaseConf.Get("ns").Map()
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	ns, _ := sys.BaseConf.Get("ns").Get(strconv.Itoa(seed.Intn(len(allNs)))).Map()
	//end
	c := NewConnection()
	c.Identity = identity
	c.CliSck = addr
	c.SvrIp = ns["ip"].(string)
	c.SvrPort = ns["port"].(string)
	c.User = userDao
	c.Device = device
	c.Token = t
	l.tToServer[t] = *c
	l.uToToken[userDao.Uid] = t
	l.online = len(l.tToServer)
	//上线通知push tip
	p := NewPushSvr()
	p.PushStatusToRoster(userDao, "online")
	//end
	info, _ := json.Marshal(c)
	res = string(info)
	return
}

func (l *loginSvr) CheckLogin(query string) (res string, islogin bool) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	//defer common.HandleError()
	local := sys.BaseConf.Get("mode").MustString()
	//本地检测
	if local[:2] == "ls" {
		r := NewRpcSvr()
		var reply []byte
		if err := r.CheckLogin(query, &reply); err == nil {
			res = string(reply)
			islogin = true
			return
		}
	}
	//远程检测
	ls, _ := sys.BaseConf.Get("ls").Map()
	for i := 0; i < len(ls); i++ {
		if local[2:] == strconv.Itoa(i) {
			continue
		}
		dd, _ := sys.BaseConf.Get("ls").Get(strconv.Itoa(i)).Map()
		//client := sys.RpcPool.GetClient(dd["ip"].(string), dd["port"].(string))
		client, err := rpc.DialHTTP("tcp", dd["ip"].(string)+":"+dd["port"].(string))
		defer client.Close()
		if client != nil && err == nil {
			//defer sys.RpcPool.PutClient(dd["ip"].(string), dd["port"].(string), client)
			var reply []byte
			err := client.Call("RpcSvr.CheckLogin", query, &reply)
			if err != nil {
				continue
			} else {
				res = string(reply)
				islogin = true
				break
			}
		}
	}
	return
}

func (l *loginSvr) Logout(query string) (isout bool) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	isout = true
	res, islogin := LoginSvr.CheckLogin(query)
	if res == "" || !islogin {
		return
	}
	c := NewConnectSvr()
	if err := json.Unmarshal([]byte(res), c); err != nil || c.User.Uid == "" {
		return
	}
	client, err := rpc.DialHTTP("tcp", c.SvrIp+":"+c.SvrPort)
	defer client.Close()
	//client := sys.RpcPool.GetClient(c.SvrIp, c.SvrPort)
	if client != nil && err == nil {
		//defer sys.RpcPool.PutClient(c.SvrIp, c.SvrPort, client)
		var reply []byte
		err := client.Call("RpcSvr.Logout", query, &reply)
		if err != nil {
			isout = false
		}
	}
	return
}
