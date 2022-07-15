package elasticsearch

import (
	"context"
	"errors"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"strings"
	"time"
)

type ElasticConfig struct {
	Endpoints           []string `json:"endpoints"`             // es 地址
	Username            string   `json:"username"`              // 用户名是用于认证的用户名
	Password            string   `json:"password"`              // 密码是用于认证的密码
	HealthCheckInterval int64    `json:"health_check_interval"` // 健康检查间隔，单位：秒（s）
	MaxRetries          int      `json:"max_retries"`           // 设置请求失败最大重试次数
	LogCfg              EsLog    `json:"log_cfg"`               // 日志相关配置
}

type EsLog struct {
	SetGzip bool `json:"set_gzip"` // 是否启用gzip 压缩
	//SetErrorLog bool `json:"set_error_log"` // 设置error 日志输出
	//SetInfoLog  bool `json:"set_info_log"`  // 设置info 日志输出
}

type EsClient struct {
	cli *elastic.Client
	cfg *ElasticConfig // 基础配置
}

func NewEsClient(cfg *ElasticConfig) (*EsClient, error) {

	if cfg == nil {
		return nil, errors.New("Es config is nil.")
	} else if len(cfg.Endpoints) == 0 {
		return nil, errors.New("Es endpoints is 0.")
	}

	client := &EsClient{
		cfg: cfg,
	}

	cli, err := elastic.NewClient(
		elastic.SetURL(convertEndpointsToUrls(cfg.Endpoints)...),
		elastic.SetBasicAuth(cfg.Username, cfg.Password),
		elastic.SetGzip(cfg.LogCfg.SetGzip),
		elastic.SetHealthcheckInterval(time.Duration(cfg.HealthCheckInterval)*time.Second),
		elastic.SetMaxRetries(cfg.MaxRetries),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC_ERROR ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "ELASTIC_INFO", log.LstdFlags)))
	if err != nil {
		logrus.Errorf("[elastic]New es client err: %v.", err)
		return nil, err
	} else {
		client.cli = cli
	}

	return client, nil
}

// InsertDoc 插入doc
func (c *EsClient) InsertDoc(indexName, docId string, doc interface{}, ctx context.Context) (*elastic.IndexResponse, error) {

	insertRes, err := c.cli.Index().
		Index(indexName). // 设置索引名称
		Id(docId).        // 设置文档id
		BodyJson(doc).    // 指定前面声明struct 对象
		Do(ctx)           // 执行请求，需要传入一个上下文对象
	if err != nil {
		return nil, err
	}

	return insertRes, nil
}

// BatchInsertDoc 批量插入doc
func (c *EsClient) BatchInsertDocs(indexName string, docIds []string, docs []interface{}, ctx context.Context) (*elastic.BulkResponse, error) {

	if len(docIds) != len(docs) && len(docIds) != 0 {
		return nil, errors.New("Len docIds and docs is not equal.")
	} else if len(docIds) == 0 { // 默认为""
		docIds = make([]string, len(docs), len(docs))
	}

	bulkReq := c.cli.Bulk()
	for i, doc := range docs {
		req := elastic.NewBulkIndexRequest().Index(indexName).Id(docIds[i]).Doc(doc)
		bulkReq = bulkReq.Add(req)
	}
	bulkRes, err := bulkReq.Do(ctx)
	if err != nil {
		return nil, err
	}

	return bulkRes, nil
}

// DeleteDoc 删除doc
func (c *EsClient) DeleteDoc(indexName, docId string, ctx context.Context) (*elastic.DeleteResponse, error) {

	// 根据id 删除一条数据
	deleteRes, err := c.cli.Delete().
		Index(indexName). // 索引名
		Id(docId).        // 文档id
		Do(ctx)           // 执行请求，需要传入一个上下文对象
	if err != nil {
		return nil, err
	}

	return deleteRes, nil
}

