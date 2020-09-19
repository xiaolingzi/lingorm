package sqlite

import (
	"database/sql"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/xiaolingzi/lingorm/internal/common"

	_ "github.com/mattn/go-sqlite3"
)

// Native struct
type Native struct {
	DatabaseConfigKey string
}

var sqlTxList map[string]*sql.Tx
var db *sql.DB

// NewNative the instance of Native
func NewNative(databaseConfigKey string) *Native {
	var m Native
	m.DatabaseConfigKey = databaseConfigKey
	return &m
}

func (m *Native) connect() {
	databaseInfo := NewConfig().GetDatabaseInfo(m.DatabaseConfigKey)

	dsn := "file:" + databaseInfo.File + "?cache=shared&_loc=auto"
	if databaseInfo.User != "" && databaseInfo.Password != "" {
		dsn += "&_auth_user=" + databaseInfo.User + "&_auth_pass=" + databaseInfo.Password
		if databaseInfo.Crypt != "" {
			dsn += "&_auth_crypt=" + databaseInfo.Crypt
		}
		if databaseInfo.Salt != "" {
			dsn += "&_auth_salt=" + databaseInfo.Salt
		}
		if databaseInfo.Timeout > 0 {
			dsn += "&_timeout=" + strconv.Itoa(databaseInfo.Timeout)
		}
	}
	var err error
	if db == nil {
		db, err = sql.Open("sqlite3", dsn)
		if err != nil {
			common.NewError().Throw(err)
		}
	}
}

// Excute excute sql
func (m *Native) Execute(query string, params map[string]interface{}, transactionKey string) (int, int) {
	tempQuery, paramList := m.convertSQL(query, params)

	var stmt *sql.Stmt
	defer func() {
		if stmt != nil {
			stmt.Close()
		}
	}()
	var res sql.Result
	var err error
	_, ok := sqlTxList[transactionKey]
	if transactionKey != "" && ok {
		tx := sqlTxList[transactionKey]
		stmt, err = tx.Prepare(tempQuery)
		if err != nil {
			common.NewError().Throw(err)
		}
		res, err = stmt.Exec(paramList...)
		if err != nil {
			common.NewError().Throw(err)
		}
	} else {
		m.connect()
		stmt, err = db.Prepare(tempQuery)
		if err != nil {
			common.NewError().Throw(err)
		}
		res, err = db.Exec(tempQuery, paramList...)
		if err != nil {
			common.NewError().Throw(err)
		}
	}
	id, _ := res.LastInsertId()
	count, _ := res.RowsAffected()
	return int(count), int(id)
}

// FetchOne return the first row that meet query criteria
func (m *Native) FetchOne(query string, params map[string]interface{}, transactionKey string) (map[string]string, error) {
	result := m.FetchAll(query, params, transactionKey)
	if len(result) > 0 {
		return result[0], nil
	}
	return nil, nil
}

// FetchAll return all the rows that meet query criteria
func (m *Native) FetchAll(query string, params map[string]interface{}, transactionKey string) []map[string]string {
	tempSQL, paramList := m.convertSQL(query, params)

	var stmt *sql.Stmt
	var rows *sql.Rows
	defer func() {
		if stmt != nil {
			stmt.Close()
		}
		if rows != nil {
			rows.Close()
		}
	}()

	var err error
	_, ok := sqlTxList[transactionKey]
	if transactionKey != "" && ok {
		tx := sqlTxList[transactionKey]
		stmt, err = tx.Prepare(tempSQL)
		if err != nil {
			common.NewError().Throw(err)
		}
		rows, err = stmt.Query(paramList...)
		if err != nil {
			common.NewError().Throw(err)
		}
	} else {
		m.connect()
		stmt, err = db.Prepare(tempSQL)
		if err != nil {
			common.NewError().Throw(err)
		}
		rows, err = stmt.Query(paramList...)
		if err != nil {
			common.NewError().Throw(err)
		}
	}

	result := m.convertRowsToMapList(rows)
	return result
}

func (m *Native) Begin() string {
	key := strconv.FormatInt(time.Now().UnixNano(), 10)
	m.connect()
	tx, err := db.Begin()
	if err != nil {
		common.NewError().Throw(err)
	}
	if sqlTxList == nil {
		sqlTxList = make(map[string]*sql.Tx)
	}
	sqlTxList[key] = tx
	return key
}

func (m *Native) Commit(transactionKey string) {
	if _, ok := sqlTxList[transactionKey]; ok {
		tx := sqlTxList[transactionKey]
		err := tx.Commit()
		if err != nil {
			common.NewError().Throw(err)
		}
		delete(sqlTxList, transactionKey)
	} else {
		common.NewError().Throw("Begin a transaction first before commit")
	}
}

func (m *Native) Rollback(transactionKey string) {
	if _, ok := sqlTxList[transactionKey]; ok {
		tx := sqlTxList[transactionKey]
		err := tx.Rollback()
		if err != nil {
			common.NewError().Throw(err)
		}
		delete(sqlTxList, transactionKey)
	} else {
		common.NewError().Throw("Begin a transaction first before rollback")
	}
}

func (m *Native) convertSQL(query string, params map[string]interface{}) (string, []interface{}) {
	if params == nil {
		return query, nil
	}
	reg, _ := regexp.Compile(":[a-zA-Z0-9_\\-]+")
	matches := reg.FindAll([]byte(query), -1)
	var paramList []interface{}
	for i := 0; i < len(matches); i++ {
		str := string(matches[i])
		key := strings.Trim(str, ":")
		query = strings.Replace(query, str, "?", 1)
		paramList = append(paramList, params[key])
	}

	return query, paramList
}

func (m *Native) convertRowsToMapList(rows *sql.Rows) []map[string]string {
	cols, _ := rows.Columns()
	vals := make([][]byte, len(cols))
	scans := make([]interface{}, len(cols))
	for k := range vals {
		scans[k] = &vals[k]
	}
	var result []map[string]string
	for rows.Next() {
		rows.Scan(scans...)
		row := make(map[string]string)
		for k, v := range vals {
			key := cols[k]
			row[key] = string(v)
		}
		result = append(result, row)
	}
	return result
}
