package arguments

import (
	"fmt"
	"strings"
)

type Configs interface {
	// This function should return the expected arguments for the programme
	Expected() map[string]string

	// This should handle the parsed arguments, however required
	SetInputs(map[string][]string)
}

func Inputs(c Configs, inputs []string) error {
	args := c.Expected()
	result := make(map[string][]string)
	var err error
	var curIndex string
	for i := 0; i < len(inputs); i++ {
		v := strings.Trim(inputs[i], " ")
		if strings.Contains(v, "=") {
			err = getSplitConfigs(v, args, &result, &curIndex)
		} else {
			err = appendConfig(&curIndex, args, &result, v)
		}
		if err != nil {
			return err
		}
	}
	c.SetInputs(result)
	return nil
}

func getSplitConfigs(v string, args map[string]string, result *map[string][]string, curIndex *string) error {
	configs := strings.Split(v, "=")
	var err error
	var errs []error
	for _, c := range configs {
		c = strings.Trim(c, " ")
		err = appendConfig(curIndex, args, result, c)
		errs = append(errs, err)
	}
	return GetErrors(errs)
}

func appendConfig(curIndex *string, args map[string]string, result *map[string][]string, arg string) error {
	removeDashes(&arg)
	if args[arg] == "" {
		if (*result)[*curIndex] == nil {
			err := fmt.Errorf("the %s argument does not exist", arg)
			return err
		}
		(*result)[*curIndex] = append((*result)[*curIndex], arg)
	} else {
		*curIndex = args[arg]
		(*result)[*curIndex] = make([]string, 0)
	}
	return nil
}

func removeDashes(input *string) {
	result := strings.Replace(*input, "-", "", 2)
	*input = result
}

func GetErrors(errs []error) error {
	var errStrings []string
	if len(errs) > 0 {
		for _, e := range errs {
			if e != nil {
				errStrings = append(errStrings, e.Error())
			}
		}
	}
	if len(errStrings) > 0 {
		return fmt.Errorf(strings.Join(errStrings, "; "))
	}
	return nil
}