// DeleteDocByCondition 删除doc（存在筛选条件）
// DeleteDocByCondsParam TermConditions/TermsConditions/RangeConditions/MustNotConditions/ShouldConditions 条件间是AND 的关系，各个筛选间是AND 的关系
// 待筛选的key 必须是原子类型，比如keyword/int等格式，text 格式可在分，所以不建议检索（分词可能不精准包括）
func (c *EsClient) DeleteDocsByCondition(indexName string, condition DeleteDocByCondsParam, ctx context.Context) (*elastic.BulkIndexByScrollResponse, error) {

	// 创建bool 查询
	boolQuery := elastic.NewBoolQuery().Must()

	// term
	terms := make([]elastic.Query, 0)
	if len(condition.TermConditions) > 0 {
		for termConditionK, termConditionV := range condition.TermConditions {
			termQuery := elastic.NewTermQuery(termConditionK, termConditionV)
			terms = append(terms, termQuery)
		}

		boolQuery.Must(terms...)
	}

	// terms
	termses := make([]elastic.Query, 0)
	if len(condition.TermsConditions) > 0 {
		for termsConditionK, termsConditionV := range condition.TermsConditions {
			termsQuery := elastic.NewTermsQuery(termsConditionK, termsConditionV...)
			termses = append(termses, termsQuery)
		}

		boolQuery.Must(termses...)
	}

	// must not
	nots := make([]elastic.Query, 0)
	if len(condition.MustNotConditions) > 0 {
		for mustNotConditionK, mustNotConditionV := range condition.MustNotConditions {
			mustNotQuery := elastic.NewTermsQuery(mustNotConditionK, mustNotConditionV...)
			nots = append(nots, mustNotQuery)
		}

		boolQuery.MustNot(nots...)
	}

	// range 筛选
	ranges := make([]elastic.Query, 0)
	if len(condition.RangeConditions) > 0 {
		for rangeConditionK, rangeConditionV := range condition.RangeConditions {
			if len(rangeConditionV) != 2 {
				return nil, errors.New("RangeConditions length is not 2.")
			} else {
				be := rangeConditionV[0] // 大于等于
				af := rangeConditionV[1] // 小于等于

				var rangeQuery *elastic.RangeQuery
				if (be == "" || be == 0 || be == nil) && af != "" {
					rangeQuery = elastic.NewRangeQuery(rangeConditionK).Lte(af)
				} else if be != "" && (af == "" || af == 0 || af == nil) {
					rangeQuery = elastic.NewRangeQuery(rangeConditionK).Gte(be)
				} else {
					rangeQuery = elastic.NewRangeQuery(rangeConditionK).Gte(be).Lte(af)
				}

				ranges = append(ranges, rangeQuery)
			}
		}

		boolQuery.Must(ranges...)
	}

	// should
	shoulds := make([]elastic.Query, 0)
	if len(condition.ShouldConditions) > 0 {
		for shouldConditionK, shouldConditionV := range condition.ShouldConditions {
			shouldQuery := elastic.NewTermsQuery(shouldConditionK, shouldConditionV...)
			shoulds = append(shoulds, shouldQuery)
		}

		boolQuery.Should(shoulds...).MinimumNumberShouldMatch(1)
	}

	deleteByQueryRes, err := c.cli.DeleteByQuery(indexName).Query(boolQuery).ProceedOnVersionConflict().Do(ctx)
	if err != nil {
		return nil, err
	}

	return deleteByQueryRes, nil
}

// UpdateDoc 更新doc
func (c *EsClient) UpdateDoc(indexName, docId string, doc interface{}, ctx context.Context) (*elastic.UpdateResponse, error) {
	updateRes, err := c.cli.Update().
		Index(indexName). // 设置索引名
		Id(docId).        // 文档id
		Doc(doc).         // 待修改部分，Map KV 格式或者JSON 格式
		Do(ctx)           // 执行请求，需要传入一个上下文对象
	if err != nil {
		return nil, err
	}

	return updateRes, nil
}

