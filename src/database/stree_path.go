package database

import (
	"fmt"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"open-devops/src/common"

	"sort"
	"strings"
)

type StreePath struct {
	Id       int64  `json:"id"`
	Level    int64  `json:"level"`
	Path     string `json:"path"`
	NodeName string `json:"node_name"`
}

func (StreePath) TableName() string {
	return "stree_path"
}

func (streePath *StreePath) CreateStreePath() error {
	return GetDb().Table(streePath.TableName()).Create(streePath).Error
}

func (streePath *StreePath) GetStreePath() (StreePath, error) {
	sp := StreePath{}
	err := GetDb().Table(streePath.TableName()).
		Where("level=? and path=? and node_name=?", streePath.Level, streePath.Path, streePath.NodeName).
		Find(&sp).Error
	return sp, err

}

func (streePath *StreePath) DeleteStreePath() (int64, error) {
	tx := GetDb().Table(streePath.TableName()).Delete(streePath)
	return tx.RowsAffected, tx.Error
}

// 带参数查询多条记录函数
func StreePathGetMany(where string, args ...interface{}) ([]StreePath, error) {
	var objs []StreePath
	err := GetDb().Table(StreePath{}.TableName()).Where(where, args...).Find(&objs).Error
	if err != nil {
		return objs, err
	}
	return objs, nil
}

func StreePathAddOne(req *common.NodeCommonReq, logger log.Logger) {
	res := strings.Split(req.Node, ".")
	if len(res) != 3 {
		level.Info(logger).Log("msg", "add.path.invalidate", "path", req.Node)
		return
	}
	// g.p.a
	g, p, a := res[0], res[1], res[2]
	// 先查g
	nodeG := &StreePath{
		Level:    1,
		Path:     "0",
		NodeName: g,
	}
	dbG, err := nodeG.GetStreePath()
	if err != nil {
		level.Error(logger).Log("msg", "check.g.failed", "path", req.Node, "err", err)
		return
	}
	// 根据查询结果判断
	switch dbG.Id {
	case 0:
		// 说明g不存在，一次插入g,p,a
		// 插入g
		err = nodeG.CreateStreePath()
		if err != nil {
			level.Error(logger).Log("msg", "g_not_exist_add_g_failed", "path", req.Node, "err", err)
			return
		}
		level.Info(logger).Log("msg", "g_not_exist_add_g_success", "path", req.Node)
		// 插入p
		nodeP := &StreePath{
			Level:    2,
			Path:     fmt.Sprintf("/%d", nodeG.Id),
			NodeName: p,
		}
		err = nodeP.CreateStreePath()
		if err != nil {
			level.Error(logger).Log("msg", "g_not_exist_add_p_failed", "path", req.Node, "err", err)
			return
		}
		level.Info(logger).Log("msg", "g_not_exist_add_p_success", "path", req.Node)
		// 插入a
		nodeA := &StreePath{
			Level:    3,
			Path:     fmt.Sprintf("/%d/%d", nodeG.Id, nodeP.Id),
			NodeName: a,
		}
		err = nodeA.CreateStreePath()
		if err != nil {
			level.Error(logger).Log("msg", "g_not_exist_add_a_failed", "path", req.Node, "err", err)
			return
		}
		level.Info(logger).Log("msg", "g_not_exist_add_a_success", "path", req.Node)
	default:
		level.Info(logger).Log("msg", "g_exist_check_p", "path", req.Node)
		// 查询p是否存在
		nodeP := &StreePath{
			Level:    2,
			Path:     fmt.Sprintf("/%d", dbG.Id),
			NodeName: p,
		}
		dbP, err := nodeP.GetStreePath()
		if err != nil {
			level.Error(logger).Log("msg", "check.g.failed", "path", req.Node, "err", err)
			return
		}
		if dbP.Id > 0 {
			// p 存在, 查询a
			nodeA := &StreePath{
				Level:    3,
				Path:     fmt.Sprintf("%s/%d", dbP.Path, dbP.Id),
				NodeName: a,
			}
			dbA, err := nodeA.GetStreePath()
			if err != nil {
				level.Error(logger).Log("msg", "g_p_exist_check_a_failed", "path", req.Node, "err", err)
				return
			}
			if dbA.Id <= 0 {
				// 说明a不存在，插入a
				err := nodeA.CreateStreePath()
				if err != nil {
					level.Error(logger).Log("msg", "g_p_exist_add_a_failed", "path", req.Node, "err", err)
					return
				}
				level.Info(logger).Log("msg", "g_p_exist_add_a_success", "path", req.Node)
				return
			}
			level.Info(logger).Log("msg", "g_p_a_exist", "path", req.Node)
			return

		}
		// p不存在
		err = nodeP.CreateStreePath()
		if err != nil {
			level.Error(logger).Log("msg", "g_exist_add_p_failed", "path", req.Node, "err", err)
			return
		}
		level.Info(logger).Log("msg", "g_exist_add_p_success", "path", req.Node)
		// 插入a
		nodeA := &StreePath{
			Level:    3,
			Path:     fmt.Sprintf("/%d/%d", dbG.Id, nodeP.Id),
			NodeName: a,
		}
		err = nodeA.CreateStreePath()
		if err != nil {
			level.Error(logger).Log("msg", "g_not_exist_add_a_failed", "path", req.Node, "err", err)
			return
		}
		level.Info(logger).Log("msg", "g_not_exist_add_a_success", "path", req.Node)
	}
}

