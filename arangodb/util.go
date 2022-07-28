package arangodb

import (
	"fmt"
)

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
