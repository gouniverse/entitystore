package entitystore

import (
	"database/sql"
	"errors"
)

// Store defines an entity store
type Store struct {
	entityTableName         string
	attributeTableName      string
	entityTrashTableName    string
	attributeTrashTableName string
	db                      *sql.DB
	dbDriverName            string
	automigrateEnabled      bool
	debugEnabled            bool
}

// StoreOption options for the vault store
type StoreOption func(*Store)

// AutoMigrate auto migrate
func (st *Store) AutoMigrate() error {
	sqls, err := st.SqlCreateTable()

	if err != nil {
		return err
	}

	for _, sql := range sqls {
		_, err := st.db.Exec(sql)
		if err != nil {
			return nil
		}
	}

	return nil
}

// EnableDebug - enables the debug option
func (st *Store) EnableDebug(debug bool) {
	st.debugEnabled = debug
}

func (st *Store) GetAttributeTableName() string {
	return st.attributeTableName
}

func (st *Store) GetAttributeTrashTableName() string {
	return st.attributeTrashTableName
}

func (st *Store) GetDB() *sql.DB {
	return st.db
}

func (st *Store) GetDebug() bool {
	return st.debugEnabled
}

func (st *Store) GetEntityTableName() string {
	return st.entityTableName
}

func (st *Store) GetEntityTrashTableName() string {
	return st.entityTrashTableName
}

