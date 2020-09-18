package steps

import (
	"github.com/cucumber/messages-go/v10"
	"strings"
)

type Data struct {
	ID  string   `json:"id"`
	Key []string `json:"key"`
}

func NewDataFrom(row *messages.PickleStepArgument_PickleTable_PickleTableRow) Data {
	if row == nil {
		return Data{}
	}

	cells := row.GetCells()
	id := cells[0].Value
	key := strings.Split(cells[1].Value, ",")

	return Data{
		ID:  id,
		Key: key,
	}
}

func NewDataSliceFrom(table *messages.PickleStepArgument_PickleTable) []Data {
	if table == nil {
		return nil
	}

	result := make([]Data, len(table.GetRows()))
	for i, row := range table.GetRows() {
		result[i] = NewDataFrom(row)
	}

	return result
}
