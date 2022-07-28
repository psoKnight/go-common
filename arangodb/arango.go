package arangodb

import (
	"context"
	"errors"
	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

type ArangoConfig struct {
	Endpoints []string `json:"endpoints"`              // es 地址
	Username  string   `json:"username"`               // 用户名是用于认证的用户名
	Password  string   `json:"password"`               // 密码是用于认证的密码
	ConnLimit int      `json:"connLimit" default:"32"` // 线程池大小
}

type ArangoClient struct {
	cli driver.Client
	cfg *ArangoConfig // 基础配置
}

// NewArangoClient 获取arangodb 客户端
func NewArangoClient(cfg *ArangoConfig) (*ArangoClient, error) {

	if cfg == nil {
		return nil, errors.New("Arango config is nil.")
	} else if len(cfg.Endpoints) == 0 {
		return nil, errors.New("Arango endpoints is 0.")
	}

	// 保证线程池大小，默认32
	if cfg.ConnLimit == 0 {
		cfg.ConnLimit = 32
	}

	client := &ArangoClient{
		cfg: cfg,
	}

	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: convertEndpointsToUrls(cfg.Endpoints),
		ConnLimit: cfg.ConnLimit,
	})
	if err != nil {
		return nil, err
	}
	cli, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(cfg.Username, cfg.Password),
	})
	if err != nil {
		return nil, err
	} else {
		client.cli = cli
	}

	return client, nil
}

