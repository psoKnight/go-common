package elasticsearch

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
)

func TestElastic(t *testing.T) {
	esCfg := &ElasticConfig{
		Endpoints:           []string{"10.117.48.122:9200"},
		Username:            "",
		Password:            "",
		HealthCheckInterval: 5,
		MaxRetries:          3,
		LogCfg:              EsLog{SetGzip: true},
	}

	// 获取elasticsearch 封装client
	client, err := NewElastic(esCfg)
	if err != nil {
		t.Errorf("New es client err: %v.", err)
		return
	}

	t.Log(client)
}

func TestInsertDoc(t *testing.T) {
	esCfg := &ElasticConfig{
		Endpoints:           []string{"10.117.48.122:9200"},
		Username:            "",
		Password:            "",
		HealthCheckInterval: 5,
		MaxRetries:          3,
		LogCfg:              EsLog{SetGzip: true},
	}

	// 获取elasticsearch 封装client
	client, err := NewElastic(esCfg)
	if err != nil {
		t.Errorf("New es client err: %v.", err)
		return
	}

	// 执行ES请求需要提供一个上下文对象
	ctx := context.Background()

	// 插入doc
	insertRes, err := client.InsertDoc("megvii", "ebg_1", staff{
		Name:   "张三",
		Gender: "男",
		Age:    25,
	}, ctx)
	if err != nil {
		t.Errorf("Es insert doc err: %v.", err)
		return
	} else {
		t.Logf("Insert res: %v.", insertRes)
	}

}

func TestUpsertDoc(t *testing.T) {
	esCfg := &ElasticConfig{
		Endpoints:           []string{"10.117.48.122:9200"},
		Username:            "",
		Password:            "",
		HealthCheckInterval: 5,
		MaxRetries:          3,
		LogCfg:              EsLog{SetGzip: true},
	}

	// 获取elasticsearch 封装client
	client, err := NewElastic(esCfg)
	if err != nil {
		t.Errorf("New es client err: %v.", err)
		return
	}

	// 执行ES请求需要提供一个上下文对象
	ctx := context.Background()

	// 更改doc
	insertOrUpdateRes, err := client.UpsertDoc("megvii", "ebg_4", staff{
		Name:   "李六",
		Gender: "女",
		Age:    25,
	}, ctx)
	if err != nil {
		t.Errorf("Es insert or update doc err: %v.", err)
		return
	} else {
		t.Logf("Insert or update res: %v.", insertOrUpdateRes)
	}

}

func TestQueryDoc(t *testing.T) {
	esCfg := &ElasticConfig{
		Endpoints:           []string{"10.117.48.122:9200"},
		Username:            "",
		Password:            "",
		HealthCheckInterval: 5,
		MaxRetries:          3,
		LogCfg:              EsLog{SetGzip: true},
	}

	// 获取elasticsearch 封装client
	client, err := NewElastic(esCfg)
	if err != nil {
		t.Errorf("New es client err: %v.", err)
		return
	}

	// 执行ES请求需要提供一个上下文对象
	ctx := context.Background()

	// 查询doc
	queryRes, err := client.QueryDoc("megvii", "ebg_1", ctx)
	if err != nil {
		t.Errorf("Es query doc err: %v.", err)
		return
	} else {
		msg := staff{}
		// 提取文档内容，原始类型是json 数据
		data, _ := queryRes.Source.MarshalJSON()
		// 将json 转成struct 结果
		json.Unmarshal(data, &msg)
		// 打印结果
		t.Logf("Query res: %v.", msg)
	}

}

func TestBatchUpdateDocs(t *testing.T) {
	esCfg := &ElasticConfig{
		Endpoints:           []string{"10.117.48.122:9200"},
		Username:            "",
		Password:            "",
		HealthCheckInterval: 5,
		MaxRetries:          3,
		LogCfg:              EsLog{SetGzip: true},
	}

	// 获取elasticsearch 封装client
	client, err := NewElastic(esCfg)
	if err != nil {
		t.Errorf("New es client err: %v.", err)
		return
	}

	// 执行ES请求需要提供一个上下文对象
	ctx := context.Background()

	// 批量查询doc
	batchQueryRes, err := client.BatchQueryDocs("megvii", []string{"ebg_1", "ebg_2", "ebg_3"}, ctx)
	if err != nil {
		t.Errorf("Es batch query doc err: %v.", err)
		return
	} else {
		// 遍历文档
		for i, doc := range batchQueryRes.Docs {
			// 转换成struct对象
			msg := staff{}
			tmp, _ := doc.Source.MarshalJSON()
			err := json.Unmarshal(tmp, &msg)
			if err != nil {
				t.Errorf("Batch query %d res: %v.", i, msg)
			} else {
				t.Logf("Batch query %d res: %v.", i, msg)
			}
		}
	}

}

