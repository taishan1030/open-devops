package database

import (
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"open-devops/src/common"
	"strings"
)

func StreePathAddTest(logger log.Logger) {
	ns := []string{
		"inf.monitor.thanos",
		"inf.monitor.kafka",
		"inf.monitor.prometheus",
		"inf.monitor.m3db",
		"inf.cicd.gray",
		"inf.cicd.deploy",
		"inf.cicd.jenkins",
		"waimai.qiangdan.queue",
		"waimai.qiangdan.worker",
		"waimai.ditu.kafka",
		"waimai.ditu.es",
		"waimai.qiangdan.es",
	}
	for _, n := range ns {
		req := &common.NodeCommonReq{
			Node: n,
		}
		StreePathAddOne(req, logger)
	}
}

// 编写查询node的测试函数
func StreePathQueryTest(logger log.Logger) {
	ns := []string{
		"a",
		"b",
		"c",
		"inf",
		"waimai",
	}
	// query_type =1
	for _, n := range ns {
		req := &common.NodeCommonReq{
			Node:      n,
			QueryType: 1,
		}
		res := StreePathQuery(req, logger)
		level.Info(logger).Log("msg", "StreePathQuery.res", "req.node", n, "num", len(res), "details", strings.Join(res, ","))
	}
	// query_type =2
	for _, n := range ns {
		req := &common.NodeCommonReq{
			Node:      n,
			QueryType: 2,
		}
		res := StreePathQuery(req, logger)
		level.Info(logger).Log("msg", "StreePathQuery.res", "req.node", n, "num", len(res), "details", strings.Join(res, ","))
	}
	// query_type =3
	ns3 := []string{
		"a.b",
		"b.a",
		"c.d",
		"inf.cicd",
		"inf.monitor",
		"waimai.ditu",
		"waimai.monitor",
		"waimai.qiangdan",
	}
	for _, n := range ns3 {
		req := &common.NodeCommonReq{
			Node:      n,
			QueryType: 3,
		}
		res := StreePathQuery(req, logger)
		level.Info(logger).Log("msg", "StreePathQuery.res", "req.node", n, "num", len(res), "details", strings.Join(res, ","))
	}
}

func StreePathDelTest(logger log.Logger) {
	ns := []string{
		"inf.cicd.jenkins",
		"inf.cicd",
		"inf",
	}
	for _, n := range ns {
		req := &common.NodeCommonReq{
			Node: n,
		}
		res := StreePathDelete(req, logger)
		level.Info(logger).Log("msg", "StreePathDelete.res", "req.node", n, "del_num", res)
	}
}
