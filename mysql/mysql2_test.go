package mysqlx

import (
	"log"
	"os"
	"testing"
	"time"
)

func TestMySQL(t *testing.T) {

	cfg := &MySQLConfig{
		Username:      "root",
		Password:      "yZY0G0Dzh5N",
		Address:       "10.171.5.193:3306",
		DatabaseName:  "test",
		MaxOpenConns:  64,
		MaxIdleConns:  4,
		LogMode:       "debug",
		Logger:        log.New(os.Stdout, "", log.LstdFlags),
		SlowThreshold: time.Duration(30) * time.Second,
	}

	// 获取mysql
	mysqlClient, err := NewMySQL(cfg)
	if err != nil {
		t.Errorf("Mysql connect failed, err: %v.", err)
		return
	}

	// 关闭mysql
	defer mysqlClient.Close()

	// 添加
	cr := &ContentRecord{
		Content:    "second content",
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}
	if err := cr.Create(mysqlClient); err != nil {
		t.Errorf("Mysql create err: %v.", err)
		return
	}

	// 根据id 排序，取第一条数据
	cr2 := &ContentRecord{}
	if err := cr2.First(mysqlClient); err != nil {
		t.Errorf("Mysql get first record err: %v.", err)
		return
	}
	t.Log(cr2)

	// 根据id 匹配
	cr3 := &ContentRecord{}
	if err := cr3.FirstById(mysqlClient, map[string]interface{}{"id": 10}); err != nil {
		t.Errorf("Mysql get record by id err: %v.", err)
		return
	}
	t.Log(cr3)

	// 获取至多5条数据
	contentRecords, err := cr3.Find(mysqlClient)
	if err != nil {
		t.Errorf("Mysql find at most 5 records err: %v.", err)
		return
	}
	for _, contentRecord := range contentRecords {
		t.Log(contentRecord)
	}

	// 根据content 删除匹配的数据
	cr4 := &ContentRecord{}
	if err = cr4.DeleteBy(mysqlClient, map[string]interface{}{"content": "second content"}); err != nil {
		t.Errorf("Mysql delete by content err: %v.", err)
		return
	}

}

type ContentRecord struct {
	Id         int64  `gorm:"column:id"`
	Content    string `gorm:"column:content"`
	CreateTime int64  `gorm:"column:create_time"`
	UpdateTime int64  `gorm:"column:update_time"`
}

/**
创建测试表test.content_record
CREATE TABLE test.content_record (
   `id` int NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '主表id',
   `content` varchar(64) NOT NULL COMMENT '名称',
   `create_time` int NOT NULL COMMENT '创建时间',
   `update_time` int NOT NULL COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
*/

// 定义表名
func (ContentRecord) TableName() string {
	return "content_record"
}

func (cr *ContentRecord) Create(m *MySQL) error {
	return m.GetClient().Create(cr).Error
}

func (cr *ContentRecord) First(m *MySQL) error {
	return m.GetClient().First(cr).Error
}

func (cr *ContentRecord) FirstById(m *MySQL, where map[string]interface{}) error {
	db := m.GetClient()
	if v, ok := where["id"]; ok {
		db = db.Where("id", v)
	}
	return db.First(cr).Error
}

func (cr *ContentRecord) Find(m *MySQL) ([]*ContentRecord, error) {
	var rows []*ContentRecord
	err := m.GetClient().Limit(5).Find(&rows).Error
	return rows, err
}

func (cr *ContentRecord) DeleteBy(m *MySQL, where map[string]interface{}) error {
	db := m.GetClient()
	if v, ok := where["content"]; ok {
		db = db.Where("content", v)
	}
	return db.Delete(cr).Error
}