func TestUpdateDoc(t *testing.T) {
	esCfg := &ElasticConfig{
		Endpoints:           []string{"10.117.48.122:9200"},
		Username:            "",
		Password:            "",
		HealthCheckInterval: 5,
		MaxRetries:          3,
		LogCfg:              EsLog{SetGzip: true},
	}

	// 获取elasticsearch 封装client
	client, err := NewElastic(esCfg)
	if err != nil {
		t.Errorf("New es client err: %v.", err)
		return
	}

	// 执行ES请求需要提供一个上下文对象
	ctx := context.Background()

	// 更新doc
	updateRes, err := client.UpdateDoc("megvii", "ebg_1", staff{
		Name:   "李四",
		Gender: "",
		Age:    0,
	}, ctx)
	if err != nil {
		t.Errorf("Es update doc err: %v.", err)
		return
	} else {
		t.Logf("Update res: %v.", updateRes)
	}
}

func TestDeleteDoc(t *testing.T) {
	esCfg := &ElasticConfig{
		Endpoints:           []string{"10.117.48.122:9200"},
		Username:            "",
		Password:            "",
		HealthCheckInterval: 5,
		MaxRetries:          3,
		LogCfg:              EsLog{SetGzip: true},
	}

	// 获取elasticsearch 封装client
	client, err := NewElastic(esCfg)
	if err != nil {
		t.Errorf("New es client err: %v.", err)
		return
	}

	// 执行ES请求需要提供一个上下文对象
	ctx := context.Background()

	// 删除doc
	deleteRes, err := client.DeleteDoc("megvii", "ebg_1", ctx)
	if err != nil {
		t.Errorf("Es delete doc err: %v.", err)
		return
	} else {
		t.Logf("Delete res: %v.", deleteRes)
	}

}

func TestDeleteIndexes(t *testing.T) {
	esCfg := &ElasticConfig{
		Endpoints:           []string{"10.117.48.122:9200"},
		Username:            "",
		Password:            "",
		HealthCheckInterval: 5,
		MaxRetries:          3,
		LogCfg:              EsLog{SetGzip: true},
	}

	// 获取elasticsearch 封装client
	client, err := NewElastic(esCfg)
	if err != nil {
		t.Errorf("New es client err: %v.", err)
		return
	}

	// 执行ES请求需要提供一个上下文对象
	ctx := context.Background()

	// 删除多个索引
	deleteIndexesRes, err := client.DeleteIndexes([]string{"megvii"}, ctx)
	if err != nil {
		t.Errorf("Es delete indexes err: %v.", err)
		return
	} else {
		t.Logf("Delete indexes res: %v.", deleteIndexesRes)
	}
}

