package client

func mergeContexts(maps ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, m := range maps {
		for k, v := range m {
			if _, ok := result[k]; !ok {
				result[k] = v
			}
		}
	}
	return result
}

func stringsMap(strings ...string) map[string]struct{} {
	result := make(map[string]struct{}, len(strings))
	for _, s := range strings {
		result[s] = struct{}{}
	}
	return result
}
