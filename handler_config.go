package main

import (
	"errors"
	"fmt"
	"github.com/vcoder4c/jirakits/cli"
	"github.com/vcoder4c/jirakits/google"
	"github.com/vcoder4c/jirakits/jira"
)

func isJiraReady() string {
	setting, err := jira.LoadSetting(projectDir())
	if err != nil {
		setting = &jira.Setting{}
		_ = jira.SaveSetting(projectDir(), setting)
	}
	if setting.URL == "" || setting.Token == "" {
		return fmt.Sprintf("Current setting: %s is not valid", setting.Print())
	}
	return ""
}

func checkJiraConfig() {
	msg := isJiraReady()
	if msg != "" {
		exitF(msg)
	}
	fmt.Println("Setting JIRA OK")
}

func configJira(ctx cli.Context) {
	args := ctx.Args()
	isChecking := args.Bool("check")
	if isChecking {
		checkJiraConfig()
		return
	}
	setting, _ := jira.LoadSetting(projectDir())
	if setting == nil {
		setting = &jira.Setting{}
	}
	if args.String("url") != "" {
		setting.URL = args.String("url")
	}
	if args.String("token") != "" {
		setting.Token = args.String("token")
	}
	err := jira.SaveSetting(projectDir(), setting)
	if err != nil {
		exit(err.Error())
	}
}

func isDriveReady() string {
	err := google.IsReady(projectDir())
	if errors.Is(err, google.CredentialIsNotExisted) {
		return fmt.Sprintf("%s: %s", err.Error(), google.GetCredentialPath(projectDir()))
	}
	if errors.Is(err, google.TokenIsNotExisted) || errors.Is(err, google.TokenIsInvalid) {
		return fmt.Sprintf("Please setup drive with command: ./%s drive --setup", Name)
	}
	if err != nil {
		return err.Error()
	}
	return ""
}

func checkDriveConfig() {
	msg := isDriveReady()
	if msg != "" {
		exit(msg)
	}
	fmt.Println("Setting Drive OK")
}

func setupDriveConfig() {
	err := google.SetupDrive(projectDir())
	if errors.Is(err, google.CredentialIsNotExisted) {
		exitF("%s: %s", err.Error(), google.GetCredentialPath(projectDir()))
	}
	if err != nil {
		exitF(err.Error())
	}
}

func configDrive(ctx cli.Context) {
	args := ctx.Args()
	isChecking := args.Bool("check")
	if isChecking {
		checkDriveConfig()
		return
	}
	isSetup := args.Bool("setup")
	if isSetup {
		setupDriveConfig()
		return
	}
	printHelp(ctx)
}

func configCheck(ctx cli.Context) {
	jiraCheck := isJiraReady()
	if jiraCheck != "" {
		fmt.Printf("Jira: %s. Use './%s help jira' for more info\n", jiraCheck, Name)
	} else {
		fmt.Print("Jira: OK\n")
	}
	driveCheck := isDriveReady()
	if driveCheck != "" {
		fmt.Printf("Drive: %s\n", driveCheck)
	} else {
		fmt.Print("Drive OK\n")
	}
	_, err := getConfig()
	if err != nil {
		fmt.Printf("Export mapping: Not found at %s\n", getConfigMappingPath())
	} else {
		fmt.Print("Export mapping: OK\n")
	}
}