// BatchUpdateDoc 批量更新doc
func (c *EsClient) BatchUpdateDocs(indexName string, docIds []string, docs []interface{}, ctx context.Context) (*elastic.BulkResponse, error) {
	if len(docIds) == 0 {
		return nil, errors.New("Len docIds is 0.")
	} else if len(docIds) != len(docs) {
		return nil, errors.New("Len docIds and docs is not equal.")
	}

	bulkReq := c.cli.Bulk()
	for i, doc := range docs {
		req := elastic.NewBulkUpdateRequest().Index(indexName).Id(docIds[i]).Doc(doc)
		bulkReq = bulkReq.Add(req)
	}
	bulkRes, err := bulkReq.Do(ctx)
	if err != nil {
		return nil, err
	}

	return bulkRes, nil
}

// UpdateDocByCondition 更新doc（存在筛选条件）
// UpdateDocByCondsParam TermConditions/TermsConditions/RangeConditions/MustNotConditions/ShouldConditions 条件间是AND 的关系，各个筛选间是AND 的关系
// 待筛选的key 必须是原子类型，比如keyword/int等格式，text 格式可在分，所以不建议检索（分词可能不精准包括）
func (c *EsClient) UpdateDocsByCondition(indexName string, scripts []string, condition UpdateDocByCondsParam, ctx context.Context) (*elastic.BulkIndexByScrollResponse, error) {
	if len(scripts) == 0 {
		return nil, errors.New("Scripts is empty.")
	}

	// 创建bool 查询
	boolQuery := elastic.NewBoolQuery().Must()

	// term
	terms := make([]elastic.Query, 0)
	if len(condition.TermConditions) > 0 {
		for termConditionK, termConditionV := range condition.TermConditions {
			termQuery := elastic.NewTermQuery(termConditionK, termConditionV)
			terms = append(terms, termQuery)
		}

		boolQuery.Must(terms...)
	}

	// terms
	termses := make([]elastic.Query, 0)
	if len(condition.TermsConditions) > 0 {
		for termsConditionK, termsConditionV := range condition.TermsConditions {
			termsQuery := elastic.NewTermsQuery(termsConditionK, termsConditionV...)
			termses = append(termses, termsQuery)
		}

		boolQuery.Must(termses...)
	}

	// must not
	nots := make([]elastic.Query, 0)
	if len(condition.MustNotConditions) > 0 {
		for mustNotConditionK, mustNotConditionV := range condition.MustNotConditions {
			mustNotQuery := elastic.NewTermsQuery(mustNotConditionK, mustNotConditionV...)
			nots = append(nots, mustNotQuery)
		}

		boolQuery.MustNot(nots...)
	}

	// range 筛选
	ranges := make([]elastic.Query, 0)
	if len(condition.RangeConditions) > 0 {
		for rangeConditionK, rangeConditionV := range condition.RangeConditions {
			if len(rangeConditionV) != 2 {
				return nil, errors.New("RangeConditions length is not 2.")
			} else {
				be := rangeConditionV[0] // 大于等于
				af := rangeConditionV[1] // 小于等于

				var rangeQuery *elastic.RangeQuery
				if (be == "" || be == 0 || be == nil) && af != "" {
					rangeQuery = elastic.NewRangeQuery(rangeConditionK).Lte(af)
				} else if be != "" && (af == "" || af == 0 || af == nil) {
					rangeQuery = elastic.NewRangeQuery(rangeConditionK).Gte(be)
				} else {
					rangeQuery = elastic.NewRangeQuery(rangeConditionK).Gte(be).Lte(af)
				}

				ranges = append(ranges, rangeQuery)
			}
		}

		boolQuery.Must(ranges...)
	}

	// should
	shoulds := make([]elastic.Query, 0)
	if len(condition.ShouldConditions) > 0 {
		for shouldConditionK, shouldConditionV := range condition.ShouldConditions {
			shouldQuery := elastic.NewTermsQuery(shouldConditionK, shouldConditionV...)
			shoulds = append(shoulds, shouldQuery)
		}

		boolQuery.Should(shoulds...).MinimumNumberShouldMatch(1)
	}

	script := strings.Join(scripts, ";")
	updateByQueryRes, err := c.cli.UpdateByQuery(indexName).Query(boolQuery).Script(elastic.NewScript(script)).ProceedOnVersionConflict().Do(ctx)
	if err != nil {
		return nil, err
	}

	return updateByQueryRes, nil
}

