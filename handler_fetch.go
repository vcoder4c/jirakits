package main

import (
	"errors"
	"fmt"
	"github.com/vcoder4c/jirakits/cli"
	"github.com/vcoder4c/jirakits/common"
	"github.com/vcoder4c/jirakits/google"
	"github.com/vcoder4c/jirakits/jira"
	"github.com/manifoldco/promptui"
	"google.golang.org/api/sheets/v4"
)

func isCreatingNew() bool {
	validate := func(input string) error {
		if common.StringInSlice(input, []string{"yes", "no"}) {
			return nil
		}
		return errors.New("Only allow yes or no")
	}
	prompt := promptui.Prompt{
		Label:    "Create new file? [yes|no]",
		Validate: validate,
	}
	result, err := prompt.Run()
	if err != nil {
		exit(err.Error())
	}
	return result == "yes"
}

func getQueryString() string {
	prompt := promptui.Prompt{
		Label:    "Your query string",
		Validate: nil,
	}
	result, err := prompt.Run()
	if err != nil {
		exit(err.Error())
	}
	return result
}

func confirmProceed(total int) bool {
	validate := func(input string) error {
		if common.StringInSlice(input, []string{"yes", "no"}) {
			return nil
		}
		return errors.New("Only allow yes or no")
	}
	prompt := promptui.Prompt{
		Label:    fmt.Sprintf("There are %v issues. Confirm to proceed? [yes|no]", total),
		Validate: validate,
	}
	result, err := prompt.Run()
	if err != nil {
		exit(err.Error())
	}
	return result == "yes"
}

func getSpreadSheetName() string {
	prompt := promptui.Prompt{
		Label:    "Name of your spreadsheet:",
		Validate: nil,
	}
	result, err := prompt.Run()
	if err != nil {
		exit(err.Error())
	}
	return result
}

func getSheetName() string {
	prompt := promptui.Prompt{
		Label:    fmt.Sprintf("Name of your spreadsheet: [%s]", common.GetCurrentDate()),
		Validate: nil,
	}
	result, err := prompt.Run()
	if err != nil {
		exit(err.Error())
	}
	if result == "" {
		return common.GetCurrentDate()
	}
	return result
}

func getExistedSpreadSheet(service *sheets.Service) *sheets.Spreadsheet {
	validate := func(input string) error {
		_, err := google.GetSpreadSheet(service, input)
		return err
	}
	prompt := promptui.Prompt{
		Label:    "Provide your spreadsheet id",
		Validate: validate,
	}
	result, err := prompt.Run()
	if err != nil {
		exit(err.Error())
	}
	spreadSheet, err := google.GetSpreadSheet(service, result)
	if err != nil {
		exit(err.Error())
	}
	return spreadSheet
}

func fetchIssues(ctx cli.Context) {
	msg := isJiraReady()
	if msg != "" {
		exit(msg)
	}
	msg = isDriveReady()
	if msg != "" {
		exit(msg)
	}
	jiraClient, err := jira.GetJiraClient(projectDir())
	if err != nil {
		exit(err.Error())
	}
	queryString := getQueryString()
	total, err := jira.GetIssuesCount(jiraClient, queryString)

	confirmYes := confirmProceed(total)
	if !confirmYes {
		return
	}

	sheetService, err := google.GetSheetService(projectDir())
	if err != nil {
		exit(err.Error())
	}
	var spreadSheet *sheets.Spreadsheet = nil
	createNew := isCreatingNew()
	if createNew {
		spreadSheetName := getSpreadSheetName()
		spreadSheet, err = google.CreateSpreadSheet(sheetService, spreadSheetName)
		if err != nil {
			exit(err.Error())
		}
	} else {
		spreadSheet = getExistedSpreadSheet(sheetService)
	}

	sheetName := getSheetName()
	// Try to create sheet name. Ignore if it's existed
	_ = google.CreateSheetOnSpreadSheet(sheetService, spreadSheet, sheetName)

	configMap, err := getConfig()
	if err != nil {
		exit(err.Error())
	}
	issues, err := jira.GetAllIssues(jiraClient, queryString)
	if err != nil {
		exit(err.Error())
	}

	rows := make([][]interface{}, 0)
	configKeys := common.GetKeys(configMap)
	rows = append(rows, configKeys)
	for _, issue := range issues {
		values, err := jira.GetValues(issue, configKeys, configMap)
		if err != nil {
			exit(err.Error())
		}
		rows = append(rows, values)
	}
	err = google.WriteToSheet(sheetService, spreadSheet, sheetName, "A1:Z10000", rows)
	if err != nil {
		exit(err.Error())
	}
	fmt.Println("Completed successfully")
}
