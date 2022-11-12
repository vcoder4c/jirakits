package google

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"os"
)

func getDriveConfig(projectDir string) (*oauth2.Config, error) {
	b, err := os.ReadFile(GetCredentialPath(projectDir))
	if err != nil {
		return nil, err
	}
	config, err := google.ConfigFromJSON(b, drive.DriveScope, sheets.SpreadsheetsScope)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func cleanToken(projectDir string) {
	_ = os.Remove(getTokenPath(projectDir))
}

func GetDriveService(projectDir string) (*drive.Service, error) {
	oauth2Config, err := getDriveConfig(projectDir)
	if err != nil {
		return nil, err
	}
	client := getClient(projectDir, oauth2Config)
	ctx := context.Background()
	service, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}
	return service, nil
}

func GetSheetService(projectDir string) (*sheets.Service, error) {
	oauth2Config, err := getDriveConfig(projectDir)
	if err != nil {
		return nil, err
	}
	client := getClient(projectDir, oauth2Config)
	ctx := context.Background()
	service, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}
	return service, nil
}

func SetupDrive(projectDir string) error {
	err := IsReady(projectDir)
	if errors.Is(err, CredentialIsNotExisted) {
		return err
	}
	cleanToken(projectDir)
	_, err = GetDriveService(projectDir)
	return err
}

func GetSpreadSheet(service *sheets.Service, spreadSheetId string) (*sheets.Spreadsheet, error) {
	req := service.Spreadsheets.Get(spreadSheetId)
	return req.Do()
}

func CreateSpreadSheet(service *sheets.Service, name string) (*sheets.Spreadsheet, error) {
	t := sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{
			Title: name,
		},
	}
	req := service.Spreadsheets.Create(&t)
	spreadSheet, err := req.Do()
	if err != nil {
		return nil, err
	}
	return spreadSheet, nil
}

func CreateSheetOnSpreadSheet(service *sheets.Service, spreadSheet *sheets.Spreadsheet, sheetName string) error {
	request := sheets.Request{
		AddSheet: &sheets.AddSheetRequest{
			Properties: &sheets.SheetProperties{
				Title: sheetName,
			},
		},
	}
	req := service.Spreadsheets.BatchUpdate(spreadSheet.SpreadsheetId, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{&request},
	})
	_, err := req.Do()
	return err
}

func WriteToSheet(service *sheets.Service, spreadSheet *sheets.Spreadsheet, sheetName string, dataRange string, rows [][]interface{}) error {
	req := &sheets.BatchUpdateValuesRequest{
		ValueInputOption: "USER_ENTERED",
	}
	rangeData := fmt.Sprintf("%s!%s", sheetName, dataRange)
	req.Data = append(req.Data, &sheets.ValueRange{
		Range:  rangeData,
		Values: rows,
	})
	_, err := service.Spreadsheets.Values.BatchUpdate(spreadSheet.SpreadsheetId, req).Do()
	return err
}
