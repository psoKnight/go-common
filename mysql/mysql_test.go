package mysqlx

import (
	"gorm.io/gorm"
	"log"
	"os"
	"testing"
	"time"
)

func TestMySQL2(t *testing.T) {

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
	cr2 := &ContentRecord2{
		Content:    "first content",
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}
	if err := mysqlClient.Create(cr2); err != nil {
		t.Errorf("Mysql create err: %v.", err)
		return
	}

	// 根据id 排序，取第一条数据
	c, err := mysqlClient.First()
	if err != nil {
		t.Errorf("Mysql get first record err: %v.", err)
		return
	}
	t.Log(c)

	// 根据id 匹配
	c2, err := mysqlClient.FirstById(10)
	if err != nil {
		t.Errorf("Mysql get record by id err: %v.", err)
		return
	}
	t.Log(c2)

	// 获取至多5条数据
	c3s, err := mysqlClient.Find()
	if err != nil {
		t.Errorf("Mysql find at most 5 records err: %v.", err)
		return
	}
	for _, contentRecord := range c3s {
		t.Log(contentRecord)
	}

	// 根据content 删除匹配的数据
	if err = mysqlClient.DeleteById(1); err != nil {
		t.Errorf("Mysql delete by content err: %v.", err)
		return
	}
}

type ContentRecord2 struct {
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
func (ContentRecord2) TableName() string {
	return "content_record"
}

func (m *MySQL) Create(c *ContentRecord2) error {
	return m.GetClient().Create(c).Error
}

func (m *MySQL) First() (ContentRecord2, error) {
	var contentRecord2 ContentRecord2
	err := m.GetClient().Model(ContentRecord2{}).First(&contentRecord2).Error
	return contentRecord2, err
}

func (m *MySQL) FirstById(id int) (ContentRecord2, error) {
	var contentRecord2 ContentRecord2
	err := m.GetClient().Model(ContentRecord2{}).Where("id = ?", id).First(&contentRecord2).Error
	if err == gorm.ErrRecordNotFound {
		return contentRecord2, nil
	}
	return contentRecord2, err
}

func (m *MySQL) Find() ([]ContentRecord2, error) {
	var rows []ContentRecord2
	err := m.GetClient().Model(ContentRecord2{}).Limit(5).Find(&rows).Error
	return rows, err
}

func (m *MySQL) DeleteById(id int) error {
	return m.GetClient().Delete(ContentRecord2{}, "id = ?", id).Error
}
