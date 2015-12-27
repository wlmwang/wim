package sys

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"sync"
)

type dbConn struct {
	lock  *sync.RWMutex
	dList map[string]*sql.DB
}

var DbConn *dbConn

func NewDbConn() (d *dbConn) {
	d = &dbConn{
		lock:  new(sync.RWMutex),
		dList: make(map[string]*sql.DB),
	}
	return
}

/**
 * 数据库连接
 *
 */
func (d *dbConn) Database() (dbmap *sql.DB, err error) {
	c := BaseConf.Get("db")
	database := c.Get("database").MustString()
	dbmap, ok := d.dList[database]
	if ok && dbmap != nil {
		return
	}
	d.lock.RLock()
	defer d.lock.RUnlock()
	dbn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", c.Get("user").MustString(), c.Get("password").MustString(), c.Get("host").MustString(), c.Get("port").MustInt(), database)
	dbmap, err = sql.Open("mysql", dbn)
	if err == nil {
		dbmap.SetMaxOpenConns(c.Get("maxOpen").MustInt())
		dbmap.SetMaxIdleConns(c.Get("maxIdle").MustInt())
		d.dList[database] = dbmap
	}
	return
}
