package gformscreator

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/viggneshvn/reddotstudios_contracts_backend/internal/contract"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/forms/v1"
	"google.golang.org/api/option"
)

type DriveService struct {
	*drive.Service
}

type FormsService struct {
	*forms.Service
}

func NewFormsService() (*FormsService, error) {
	// Load the service account credentials from the JSON file.
	ctx := context.Background()
	sa := option.WithCredentialsFile("/etc/secrets/credentials.json")

	formsService, err := forms.NewService(ctx, sa)
	if err != nil {
		return nil, fmt.Errorf("failed while creating a new service with err : %w", err)
	}

	return &FormsService{formsService}, nil
}

func NewDriveService() (*DriveService, error) {
	// Load the service account credentials from the JSON file.
	ctx := context.Background()
	sa := option.WithCredentialsFile("/etc/secrets/credentials.json")

	driveService, err := drive.NewService(ctx, sa)
	if err != nil {
		return nil, fmt.Errorf("failed while creating a new service with err : %w", err)
	}

	return &DriveService{driveService}, nil
}

func (driveService *DriveService) uploadImageToDrive(imgPath string) (*string, error) {

	// Read the image file's content.
	fileContent, err := ioutil.ReadFile(imgPath)
	if err != nil {
		return nil, fmt.Errorf("failed while reading image file with err : %w", err)
	}

	// Create a new Google Drive file.
	driveFile := &drive.File{
		Name:     filepath.Base(imgPath),
		Parents:  []string{"1UX-0xXQPRbV5aj1G_NNX06gyODgakvQP"},
		MimeType: "image/jpeg",
	}

	// Upload the image to Google Drive.
	uploadedFile, err := driveService.Files.Create(driveFile).Media(bytes.NewReader(fileContent)).Do()
	if err != nil {
		return nil, fmt.Errorf("failed while creating image file in drive with err : %w", err)
	}

	// Set the permissions for the file to be public
	permission := &drive.Permission{
		Type:               "anyone",
		Role:               "reader",
		AllowFileDiscovery: false,
	}

	_, err = driveService.Permissions.Create(uploadedFile.Id, permission).Do()
	if err != nil {
		return nil, fmt.Errorf("failed while creating public permissions for the image file in drive with err : %w", err)
	}

	file, err := driveService.Files.Get(uploadedFile.Id).Fields("webContentLink").Do()
	if err != nil {
		return nil, fmt.Errorf("failed while getting the web content link for the image file in drive with err : %w", err)
	}

	return &file.WebContentLink, nil
}

func (driveService *DriveService) moveFormFileToSharedDirectory(form *forms.Form) error {

	// Retrieve the file metadata
	formFileMetadata, err := driveService.Files.Get(form.FormId).Do()
	if err != nil {
		return fmt.Errorf("failure to get the form file metadata with error : %w", err)
	}

	// Move the form file to a different folder
	oldParents := strings.Join(formFileMetadata.Parents[:], ",")
	newParentFolderID := "1aMZeE6MnjTmtsxwD4T2Xye4sSgbIVn12" // Replace with the ID of the new parent folder
	_, err = driveService.Files.Update(formFileMetadata.Id, nil).AddParents(newParentFolderID).RemoveParents(oldParents).Do()
	if err != nil {
		return fmt.Errorf("failure to update the parent of the form with error : %w", err)
	}

	return nil
}

func (driveService *DriveService) cleanUpImageFiles(parentID string) error {
	query := fmt.Sprintf("'%s' in parents and trashed=false", parentID)
	files, err := driveService.Files.List().Q(query).Fields("files(id, name)").Do()
	if err != nil {
		return fmt.Errorf("failure to list the files in the parent directory : %s with error : %w", parentID, err)
	}

	for _, file := range files.Files {
		err := driveService.Files.Delete(file.Id).Do()
		if err != nil {
			return fmt.Errorf("failure to delete the file : %v with error : %w", file, err)
		}

	}
	return nil
}

func (driveService *DriveService) deleteFile(fileID string) error {
	err := driveService.Files.Delete(fileID).Do()
	if err != nil {
		return fmt.Errorf("failed to delete the file with ID %s: %w", fileID, err)
	}
	return nil
}

func (formsService *FormsService) CreateFormWithTitleAndDescription(contract *contract.Contract, title *string, description *string) (*forms.Form, error) {
	// Create a new Google Form
	form := &forms.Form{
		Info: &forms.Info{
			DocumentTitle: fmt.Sprintf("Contract-%s-%s", strings.ReplaceAll(contract.ClientDetails.ClientName, " ", "_"), strings.ReplaceAll(contract.EventDetails.EventName, " ", "_")),
			Title:         *title,
		},
	}

	// Insert the Google Form into the user's Google Forms account
	form, err := formsService.Forms.Create(form).Do()
	if err != nil {
		return nil, fmt.Errorf("failure to create a google form with error : %w", err)
	}

	_, err = formsService.Forms.BatchUpdate(form.FormId, &forms.BatchUpdateFormRequest{Requests: []*forms.Request{
		{
			UpdateFormInfo: &forms.UpdateFormInfoRequest{
				Info: &forms.Info{
					Description: *description,
				},
				UpdateMask: "description",
			},
		},
	}},
	).Do()
	if err != nil {
		return nil, fmt.Errorf("failure to update the description for the google form with error : %w", err)
	}

	return form, nil
}