func StreePathQuery(req *common.NodeCommonReq, logger log.Logger) (res []string) {

	switch req.QueryType {
	case 1:
		nodeG := &StreePath{
			Level:    1,
			Path:     "0",
			NodeName: req.Node,
		}
		// 判断g是否存在
		dbG, err := nodeG.GetStreePath()
		if err != nil {
			level.Error(logger).Log("msg", "query_g_failed", "path", req.Node, "err", err)
			return
		}
		if dbG.Id <= 0 {
			// 说明要查询的g不存在
			return
		}
		pathP := fmt.Sprintf("/%d", dbG.Id)
		// 根据g查询 所有p的列表 node=g query_type=1
		whereStr := "level=? and path=?"
		ps, err := StreePathGetMany(whereStr, 2, pathP)
		if err != nil {
			level.Error(logger).Log("msg", "query_ps_failed", "path", req.Node, "err", err)
			return
		}
		for _, i := range ps {
			res = append(res, i.NodeName)
		}
		sort.Strings(res)
		return
	case 2:
		/*
			编写query_type=2的查询 根据g查询 所有g.p.a的列表
			先查 g ，再查p 最后查a ，中间有一步没有都返回空
		*/
		nodeG := &StreePath{
			Level:    1,
			Path:     "0",
			NodeName: req.Node,
		}
		// 判断g是否存在
		dbG, err := nodeG.GetStreePath()
		if err != nil {
			level.Error(logger).Log("msg", "query_g_failed", "path", req.Node, "err", err)
			return
		}
		if dbG.Id <= 0 {
			// 说明要查询的g不存在
			return
		}
		pathP := fmt.Sprintf("/%d", dbG.Id)
		// 根据g查询 所有p的列表 node=g query_type=1
		whereStr := "level=? and path=?"
		ps, err := StreePathGetMany(whereStr, 2, pathP)
		if err != nil {
			level.Error(logger).Log("msg", "query_ps_failed", "path", req.Node, "err", err)
			return
		}
		if len(ps) == 0 {
			//	 说明g下面没有p
			return
		}
		for _, p := range ps {
			pathA := fmt.Sprintf("%s/%d", p.Path, p.Id)
			as, err := StreePathGetMany(whereStr, 3, pathA)
			if err != nil {
				level.Error(logger).Log("msg", "query_as_failed", "path", req.Node, "err", err)
				continue
			}
			if len(as) == 0 {
				// 说明该p下没有a
				continue
			}
			for _, a := range as {
				fullPath := fmt.Sprintf("%s.%s.%s", dbG.NodeName, p.NodeName, a.NodeName)
				res = append(res, fullPath)
			}
		}
		sort.Strings(res)
		return
	case 3:
		/*
			编写query_type=3的查询 根据g.p查询 所有g.p.a的列表 node=g.p query_type=3
			先查询 g 和p，不存在直接返回空
			查p时需要带上p.name查询
		*/
		gps := strings.Split(req.Node, ".")
		g, p := gps[0], gps[1]
		nodeG := &StreePath{
			Level:    1,
			Path:     "0",
			NodeName: g,
		}
		// 判断g是否存在
		dbG, err := nodeG.GetStreePath()
		if err != nil {
			level.Error(logger).Log("msg", "query_g_failed", "path", req.Node, "err", err)
			return
		}
		if dbG.Id <= 0 {
			// 说明要查询的g不存在
			return
		}
		//g存在，这里不需要查全量的p，只查询匹配这个node_name的p
		NodeP := &StreePath{
			Level:    2,
			Path:     fmt.Sprintf("/%d", dbG.Id),
			NodeName: p,
		}
		dbP, err := NodeP.GetStreePath()
		if err != nil {
			level.Error(logger).Log("msg", "query_p_failed", "path", req.Node, "err", err)
			return
		}
		if dbP.Id <= 0 {
			// 说明p不存在
			return
		}
		pathA := fmt.Sprintf("%s/%d", dbP.Path, dbP.Id)
		whereStr := "level=? and path=? "
		as, err := StreePathGetMany(whereStr, 3, pathA)
		if err != nil {
			level.Error(logger).Log("msg", "query_as_failed", "path", req.Node, "err", err)
			return
		}
		for _, a := range as {
			fullPath := fmt.Sprintf("%s.%s.%s", dbG.NodeName, dbP.NodeName, a.NodeName)
			res = append(res, fullPath)
		}
		sort.Strings(res)
		return
	}
	return
}