func TestQueryDocsByConditions(t *testing.T) {
	esCfg := &ElasticConfig{
		Endpoints:           []string{"10.117.48.122:9200"},
		Username:            "",
		Password:            "",
		HealthCheckInterval: 5,
		MaxRetries:          3,
		LogCfg:              EsLog{SetGzip: true},
	}

	// 获取elasticsearch 封装client
	client, err := NewElastic(esCfg)
	if err != nil {
		t.Errorf("New es client err: %v.", err)
		return
	}

	// 执行ES请求需要提供一个上下文对象
	ctx := context.Background()

	// 查询doc（存在筛选条件）
	queryDocByConditionRes, err := client.QueryDocsByConditions("megvii", QueryDocByCondsParam{
		TermConditions:    nil,                                                                         // TermConditions:  map[string]interface{}{"name": "张三"}
		TermsConditions:   map[string][]interface{}{"gender.keyword": {"男", "女"}, "age": {20, 18, 25}}, // TermConditions:  map[string]interface{}{"name": "张三"}
		MustNotConditions: nil,                                                                         // map[string][]interface{}{"gender.keyword": {"男"}}
		RangeConditions:   map[string][]interface{}{"age": {18, 29}},                                   // map[string][]interface{}{"age": {18, 29}}
		ShouldConditions:  map[string][]interface{}{"age": {18, 25}},                                   // map[string][]interface{}{"age": {18, 25}}
	}, "age", false, 0, 10, ctx)
	if err != nil {
		t.Errorf("Es query doc by condition err: %v.", err)
		return
	} else {
		t.Logf("Query by condition res: %v.", queryDocByConditionRes)

		if queryDocByConditionRes.TotalHits() > 0 {
			// 查询结果不为空，则遍历结果
			var staffs staff
			// 通过Each方法，将es 结果的json 结构转换成struct 对象
			for i, item := range queryDocByConditionRes.Each(reflect.TypeOf(staffs)) {
				if staff, ok := item.(staff); ok {
					t.Logf("staff %d: %v.", i, staff)
				}
			}
		}
	}
}

func TestUpdateDocsByConditions(t *testing.T) {
	esCfg := &ElasticConfig{
		Endpoints:           []string{"10.117.48.122:9200"},
		Username:            "",
		Password:            "",
		HealthCheckInterval: 5,
		MaxRetries:          3,
		LogCfg:              EsLog{SetGzip: true},
	}

	// 获取elasticsearch 封装client
	client, err := NewElastic(esCfg)
	if err != nil {
		t.Errorf("New es client err: %v.", err)
		return
	}

	// 执行ES请求需要提供一个上下文对象
	ctx := context.Background()

	// 更新doc（存在筛选条件）
	updateDocByConditionRes, err := client.UpdateDocsByConditions("megvii", []string{"ctx._source['temp']='测试7'"}, UpdateDocByCondsParam{
		TermConditions:    map[string]interface{}{"name.keyword": "张三"},           // map[string]interface{}{"name": "张三"}
		TermsConditions:   map[string][]interface{}{"gender.keyword": {"男", "女"}}, // map[string][]interface{}{"gender":{"男","女"}}
		RangeConditions:   map[string][]interface{}{"age": {"", 50}},              // map[string][]interface{}{"age": {"",29}}
		MustNotConditions: map[string][]interface{}{"gender.keyword": {"男"}},      // map[string][]interface{}{"gender.keyword": {"男"}}
	}, ctx)
	if err != nil {
		t.Errorf("Es update doc by condition err: %v.", err)
		return
	} else {
		t.Logf("Update by condition res: %v.", updateDocByConditionRes)
	}
}

func TestDeleteDocsByConditions(t *testing.T) {
	esCfg := &ElasticConfig{
		Endpoints:           []string{"10.117.48.122:9200"},
		Username:            "",
		Password:            "",
		HealthCheckInterval: 5,
		MaxRetries:          3,
		LogCfg:              EsLog{SetGzip: true},
	}

	// 获取elasticsearch 封装client
	client, err := NewElastic(esCfg)
	if err != nil {
		t.Errorf("New es client err: %v.", err)
		return
	}

	// 执行ES请求需要提供一个上下文对象
	ctx := context.Background()
	// 删除doc（存在筛选条件）
	deleteDocByConditionRes, err := client.DeleteDocsByConditions("megvii", DeleteDocByCondsParam{
		TermConditions:    map[string]interface{}{"name.keyword": "王五"},           // map[string]interface{}{"name": "张三"}
		TermsConditions:   map[string][]interface{}{"gender.keyword": {"男", "女"}}, // map[string][]interface{}{"gender":{"男","女"}}
		RangeConditions:   map[string][]interface{}{"age": {"", 30}},              // map[string][]interface{}{"age": {"",29}}
		MustNotConditions: map[string][]interface{}{"gender.keyword": {"男"}},      // map[string][]interface{}{"gender.keyword": {"男"}}
	}, ctx)
	if err != nil {
		t.Errorf("Es delete doc by condition err: %v.", err)
		return
	} else {
		t.Logf("Delete by condition res: %v.", deleteDocByConditionRes)
	}
}

