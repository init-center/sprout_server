package models

type ConfigItem struct {
	Key     string `db:"config_key" json:"key"`
	Value   string `db:"config_value" json:"value"`
	Explain string `db:"config_explain" json:"explain"`
}

type Configs = []ConfigItem

type ConfigList struct {
	Page Page    `json:"page"`
	List Configs `json:"list"`
}
