package router

type Params map[string]string

func (p Params) append(key, value string) {
	p[key] = value
}

func (p Params) Has(key string) bool {
	_, found := p[key]
	return found
}

func (p Params) Get(key string) string {
	return p[key]
}
