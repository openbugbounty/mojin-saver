package command

type Tags map[string]interface{}

func (t Tags) AddTag(name string, val interface{}) {
	switch val.(type) {
	case string:
		if v, ok := t[name]; ok {
			t[name] = duplicateAppend(v, val.(string))
		} else {
			t[name] = val
		}
	case []string:
		slice := val.([]string)
		for _, v := range slice {
			t.AddTag(name, v)
		}
	default:
	}
}

func duplicateAppend(input interface{}, val string) interface{} {
	switch input.(type) {
	case string:
		if input.(string) == val {
			return input
		}
		return []string{input.(string), val}
	case []string:
		m := make(map[string]struct{})
		for _, v := range input.([]string) {
			m[v] = struct{}{}
		}
		if _, ok := m[val]; ok {
			return input
		}
		return append(input.([]string), val)
	default:
	}

	return []string{}
}
