package clickhouse

import (
	"context"
	"fmt"
	"testing"
)

func TestNewClickhouse(t *testing.T) {

	cfg := &ClickhouseConfig{
		Addrs:    []string{"10.117.48.122:9000"},
		Database: "default",
		Username: "default",
		Password: "megvii2022",
	}

	client, err := NewClickhouse(cfg)
	if err != nil {
		t.Errorf("New clickhouse client err: %v.", err)
		return
	}
	t.Log(client)
	t.Log(client.cfg)
}

func TestExecCQL(t *testing.T) {

	cfg := &ClickhouseConfig{
		Addrs:    []string{"10.117.48.122:9000"},
		Database: "default",
		Username: "default",
		Password: "megvii2022",
	}

	client, err := NewClickhouse(cfg)
	if err != nil {
		t.Errorf("New clickhouse client err: %v.", err)
		return
	}

	ctx := context.Background()
	err = client.ExecCQL("CREATE database IF NOT EXISTS megvii;", ctx)
	if err != nil {
		t.Errorf("Exec cql err: %v.", err)
		return
	}

	err = client.ExecCQL("CREATE TABLE IF NOT EXISTS megvii.staff (`name` String, `age` UInt32, `gender` UInt8) ENGINE = MergeTree() PARTITION BY gender ORDER BY age;", ctx)
	if err != nil {
		t.Errorf("Exec cql 2 err: %v.", err)
		return
	}
}

func TestInsert(t *testing.T) {
	cfg := &ClickhouseConfig{
		Addrs:    []string{"10.117.48.122:9000"},
		Database: "default",
		Username: "default",
		Password: "megvii2022",
	}

	client, err := NewClickhouse(cfg)
	if err != nil {
		t.Errorf("New clickhouse client err: %v.", err)
		return
	}

	ctx := context.Background()

	insert := fmt.Sprintf("INSERT INTO megvii.staff (name, age, gender) VALUES ('%s', %d, %d)", "stick", 1, 2)
	t.Log(insert)
	err = client.ExecCQL(insert, ctx)
	if err != nil {
		t.Errorf("Insert err: %v.", err)
		return
	}
}

func TestBatchInsert(t *testing.T) {
	cfg := &ClickhouseConfig{
		Addrs:    []string{"10.117.48.122:9000"},
		Database: "default",
		Username: "default",
		Password: "megvii2022",
	}

	client, err := NewClickhouse(cfg)
	if err != nil {
		t.Errorf("New clickhouse client err: %v.", err)
		return
	}

	ctx := context.Background()
	batch, err := client.cli.PrepareBatch(ctx, "INSERT INTO megvii.staff")
	if err != nil {
		t.Errorf("Prepare batch err: %v.", err)
		return
	}

	for i := 1; i <= 10000; i++ {
		err := batch.AppendStruct(&staff{
			Name:   "旷小刚",
			Gender: 1,
			Age:    uint32(i),
		})
		if err != nil {
			t.Errorf("Append struct %d err: %v.", i, err)
			continue
		}
	}

	err = batch.Send()
	if err != nil {
		t.Errorf("Send err: %v.", err)
		return
	}
}

func TestSelect(t *testing.T) {

	cfg := &ClickhouseConfig{
		Addrs:    []string{"10.117.48.122:9000"},
		Database: "default",
		Username: "default",
		Password: "megvii2022",
	}

	client, err := NewClickhouse(cfg)
	if err != nil {
		t.Errorf("New clickhouse client err: %v.", err)
		return
	}

	var results []staff

	ctx := context.Background()
	err = client.SelectCQL("SELECT name, gender, age FROM megvii.staff LIMIT 3", &results, ctx)
	if err != nil {
		t.Errorf("Select err: %v.", err)
		return
	}

	for _, v := range results {
		t.Log(v)
	}
}

func TestSelectRows(t *testing.T) {

	cfg := &ClickhouseConfig{
		Addrs:    []string{"10.117.48.122:9000"},
		Database: "default",
		Username: "default",
		Password: "megvii2022",
	}

	client, err := NewClickhouse(cfg)
	if err != nil {
		t.Errorf("New clickhouse client err: %v.", err)
		return
	}

	ctx := context.Background()
	rows, err := client.cli.Query(ctx, "SELECT name, gender, age FROM megvii.staff PREWHERE age >= $1 and age < $2", 50, 300)
	if err != nil {
		t.Errorf("Query err: %v.", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		s := &staff{}
		if err := rows.ScanStruct(s); err != nil {
			return
		}
		t.Log(s)
	}
}

func TestGetCount(t *testing.T) {

	cfg := &ClickhouseConfig{
		Addrs:    []string{"10.117.48.122:9000"},
		Database: "default",
		Username: "default",
		Password: "megvii2022",
	}

	client, err := NewClickhouse(cfg)
	if err != nil {
		t.Errorf("New clickhouse client err: %v.", err)
		return
	}

	ctx := context.Background()
	count, err := client.GetCQLCount("SELECT COUNT() c FROM megvii.staff", ctx)
	if err != nil {
		t.Errorf("Get count err: %v.", err)
		return
	}

	t.Log(count)
}

type staff struct {
	Name   string `json:"name" ch:"name"`
	Gender uint8  `json:"gender" ch:"gender"`
	Age    uint32 `json:"age" ch:"age"`
}
