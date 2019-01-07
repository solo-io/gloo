package cliutil

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/AlecAivazis/survey.v1"
)

func GetYesInput(msg string) (bool, error) {
	var yesAnswer string

	if err := GetStringInputDefault(
		msg,
		&yesAnswer,
		"N",
	); err != nil {
		return false, err
	}
	return strings.ToLower(yesAnswer) != "y", nil
}

func GetStringInput(msg string, value *string) error {
	return GetStringInputDefault(msg, value, "")
}

func GetStringInputDefault(msg string, value *string, defaultValue string) error {
	prompt := &survey.Input{Message: msg, Default: defaultValue}
	if err := survey.AskOne(prompt, value, nil); err != nil {
		return err
	}
	return nil
}

func GetUint32Input(msg string, value *uint32) error {
	return GetUint32InputDefault(msg, value, 0)
}

func GetUint32InputDefault(msg string, value *uint32, defaultValue uint32) error {
	var strValue string
	prompt := &survey.Input{Message: msg, Default: strconv.Itoa(int(defaultValue))}
	if err := survey.AskOne(prompt, &strValue, nil); err != nil {
		return err
	}
	val, err := strconv.Atoi(strValue)
	if err != nil {
		return err
	}
	*value = uint32(val)
	return nil
}

func GetBoolInput(msg string, value *bool) error {
	return GetBoolInputDefault(msg, value, false)
}

func GetBoolInputDefault(msg string, value *bool, defaultValue bool) error {
	var strValue string
	defaultValueStr := "N"
	if defaultValue {
		defaultValueStr = "y"
	}
	prompt := &survey.Input{Message: msg + " [y/N]: ", Default: defaultValueStr}
	if err := survey.AskOne(prompt, &strValue, nil); err != nil {
		return err
	}
	*value = strings.ToLower(defaultValueStr) == "y"
	return nil
}

func GetStringSliceInput(msg string, value *[]string) error {
	prompt := &survey.Input{Message: msg}
	var lastChoice string
	for {
		if err := survey.AskOne(prompt, &lastChoice, nil); err != nil {
			return err
		}

		if lastChoice == "" {
			break
		}
		*value = append(*value, lastChoice)
	}
	return nil
}

func ChooseFromList(message string, choice *string, options []string) error {
	if len(options) == 0 {
		return fmt.Errorf("No options to select from (for prompt: %v)", message)
	}

	question := &survey.Select{
		Message: message,
		Options: options,
	}

	if err := survey.AskOne(question, choice, survey.Required); err != nil {
		// this should not error
		fmt.Println("error with input")
		return err
	}

	return nil
}

func ChooseBool(message string, target *bool) error {

	yes, no := "yes", "no"

	question := &survey.Select{
		Message: message,
		Options: []string{yes, no},
	}

	var choice string
	if err := survey.AskOne(question, &choice, survey.Required); err != nil {
		return err
	}

	*target = choice == yes
	return nil
}
