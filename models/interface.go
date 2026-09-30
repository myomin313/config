package models

type ColumnDropper interface {
	DropColumns() []string
}