// UpsertDoc 更改（更新或插入）doc
func (c *EsClient) UpsertDoc(indexName, docId string, doc interface{}, ctx context.Context) (interface{}, error) {
	if docId == "" {
		// 不存在，执行插入操作
		return c.InsertDoc(indexName, docId, doc, ctx)
	}

	exist, err := c.cli.Exists().
		Index(indexName). // 设置索引名
		Id(docId).        // 文档id
		Do(ctx)           // 执行请求，需要传入一个上下文对象
	if err != nil {
		return nil, err
	}

	if exist {
		// 存在，执行更新操作
		return c.UpdateDoc(indexName, docId, doc, ctx)
	} else {
		// 不存在，执行插入操作
		return c.InsertDoc(indexName, docId, doc, ctx)
	}

}

// BatchUpsertDoc 批量更改（更新或插入）doc
func (c *EsClient) BatchUpsertDocs(indexName string, docIds []string, docs []interface{}, ctx context.Context) (*elastic.BulkResponse, error) {

	if len(docIds) == 0 {
		return nil, errors.New("Len docIds is 0.")
	} else if len(docIds) != len(docs) {
		return nil, errors.New("Len docIds and docs is not equal.")
	}

	bulkReq := c.cli.Bulk()
	for i, doc := range docs {
		req := elastic.NewBulkUpdateRequest().Index(indexName).Id(docIds[i]).Doc(doc).DocAsUpsert(true)
		bulkReq = bulkReq.Add(req)
	}
	bulkRes, err := bulkReq.Do(ctx)
	if err != nil {
		return nil, err
	}

	return bulkRes, nil
}

// QueryDoc 查询doc
func (c *EsClient) QueryDoc(indexName, docId string, ctx context.Context) (*elastic.GetResult, error) {
	// 根据id 查询文档
	queryRes, err := c.cli.Get().
		Index(indexName). // 指定索引名
		Id(docId).        // 设置文档id
		Do(ctx)           // 执行请求
	if err != nil {
		return nil, err
	}
	return queryRes, nil
}

