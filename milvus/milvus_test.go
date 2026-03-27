package milvusx

import (
	"context"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"math/rand"
	"testing"
	"time"
)

func TestMilvusRepo(t *testing.T) {

	cfg := &MilvusConfig{
		Address:  "rw-535edc.milvus.qihudb.net:19032",
		Username: "milvus_wr",
		Password: "5fa2d4b73a5d7365",
		Timeout:  5,
		DbName:   "video_similar",
	}
	cli, err := NewMilvus(cfg)
	if err != nil {
		t.Errorf("Milvus connect failed, err: %v", err)
		return
	}
	defer cli.Close()

	t.Log("Connect milvus success")

	// 检查collection 是否存在
	isHas, err := cli.GetClient().HasCollection(context.Background(), testCollection)
	if err != nil {
		t.Errorf("Check collection failed, err: %v", err)
		return
	}
	t.Logf("Collection exists: %v", isHas)

	repo := NewMilvusRepo(cli, testCollection)

	// 插入数据
	rec := &VectorRecord{
		VideoID:  randomInt16(),
		FrameNum: randomInt16(),
		Vector:   randomVector(128),
	}

	if err := repo.Insert(context.Background(), rec); err != nil {
		t.Errorf("Insert failed, err: %v", err)
		return
	}

	t.Log("Insert success")

	// 搜索
	results, err := repo.Search(context.Background(), rec.Vector, 10)
	if err != nil {
		t.Errorf("Search failed, err: %v", err)
		return
	}

	if len(results) == 0 {
		t.Error("Search returned 0 results")
		return
	}

	for i, r := range results {
		t.Logf("Result %d, id: %d, video_id:%d, frame_num:%d, score:%f", i, r.ID, r.VideoID, r.FrameNum, r.Score)
	}
	/**
	示例：
	  milvus_test.go:66: Result 0, id: 465037424767021753, video_id:28305, frame_num:26893, score:44.411415
	    milvus_test.go:66: Result 1, id: 465037424767021658, video_id:1001, frame_num:1, score:40.929562
	    milvus_test.go:66: Result 2, id: 465037424767021693, video_id:1001, frame_num:2, score:40.929562
	    milvus_test.go:66: Result 3, id: 465037424767021695, video_id:1001, frame_num:3, score:40.929562
	    milvus_test.go:66: Result 4, id: 465037424767021697, video_id:1001, frame_num:4, score:40.929562
	    milvus_test.go:66: Result 5, id: 465037424767021742, video_id:16356, frame_num:23118, score:40.929562
	    milvus_test.go:66: Result 6, id: 465037424767021743, video_id:32750, frame_num:13408, score:40.929562
	    milvus_test.go:66: Result 7, id: 465037424767021744, video_id:16517, frame_num:24917, score:40.929562
	    milvus_test.go:66: Result 8, id: 465037424767021752, video_id:16457, frame_num:7717, score:35.160881
	    milvus_test.go:66: Result 9, id: 465037424767021749, video_id:23356, frame_num:27002, score:34.239025

	*/
}

type VectorRecord struct {
	VideoID  int16
	FrameNum int16
	Vector   []float32
}

type SearchResult struct {
	ID       int64
	VideoID  int16
	FrameNum int16
	Score    float32
}

type VectorRepository interface {
	// 插入向量数据
	Insert(ctx context.Context, rec *VectorRecord) error
	// 根据向量检索topK 数据
	Search(ctx context.Context, vec []float32, topK int) ([]SearchResult, error)
}

type MilvusRepo struct {
	cli        *Milvus
	collection string
}

func NewMilvusRepo(cli *Milvus, collection string) *MilvusRepo {
	return &MilvusRepo{
		cli:        cli,
		collection: collection,
	}
}

func (r *MilvusRepo) Insert(ctx context.Context, rec *VectorRecord) error {
	_, err := r.cli.GetClient().Insert(
		ctx,
		r.collection,
		"",
		entity.NewColumnInt16("video_id", []int16{rec.VideoID}),
		entity.NewColumnInt16("frame_num", []int16{rec.FrameNum}),
		entity.NewColumnFloatVector("frame_vector", len(rec.Vector), [][]float32{rec.Vector}),
	)
	return err
}

func (r *MilvusRepo) Search(ctx context.Context, vec []float32, topK int) ([]SearchResult, error) {
	sp, _ := entity.NewIndexIvfFlatSearchParam(128)

	results, err := r.cli.GetClient().Search(
		ctx,
		r.collection,
		nil,
		"",
		[]string{"id", "video_id", "frame_num"},
		[]entity.Vector{entity.FloatVector(vec)},
		"frame_vector",
		entity.IP,
		topK,
		sp,
	)
	if err != nil {
		return nil, err
	}

	var out []SearchResult
	for _, result := range results {
		idCol := result.Fields.GetColumn("id").(*entity.ColumnInt64)
		videoCol := result.Fields.GetColumn("video_id").(*entity.ColumnInt16)
		frameCol := result.Fields.GetColumn("frame_num").(*entity.ColumnInt16)
		scores := result.Scores

		for i := 0; i < result.ResultCount; i++ {
			out = append(out, SearchResult{
				ID:       idCol.Data()[i],
				VideoID:  videoCol.Data()[i],
				FrameNum: frameCol.Data()[i],
				Score:    scores[i],
			})
		}
	}

	return out, nil
}

// 测试集合
const testCollection = "gosdk_insert_example"

func randomVector(dim int) []float32 {
	vec := make([]float32, dim)
	for i := range vec {
		vec[i] = rand.Float32() // 0.0 ~ 1.0
	}
	return vec
}

// 初始化随机种子
func init() {
	rand.Seed(time.Now().UnixNano())
}

// 返回随机 int16 (0~32767)
func randomInt16() int16 {
	return int16(rand.Intn(1 << 15))
}
