package Schema

type Field struct {
	Name   string                 `json:"name"`
	Type   string                 `json:"type"`
	Pos    uint                   `json:"pos"`
	Length uint                   `json:"length"`
	Tags   map[string]interface{} `json:"tags"`
}

type Schema struct {
	Type   string  `json:"type"`
	Name   string  `json:"name"`
	Fields []Field `json:"fields"`
}
