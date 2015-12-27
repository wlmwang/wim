package dao

import (
	"strconv"
	"time"
	"zim/sys"
)

type TagDao struct {
	baseDao
	Tid      string `json:"tid" db:"tid"`
	Tname    string `json:"tname" db:"tname"`
	Uid      string `json:"uid" db:"uid"`
	timeline int64  `db:"timeline"`
}

func NewTagDao() (t *TagDao) {
	t = &TagDao{}
	t.SetTableName("zim_tag")
	return
}

func (t *TagDao) AddTag(uid, tname string) (id string, err error) {
	dbmap, err := sys.DbConn.Database()
	if err != nil {
		return
	}
	sql := "INSERT INTO `" + t.GetTableName() + "`(`uid`,`tname`,`timeline`) VALUES(?,?,?)"
	stmt, _ := dbmap.Prepare(sql)
	defer stmt.Close()
	res, err := stmt.Exec(uid, tname, time.Now().Unix())
	if err == nil {
		if i, err := res.LastInsertId(); err == nil {
			t.Tid, id = strconv.FormatInt(i, 10), strconv.FormatInt(i, 10)
		}
	}
	return
}

func (t *TagDao) DelTag(uid, tname string) (err error) {
	dbmap, err := sys.DbConn.Database()
	if err != nil {
		return
	}
	sql := "DELETE FROM `" + t.GetTableName() + "` WHERE `uid`=? AND `tname`=?"
	stmt, _ := dbmap.Prepare(sql)
	defer stmt.Close()
	_, err = stmt.Exec(uid, tname)
	return
}

func (t *TagDao) GetTagByUT(uid, tname string) (err error) {
	dbmap, err := sys.DbConn.Database()
	if err != nil {
		return
	}
	sql := "SELECT `tid`,`uid`,`tname` FROM `" + t.GetTableName() + "` WHERE `uid`='" + uid + "' AND tname='" + tname + "'"
	urow, err := dbmap.Query(sql)
	defer urow.Close()
	if err == nil {
		urow.Next()
		err = urow.Scan(&t.Tid, &t.Uid, &t.Tname)
	}
	return
}

func (t *TagDao) GetTag(uid string) (tm []string, err error) {
	tm = make([]string, 0)
	sql := "SELECT `tname` FROM `" + t.GetTableName() + "` WHERE AND `uid`=" + uid
	dbmap, err := sys.DbConn.Database()
	if err != nil {
		return
	}
	urow, err := dbmap.Query(sql)
	defer urow.Close()
	for urow.Next() {
		tname := ""
		if err = urow.Scan(&tname); err == nil {
			tm = append(tm, tname)
		}
	}
	return
}

func (t *TagDao) GetUid(tname string) (tm []string, err error) {
	dbmap, err := sys.DbConn.Database()
	if err != nil {
		return
	}
	tm = make([]string, 0)
	sql := "SELECT `uid` FROM `" + t.GetTableName() + "` WHERE `tname`='" + tname + "'"
	urow, err := dbmap.Query(sql)
	defer urow.Close()
	for urow.Next() {
		if err = urow.Scan(&t.Uid); err == nil {
			tm = append(tm, t.Uid)
		}
	}
	return
}
