package elasticsearch

import (
	"fmt"
)

type UpdateDocByCondsParam struct {
	TermConditions    map[string]interface{}   // 精准匹配
	TermsConditions   map[string][]interface{} // 单个term 多条件精准匹配
	RangeConditions   map[string][]interface{} // 范围匹配
	MustNotConditions map[string][]interface{} // 不需匹配
	ShouldConditions  map[string][]interface{} // 匹配至少1个
}

type QueryDocByCondsParam struct {
	TermConditions    map[string]interface{}   // 精准匹配
	TermsConditions   map[string][]interface{} // 单个term 多条件精准匹配
	RangeConditions   map[string][]interface{} // 范围匹配
	MustNotConditions map[string][]interface{} // 不需匹配
	ShouldConditions  map[string][]interface{} // 匹配至少1个
}

type DeleteDocByCondsParam struct {
	TermConditions    map[string]interface{}   // 精准匹配
	TermsConditions   map[string][]interface{} // 单个term 多条件精准匹配
	RangeConditions   map[string][]interface{} // 范围匹配
	MustNotConditions map[string][]interface{} // 不需匹配
	ShouldConditions  map[string][]interface{} // 匹配至少1个
}

func convertEndpointsToUrls(endpoints []string) []string {
	if len(endpoints) == 0 {
		return endpoints
	}

	urls := make([]string, 0, len(endpoints))
	for _, endpoint := range endpoints {
		if endpoint != "" {
			urls = append(urls, fmt.Sprintf("http://%s", endpoint))
		}
	}

	return urls
}