// ConnectDatabase 获取db
func (c *ArangoClient) ConnectDatabase(dbName string, ctx context.Context) (driver.Database, error) {
	db, err := c.cli.Database(ctx, dbName)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// DatabaseExist 判断db 是否存在
func (c *ArangoClient) DatabaseExist(dbName string, ctx context.Context) (bool, error) {
	exist, err := c.cli.DatabaseExists(ctx, dbName)
	if err != nil {
		return false, err
	}
	return exist, nil
}

// CreateDatabase 创建db
func (c *ArangoClient) CreateDatabase(dbName string, ctx context.Context, options *driver.CreateDatabaseOptions) (driver.Database, error) {
	db, err := c.cli.CreateDatabase(ctx, dbName, options)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// DeleteDatabase 删除db
func (c *ArangoClient) DeleteDatabase(dbName string, ctx context.Context) error {
	db, err := c.ConnectDatabase(dbName, ctx)
	if err != nil {
		return err
	}
	err = db.Remove(ctx)
	if err != nil {
		return err
	}
	return nil
}

// ConnectColletction 获取collection
func (c *ArangoClient) ConnectColletction(dbName string, colName string, ctx context.Context) (driver.Collection, error) {
	db, err := c.cli.Database(ctx, dbName)
	if err != nil {
		return nil, err
	}
	col, err := db.Collection(ctx, colName)
	if err != nil {
		return nil, err
	}
	return col, nil
}

// ColletctionExist 判断collection 是否存在
func (c *ArangoClient) ColletctionExist(dbName string, colName string, ctx context.Context) (bool, error) {
	db, err := c.cli.Database(ctx, dbName)
	if err != nil {
		return false, err
	}
	exist, err := db.CollectionExists(ctx, colName)
	if err != nil {
		return false, err
	}
	return exist, nil
}

// CreateColletction 创建集合
func (c *ArangoClient) CreateColletction(dbName string, colName string, ctx context.Context, options *driver.CreateCollectionOptions) (driver.Collection, error) {
	db, err := c.cli.Database(ctx, dbName)
	if err != nil {
		return nil, err
	}
	col, err := db.CreateCollection(ctx, colName, options)
	if err != nil {
		return nil, err
	}
	return col, nil
}

// DeleteColletction 删除集合
func (c *ArangoClient) DeleteColletction(dbName string, colName string, ctx context.Context) error {
	db, err := c.cli.Database(ctx, dbName)
	if err != nil {
		return err
	}
	col, err := db.Collection(ctx, colName)
	if err != nil {
		return err
	}
	err = col.Remove(ctx)
	if err != nil {
		return err
	}
	return nil
}

// CreateDocument 创建文档
func (c *ArangoClient) CreateDocument(dbName string, colName string, doc interface{}, ctx context.Context) (driver.DocumentMeta, error) {
	col, err := c.ConnectColletction(dbName, colName, ctx)
	if err != nil {
		return driver.DocumentMeta{}, err
	}
	meta, err := col.CreateDocument(ctx, doc)
	if err != nil {
		return meta, err
	}
	return meta, nil
}

// CreateDocuments 创建多个文档
func (c *ArangoClient) CreateDocuments(dbName string, colName string, docs []interface{}, ctx context.Context) (driver.DocumentMetaSlice, driver.ErrorSlice, error) {
	col, err := c.ConnectColletction(dbName, colName, ctx)
	if err != nil {
		return driver.DocumentMetaSlice{}, driver.ErrorSlice{}, err
	}
	metas, errs, err := col.CreateDocuments(ctx, docs)
	if err != nil {
		return metas, errs, err
	}
	return metas, errs, nil
}

// DeleteDocument 删除文档
func (c *ArangoClient) DeleteDocument(dbName string, colName string, key string, ctx context.Context) (driver.DocumentMeta, error) {
	if key == "" {
		return driver.DocumentMeta{}, errors.New("Key is '', please check.")
	}

	col, err := c.ConnectColletction(dbName, colName, ctx)
	if err != nil {
		return driver.DocumentMeta{}, err
	}
	meta, err := col.RemoveDocument(ctx, key)
	if err != nil {
		return meta, err
	}
	return meta, nil
}

// DeleteDocuments 删除多个文档
func (c *ArangoClient) DeleteDocuments(dbName string, colName string, keys []string, ctx context.Context) (driver.DocumentMetaSlice, driver.ErrorSlice, error) {
	// 处理非空keys
	for _, key := range keys {
		if key == "" {
			return driver.DocumentMetaSlice{}, driver.ErrorSlice{}, errors.New("Keys exist '', please check.")
		}
	}

	col, err := c.ConnectColletction(dbName, colName, ctx)
	if err != nil {
		return driver.DocumentMetaSlice{}, driver.ErrorSlice{}, err
	}
	metas, errs, err := col.RemoveDocuments(ctx, keys)
	if err != nil {
		return metas, errs, err
	}
	return metas, errs, nil
}

// UpdateDocument 更新文档
func (c *ArangoClient) UpdateDocument(dbName string, colName string, key string, updateDoc interface{}, ctx context.Context) (driver.DocumentMeta, error) {
	if key == "" {
		return driver.DocumentMeta{}, errors.New("Key is '', please check.")
	}

	col, err := c.ConnectColletction(dbName, colName, ctx)
	if err != nil {
	}
	meta, err := col.UpdateDocument(ctx, key, updateDoc)
	if err != nil {
		return driver.DocumentMeta{}, err
	}
	if err != nil {
		return meta, err
	}
	return meta, nil
}

// UpdateDocuments 删除多个文档
func (c *ArangoClient) UpdateDocuments(dbName string, colName string, keys []string, updateDocs []interface{}, ctx context.Context) (driver.DocumentMetaSlice, driver.ErrorSlice, error) {
	// 处理非空keys
	for _, key := range keys {
		if key == "" {
			return driver.DocumentMetaSlice{}, driver.ErrorSlice{}, errors.New("Keys exist '', please check.")
		}
	}

	col, err := c.ConnectColletction(dbName, colName, ctx)
	if err != nil {
		return driver.DocumentMetaSlice{}, driver.ErrorSlice{}, err
	}
	metas, errs, err := col.UpdateDocuments(ctx, keys, updateDocs)
	if err != nil {
		return metas, errs, err
	}
	return metas, errs, nil
}

// 查询文档
func (c *ArangoClient) QueryDocuments(dbName string, query string, ctx context.Context) ([]interface{}, error) {
	if query == "" {
		return []interface{}{}, errors.New("Query is '', please check.")
	}

	db, err := c.ConnectDatabase(dbName, ctx)
	if err != nil {
		return []interface{}{}, err
	}

	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		return []interface{}{}, err
	}
	defer cursor.Close() // 关闭游标

	docs := make([]interface{}, 0)
	for {
		var doc interface{}
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			continue
		}

		// 去除nil 数据
		if doc != nil {
			docs = append(docs, doc)
		}
	}

	return docs, nil
}

// QueryDocumentsBindVariables 查询文档(带有变量)
func (c *ArangoClient) QueryDocumentsBindVariables(dbName string, query string, bindVars map[string]interface{}, ctx context.Context) ([]interface{}, error) {
	if query == "" {
		return []interface{}{}, errors.New("Query is '', please check.")
	}

	db, err := c.ConnectDatabase(dbName, ctx)
	if err != nil {
		return []interface{}{}, err
	}

	cursor, err := db.Query(ctx, query, bindVars)
	if err != nil {
		return []interface{}{}, err
	}
	defer cursor.Close() // 关闭游标

	docs := make([]interface{}, 0)
	for {
		var doc interface{}
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			continue
		}

		// 去除nil 数据
		if doc != nil {
			docs = append(docs, doc)
		}
	}

	return docs, nil
}