func (formsService *FormsService) CreateImageItem(form *forms.Form, title string, sourceURI *string, index int64) error {

	_, err := formsService.Forms.BatchUpdate(form.FormId, &forms.BatchUpdateFormRequest{Requests: []*forms.Request{
		{
			CreateItem: &forms.CreateItemRequest{
				Item: &forms.Item{
					Title: title,
					ImageItem: &forms.ImageItem{
						Image: &forms.Image{
							SourceUri: *sourceURI,
						},
					},
				},
				Location: &forms.Location{Index: index, ForceSendFields: []string{"Index"}},
			},
		},
	}},
	).Do()
	if err != nil {
		return fmt.Errorf("failure to create an image item in the form with title : %s with sourceURI : %s at index %d with error : %w", title, *sourceURI, index, err)
	}

	return nil
}

func (formsService *FormsService) CreateSignatureItem(form *forms.Form, title string, index int64) error {

	_, err := formsService.Forms.BatchUpdate(form.FormId, &forms.BatchUpdateFormRequest{Requests: []*forms.Request{
		{
			CreateItem: &forms.CreateItemRequest{
				Item: &forms.Item{
					Title: title,
					QuestionItem: &forms.QuestionItem{
						Question: &forms.Question{
							TextQuestion: &forms.TextQuestion{
								Paragraph: false,
							},
							Required: true,
						},
					},
				},
				Location: &forms.Location{Index: index, ForceSendFields: []string{"Index"}},
			},
		},
	}},
	).Do()
	if err != nil {
		return fmt.Errorf("failure to create a signature item in the form with title : %s at index %d with error : %w", title, index, err)
	}

	return nil
}

func CreateGoogleForm(contract *contract.Contract) error {
	logger := logrus.New()
	var form *forms.Form
	var toBeDeletedforms []*forms.Form

	ds, err := NewDriveService()
	if err != nil {
		return fmt.Errorf("failure to a new drive service with error : %w", err)
	}

	fs, err := NewFormsService()
	if err != nil {
		return fmt.Errorf("failure to a new forms service with error : %w", err)
	}

	formTitle := "RED DOT STUDIOS SERVICES AGREEMENT"
	formDescription := fmt.Sprintf("This Agreement was made and entered into on %s between Red Dot Studios, a Massachusetts based photography and videography Service and %s(\"Client\").", time.Now().Format("01/02/2006"), contract.ClientDetails.ClientName)

	defer func() {
		if err != nil {
			logger.WithError(err).Errorf("Deleting the form since the error was not nil")
			for _, toBeDeletedform := range toBeDeletedforms {
				ds.deleteFile(toBeDeletedform.FormId)
			}
		}
		ds.cleanUpImageFiles("1UX-0xXQPRbV5aj1G_NNX06gyODgakvQP")
	}()

	var jobComplete bool
	var retryCount int

	for retryCount < 3 && !jobComplete {
		contractsImageURL, err := ds.uploadImageToDrive("img/image-contract.jpg")
		if err != nil {
			retryCount++
			logger.WithField("retryCount", retryCount).WithError(err).Errorf("failed to upload image to drive, retrying ...")
			time.Sleep(5)
			continue
		}

		termsImageURL, err := ds.uploadImageToDrive("img/image-terms.jpg")
		if err != nil {
			retryCount++
			logger.WithField("retryCount", retryCount).WithError(err).Errorf("failed to upload image to drive, retrying ...")
			time.Sleep(5)
			continue
		}

		form, err = fs.CreateFormWithTitleAndDescription(contract, &formTitle, &formDescription)
		if err != nil {
			retryCount++
			logger.WithField("retryCount", retryCount).WithError(err).Errorf("failed to create a form, retrying ...")
			time.Sleep(5)
			continue
		}

		err = fs.CreateImageItem(form, "The Client hereby agree as follows:", contractsImageURL, 0)
		if err != nil {
			toBeDeletedforms = append(toBeDeletedforms, form)
			retryCount++
			logger.WithField("retryCount", retryCount).WithError(err).Errorf("failed to create image item, retrying ...")
			time.Sleep(5)
			continue
		}

		err = fs.CreateImageItem(form, "Terms and Conditions", termsImageURL, 1)
		if err != nil {
			toBeDeletedforms = append(toBeDeletedforms, form)
			retryCount++
			logger.WithField("retryCount", retryCount).WithError(err).Errorf("failed to create image item, retrying ...")
			time.Sleep(5)
			continue
		}

		err = fs.CreateSignatureItem(form, "Digital Signature (Printed Name):", 2)
		if err != nil {
			toBeDeletedforms = append(toBeDeletedforms, form)
			retryCount++
			logger.WithField("retryCount", retryCount).WithError(err).Errorf("failed to create signature item, retrying ...")
			time.Sleep(5)
			continue
		}

		jobComplete = true
	}
	if err != nil {
		return fmt.Errorf("failure while creating form item with error : %w", err)
	}

	err = ds.moveFormFileToSharedDirectory(form)
	if err != nil {
		return fmt.Errorf("failure to move form file to shared directory with error : %w", err)
	}

	return nil
}
