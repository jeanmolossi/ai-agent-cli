package carbon

import (
	"github.com/dromara/carbon/v2"
)

type Carbon = carbon.Carbon

func SetTimezone(timezone string) {
	carbon.SetTimezone(timezone)
}
