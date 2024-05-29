package messages

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

var Messages *Templates

type Templates struct {
	ElectricityStatusOn     string   `json:"electricity_status_on"`
	ElectricityStatusOff    string   `json:"electricity_status_off"`
	ElectricityTurnOn       string   `json:"electricity_turn_on"`
	ElectricityTurnOff      string   `json:"electricity_turn_off"`
	ScheduleMessage         string   `json:"schedule_message"`
	ScheduleMessageContinue string   `json:"schedule_message_continue"`
	MinutesForms            []string `json:"minutes_forms"`
	HoursForms              []string `json:"hours_forms"`
	And                     string   `json:"and"`
}

func LoadMessages(localePath string) error {
	raw, err := os.ReadFile(localePath)
	err = json.Unmarshal(raw, &Messages)
	if err != nil {
		return err
	}

	return validate()
}

func validate() error {
	v := reflect.ValueOf(*Messages)
	tType := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := tType.Field(i).Name

		if field.IsZero() {
			return fmt.Errorf("field %s is empty", fieldName)
		}
	}

	return nil
}
