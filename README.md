# go-common

[![Go Report Card](https://goreportcard.com/badge/github.com/psoKnight/go-common)](https://goreportcard.com/report/github.com/psoKnight/go-common)
[![Master](https://github.com/psoKnight/go-common/workflows/Go/badge.svg)](https://github.com/psoKnight/go-common/actions)

Golang 公共封装库。

# 依赖管理

`go-common` 工程使用`Go Modules`管理第三方依赖，详情请参考[Go Modules 官方文档](https://blog.golang.org/using-go-modules)。

# 源代码目录结构说明

```$xslt
├── arangodb    // ArangoDB 基础封装库
├── cache   // GCache LRU 基础封装库
├── clickhouse  // ClickHouse 基础封装库
├── elasticsearch    // ElasticSearch 基础封装库
├── etcd    // ETCD 基础封装库
├── grpcz   // GRPC 基础封装库
├── http    // HTTP 基础封装库
│     ├── http-native // 原生操作http 封装（net/http）
│     ├── http-go // HTTP-GO v1.0.0
├── kafka   // Kafka 基础封装库
├── log // Log 基础封装库
├── mqtt    // MQTT 基础封装库
├── mysql   // MySQL 基础封装库
├── redis   // Redis 基础封装库
├── rocketmq    // RocketMQ 基础封装库
├── seaweed-fs  // SeaweedFS 基础封装库
├── utils   // 工具/方法包
│   ├── copy_struct_by_tag_util // 根据JSON tag 拷贝结构体 
│   ├── date_util   // 日期转换
│   ├── discard_json_comments   // 去除JSON 注释  
│   ├── feature // 特征计算相关
│   │     ├── L2_distance // 使用L2（Euclidean）distance 计算图片特征值
│   ├── geohash_util    // GeoHash 算法
│   ├── pool_util   //  并发协程池
│   ├── rand_util   // 随机数 
│   ├── semaphore_util   // 并发信号量 
```

## 包规范

    如果包的命名尾缀带"_"表示未经过测试，示例：hadoop_。

# Git tag 记录

- v1.0.0：增加arangodb、cache、clickhouse、elasticsearch、etcd、grpc、http、kafka、log、mqtt（包括mqtt 和alibaba
  mqtt）、mysql、redis、rocketmq、seaweed-fs、utils 基础封装库；