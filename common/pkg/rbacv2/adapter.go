package rbacv2

import (
	"fmt"
	"runtime"

	_ "tinypro/common/lib/db/mysql"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	m2 "github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"

	"tinypro/common/models"
)

// Adapter represents the Xorm adapter for policy storage.
type Adapter struct {
	//driverName     string
	//dataSourceName string
	//dbSpecified    bool
	o orm.Ormer
}

// finalizer is the destructor for Adapter.
func finalizer(a *Adapter) {
}

// NewAdapter is the constructor for Adapter.
// dbSpecified is an optional bool parameter. The default value is false.
// It's up to whether you have specified an existing DB in dataSourceName.
// If dbSpecified == true, you need to make sure the DB in dataSourceName exists.
// If dbSpecified == false, the adapter will automatically create a DB named "casbin".
func NewAdapter() *Adapter {
	a := &Adapter{}
	//a.driverName = driverName
	//a.dataSourceName = dataSourceName
	//
	//if len(dbSpecified) == 0 {
	//	a.dbSpecified = false
	//} else if len(dbSpecified) == 1 {
	//	a.dbSpecified = dbSpecified[0]
	//} else {
	//	panic(errors.New("invalid parameter: dbSpecified"))
	//}

	a.o = orm.NewOrm()

	// Call the destructor when the object is released.
	runtime.SetFinalizer(a, finalizer)

	return a
}

func loadPolicyLine(line models.CasbinRule, model m2.Model) {
	lineText := line.PType
	if line.V0 != "" {
		lineText += ", " + line.V0
	}
	if line.V1 != "" {
		lineText += ", " + line.V1
	}
	if line.V2 != "" {
		lineText += ", " + line.V2
	}
	if line.V3 != "" {
		lineText += ", " + line.V3
	}
	if line.V4 != "" {
		lineText += ", " + line.V4
	}
	if line.V5 != "" {
		lineText += ", " + line.V5
	}

	persist.LoadPolicyLine(lineText, model)
}

// LoadPolicy loads policy from database.
func (a *Adapter) LoadPolicy(model m2.Model) error {
	var page = 0
	var pSize = 1000
	var lines []models.CasbinRule
	var total int64
	var err error
	for {
		var t []models.CasbinRule
		total, err = a.o.QueryTable(models.CASBIN_RULE_TABLENAME).Offset(page * pSize).Limit(pSize).All(&t)
		logs.Info("t: %v total: %v err: %v", len(t), total, err)
		if err != nil && err != orm.ErrNoRows {
			return err
		}

		page++
		lines = append(lines, t...)

		if total == 0 || total < int64(pSize) {
			break
		}
	}

	for _, line := range lines {
		loadPolicyLine(line, model)
	}

	return nil
}

func savePolicyLine(pType string, rule []string) models.CasbinRule {
	line := models.CasbinRule{}

	line.PType = pType
	if len(rule) > 0 {
		line.V0 = rule[0]
	}
	if len(rule) > 1 {
		line.V1 = rule[1]
	}
	if len(rule) > 2 {
		line.V2 = rule[2]
	}
	if len(rule) > 3 {
		line.V3 = rule[3]
	}
	if len(rule) > 4 {
		line.V4 = rule[4]
	}
	if len(rule) > 5 {
		line.V5 = rule[5]
	}

	return line
}

// SavePolicy saves policy to database.
func (a *Adapter) SavePolicy(model m2.Model) error {
	var lines []models.CasbinRule

	for pType, ast := range model["p"] {
		for _, rule := range ast.Policy {
			line := savePolicyLine(pType, rule)
			lines = append(lines, line)
		}
	}

	for pType, ast := range model["g"] {
		for _, rule := range ast.Policy {
			line := savePolicyLine(pType, rule)
			lines = append(lines, line)
		}
	}

	tx, _ := a.o.Begin()
	_, err := tx.Raw(fmt.Sprintf("TRUNCATE %s", models.CASBIN_RULE_TABLENAME)).Exec()
	if err != nil {
		_ = tx.Rollback()
		logs.Error("[SavePolicy] TRUNCATE err: %v", err)
		return err
	}

	_, err = tx.InsertMulti(len(lines), lines)
	if err != nil {
		_ = tx.Rollback()
		logs.Error("[SavePolicy] InsertMulti err: %v", err)
		return err
	}
	_ = tx.Commit()

	return err
}

// AddPolicy adds a policy rule to the storage.
func (a *Adapter) AddPolicy(sec string, pType string, rule []string) error {
	line := savePolicyLine(pType, rule)
	_, err := a.o.Insert(&line)
	return err
}

// RemovePolicy removes a policy rule from the storage.
func (a *Adapter) RemovePolicy(sec string, pType string, rule []string) error {
	line := savePolicyLine(pType, rule)
	_, err := a.o.Delete(&line, "p_type", "v0", "v1", "v2", "v3", "v4", "v5")
	return err
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
func (a *Adapter) RemoveFilteredPolicy(sec string, pType string, fieldIndex int, fieldValues ...string) error {
	line := models.CasbinRule{}

	line.PType = pType
	filter := []string{}
	filter = append(filter, "p_type")
	if fieldIndex <= 0 && 0 < fieldIndex+len(fieldValues) {
		line.V0 = fieldValues[0-fieldIndex]
		filter = append(filter, "v0")
	}
	if fieldIndex <= 1 && 1 < fieldIndex+len(fieldValues) {
		line.V1 = fieldValues[1-fieldIndex]
		filter = append(filter, "v1")
	}
	if fieldIndex <= 2 && 2 < fieldIndex+len(fieldValues) {
		line.V2 = fieldValues[2-fieldIndex]
		filter = append(filter, "v2")
	}
	if fieldIndex <= 3 && 3 < fieldIndex+len(fieldValues) {
		line.V3 = fieldValues[3-fieldIndex]
		filter = append(filter, "v3")
	}
	if fieldIndex <= 4 && 4 < fieldIndex+len(fieldValues) {
		line.V4 = fieldValues[4-fieldIndex]
		filter = append(filter, "v4")
	}
	if fieldIndex <= 5 && 5 < fieldIndex+len(fieldValues) {
		line.V5 = fieldValues[5-fieldIndex]
		filter = append(filter, "v5")
	}

	_, err := a.o.Delete(&line, filter...)
	return err
}