func TestBatchInsertDocs(t *testing.T) {
	esCfg := &ElasticConfig{
		Endpoints:           []string{"10.117.48.122:9200"},
		Username:            "",
		Password:            "",
		HealthCheckInterval: 5,
		MaxRetries:          3,
		LogCfg:              EsLog{SetGzip: true},
	}

	// 获取elasticsearch 封装client
	client, err := NewElastic(esCfg)
	if err != nil {
		t.Errorf("New es client err: %v.", err)
		return
	}

	// 执行ES请求需要提供一个上下文对象
	ctx := context.Background()
	// BatchInsertDoc 批量插入doc
	BatchInsertDocRes, err := client.BatchInsertDocs("megvii", []string{}, []interface{}{staff{
		Name:   "宋一",
		Gender: "女",
		Age:    21,
	}, staff{
		Name:   "宋九",
		Gender: "男",
		Age:    23,
	}}, ctx)
	if err != nil {
		t.Errorf("Es batch insert doc err: %v.", err)
		return
	} else {
		t.Logf("Batch insert doc res: %v.", BatchInsertDocRes)
	}
}

func TestBatchUpsertDocs(t *testing.T) {
	esCfg := &ElasticConfig{
		Endpoints:           []string{"10.117.48.122:9200"},
		Username:            "",
		Password:            "",
		HealthCheckInterval: 5,
		MaxRetries:          3,
		LogCfg:              EsLog{SetGzip: true},
	}

	// 获取elasticsearch 封装client
	client, err := NewElastic(esCfg)
	if err != nil {
		t.Errorf("New es client err: %v.", err)
		return
	}

	// 执行ES请求需要提供一个上下文对象
	ctx := context.Background()

	// BatchUpsertDocs 批量插入doc
	batchUpsertDocsRes, err := client.BatchUpsertDocs("megvii", []string{"ebg_11", "ebg_21"}, []interface{}{staff{
		Name:   "张三",
		Gender: "男",
		Age:    21,
	}, staff{
		Name:   "李四",
		Gender: "男",
		Age:    23,
	}}, ctx)
	if err != nil {
		t.Errorf("Es batch upsert doc err: %v.", err)
		return
	} else {
		t.Logf("Batch upsert doc res: %v.", batchUpsertDocsRes)
	}
}

func TestDeleteIndex(t *testing.T) {
	esCfg := &ElasticConfig{
		Endpoints:           []string{"10.117.48.122:9200"},
		Username:            "",
		Password:            "",
		HealthCheckInterval: 5,
		MaxRetries:          3,
		LogCfg:              EsLog{SetGzip: true},
	}

	// 获取elasticsearch 封装client
	client, err := NewElastic(esCfg)
	if err != nil {
		t.Errorf("New es client err: %v.", err)
		return
	}

	// 执行ES请求需要提供一个上下文对象
	ctx := context.Background()

	// 删除单个索引
	deleteIndexRes, err := client.DeleteIndex("empty", ctx)
	if err != nil {
		t.Errorf("Es delete index err: %v.", err)
		return
	} else {
		t.Logf("Delete index res: %v.", deleteIndexRes)
	}
}

func TestExistIndex(t *testing.T) {
	esCfg := &ElasticConfig{
		Endpoints:           []string{"10.117.48.122:9200"},
		Username:            "",
		Password:            "",
		HealthCheckInterval: 5,
		MaxRetries:          3,
		LogCfg:              EsLog{SetGzip: true},
	}

	// 获取elasticsearch 封装client
	client, err := NewElastic(esCfg)
	if err != nil {
		t.Errorf("New es client err: %v.", err)
		return
	}

	// 执行ES请求需要提供一个上下文对象
	ctx := context.Background()

	// 判断索引是否存在
	ExistIndexRes, err := client.ExistIndex("megvii", ctx)
	if err != nil {
		t.Errorf("Es exist index err: %v.", err)
		return
	} else {
		t.Logf("Exist index res: %v.", ExistIndexRes)
	}
}

type staff struct {
	Name   string `json:"name,omitempty"`
	Gender string `json:"gender,omitempty"`
	Age    int    `json:"age,omitempty"`
}