// UpdateDocByCondition 查询doc（存在筛选条件）
// QueryDocByCondsParam TermConditions/TermsConditions/RangeConditions/MustNotConditions/ShouldConditions 条件间是AND 的关系，各个筛选间是AND 的关系
// 待筛选的key 必须是原子类型，比如keyword/int等格式，text 格式可在分，所以不建议检索（分词可能不精准包括）
func (c *EsClient) QueryDocsByCondition(indexName string, condition QueryDocByCondsParam, sortField string, sortByASC bool, from, size int, ctx context.Context) (*elastic.SearchResult, error) {
	// 创建bool 查询
	boolQuery := elastic.NewBoolQuery().Must()

	// term
	terms := make([]elastic.Query, 0)
	if len(condition.TermConditions) > 0 {
		for termConditionK, termConditionV := range condition.TermConditions {
			termQuery := elastic.NewTermQuery(termConditionK, termConditionV)
			terms = append(terms, termQuery)
		}

		boolQuery.Must(terms...)
	}

	// terms
	termses := make([]elastic.Query, 0)
	if len(condition.TermsConditions) > 0 {
		for termsConditionK, termsConditionV := range condition.TermsConditions {
			termsQuery := elastic.NewTermsQuery(termsConditionK, termsConditionV...)
			termses = append(termses, termsQuery)
		}

		boolQuery.Must(termses...)
	}

	// must not
	nots := make([]elastic.Query, 0)
	if len(condition.MustNotConditions) > 0 {
		for mustNotConditionK, mustNotConditionV := range condition.MustNotConditions {
			mustNotQuery := elastic.NewTermsQuery(mustNotConditionK, mustNotConditionV...)
			nots = append(nots, mustNotQuery)
		}

		boolQuery.MustNot(nots...)
	}

	// range 筛选
	ranges := make([]elastic.Query, 0)
	if len(condition.RangeConditions) > 0 {
		for rangeConditionK, rangeConditionV := range condition.RangeConditions {
			if len(rangeConditionV) != 2 {
				return nil, errors.New("RangeConditions length is not 2.")
			} else {
				be := rangeConditionV[0] // 大于等于
				af := rangeConditionV[1] // 小于等于

				var rangeQuery *elastic.RangeQuery
				if (be == "" || be == 0 || be == nil) && af != "" {
					rangeQuery = elastic.NewRangeQuery(rangeConditionK).Lte(af)
				} else if be != "" && (af == "" || af == 0 || af == nil) {
					rangeQuery = elastic.NewRangeQuery(rangeConditionK).Gte(be)
				} else {
					rangeQuery = elastic.NewRangeQuery(rangeConditionK).Gte(be).Lte(af)
				}

				ranges = append(ranges, rangeQuery)
			}
		}

		boolQuery.Must(ranges...)
	}

	// should
	shoulds := make([]elastic.Query, 0)
	if len(condition.ShouldConditions) > 0 {
		for shouldConditionK, shouldConditionV := range condition.ShouldConditions {
			shouldQuery := elastic.NewTermsQuery(shouldConditionK, shouldConditionV...)
			shoulds = append(shoulds, shouldQuery)
		}

		boolQuery.Should(shoulds...).MinimumNumberShouldMatch(1)
	}

	query := c.cli.Search().Index(indexName).Query(boolQuery)
	if sortField != "" {
		query.Sort(sortField, sortByASC)
	}
	searchRes, err := query.From(from).Size(size).Do(ctx)
	if err != nil {
		return nil, err
	}

	return searchRes, nil
}

// BatchQueryDoc 批量查询doc
func (c *EsClient) BatchQueryDocs(indexName string, docIds []string, ctx context.Context) (*elastic.MgetResponse, error) {
	if len(docIds) == 0 {
		return nil, errors.New("Len docIds is 0.")
	}
	multiGet := c.cli.MultiGet() // 通过NewMultiGetItem 配置查询条件
	for _, id := range docIds {
		multiGet.Add(elastic.NewMultiGetItem().Index(indexName).Id(id))
	}
	multiGetRes, err := multiGet.Do(ctx)
	if err != nil {
		return nil, err
	}
	return multiGetRes, nil
}

// DeleteIndexes 删除多个索引
func (c *EsClient) DeleteIndexes(indexNames []string, ctx context.Context) (*elastic.IndicesDeleteResponse, error) {
	deleteIndexRes, err := c.cli.DeleteIndex(indexNames...).Do(ctx) // 执行请求，需要传入一个上下文对象
	if err != nil {
		return nil, err
	}

	return deleteIndexRes, nil
}

// DeleteIndex 删除单个索引
func (c *EsClient) DeleteIndex(indexName string, ctx context.Context) (*elastic.IndicesDeleteResponse, error) {
	deleteIndexRes, err := c.cli.DeleteIndex(indexName).Do(ctx) // 执行请求，需要传入一个上下文对象
	if err != nil {
		return nil, err
	}

	return deleteIndexRes, nil
}

// ExistIndex 索引是否存在
func (c *EsClient) ExistIndex(indexName string, ctx context.Context) (bool, error) {
	ExistIndexRes, err := c.cli.IndexExists(indexName).Do(ctx)
	if err != nil {
		return false, err
	}

	return ExistIndexRes, nil
}