func (st *Store) SqlCreateTable() ([]string, error) {

	sqlMysql1 := `
	CREATE TABLE IF NOT EXISTS ` + st.entityTableName + ` (
		id varchar(40) NOT NULL PRIMARY KEY,
		entity_type varchar(40) NOT NULL,
		entity_handle varchar(60) DEFAULT '',
		created_at datetime NOT NULL,
		updated_at datetime NOT NULL
	 );
	`

	sqlMysql2 := `
	CREATE TABLE IF NOT EXISTS ` + st.attributeTableName + ` (
		id varchar(40) NOT NULL PRIMARY KEY,
		entity_id varchar(40) NOT NULL,
		attribute_key varchar(255) NOT NULL,
		attribute_value text,
		created_at datetime NOT NULL,
		updated_at datetime NOT NULL
	);
	`

	sqlMysql3 := `
	CREATE TABLE IF NOT EXISTS ` + st.entityTrashTableName + ` (
		id varchar(40) NOT NULL PRIMARY KEY,
		entity_type varchar(40) NOT NULL,
		entity_handle varchar(60) DEFAULT '',
		created_at datetime NOT NULL,
		updated_at datetime NOT NULL,
		deleted_at datetime NOT NULL,
		deleted_by varchar(40)
	);
	`

	sqlMysql4 := `
	CREATE TABLE IF NOT EXISTS ` + st.attributeTrashTableName + ` (
		id varchar(40) NOT NULL PRIMARY KEY,
		entity_id varchar(40) NOT NULL,
		attribute_key varchar(255) NOT NULL,
		attribute_value text,
		created_at datetime NOT NULL,
		updated_at datetime NOT NULL,
		deleted_at datetime NOT NULL,
		deleted_by varchar(40)
	);
	`

	sqlPostgres1 := `
	CREATE TABLE IF NOT EXISTS ` + st.attributeTableName + ` (
		"id" varchar(40) NOT NULL PRIMARY KEY,
		"entity_id" varchar(40) NOT NULL,
		"attribute_key" varchar(255) NOT NULL,
		"attribute_value" text,
		"created_at" timestamptz(6) NOT NULL,
		"updated_at" timestamptz(6) NOT NULL
	);
	`

	sqlPostgres2 := `
	CREATE TABLE IF NOT EXISTS ` + st.entityTableName + ` (
	   "id" varchar(40) NOT NULL PRIMARY KEY,
	   "entity_type" varchar(40) NOT NULL,
	   "entity_handle" varchar(60) DEFAULT '',
	   "created_at" timestamptz(6),
	   "updated_at" timestamptz(6)
	);
	`

	sqlPostgres3 := `
	CREATE TABLE IF NOT EXISTS ` + st.entityTrashTableName + ` (
		"id" varchar(40) NOT NULL PRIMARY KEY,
		"entity_type" varchar(40) NOT NULL,
		"entity_handle" varchar(60) DEFAULT '',
		"created_at" timestamptz(6) NOT NULL,
		"updated_at" timestamptz(6) NOT NULL,
		"deleted_at" timestamptz(6) NOT NULL,
		"deleted_by" varchar(40)
	);
	`

	sqlPostgres4 := `
	CREATE TABLE IF NOT EXISTS ` + st.attributeTrashTableName + ` (
		"id" varchar(40) NOT NULL PRIMARY KEY,
		"entity_id" varchar(40) NOT NULL,
		"attribute_key" varchar(255) NOT NULL,
		"attribute_value" text,
		"created_at" timestamptz(6) NOT NULL,
		"updated_at" timestamptz(6) NOT NULL,
		"deleted_at" timestamptz(6) NOT NULL,
		"deleted_by" varchar(40)
	);
	`

	sqlSqlite1 := `
	CREATE TABLE IF NOT EXISTS "` + st.attributeTableName + `" (
		"id" varchar(40) NOT NULL PRIMARY KEY,
		"entity_id" varchar(40) NOT NULL,
		"attribute_key" varchar(255) NOT NULL,
		"attribute_value" text,
		"created_at" datetime NOT NULL,
		"updated_at" datetime NOT NULL
	);
	`
	sqlSqlite2 := `
	CREATE TABLE IF NOT EXISTS "` + st.entityTableName + `" (
	   "id" varchar(40) NOT NULL PRIMARY KEY,
	   "entity_type" varchar(40) NOT NULL,
	   "entity_handle" varchar(60) DEFAULT '',
	   "created_at" datetime NOT NULL,
	   "updated_at" datetime NOT NULL
	);
	`

	sqlSqlite3 := `
	CREATE TABLE IF NOT EXISTS "` + st.entityTrashTableName + `" (
		"id" varchar(40) NOT NULL PRIMARY KEY,
		"entity_type" varchar(40) NOT NULL,
		"entity_handle" varchar(60) DEFAULT '',
		"created_at" datetime NOT NULL,
		"updated_at" datetime NOT NULL,
		"deleted_at" datetime NOT NULL,
		"deleted_by" varchar(40)
	);
	`

	sqlSqlite4 := `
	CREATE TABLE IF NOT EXISTS "` + st.attributeTrashTableName + `" (
		"id" varchar(40) NOT NULL PRIMARY KEY,
		"entity_id" varchar(40) NOT NULL,
		"attribute_key" varchar(255) NOT NULL,
		"attribute_value" text,
		"created_at" datetime NOT NULL,
		"updated_at" datetime NOT NULL,
		"deleted_at" datetime NOT NULL,
		"deleted_by" varchar(40)
	);
	`

	sqls := []string{}

	if st.dbDriverName == "mysql" {
		sqls = append(sqls, sqlMysql1)
		sqls = append(sqls, sqlMysql2)
		sqls = append(sqls, sqlMysql3)
		sqls = append(sqls, sqlMysql4)
	} else if st.dbDriverName == "postgres" {
		sqls = append(sqls, sqlPostgres1)
		sqls = append(sqls, sqlPostgres2)
		sqls = append(sqls, sqlPostgres3)
		sqls = append(sqls, sqlPostgres4)
	} else if st.dbDriverName == "sqlite" {
		sqls = append(sqls, sqlSqlite1)
		sqls = append(sqls, sqlSqlite2)
		sqls = append(sqls, sqlSqlite3)
		sqls = append(sqls, sqlSqlite4)
	} else {
		return nil, errors.New("unsupported driver " + st.dbDriverName)
	}

	return sqls, nil
}
