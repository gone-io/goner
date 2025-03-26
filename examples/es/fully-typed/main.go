package main

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gone-io/gone/v2"
	"github.com/gone-io/goner/es"
	"github.com/gone-io/goner/viper"
)

type esUser struct {
	gone.Flag

	esClient *elasticsearch.TypedClient `gone:"*"`
	logger   gone.Logger                `gone:"*"`
}

func (s *esUser) Use() {
	// 创建索引
	create, err := s.esClient.Indices.Create("my_index").Do(context.TODO())
	if err != nil {
		s.logger.Errorf("Indices.Create err:%v", err)
		return
	}
	s.logger.Infof("created result: %#+v", create)

	// 创建文档
	document := struct {
		Name string `json:"name"`
	}{
		"go-elasticsearch",
	}
	index, err := s.esClient.Index("my_index").Document(document).Do(context.TODO())
	if err != nil {
		s.logger.Errorf("Index err:%v", err)
		return
	}
	s.logger.Infof("Index result: %#+v", index)

	// 查询文档
	get, err := s.esClient.Get("my_index", index.Id_).Do(context.TODO())
	if err != nil {
		s.logger.Errorf("Get err:%v", err)
		return
	}
	s.logger.Infof("Get result: %#+v", get)

	// 删除索引
	response, err := s.esClient.Indices.Delete("my_index").Do(context.TODO())
	if err != nil {
		s.logger.Errorf("Delete err:%v", err)
		return
	}
	s.logger.Infof("Delete result: %#+v", response)
}

func main() {
	gone.
		NewApp(
			viper.Load,         //使用viper读取本地配置文件
			es.LoadTypedClient, //使用*elasticsearch.TypedClient
		).
		Run(func(esUser *esUser) {
			esUser.Use()
		})
}
