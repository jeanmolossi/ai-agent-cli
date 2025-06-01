package console

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strings"

	"github.com/urfave/cli/v3"
	"github.com/xrash/smetrics"
)

func commandNotFound(_ context.Context, cmd *cli.Command, command string) {
	var (
		msgTxt     = fmt.Sprintf("Command '%s' is not defined.", command)
		suggestion string
	)

	if alternatives := findAlternatives(command, func() (collection []string) {
		for i := range cmd.Commands {
			collection = append(collection, cmd.Commands[i].Names()...)
		}

		return
	}()); len(alternatives) > 0 {
		if len(alternatives) == 1 {
			msgTxt = msgTxt + " Did you mean this?"
		} else {
			msgTxt = msgTxt + " Did you mean one of these?"
		}
		suggestion = "\n  " + strings.Join(alternatives, "\n  ")
	}

	slog.Error(msgTxt)
	slog.Info(suggestion)
}

func findAlternatives(name string, collection []string) (result []string) {
	var (
		threshold       = 1e3
		alternatives    = make(map[string]float64)
		collectionParts = make(map[string][]string)
	)

	for i := range collection {
		collectionParts[collection[i]] = strings.Split(collection[i], ":")
	}

	for i, sub := range strings.Split(name, ":") {
		for collectionName, parts := range collectionParts {
			exists := alternatives[collectionName] != 0

			if len(parts) <= i {
				if exists {
					alternatives[collectionName] += threshold
				}

				continue
			}

			lev := smetrics.WagnerFischer(sub, parts[i], 1, 1, 1)

			if float64(lev) <= float64(len(sub))/3 || strings.Contains(parts[i], sub) {
				if exists {
					alternatives[collectionName] += float64(lev)
				} else {
					alternatives[collectionName] = float64(lev)
				}
			} else if exists {
				alternatives[collectionName] += threshold
			}
		}
	}

	for _, item := range collection {
		lev := smetrics.WagnerFischer(name, item, 1, 1, 1)

		if float64(lev) <= float64(len(name))/3 || strings.Contains(item, name) {
			if alternatives[item] != 0 {
				alternatives[item] -= float64(lev)
			} else {
				alternatives[item] = float64(lev)
			}
		}
	}

	type scoredItem struct {
		name  string
		score float64
	}

	var sortedAlternatives []scoredItem

	for item, score := range alternatives {
		if score < 2*threshold {
			sortedAlternatives = append(sortedAlternatives, scoredItem{item, score})
		}
	}

	sort.Slice(sortedAlternatives, func(i, j int) bool {
		if sortedAlternatives[i].score == sortedAlternatives[j].score {
			return sortedAlternatives[i].name < sortedAlternatives[j].name
		}

		return sortedAlternatives[i].score < sortedAlternatives[j].score
	})

	for _, item := range sortedAlternatives {
		result = append(result, item.name)
	}

	return result
}

func onUsageError(_ context.Context, _ *cli.Command, err error, _ bool) error {
	if flag, ok := strings.CutPrefix(err.Error(), "flag provided but not defined: -"); ok {
		slog.Error("The '%s' option does not exist.", flag)
		return nil
	}

	if flag, ok := strings.CutPrefix(err.Error(), "flag needs an argument: "); ok {
		slog.Error("The '%s' option requires a value.", flag)
		return nil
	}

	if errMsg := err.Error(); strings.HasPrefix(errMsg, "invalid value") && strings.Contains(errMsg, "for flag -") {
		var value, flag string
		if _, parseErr := fmt.Sscanf(errMsg, "invalid value %q for flag -%s", &value, &flag); parseErr == nil {
			slog.Error("Invalid value '%s' for option '%s'.", value, strings.TrimSuffix(flag, ":"))
			return nil
		}
	}

	return err
}
