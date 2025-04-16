package cassandraadapter

import (
	"sort"
	"strconv"
	"strings"

	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/gocql/gocql"
)

// Adapter represents the Cassandra adapter for policy storage.
type Adapter struct {
	session   *gocql.Session
	tableName string
}

// NewAdapter is the constructor for Adapter.
func NewAdapter(session *gocql.Session) *Adapter {
	return &Adapter{
		session:   session,
		tableName: "casbin_policy",
	}
}

func (a *Adapter) close() {
	a.session.Close()
}

func (a *Adapter) createTable() {
	query := `CREATE TABLE IF NOT EXISTS ` + a.tableName + ` (
		no text PRIMARY KEY,
		ptype text,
		v1 text,
		v2 text,
		v3 text,
		v4 text
	)`
	if err := a.session.Query(query).Exec(); err != nil {
		panic(err)
	}
}

func (a *Adapter) dropTable() {
	if err := a.session.Query("DROP TABLE IF EXISTS " + a.tableName).Exec(); err != nil {
		panic(err)
	}
}

func loadPolicyLine(line string, model model.Model) {
	if line == "" {
		return
	}
	tokens := strings.Split(line, ", ")
	key := tokens[0]
	sec := key[:1]
	model[sec][key].Policy = append(model[sec][key].Policy, tokens[1:])
}

// LoadPolicy loads policy from database.
func (a *Adapter) LoadPolicy(model model.Model) error {
	var (
		no    string
		ptype string
		v1    string
		v2    string
		v3    string
		v4    string
	)

	lines := make(map[int]string)
	iter := a.session.Query(`SELECT no, ptype, v1, v2, v3, v4 FROM ` + a.tableName).Iter()
	for iter.Scan(&no, &ptype, &v1, &v2, &v3, &v4) {
		line := ptype
		if v1 != "" {
			line += ", " + v1
		}
		if v2 != "" {
			line += ", " + v2
		}
		if v3 != "" {
			line += ", " + v3
		}
		if v4 != "" {
			line += ", " + v4
		}
		i, _ := strconv.Atoi(no)
		lines[i] = line
	}

	var keys []int
	for k := range lines {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		loadPolicyLine(lines[k], model)
	}

	return iter.Close()
}

func (a *Adapter) writeTableLine(no int, ptype string, rule []string) error {
	values := "'" + strconv.Itoa(no) + "','" + ptype + "'"
	for i := range rule {
		values += ",'" + rule[i] + "'"
	}
	for i := 0; i < 4-len(rule); i++ {
		values += ",''"
	}
	query := "INSERT INTO " + a.tableName + " (no, ptype, v1, v2, v3, v4) VALUES(" + values + ")"
	return a.session.Query(query).Exec()
}

// SavePolicy saves policy to database.
func (a *Adapter) SavePolicy(model model.Model) error {
	a.dropTable()
	a.createTable()

	no := 0
	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			if err := a.writeTableLine(no, ptype, rule); err != nil {
				return err
			}
			no++
		}
	}

	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			if err := a.writeTableLine(no, ptype, rule); err != nil {
				return err
			}
			no++
		}
	}
	return nil
}

// AddPolicy adds a policy rule to the storage.
func (a *Adapter) AddPolicy(sec string, ptype string, rule []string) error {
	// Generate a new 'no' value
	var maxNo string
	var current int
	iter := a.session.Query("SELECT no FROM " + a.tableName).Iter()
	for iter.Scan(&maxNo) {
		if val, err := strconv.Atoi(maxNo); err == nil && val > current {
			current = val
		}
	}
	no := current + 1
	return a.writeTableLine(no, ptype, rule)
}

// RemovePolicy removes a policy rule from the storage.
func (a *Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	whereClause := "ptype = '" + ptype + "'"
	for i, val := range rule {
		whereClause += " AND v" + strconv.Itoa(i+1) + " = '" + val + "'"
	}
	return a.session.Query("DELETE FROM " + a.tableName + " WHERE " + whereClause).Exec()
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
func (a *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	whereClause := "ptype = '" + ptype + "'"
	for i, val := range fieldValues {
		if val == "" {
			continue
		}
		if i+fieldIndex > 3 {
			break
		}
		whereClause += " AND v" + strconv.Itoa(i+fieldIndex+1) + " = '" + val + "'"
	}
	return a.session.Query("DELETE FROM " + a.tableName + " WHERE " + whereClause).Exec()
}

// Ensure Adapter implements persist.Adapter
var _ persist.Adapter = (*Adapter)(nil)
