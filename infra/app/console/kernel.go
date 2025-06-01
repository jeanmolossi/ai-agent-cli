package console

import "github.com/jeanmolossi/ai-agent-cli/app/contracts/console"

type Kernel struct{}

func (k Kernel) Commands() []console.Command {
	return []console.Command{}
}

// func (k Kernel) Schedule() []schedule.Event{
//     return []schedule.Event{}
// }
