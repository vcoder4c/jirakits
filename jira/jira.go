package jira

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/thedevsaddam/gojsonq/v2"
)

func GetJiraClient(projectDir string) (*jira.Client, error) {
	setting, err := LoadSetting(projectDir)
	if err != nil {
		return nil, err
	}
	if setting.URL == "" || setting.Token == "" {
		return nil, errors.New("Jira URL and Token must be configured")
	}
	tp := jira.BearerAuthTransport{Token: setting.Token}
	jiraClient, err := jira.NewClient(tp.Client(), setting.URL)
	if err != nil {
		return nil, err
	}
	return jiraClient, nil
}

func GetIssuesCount(client *jira.Client, searchString string) (int, error) {
	last := 0
	opt := &jira.SearchOptions{
		MaxResults: 10, // Max results can go up to 1000
		StartAt:    last,
	}
	_, resp, err := client.Issue.Search(searchString, opt)
	if err != nil {
		return 0, err
	}

	return resp.Total, nil
}

// getAllIssues will implement pagination of api and get all the issues.
// Jira API has limitation as to maxResults it can return at one time.
// You may have use case where you need to get all the issues according to jql
// This is where this example comes in.
func GetAllIssues(client *jira.Client, searchString string) ([]jira.Issue, error) {
	last := 0
	var issues []jira.Issue
	for {
		opt := &jira.SearchOptions{
			MaxResults: 1000, // Max results can go up to 1000
			StartAt:    last,
		}

		chunk, resp, err := client.Issue.Search(searchString, opt)
		if err != nil {
			return nil, err
		}

		total := resp.Total
		if issues == nil {
			issues = make([]jira.Issue, 0, total)
		}
		issues = append(issues, chunk...)
		last = resp.StartAt + len(chunk)
		if last >= total {
			return issues, nil
		}
	}
}

func GetIssue(client *jira.Client, issueId string) (*jira.Issue, error) {
	issue, _, err := client.Issue.Get(issueId, nil)
	if err != nil {
		return nil, err
	}
	return issue, nil
}

func getField(data []byte, field string) interface{} {
	return gojsonq.New().FromString(string(data)).Find(field)
}

func GetValues(issue jira.Issue, keys []interface{}, keyMap map[string]string) ([]interface{}, error) {
	data, err := json.Marshal(issue)
	if err != nil {
		return nil, err
	}
	result := make([]interface{}, len(keys))
	for idx, key := range keys {
		field := keyMap[fmt.Sprintf("%v", key)]
		value := getField(data, field)
		if value == nil {
			result[idx] = "N/A"
		} else {
			result[idx] = fmt.Sprintf("%v", value)
		}
	}
	return result, nil
}
