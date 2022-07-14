package mysqlx

import (
	"log"
	"os"
	"testing"
	"time"
)

func TestMySQL(t *testing.T) {

	// 获取mysql
	cfg := &Config{
		Username:      "root",
		Password:      "nXfYAxwc-4q",
		Address:       "10.171.4.214:3306",
		DatabaseName:  "test",
		MaxOpenConns:  64,
		MaxIdleConns:  4,
		LogMode:       "debug",
		Logger:        log.New(os.Stdout, "", log.LstdFlags),
		SlowThreshold: time.Duration(30) * time.Second,
	}

	mysqlClient, errNM := NewMySQL(cfg)
	if errNM != nil {
		t.Errorf("Mysql connect failed, err: %v.", errNM)
		return
	}

	// 关闭mysql
	defer mysqlClient.Close()

	// 添加
	cr := &ContentRecord{
		Content:    "first content",
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}
	if errC := cr.Create(mysqlClient); errC != nil {
		t.Errorf("Mysql create err: %v.", errC)
		return
	}

	// 根据id 排序 取第一条数据
	cr2 := &ContentRecord{}
	if errF := cr2.First(mysqlClient); errF != nil {
		t.Errorf("Mysql get first record err: %v.", errF)
		return
	}
	t.Log(cr2)

	// 根据`id`匹配
	cr3 := &ContentRecord{}
	if errFBI := cr3.FirstById(mysqlClient, map[string]interface{}{"id": 2}); errFBI != nil {
		t.Errorf("Mysql get record by `id` err: %v.", errFBI)
		return
	}
	t.Log(cr3)

	// 获取至多5条数据
	contentRecords, errF := cr3.Find(mysqlClient)
	if errF != nil {
		t.Errorf("Mysql find at most 5 records err: %v.", errNM)
		return
	}
	for _, contentRecord := range contentRecords {
		t.Log(contentRecord)
	}

	// 根据`content`删除匹配的数据
	cr4 := &ContentRecord{}
	errDB := cr4.DeleteBy(mysqlClient, map[string]interface{}{"content": "second content"})
	if errDB != nil {
		t.Errorf("Mysql delete by `content` err: %v.", errDB)
		return
	}

}

type ContentRecord struct {
	Id         int64  `gorm:"column:id"`
	Content    string `gorm:"column:content"`
	CreateTime int64  `gorm:"column:create_time"`
	UpdateTime int64  `gorm:"column:update_time"`
}

// 定义表名
func (ContentRecord) TableName() string {
	return "content_record"
}

func (m *ContentRecord) Create(client *Mysql) error {
	return client.DB.Create(m).Error
}

func (m *ContentRecord) First(client *Mysql) error {
	return client.DB.First(m).Error
}

func (m *ContentRecord) FirstById(client *Mysql, where map[string]interface{}) error {
	db := client.DB
	if v, ok := where["id"]; ok {
		db = db.Where("id", v)
	}
	return db.First(m).Error
}

func (m *ContentRecord) Find(client *Mysql) ([]*ContentRecord, error) {
	var rows []*ContentRecord
	err := client.DB.Limit(5).Find(&rows).Error
	return rows, err
}

func (m *ContentRecord) DeleteBy(client *Mysql, where map[string]interface{}) error {
	db := client.DB
	if v, ok := where["content"]; ok {
		db = db.Where("content", v)
	}
	return db.Delete(m).Error
}
