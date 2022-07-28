package arangodb

import (
	"context"
	"testing"
)

func TestArango(t *testing.T) {
	arangoCfg := &ArangoConfig{
		Endpoints: []string{"127.0.0.1:8529"},
		Username:  "root",
		Password:  "root",
		ConnLimit: 64,
	}

	// 获取arangodb 封装client
	client, err := NewArangoClient(arangoCfg)
	if err != nil {
		t.Errorf("New arango client err: %v.", err)
		return
	}
	ctx := context.Background()

	database, err := client.ConnectDatabase("megvii", ctx)
	if err != nil {
		t.Errorf("Conn database err: %v.", err)
		return
	}

	col, err := client.ConnectColletction("megvii", "staff", ctx)
	if err != nil {
		t.Errorf("Conn database err: %v.", err)
		return
	}

	t.Log(database, col)
	t.Logf("Conn limit: %d.", client.cfg.ConnLimit)
}

func TestDatabaseExist(t *testing.T) {

	arangoCfg := &ArangoConfig{
		Endpoints: []string{"127.0.0.1:8529"},
		Username:  "root",
		Password:  "root",
	}

	// 获取arangodb 封装client
	client, err := NewArangoClient(arangoCfg)
	if err != nil {
		t.Errorf("New arango client err: %v.", err)
		return
	}

	ctx := context.Background()
	dbExist, err := client.DatabaseExist("megvii", ctx)
	if err != nil {
		t.Errorf("Database exist err: %v.", err)
		return
	}
	t.Logf("Databse exist: %t.", dbExist)
}

func TestColletctionExist(t *testing.T) {

	arangoCfg := &ArangoConfig{
		Endpoints: []string{"127.0.0.1:8529"},
		Username:  "root",
		Password:  "root",
	}

	// 获取arangodb 封装client
	client, err := NewArangoClient(arangoCfg)
	if err != nil {
		t.Errorf("New arango client err: %v.", err)
		return
	}

	ctx := context.Background()
	colExist, err := client.ColletctionExist("megvii", "staff", ctx)
	if err != nil {
		t.Errorf("Collection exist err: %v.", err)
		return
	}
	t.Logf("Collection exist: %t.", colExist)
}

func TestCreateDatabase(t *testing.T) {
	arangoCfg := &ArangoConfig{
		Endpoints: []string{"127.0.0.1:8529"},
		Username:  "root",
		Password:  "root",
	}

	// 获取arangodb 封装client
	client, err := NewArangoClient(arangoCfg)
	if err != nil {
		t.Errorf("New arango client err: %v.", err)
		return
	}

	ctx := context.Background()
	database, err := client.CreateDatabase("test", ctx, nil)
	if err != nil {
		t.Errorf("Create arango database err: %v.", err)
		return
	}
	t.Logf("Create databse: %s.", database.Name())
}

func TestDeleteDatabase(t *testing.T) {
	arangoCfg := &ArangoConfig{
		Endpoints: []string{"127.0.0.1:8529"},
		Username:  "root",
		Password:  "root",
	}

	// 获取arangodb 封装client
	client, err := NewArangoClient(arangoCfg)
	if err != nil {
		t.Errorf("New arango client err: %v.", err)
		return
	}

	ctx := context.Background()
	err = client.DeleteDatabase("test", ctx)
	if err != nil {
		t.Errorf("Create arango database err: %v.", err)
		return
	}
}

func TestCreateCollection(t *testing.T) {
	arangoCfg := &ArangoConfig{
		Endpoints: []string{"127.0.0.1:8529"},
		Username:  "root",
		Password:  "root",
	}

	// 获取arangodb 封装client
	client, err := NewArangoClient(arangoCfg)
	if err != nil {
		t.Errorf("New arango client err: %v.", err)
		return
	}

	ctx := context.Background()
	col, err := client.CreateColletction("megvii", "test", ctx, nil)
	if err != nil {
		t.Errorf("Create arango collection err: %v.", err)
		return
	}
	t.Logf("Create collection: %s.", col.Name())
}

func TestDeleteCollection(t *testing.T) {
	arangoCfg := &ArangoConfig{
		Endpoints: []string{"127.0.0.1:8529"},
		Username:  "root",
		Password:  "root",
	}

	// 获取arangodb 封装client
	client, err := NewArangoClient(arangoCfg)
	if err != nil {
		t.Errorf("New arango client err: %v.", err)
		return
	}

	ctx := context.Background()
	err = client.DeleteColletction("megvii", "test", ctx)
	if err != nil {
		t.Errorf("Create arango collection err: %v.", err)
		return
	}
}

func TestCreateDocument(t *testing.T) {
	arangoCfg := &ArangoConfig{
		Endpoints: []string{"127.0.0.1:8529"},
		Username:  "root",
		Password:  "root",
	}

	// 获取arangodb 封装client
	client, err := NewArangoClient(arangoCfg)
	if err != nil {
		t.Errorf("New arango client err: %v.", err)
		return
	}

	ctx := context.Background()
	document, err := client.CreateDocument("megvii", "staff", staff{
		Name:   "李四",
		Gender: "男",
		Age:    20,
	}, ctx)
	if err != nil {
		t.Errorf("Create arango collection err: %v.", err)
		return
	}
	t.Logf("%+v", document)
}

