/*
 * Copyright (c) 2019, Tranquil Data, Inc. All rights reserved.
 */

package sqlparser

func (*InsertPostgres) iStatement() {}
func (*UpdatePostgres) iStatement() {}
func (*DeletePostgres) iStatement() {}

type InsertPostgres struct {
	Action     string
	Comments   Comments
	Ignore     string
	Table      TableName
	Partitions Partitions
	Columns    Columns
	Rows       InsertRows
	OnDup      OnDup
	Returning  Returning
}

func (node *InsertPostgres) Format(buf *TrackedBuffer) {
	buf.astPrintf(node, "%s %v%sinto %v%v%v %v%v%v",
		node.Action,
		node.Comments, node.Ignore,
		node.Table, node.Partitions, node.Columns, node.Rows, node.OnDup,
		node.Returning)
}

func (node *InsertPostgres) walkSubtree(visit Visit) error {
	if node == nil {
		return nil
	}
	return Walk(
		visit,
		node.Comments,
		node.Table,
		node.Columns,
		node.Rows,
		node.OnDup,
		node.Returning,
	)
}

type UpdatePostgres struct {
	Comments   Comments
	Ignore     string
	TableExprs TableExprs
	Exprs      UpdateExprs
	Where      *Where
	OrderBy    OrderBy
	Limit      *Limit
	Returning  Returning
}

// Format formats the node.
func (node *UpdatePostgres) Format(buf *TrackedBuffer) {
	buf.astPrintf(node, "update %v%s%v set %v%v%v%v%v",
		node.Comments, node.Ignore, node.TableExprs,
		node.Exprs, node.Where, node.OrderBy, node.Limit,
		node.Returning)
}

type DeletePostgres struct {
	Comments   Comments
	Targets    TableNames
	TableExprs TableExprs
	Partitions Partitions
	Where      *Where
	OrderBy    OrderBy
	Limit      *Limit
	Returning  Returning
}

// Format formats the node.
func (node *DeletePostgres) Format(buf *TrackedBuffer) {
	buf.astPrintf(node, "delete %v", node.Comments)
	if node.Targets != nil {
		buf.astPrintf(node, "%v ", node.Targets)
	}
	buf.astPrintf(node, "from %v%v%v%v%v%v", node.TableExprs, node.Partitions,
		node.Where, node.OrderBy, node.Limit, node.Returning)
}

// Returning represents a RETURNING clause.
type Returning SelectExprs

// Format formats the node.
func (node Returning) Format(buf *TrackedBuffer) {
	if node == nil {
		return
	}
	buf.astPrintf(node, " returning %v", SelectExprs(node))
}

// JoinTableExpr.Join
const (
	// The following terms are Added for handling Posgres grammar. Some of them
	// are not needed (for example CROSS JOIN is the same as INNER JOIN), but
	//are included to render the query as close to the original one as possible.
	FullJoinStr       = "full join"
	FullOuterJoinStr  = "full outer join"
	LeftOuterJoinStr  = "left outer join"
	RightOuterJoinStr = "right outer join"
	InnerJoinStr      = "inner join"
	CrossJoinStr      = "cross join"
)

// Select.Distinct
const (
	AllStr = "all "
)

func (*ParenExpr) iExpr() {}

// ParenExpr represents a parenthesized boolean expression.
type ParenExpr struct {
	Expr Expr
}

// Format formats the node.
func (node *ParenExpr) Format(buf *TrackedBuffer) {
	buf.astPrintf(node, "(%v)", node.Expr)
}
