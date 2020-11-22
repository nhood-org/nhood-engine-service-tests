package steps

import (
	"github.com/cucumber/messages-go/v10"
	"strings"
)

type Data struct {
	UUID      string   `json:"uuid"`
	Key       []string `json:"key"`
	Reference string   `json:"reference"`
}

func NewDataFrom(row *messages.PickleStepArgument_PickleTable_PickleTableRow) Data {
	if row == nil {
		return Data{}
	}

	cells := row.GetCells()
	key := strings.Split(cells[0].Value, ",")
	reference := cells[1].Value

	return Data{
		Key:       key,
		Reference: reference,
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