// 传入g，如果g下有p就不让删g
// 传入g.p，如果p下有a就不让删p
// 传入g.p.a，直接删
func StreePathDelete(req *common.NodeCommonReq, logger log.Logger) (delNum int64) {
	path := strings.Split(req.Node, ".")
	pLevel := len(path)
	//	  传入g，如果g下有p就不让删g
	nodeG := &StreePath{
		Level:    1,
		Path:     "0",
		NodeName: path[0],
	}
	dbG, err := nodeG.GetStreePath()
	if err != nil {
		level.Error(logger).Log("msg", "query_g_failed", "path", req.Node, "err", err)
		return
	}
	if dbG.Id <= 0 {
		// 说明要删除的g不存在
		return 0
	}
	pathP := fmt.Sprintf("/%d", dbG.Id)
	switch pLevel {
	case 1:
		whereStr := "level=? and path=?"
		ps, err := StreePathGetMany(whereStr, 2, pathP)
		if err != nil {
			level.Error(logger).Log("msg", "query_ps_failed", "path", req.Node, "err", err)
			return
		}
		if len(ps) > 0 {
			level.Warn(logger).Log("msg", "del_g_reject", "path", req.Node, "reason", "g_has_ps", "ps_num", len(ps))
			return
		}
		delNum, err = dbG.DeleteStreePath()
		if err != nil {
			level.Error(logger).Log("msg", "del_g_failed", "path", req.Node, "err", err)
			return
		}
		level.Info(logger).Log("msg", "del_g_success", "path", req.Node)
		return
	case 2:
		nodeP := &StreePath{
			Level:    2,
			Path:     pathP,
			NodeName: path[1],
		}
		dbP, err := nodeP.GetStreePath()
		if err != nil {
			level.Error(logger).Log("msg", "query_p_failed", "path", req.Node, "err", err)
			return
		}
		if dbP.Id <= 0 {
			// 说明p不存在
			return
		}
		pathA := fmt.Sprintf("%s/%d", dbP.Path, dbP.Id)
		whereStr := "level=? and path=?"
		as, err := StreePathGetMany(whereStr, 3, pathA)
		if err != nil {
			level.Error(logger).Log("msg", "query_as_failed", "path", req.Node, "err", err)
			return
		}
		if len(as) > 0 {
			level.Warn(logger).Log("msg", "del_g_p_reject", "path", req.Node, "reason", "p_has_as", "as_num", len(as))
			return
		}
		delNum, err = dbP.DeleteStreePath()
		if err != nil {
			level.Error(logger).Log("msg", "del_p_failed", "path", req.Node, "err", err)
			return
		}
		level.Info(logger).Log("msg", "del_p_success", "path", req.Node)
		return
	case 3:
		nodeP := &StreePath{
			Level:    2,
			Path:     pathP,
			NodeName: path[1],
		}
		dbP, err := nodeP.GetStreePath()
		if err != nil {
			level.Error(logger).Log("msg", "query_p_failed", "path", req.Node, "err", err)
			return
		}
		if dbP.Id <= 0 {
			// 说明p不存在
			return
		}
		nodeA := &StreePath{
			Level:    3,
			Path:     fmt.Sprintf("%s/%d", dbP.Path, dbP.Id),
			NodeName: path[2],
		}
		dbA, err := nodeA.GetStreePath()
		if err != nil {
			level.Error(logger).Log("msg", "query_a_failed", "path", req.Node, "err", err)
			return
		}
		if dbA.Id <= 0 {
			return
		}
		delNum, err = dbA.DeleteStreePath()
		if err != nil {
			level.Error(logger).Log("msg", "del_a_failed", "path", req.Node, "err", err)
			return
		}
		level.Info(logger).Log("msg", "del_a_success", "path", req.Node)
		return
	}
	return 0
}