func TestCreateDocuments(t *testing.T) {
	arangoCfg := &ArangoConfig{
		Endpoints: []string{"127.0.0.1:8529"},
		Username:  "root",
		Password:  "root",
	}

	// 获取arangodb 封装client
	client, err := NewArangoClient(arangoCfg)
	if err != nil {
		t.Errorf("New arango client err: %v.", err)
		return
	}

	ctx := context.Background()
	documents, _, err := client.CreateDocuments("megvii", "staff", []interface{}{staff{
		Name:   "王六",
		Gender: "女",
		Age:    18,
	}}, ctx)
	if err != nil {
		t.Errorf("Create documents err: %v.", err)
		return
	}
	t.Logf("%+v", documents)
}

func TestDeleteDocument(t *testing.T) {
	arangoCfg := &ArangoConfig{
		Endpoints: []string{"127.0.0.1:8529"},
		Username:  "root",
		Password:  "root",
	}

	// 获取arangodb 封装client
	client, err := NewArangoClient(arangoCfg)
	if err != nil {
		t.Errorf("New arango client err: %v.", err)
		return
	}

	ctx := context.Background()
	document, err := client.DeleteDocument("megvii", "staff", "1658988706", ctx)
	if err != nil {
		t.Errorf("Delete document err: %v.", err)
		return
	}
	t.Logf("%+v", document)
}

func TestDeleteDocuments(t *testing.T) {
	arangoCfg := &ArangoConfig{
		Endpoints: []string{"127.0.0.1:8529"},
		Username:  "root",
		Password:  "root",
	}

	// 获取arangodb 封装client
	client, err := NewArangoClient(arangoCfg)
	if err != nil {
		t.Errorf("New arango client err: %v.", err)
		return
	}

	ctx := context.Background()
	documents, _, err := client.DeleteDocuments("megvii", "staff", []string{"1658988706", ""}, ctx)
	if err != nil {
		t.Errorf("Delete documents err: %v.", err)
		return
	}
	t.Logf("%+v", documents)
}

func TestUpdateDocument(t *testing.T) {
	arangoCfg := &ArangoConfig{
		Endpoints: []string{"127.0.0.1:8529"},
		Username:  "root",
		Password:  "root",
	}

	// 获取arangodb 封装client
	client, err := NewArangoClient(arangoCfg)
	if err != nil {
		t.Errorf("New arango client err: %v.", err)
		return
	}

	ctx := context.Background()
	document, err := client.UpdateDocument("megvii", "staff", "", staff{
		Name: "薛一",
	}, ctx)
	if err != nil {
		t.Errorf("Update document err: %v.", err)
		return
	}
	t.Logf("%+v", document)
}

func TestUpdateDocuments(t *testing.T) {
	arangoCfg := &ArangoConfig{
		Endpoints: []string{"127.0.0.1:8529"},
		Username:  "root",
		Password:  "root",
	}

	// 获取arangodb 封装client
	client, err := NewArangoClient(arangoCfg)
	if err != nil {
		t.Errorf("New arango client err: %v.", err)
		return
	}

	ctx := context.Background()
	documents, _, err := client.UpdateDocuments("megvii", "staff", []string{"1658988707"}, []interface{}{staff{
		Age: 10,
	}}, ctx)
	if err != nil {
		t.Errorf("Update documents err: %v.", err)
		return
	}
	t.Logf("%+v", documents)
}

func TestQueryDocuments(t *testing.T) {
	arangoCfg := &ArangoConfig{
		Endpoints: []string{"127.0.0.1:8529"},
		Username:  "root",
		Password:  "root",
	}

	// 获取arangodb 封装client
	client, err := NewArangoClient(arangoCfg)
	if err != nil {
		t.Errorf("New arango client err: %v.", err)
		return
	}

	ctx := context.Background()
	docs, err := client.QueryDocuments("megvii", "FOR d IN staff RETURN d", ctx)
	if err != nil {
		t.Errorf("Query documents err: %v.", err)
		return
	}
	for i, doc := range docs {
		t.Logf("Doc %d: %+v.", i+1, doc)
	}
}

func TestQueryDocumentsBindVariables(t *testing.T) {
	arangoCfg := &ArangoConfig{
		Endpoints: []string{"127.0.0.1:8529"},
		Username:  "root",
		Password:  "root",
	}

	// 获取arangodb 封装client
	client, err := NewArangoClient(arangoCfg)
	if err != nil {
		t.Errorf("New arango client err: %v.", err)
		return
	}

	ctx := context.Background()
	docs, err := client.QueryDocumentsBindVariables("megvii", "FOR d IN staff FILTER d.name == @name RETURN d", map[string]interface{}{"name": "张三"}, ctx)
	if err != nil {
		t.Errorf("Query documents err: %v.", err)
		return
	}
	for i, doc := range docs {
		t.Logf("Doc %d: %+v.", i+1, doc)
	}
}

type staff struct {
	Name   string `json:"name,omitempty"`
	Gender string `json:"gender,omitempty"`
	Age    int    `json:"age,omitempty"`

	Key string `json:"_key,omitempty"` // _key
}
