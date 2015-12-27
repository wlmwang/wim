// Copyright 2014 G&W. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package service

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
	"zim/dao"
	"zim/sys"
)

const (
	writeWait      = 30 * time.Second //Time allowed to write a message to the peer.
	pongWait       = 10 * time.Second //Time allowed to read the next pong message from the peer.
	maxMessageSize = 1024             //Maximum message size allowed from peer.
)

type connection struct {
	Identity string       `json:"identity"` //客户端identity
	CliSck   string       `json:"clisck"`   //客户端ip:port
	SvrIp    string       `json:"svr_ip"`   //服务端ip
	SvrPort  string       `json:"svr_port"` //服务端port
	Token    string       `json:"token"`    //标示该连接的随机字符。
	Device   string       `json:"device"`   //设备类型（web、ios、android、winphone）
	disabled int          `json:"disabled"` //0打开 1断开
	User     *dao.UserDao `json:"user"`     //用户句柄
}
type connectSvr struct {
	ws *websocket.Conn //websocket connection
	connection
}

func NewConnection() *connection {
	return &connection{
		User: new(dao.UserDao),
	}
}

func NewConnectSvr() *connectSvr {
	return &connectSvr{
		connection: connection{User: new(dao.UserDao)},
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  maxMessageSize,
	WriteBufferSize: maxMessageSize,
}

/**
 * 连接
 */
func (c *connectSvr) Connect(w http.ResponseWriter, r *http.Request) (int, error) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if _, ok := err.(websocket.HandshakeError); ok {
		return 4007, errors.New(sys.LangConf.Get("4007").MustString())
	} else if err != nil {
		return 4008, errors.New(sys.LangConf.Get("4008").MustString())
	}
	//登陆信息
	if res, islogin := LoginSvr.CheckLogin("c/" + r.FormValue("token")); res != "" && islogin {
		if err := json.Unmarshal([]byte(res), c); err != nil {
			return 4011, errors.New(sys.LangConf.Get("4011").MustString())
		}
	}
	c.ws = ws
	c.Token = "c/" + c.Token
	//添加连接池
	//ConnectHub.register <- c
	ConnectHub.add(c)
	return 0, err
}

/**
 * 循环读取消息
 */
func (c *connectSvr) Reader() (err error) {
	defer func() {
		token, _ := c.getToken()
		ConnectHub.close(token)
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		requestSrv := NewRequestSrv()
		requestSrv.parse(c, message)
	}
	return err
}

func (c *connectSvr) discon() error {
	c.disabled = 1
	c.ws.Close()
	return nil
}

func (c *connectSvr) getToken() (string, error) {
	if c.disabled == 1 {
		return c.Token, errors.New(sys.LangConf.Get("3000").MustString())
	}
	return c.Token, nil
}

func (c *connectSvr) getDevice() (string, error) {
	if c.disabled == 1 {
		return "", errors.New(sys.LangConf.Get("3000").MustString())
	}
	return c.Device, nil
}

func (c *connectSvr) setUser(user *dao.UserDao) (bool, error) {
	if c.disabled == 1 {
		return false, errors.New(sys.LangConf.Get("3000").MustString())
	}
	c.User = user
	return true, nil
}

func (c *connectSvr) getUser() (*dao.UserDao, error) {
	if c.disabled == 1 {
		return nil, errors.New(sys.LangConf.Get("3000").MustString())
	}
	if user := c.User; user == nil || user.Uid == "" {
		return nil, errors.New(sys.LangConf.Get("4011").MustString())
	}
	return c.User, nil
}

func (c *connectSvr) sendText(message []byte) error {
	if _, err := c.getToken(); err != nil {
		return errors.New(sys.LangConf.Get("3000").MustString())
	}
	return c.write(websocket.TextMessage, message)
}

func (c *connectSvr) sendBinary(message []byte) error {
	if _, err := c.getToken(); err != nil {
		return errors.New(sys.LangConf.Get("3000").MustString())
	}
	return c.write(websocket.BinaryMessage, message)
}

func (c *connectSvr) ping() error {
	return c.write(websocket.PingMessage, []byte{})
}

/**
 * write writes a message with the given message type and payload.
 */
func (c *connectSvr) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}
