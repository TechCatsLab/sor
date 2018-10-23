package service

import (
	"database/sql"

	"github.com/TechCatsLab/sor/category/config"

	"github.com/TechCatsLab/sor/category/model/mysql"
)

type TransactService struct {
	db   *sql.DB
	SQLS []string
}

func NewCategoryService(c *config.Config, db *sql.DB) *TransactService {
	ts := &TransactService{
		db: db,
		SQLS: []string{
			`CREATE DATABASE IF NOT EXISTS ` + c.CategoryDB,
			`CREATE TABLE IF NOT EXISTS ` + c.CategoryDB + `.` + c.CategoryTable + `(
				categoryId INT(11) NOT NULL AUTO_INCREMENT COMMENT '类别id',
				parentId INT(11) DEFAULT NULL  COMMENT '父类别id',
				name VARCHAR(50) DEFAULT NULL COMMENT '类别名称',
				status TINYINT(1) DEFAULT '1' COMMENT '状态1-在售，2-废弃',
				createTime DATETIME DEFAULT current_timestamp COMMENT '创建时间',
				PRIMARY KEY (categoryId),
				INDEX(parentId)
			)ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8`,
			`INSERT INTO ` + c.CategoryDB + `.` + c.CategoryTable + `(parentId,name) VALUES (?,?)`,
			`UPDATE ` + c.CategoryDB + `.` + c.CategoryTable + `SET status = ? WHERE categoryId = ? LIMIT 1`,
			`UPDATE ` + c.CategoryDB + `.` + c.CategoryTable + `SET name = ? WHERE categoryId = ? LIMIT 1`,
			`SELECT * FROM ` + c.CategoryDB + `.` + c.CategoryTable + ` WHERE parentId = ?`,
		},
	}
	return ts
}
func (ts *TransactService) CreateDB() error {
	return mysql.CreateDB(ts.db, ts.SQLS[0])
}

func (ts *TransactService) CreateTable() error {
	return mysql.CreateTable(ts.db, ts.SQLS[1])
}

//返回插入的编号
func (ts *TransactService) Insert(parentId uint, name string) (uint, error) {
	return mysql.InsertCategory(ts.db, ts.SQLS[2], parentId, name)
}

func (ts *TransactService) ChangeCategoryStatus(categoryId uint, status int8) error {
	return mysql.ChangeCategoryStatus(ts.db, ts.SQLS[3], categoryId, status)
}

func (ts *TransactService) ChangeCategoryName(categoryId uint, name string) error {
	return mysql.ChangeCategoryName(ts.db, ts.SQLS[4], categoryId, name)
}

//返回父级目录为parentId的目录
func (ts *TransactService) LisitChirldrenByParentId(parentId uint) ([]*mysql.Category, error) {
	categorys, err := mysql.LisitChirldrenByParentId(ts.db, ts.SQLS[5], parentId)
	if err != nil {
		return nil, err
	}

	return categorys, nil
}

var product = []string{
	`CREATE DATABASE IF NOT EXISTS mall`,
	`CREATE TABLE IF NOT EXISTS mall.product(
		productId INT(11) NOT NULL AUTO_INCREMENT COMMENT '商品id',
		categoryId INT(11) NOT NULL COMMENT '分类id，对应分类表的主键',
		name VARCHAR(100) NOT NULL COMMENT '商品名称',
		mainImage VARCHAR(500) NOT NULL COMMENT '商品图片，相对地址',
		detail TEXT COMMENT '商品详情',
		price DECIMAL(20,2) NOT NULL COMMENT '价格，单位-元保留两位小数',
		stock INT(11) NOT NULL COMMENT '库存数量',
		status TINYINT(1) NOT NULL COMMENT '状态，1-在售 2-下架 3-删除',
		createTime DATETIME DEFAULT NULL COMMENT '创建时间',
		updateTime DATATIME DEFAULT NULL COMMENT '更新时间',
		PRIMARY KEY (productId)
	)ENGINE = InnoDB AUTO_INCREMENT= 2 DEFAULT CHARSET=utf8`,
}
