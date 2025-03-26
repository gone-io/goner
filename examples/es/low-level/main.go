package main

import (
	"bytes"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/es"
	"github.com/gone-io/goner/viper"
	"io"
)

type esUser struct {
	gone.Flag

	esClient *elasticsearch.Client `gone:"*"`
	logger   gone.Logger           `gone:"*"`
}

func (s *esUser) Use() {
	// 创建索引
	create, err := s.esClient.Indices.Create("my_index")
	if err != nil {
		s.logger.Errorf("Indices.Create err:%v", err)
		return
	}
	s.logger.Infof("created result: %v", create)

	// 创建文档
	document := struct {
		Name string `json:"name"`
	}{
		"go-elasticsearch",
	}
	data, _ := json.Marshal(document)
	index, err := s.esClient.Index("my_index", bytes.NewReader(data))
	if err != nil {
		s.logger.Errorf("Index err:%v", err)
		return
	}
	s.logger.Infof("Index result: %v", index)

	type ID struct {
		ID string `json:"_id"`
	}
	var id ID
	all, err := io.ReadAll(index.Body)
	if err != nil {
		s.logger.Errorf("ReadAll err:%v", err)
		return
	}

	err = json.Unmarshal(all, &id)
	if err != nil {
		s.logger.Errorf("Unmarshal err:%v", err)
		return
	}

	// 查询文档
	get, err := s.esClient.Get("my_index", id.ID)
	if err != nil {
		s.logger.Errorf("Get err:%v", err)
	}
	s.logger.Infof("Get result: %v", get)

	// 删除索引
	response, err := s.esClient.Indices.Delete([]string{"my_index"})
	if err != nil {
		s.logger.Errorf("Delete err:%v", err)
		return
	}
	s.logger.Infof("Delete result: %v", response)
}

func main() {
	gone.
		NewApp(
			viper.Load, //使用viper读取本地配置文件
			es.Load,    //使用*elasticsearch.Client
		).
		Run(func(esUser *esUser) {
			esUser.Use()
		})
}
