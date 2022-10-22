package gdrive

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type GDriveWebAPI struct {
	d *drive.Service
}

var (
	ErrFileNotFound = errors.New("file not found")
)

func New(apiJSONFilePath string) *GDriveWebAPI {
	d, err := drive.NewService(context.Background(), option.WithCredentialsFile(apiJSONFilePath))
	if err != nil {
		panic(err)
	}

	return &GDriveWebAPI{
		d: d,
	}
}

func (g *GDriveWebAPI) UploadCSVFile(ctx context.Context, name string, data []byte) (string, error) {
	fileId, err := g.getFileIdByName(ctx, name)
	if err != nil {
		if !errors.Is(err, ErrFileNotFound) {
			return "", fmt.Errorf("GDriveWebAPI.UploadCSVFile: g.getFileIdByName: %w", err)
		}

		id, err := g.createFile(ctx, name, data)
		if err != nil {
			return "", fmt.Errorf("GDriveWebAPI.UploadCSVFile: g.createFile: %w", err)
		}

		return g.getFileURL(id), nil
	}

	err = g.updateFile(ctx, fileId, data)
	if err != nil {
		return "", fmt.Errorf("GDriveWebAPI.UploadCSVFile: g.updateFile: %w", err)
	}

	return g.getFileURL(fileId), nil
}

func (g *GDriveWebAPI) DeleteFile(ctx context.Context, name string) error {
	fileId, err := g.getFileIdByName(ctx, name)
	if err != nil {
		return fmt.Errorf("GDriveWebAPI.DeleteFile: g.getFileIdByName: %w", err)
	}

	err = g.d.Files.Delete(fileId).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("GDriveWebAPI.DeleteFile: g.d.Files.Delete: %w", err)
	}

	return nil
}

func (g *GDriveWebAPI) GetAllFilenames(ctx context.Context) ([]string, error) {
	files, err := g.getAllFiles(ctx)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(files))
	for _, file := range files {
		names = append(names, file.Name)
	}

	return names, nil
}

// createFile creates a csv file in Google Drive with public read access and returns its ID and URL
func (g *GDriveWebAPI) createFile(ctx context.Context, name string, content []byte) (string, error) {
	file := &drive.File{
		Name:     name,
		MimeType: "text/csv",
	}

	permissions := &drive.Permission{
		Type: "anyone",
		Role: "reader",
	}

	_, err := g.d.Files.Create(file).Context(ctx).Media(bytes.NewReader(content)).Do()
	if err != nil {
		return "", err
	}

	fileId, err := g.getFileIdByName(ctx, name)
	if err != nil {
		return "", err
	}

	_, err = g.d.Permissions.Create(fileId, permissions).Context(ctx).Do()
	if err != nil {
		return "", err
	}

	return fileId, nil
}

func (g *GDriveWebAPI) updateFile(ctx context.Context, id string, content []byte) error {
	_, err := g.d.Files.Update(id, &drive.File{}).Context(ctx).Media(bytes.NewReader(content)).Do()

	return err
}

func (g *GDriveWebAPI) getFileURL(id string) string {
	return fmt.Sprintf("https://drive.google.com/file/d/%s/view?usp=sharing", id)
}

func (g *GDriveWebAPI) getAllFiles(ctx context.Context) ([]*drive.File, error) {
	r, err := g.d.Files.List().Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return r.Files, nil
}

func (g *GDriveWebAPI) getFileIdByName(ctx context.Context, name string) (string, error) {
	files, err := g.getAllFiles(ctx)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if file.Name == name {
			return file.Id, nil
		}
	}

	return "", ErrFileNotFound
}
