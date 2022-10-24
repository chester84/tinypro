package types

import (
	"fmt"
)

const (
	HtmlSelected = "selected"
	HtmlChecked  = "checked"

	Html404 = `<html>
<head><title>404 Not Found</title></head>
<body>
<center><h1>404 Not Found</h1></center>
</body>
</html>`

	Html401 = `<html>
<head><title>401 Authorization Required</title></head>
<body>
<center><h1>401 Authorization Required</h1></center>
</body>
</html>`

	Html403 = `<html>
<head><title>403 Forbidden</title></head>
<body>
<center><h1>403 Forbidden</h1></center>
</body>
</html>`
)

var sortFieldMap = map[string]bool{
	"id":     true,
	"ctime":  true,
	"heat":   true,
	"weight": true,

	"created_at": true,
	"fetched_at": true,
	"last_op_at": true,
	"modify_at":  true,
	"updated_at": true,

	"register_at": true,

	"transaction_at": true,
}

func SQLOrderBy(sortField, orderBy string) (subSQL string) {
	if orderBy != "ASC" {
		orderBy = "DESC"
	} else {
		orderBy = "ASC"
	}
	if sortFieldMap[sortField] {
		subSQL = fmt.Sprintf(`ORDER BY %s %s`, sortField, orderBy)
	}

	return
}

func AdmBoxDivConf() []string {
	return []string{
		`box-default`,
		`box-primary`,
		`box-info`,
		`box-danger`,
		`box-warning`,
		`box-success`,
	}
}
