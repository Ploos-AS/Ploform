package formschema

type Field struct {
 ID string `json:"id"`
 Type string `json:"type"`
 Label string `json:"label"`
 Required bool `json:"required,omitempty"`
 Options []string `json:"options,omitempty"`
}

type Schema struct {
 Fields []Field `json:"fields"`
}
