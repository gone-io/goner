# 使用 github.com/gone-io/goner/es 示例

## 环境准备

### 1. 本地部署 Elasticsearch

使用官方提供的快速部署脚本：

```bash
curl -fsSL https://elastic.co/start-local | sh
```

脚本会自动下载并启动最新版本的 Elasticsearch。启动成功后，可以通过 http://localhost:9200 访问。

### 2. 配置 API Key

1. 访问 Elasticsearch 控制台，创建 API Key
2. 将生成的 API Key 复制到配置文件 `config/default.yaml` 中：

![API Key 配置](images/img.png)

配置文件示例：
```yaml
es:
  addresses: http://localhost:9200   # Elasticsearch 节点地址
  aPIKey: "your-api-key-here"      # 将这里替换为你的 API Key
```

## 运行示例

本示例提供了两种使用方式的演示：低级 API 和完全类型化 API。

### 1. 低级 API 示例

运行命令：
```bash
go run ./low-level/main.go 
```

执行结果：
```log
2025/03/26 16:22:16 created result: [200 OK] {"acknowledged":true,"shards_acknowledged":true,"index":"my_index"}
2025/03/26 16:22:16 Index result: [201 Created] {"_index":"my_index","_id":"KH6L0ZUBb_HhHUd_hFDu","_version":1,"result":"created","_shards":{"total":2,"successful":1,"failed":0},"_seq_no":0,"_primary_term":1}
2025/03/26 16:22:16 Get result: [200 OK] {"_index":"my_index","_id":"KH6L0ZUBb_HhHUd_hFDu","_version":1,"_seq_no":0,"_primary_term":1,"found":true,"_source":{"name":"go-elasticsearch"}}
2025/03/26 16:22:16 Delete result: [200 OK] {"acknowledged":true}
```

示例演示了以下操作：
1. 创建索引
2. 添加文档
3. 查询文档
4. 删除索引

### 2. 完全类型化 API 示例

运行命令：
```bash
go run ./fully-typed/main.go
```

执行结果：
```log
2025/03/26 16:23:01 created result: &create.Response{Acknowledged:true, Index:"my_index", ShardsAcknowledged:true}
2025/03/26 16:23:01 Index result: &index.Response{ForcedRefresh:(*bool)(nil), Id_:"KX6M0ZUBb_HhHUd_NVAG", Index_:"my_index", PrimaryTerm_:(*int64)(0x1400000f038), Result:result.Result{Name:"created"}, SeqNo_:(*int64)(0x1400000f030), Shards_:types.ShardStatistics{Failed:0x0, Failures:[]types.ShardFailure(nil), Skipped:(*uint)(nil), Successful:0x1, Total:0x2}, Version_:1}
2025/03/26 16:23:01 Get result: &get.Response{Fields:map[string]json.RawMessage{}, Found:true, Id_:"KX6M0ZUBb_HhHUd_NVAG", Ignored_:[]string(nil), Index_:"my_index", PrimaryTerm_:(*int64)(0x140000a2448), Routing_:(*string)(nil), SeqNo_:(*int64)(0x140000a2440), Source_:json.RawMessage{0x7b, 0x22, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x3a, 0x22, 0x67, 0x6f, 0x2d, 0x65, 0x6c, 0x61, 0x73, 0x74, 0x69, 0x63, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x22, 0x7d}, Version_:(*int64)(0x140000a2438)}
2025/03/26 16:23:01 Delete result: &delete.Response{Acknowledged:true, Shards_:(*types.ShardStatistics)(nil)}
```

完全类型化 API 提供了更好的类型安全性和 IDE 支持，推荐在生产环境中使用。

## 常见问题

1. 如果遇到连接问题，请确保：
   - Elasticsearch 服务已正常启动
   - API Key 配置正确
   - 防火墙未阻止相关端口

2. 如需了解更多配置选项，请参考 `es/README_CN.md` 文档。