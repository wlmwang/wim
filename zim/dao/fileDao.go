package dao

import (
//"time"
//"zim/sys"
)

type fileDao struct {
	baseDao
	fid      int64
	Path     string //地址
	Sort     string //分类
	Appid    string
	uid      string
	timeline int64
}

func NewFileDao() (f *fileDao) {
	f = &fileDao{}
	f.SetTableName("zim_file")
	return
}

/*
func (f *fileDao) Save(sort, filepath, uid string) (fid int64, err error) {
	//插入数据
	f.Path, f.Sort, f.uid = filepath, sort, uid
	dbmap, err := sys.DbConn.Database()
	if err != nil {
		return
	}
	stmt, err := dbmap.Prepare("INSERT INTO `" + f.GetTableName() + "`(`path`,`sort`,`uid`,`timeline`) VALUES(?,?,?,?)")
	defer stmt.Close()
	res, _ := stmt.Exec(filepath, sort, uid, time.Now().Unix())
	if f.fid, err = res.LastInsertId(); err == nil {
		fid = f.fid
	}
	return
}
*/
