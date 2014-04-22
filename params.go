package bourbon

// Params is a map of parameter names and their values from the request URL.
type Params map[string]string

func createParams(c *context) Params {
	params := make(Params)

	keys := c.route.Regexp.FindAllStringSubmatch(c.route.Pattern, -1)
	if len(keys) == 0 {
		return params
	}

	vals := c.route.Regexp.FindAllStringSubmatch(c.r.URL.Path, -1)
	for i := 1; i < len(keys[0]); i++ {
		key := keys[0][i]
		val := vals[0][i]
		key = key[1 : len(key)-1]
		params[key] = val
	}

	return params
}
