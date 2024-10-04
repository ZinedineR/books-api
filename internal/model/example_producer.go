package model

type ExampleMessage struct {
	Id                 string     `json:"id" form:"id"`
	TriggerBy          string     `form:"trigger_by" json:"trigger_by"`
	Module             string     `form:"module" json:"module"`
	Description        string     `form:"description" json:"description"`
	Email              string     `json:"email" form:"email"`
	Name               string     `bson:"name" json:"name,omitempty"`
	Cc                 []string   `json:"cc" form:"cc"`
	Variable           []Variable `json:"variable" form:"variable"`
	Host               string     `bson:"host" json:"host,omitempty"`
	AttachmentFilePath string     `json:"attachmentFilePath" form:"attachmentFilePath"`
}

type Variable struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func (model *ExampleMessage) VariableToMap() map[string]interface{} {
	variableMap := make(map[string]interface{})

	for _, v := range model.Variable {
		variableMap[v.Key] = v.Value
	}

	return variableMap
}
