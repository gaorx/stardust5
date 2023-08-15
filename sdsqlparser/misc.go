package sdsqlparser

import (
	"github.com/blastrain/vitess-sqlparser/sqlparser"
)

func SqlWithLimit(selectSql string, limit, offset string) (string, bool) {
	stmt0, err := sqlparser.Parse(selectSql)
	if err != nil {
		return "", false
	}
	stmt, ok := stmt0.(*sqlparser.Select)
	if !ok {
		return "", false
	}
	stmt.SetLimit(&sqlparser.Limit{
		Rowcount: sqlparser.NewValArg([]byte(limit)),
		Offset:   sqlparser.NewValArg([]byte(offset)),
	})
	return sqlparser.String(stmt), true
}

func SqlForCount(selectSql string) (string, bool) {
	stmt0, err := sqlparser.Parse(selectSql)
	if err != nil {
		return "", false
	}
	stmt, ok := stmt0.(*sqlparser.Select)
	if !ok {
		return "", false
	}
	stmt.SelectExprs = sqlparser.SelectExprs{
		&sqlparser.AliasedExpr{
			Expr: &sqlparser.FuncExpr{
				Name: sqlparser.NewColIdent("COUNT"),
				Exprs: sqlparser.SelectExprs{
					&sqlparser.StarExpr{},
				},
			},
			As: sqlparser.NewColIdent(""),
		},
	}
	return sqlparser.String(stmt), true
}
