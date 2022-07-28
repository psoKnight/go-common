package arangodb

import (
	"context"
	"errors"
	"fmt"
	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

type ArangoConfig struct {
	Endpoints []string `json:"endpoints"` // arangodb 地址
	Username  string   `json:"username"`  // 用户名是用于认证的用户名
	Password  string   `json:"password"`  // 密码是用于认证的密码
	ConnLimit int      `json:"connLimit"` // 线程池大小
}

type ArangoClient struct {
	cli driver.Client
	cfg *ArangoConfig // 基础配置
}

// NewArangoDB 获取arangodb 客户端
func NewArangoDB(cfg *ArangoConfig) (*ArangoClient, error) {

	if cfg == nil {
		return nil, errors.New("[arango]config is nil")
	} else if len(cfg.Endpoints) == 0 {
		return nil, errors.New("[arango]endpoints is 0")
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

// GetClient 获取arangodb client
func (c *ArangoClient) GetClient() driver.Client {
	return c.cli
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

// ConnectCollection 获取collection
func (c *ArangoClient) ConnectCollection(dbName string, colName string, ctx context.Context) (driver.Collection, error) {
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

// CollectionExist 判断collection 是否存在
func (c *ArangoClient) CollectionExist(dbName string, colName string, ctx context.Context) (bool, error) {
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

// CreateCollection 创建集合
func (c *ArangoClient) CreateCollection(dbName string, colName string, ctx context.Context, options *driver.CreateCollectionOptions) (driver.Collection, error) {
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

// DeleteCollection 删除集合
func (c *ArangoClient) DeleteCollection(dbName string, colName string, ctx context.Context) error {
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
	col, err := c.ConnectCollection(dbName, colName, ctx)
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
	col, err := c.ConnectCollection(dbName, colName, ctx)
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
		return driver.DocumentMeta{}, errors.New("[arango]key is '', please check")
	}

	col, err := c.ConnectCollection(dbName, colName, ctx)
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
			return driver.DocumentMetaSlice{}, driver.ErrorSlice{}, errors.New("[arango]keys exist '', please check")
		}
	}

	col, err := c.ConnectCollection(dbName, colName, ctx)
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
		return driver.DocumentMeta{}, errors.New("[arango]key is '', please check")
	}

	col, err := c.ConnectCollection(dbName, colName, ctx)
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
			return driver.DocumentMetaSlice{}, driver.ErrorSlice{}, errors.New("[arango]keys exist '', please check")
		}
	}

	col, err := c.ConnectCollection(dbName, colName, ctx)
	if err != nil {
		return driver.DocumentMetaSlice{}, driver.ErrorSlice{}, err
	}
	metas, errs, err := col.UpdateDocuments(ctx, keys, updateDocs)
	if err != nil {
		return metas, errs, err
	}
	return metas, errs, nil
}

// QueryDocuments 查询文档
func (c *ArangoClient) QueryDocuments(dbName string, query string, ctx context.Context) ([]interface{}, error) {
	if query == "" {
		return []interface{}{}, errors.New("[arango]query is '', please check")
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
		return []interface{}{}, errors.New("[arango]query is '', please check")
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

// ExecuteAQL 执行AQL
func (c *ArangoClient) ExecuteAQL(dbName string, query string, ctx context.Context) error {
	if query == "" {
		return errors.New("[arango]query is '', please check")
	}

	db, err := c.ConnectDatabase(dbName, ctx)
	if err != nil {
		return err
	}

	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		return err
	}
	defer cursor.Close() // 关闭游标

	return nil
}

// ExecuteAQLBindVariables 执行AQL(带有变量)
func (c *ArangoClient) ExecuteAQLBindVariables(dbName string, query string, bindVars map[string]interface{}, ctx context.Context) error {
	if query == "" {
		return errors.New("[arango]query is '', please check")
	}

	db, err := c.ConnectDatabase(dbName, ctx)
	if err != nil {
		return err
	}

	cursor, err := db.Query(ctx, query, bindVars)
	if err != nil {
		return err
	}
	defer cursor.Close() // 关闭游标

	return nil
}

// GetCollectionCount 获取集合总数
func (c *ArangoClient) GetCollectionCount(dbName string, colName string, ctx context.Context) (int64, error) {

	db, err := c.ConnectDatabase(dbName, ctx)
	if err != nil {
		return 0, err
	}

	col, err := db.Collection(ctx, colName)
	if err != nil {
		return 0, err
	}

	count, err := col.Count(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetQueryCount 获取query 查询总数
func (c *ArangoClient) GetQueryCount(dbName string, query string, ctx context.Context) (int64, error) {
	if query == "" {
		return 0, errors.New("[arango]query is '', please check")
	}

	db, err := c.ConnectDatabase(dbName, ctx)
	if err != nil {
		return 0, err
	}

	ctx = driver.WithQueryCount(ctx)

	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		return 0, err
	}
	defer cursor.Close() // 关闭游标

	return cursor.Count(), nil
}

// GetQueryBindVariablesCount 获取query 查询总数(带有变量)
func (c *ArangoClient) GetQueryBindVariablesCount(dbName string, query string, bindVars map[string]interface{}, ctx context.Context) (int64, error) {
	if query == "" {
		return 0, errors.New("[arango]query is '', please check")
	}

	db, err := c.ConnectDatabase(dbName, ctx)
	if err != nil {
		return 0, err
	}

	ctx = driver.WithQueryCount(ctx)

	cursor, err := db.Query(ctx, query, bindVars)
	if err != nil {
		return 0, err
	}
	defer cursor.Close() // 关闭游标

	return cursor.Count(), nil
}

// GetRelationBetweenTwoIds 获取两个ids 间的关系
/*
isAny
	true: 双向查询 id1->id2 || id2->id1
	false: outbound 查询 id->id2
*/
func (c *ArangoClient) GetRelationBetweenTwoIds(dbName, relationCol, id1, id2 string, isAny bool, ctx context.Context) ([]interface{}, error) {

	if id1 == "" || id2 == "" {
		return []interface{}{}, errors.New("[arango]id is '', please check")
	}

	bindVars := make(map[string]interface{})
	bindVars["id1"] = id1
	bindVars["id2"] = id2
	bindVars["@col"] = relationCol

	filter := ""
	if isAny {
		// 双向查询 id1->id2 || id2->id1
		filter = "FILTER (d._from == @id1 AND d._to == @id2) OR (d._from == @id2 AND d._to == @id1)"
	} else {
		// outbound 查询 id->id2
		filter = "FILTER d._from == @id1 AND d._to == @id2"
	}

	query := fmt.Sprintf("FOR d IN @@col %s RETURN d", filter)

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

// GetRelationCountBetweenTwoIds 获取两个ids 间的关系总数
/*
isAny
	true: 双向查询 id1->id2 || id2->id1
	false: outbound 查询 id->id2
*/
func (c *ArangoClient) GetRelationCountBetweenTwoIds(dbName, relationCol, id1, id2 string, isAny bool, ctx context.Context) (int64, error) {

	if id1 == "" || id2 == "" {
		return 0, errors.New("[arango]id is '', please check")
	}

	bindVars := make(map[string]interface{})
	bindVars["id1"] = id1
	bindVars["id2"] = id2
	bindVars["@col"] = relationCol

	filter := ""
	if isAny {
		// 双向查询 id1->id2 || id2->id1
		filter = "FILTER (d._from == @id1 AND d._to == @id2) OR (d._from == @id2 AND d._to == @id1)"
	} else {
		// outbound 查询 id->id2
		filter = "FILTER d._from == @id1 AND d._to == @id2"
	}

	query := fmt.Sprintf("FOR d IN @@col %s RETURN d", filter)

	db, err := c.ConnectDatabase(dbName, ctx)
	if err != nil {
		return 0, err
	}

	ctx = driver.WithQueryCount(ctx)

	cursor, err := db.Query(ctx, query, bindVars)
	if err != nil {
		return 0, err
	}
	defer cursor.Close() // 关闭游标

	return cursor.Count(), nil
}

// GetRelationById 获取两个ids 间的关系
/*
direct：路径方向
	out: 外向查询 id->
	in: 内向查询 id<-
	any: 双向查询 id-> || id<-
*/
func (c *ArangoClient) GetRelationById(dbName, relationCol, id, direct string, ctx context.Context) ([]interface{}, error) {

	if id == "" {
		return []interface{}{}, errors.New("[arango]id is '', please check")
	}

	bindVars := make(map[string]interface{})
	bindVars["id"] = id
	bindVars["@col"] = relationCol

	tag := ""
	if direct == "out" {
		// 外向查询 id->
		tag = "OUTBOUND"
	} else if direct == "in" {
		// 内向查询 id<-
		tag = "INBOUND"
	} else {
		// ANY 双向查询 id-> || id<-
		tag = "ANY"
	}

	query := fmt.Sprintf("FOR v, e IN %s @id @@col RETURN e", tag)

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

// GetRelationCountById 获取两个ids 间的关系
/*
direct：路径方向
	out: 外向查询 id->
	in: 内向查询 id<-
	any: 双向查询 id-> || id<-
*/
func (c *ArangoClient) GetRelationCountById(dbName, relationCol, id, direct string, ctx context.Context) (int64, error) {

	if id == "" {
		return 0, errors.New("[arango]id is '', please check")
	}

	bindVars := make(map[string]interface{})
	bindVars["id"] = id
	bindVars["@col"] = relationCol

	tag := ""
	if direct == "out" {
		// 外向查询 id->
		tag = "OUTBOUND"
	} else if direct == "in" {
		// 内向查询 id<-
		tag = "INBOUND"
	} else {
		// ANY 双向查询 id-> || id<-
		tag = "ANY"
	}

	query := fmt.Sprintf("FOR v, e IN %s @id @@col RETURN e", tag)

	db, err := c.ConnectDatabase(dbName, ctx)
	if err != nil {
		return 0, err
	}

	ctx = driver.WithQueryCount(ctx)

	cursor, err := db.Query(ctx, query, bindVars)
	if err != nil {
		return 0, err
	}
	defer cursor.Close() // 关闭游标

	return cursor.Count(), nil
}
