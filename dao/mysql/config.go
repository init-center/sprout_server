package mysql

import (
	"sprout_server/models"
	"sprout_server/models/queryfields"
)

func CreateConfig(p *models.ParamsAddConfig) (err error) {
	sqlStr := `INSERT INTO t_global_config(config_key, config_value, config_explain) VALUES(?, ?, ?)`
	_, err = db.Exec(sqlStr, p.Key, p.Value, p.Explain)
	return
}

func CheckConfigExistByKey(key string) (bool, error) {
	sqlStr := `SELECT count(config_key) FROM t_global_config WHERE config_key = ?`
	var count int
	if err := db.Get(&count, sqlStr, key); err != nil {
		return false, err
	}
	return count > 0, nil
}

func UpdateConfig(p *models.ParamsUpdateConfig, u *models.UriUpdateConfig) (err error) {
	sqlStr := `UPDATE t_global_config SET config_key = ?, config_value = ?, config_explain = ? WHERE config_key = ?`
	_, err = db.Exec(sqlStr, p.Key, p.Value, p.Explain, u.Key)
	return
}

func DeleteConfig(key string) (err error) {
	sqlStr := `DELETE FROM t_global_config WHERE config_key = ?`
	_, err = db.Exec(sqlStr, key)
	if err != nil {
		return
	}
	return
}

func GetConfigByKey(key string) (config models.ConfigItem, err error) {
	sqlStr := `SELECT config_key, config_value, config_explain FROM t_global_config WHERE config_key = ?`
	err = db.Get(&config, sqlStr, key)
	return config, err

}

func GetConfigs(queryFields *queryfields.ConfigQueryFields) (configs models.ConfigList, err error) {
	sqlStr := `SELECT config_key, config_value, config_explain FROM t_global_config `

	sqlStr = dynamicConcatConfigSql(sqlStr, queryFields)
	sqlStr += ` ORDER BY id DESC `
	var limit = queryFields.Limit
	if queryFields.Page != 0 && queryFields.Limit != 0 {
		sqlStr += ` LIMIT ? OFFSET ?`
		err = db.Select(&configs.List, sqlStr, queryFields.Key, queryFields.Value, queryFields.Explain, limit, (queryFields.Page-1)*limit)
	} else {
		err = db.Select(&configs.List, sqlStr, queryFields.Key, queryFields.Value, queryFields.Explain)
	}

	if err != nil {
		return
	}

	if len(configs.List) == 0 {
		configs.List = make([]models.ConfigItem, 0, 0)
	}

	// get config count
	countSqlStr := ` SELECT COUNT(id) FROM t_global_config `
	countSqlStr = dynamicConcatConfigSql(countSqlStr, queryFields)
	err = db.Get(&configs.Page.Count, countSqlStr, queryFields.Key, queryFields.Value, queryFields.Explain)
	if err != nil {
		return
	}
	configs.Page.CurrentPage = queryFields.Page
	configs.Page.Size = queryFields.Limit

	return

}

func dynamicConcatConfigSql(sqlStr string, queryFields *queryfields.ConfigQueryFields) string {
	if queryFields.Key != "" {
		sqlStr += ` WHERE config_key = ? `
	} else {
		sqlStr += ` WHERE LENGTH(?) = 0 `
	}

	if queryFields.Value != "" {
		sqlStr += ` AND config_value LIKE CONCAT("%", ?, "%") `
	} else {
		sqlStr += ` AND LENGTH(?) = 0 `
	}

	if queryFields.Explain != "" {
		sqlStr += ` AND config_explain LIKE CONCAT("%", ?, "%") `
	} else {
		sqlStr += ` AND LENGTH(?) = 0 `
	}

	return sqlStr
}
