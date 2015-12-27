package dao

import (
	"encoding/json"
	"time"
	"zim/sys"
)

type RequestDao struct {
	baseDao
	rid      int64  `db:"rid"`                     //rid
	Cmd      string `json:"cmd" db:"cmd"`          //action
	Fuid     string `json:"fuid" db:"fuid"`        //Sender Uid
	Fname    string `json:"fname" db:"fname"`      //Sender Uid
	To       string `json:"to" db:"to"`            //Receiver标签（u/*、g/*、t/*,b/*）
	Tuid     string `json:"tuid" db:"tuid"`        //tuid
	Expired  int    `json:"expired" db:"expired"`  //消息有效期，单位秒。默认86400
	Stime    int64  `json:"stime" db:"stime"`      //发送时间
	SeqCli   int64  `json:"seq_cli"  db:"seq_cli"` //seq_cli （cmd=receive时，用到）
	timeline int64  `db:"timeline"`

	Message map[string]string `json:"message" db:"message"` //消息
	Option  map[string]string `json:"option" db:"option"`   //附加信息
}

func NewRequestDao() (r *RequestDao) {
	r = &RequestDao{
		Message: make(map[string]string),
		Option:  make(map[string]string),
	}
	r.SetTableName("zim_request")
	return
}

func (r *RequestDao) Save() (id int64, err error) {
	//插入数据
	var message, option interface{}
	if r.Message != nil {
		message, _ = json.Marshal(r.Message)
	}
	if r.Option != nil {
		option, _ = json.Marshal(r.Option)
	}
	dbmap, err := sys.DbConn.Database()
	if err != nil {
		return
	}
	stmt, err := dbmap.Prepare("INSERT INTO `" + r.GetTableName() + "`(`cmd`,`fuid`,`fname`,`to`,`message`,`option`,`stime`,`expired`,`timeline`) VALUES(?,?,?,?,?,?,?,?,?)")
	defer stmt.Close()
	res, _ := stmt.Exec(r.Cmd, r.Fuid, r.Fname, r.To, message, option, r.Stime, r.Expired, time.Now().Unix())
	if r.rid, err = res.LastInsertId(); err == nil {
		id = r.rid
	}
	return
}
