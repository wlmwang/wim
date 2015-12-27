package dao

import (
	"encoding/json"
	"strconv"
	"time"
	"zim/sys"
)

type ReceiveDao struct {
	baseDao
	Cmd     string    `json:"cmd" db:"cmd"` //action
	SeqCli  int64     `json:"seq_cli"`      //seq_cli
	SeqSvr  int64     `json:"seq_svr"`      //seq_svr
	Message []SendDao `json:"message"`      //message
}

func NewReceiveDao() (r *ReceiveDao) {
	r = &ReceiveDao{
		Message: make([]SendDao, 0),
	}
	r.SetTableName("zim_msgstore")
	return
}

func (r *ReceiveDao) Get(req *RequestDao) (err error) {
	u := NewUserDao()
	news, err := u.haveMessage(req.Tuid, req.SeqCli)
	if req.SeqCli > u.SeqCli && req.SeqCli <= u.SeqSvr {
		u.setSeq(req.SeqCli, req.Tuid, "c")
	}
	if !news {
		return
	}
	sql := "SELECT `mid`,`rid`,`seq`,`cmd`,`fuid`,`fname`,`to`,`tuid`,`message`,`option`,`stime`,`expired` FROM `" + r.GetTableName() + "` " + "WHERE `tuid`='" + req.Tuid + "' AND `seq`> " + strconv.FormatInt(req.SeqCli, 10) + " ORDER BY `seq` ASC LIMIT 20"
	dbmap, err := sys.DbConn.Database()
	if err != nil {
		return
	}
	urow, err := dbmap.Query(sql)
	defer urow.Close()
	if err != nil {
		return
	}
	var seqCli int64 = 0
	for urow.Next() {
		sd := NewSendDao()
		var message, option string
		if err := urow.Scan(&sd.mid, &sd.rid, &sd.Seq, &sd.Cmd, &sd.Fuid, &sd.Fname, &sd.To, &sd.Tuid, &message, &option, &sd.Stime, &sd.Expired); err == nil {
			if sd.Expired > 0 && sd.Stime+int64(sd.Expired) < time.Now().Unix() {
				continue
			}
			if message != "" {
				json.Unmarshal([]byte(message), &sd.Message)
			}
			if option != "" {
				json.Unmarshal([]byte(option), &sd.Option)
			}
			r.Message = append(r.Message, *sd)
			seqCli = sd.Seq
		}
	}
	r.SeqCli, r.SeqSvr = seqCli, u.SeqSvr
	r.Cmd = "message"
	return
}
