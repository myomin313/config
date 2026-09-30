package config

import "fmt"

func DropColumns(model interface{}, columns []string) error {
	for _, column := range columns {
		if DB.Migrator().HasColumn(model, column) {
			if err := DB.Migrator().DropColumn(model, column); err != nil {
				return err
			}
			fmt.Println("Dropped column:", column)
		}
	}
	return nil
}
