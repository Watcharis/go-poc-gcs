package gcs

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type googleCloudStorageRepository struct {
	credentialJSON []byte
}

type GoogleCloudStorageRepository interface {
	Download(ctx context.Context, bucket string, object string, destination string) error
	Upload(ctx context.Context, bucket string, object string) error
	ReadFile(ctx context.Context, bucket string, object string) error
	GenerateV4GetObjectSignedURL(ctx context.Context, bucket string, object string) error
}

func NewGoogleCloudStorageClient(ctx context.Context, credentialJSON []byte) (*storage.Client, error) {
	storageClient, err := storage.NewClient(ctx, option.WithCredentialsJSON(credentialJSON))
	if err != nil {
		return nil, err
	}
	return storageClient, nil
}

func NewGoogleCloudStorageRepository(credentialJSON []byte) GoogleCloudStorageRepository {
	return &googleCloudStorageRepository{
		credentialJSON: credentialJSON,
	}
}

func (gcs *googleCloudStorageRepository) Download(ctx context.Context, bucket string, object string, destination string) error {

	client, err := NewGoogleCloudStorageClient(ctx, gcs.credentialJSON)
	if err != nil {
		return err
	}
	defer client.Close()

	f, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("os.Create: %w", err)
	}

	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("Object(%q).NewReader: %w", object, err)
	}
	defer rc.Close()

	if _, err := io.Copy(f, rc); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}

	if err = f.Close(); err != nil {
		return fmt.Errorf("f.Close: %w", err)
	}

	return nil
}

func (gcs *googleCloudStorageRepository) Upload(ctx context.Context, bucket string, object string) error {

	client, err := NewGoogleCloudStorageClient(ctx, gcs.credentialJSON)
	if err != nil {
		return err
	}
	defer client.Close()

	// Open local file.
	f, err := os.Open(object)
	if err != nil {
		return fmt.Errorf("os.Open: %w", err)
	}
	defer f.Close()

	o := client.Bucket(bucket).Object(object)

	// o = o.If(storage.Conditions{DoesNotExist: true})

	wc := o.NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %w", err)
	}

	return nil
}

func (gcs *googleCloudStorageRepository) ReadFile(ctx context.Context, bucket string, object string) error {

	client, err := NewGoogleCloudStorageClient(ctx, gcs.credentialJSON)
	if err != nil {
		return err
	}
	defer client.Close()

	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return err
	}

	defer rc.Close()

	// slurp, err := io.ReadAll(rc)
	// if err != nil {
	// 	return err
	// }

	scanner := bufio.NewScanner(rc)
	row := 0
	for scanner.Scan() {
		line := scanner.Text()
		if row >= 2 {
			break
		}
		fmt.Println("line :", line)
		fmt.Println("row :", row)
		row += 1
	}

	return nil
}

func (gcs *googleCloudStorageRepository) GenerateV4GetObjectSignedURL(ctx context.Context, bucket string, object string) error {

	client, err := NewGoogleCloudStorageClient(ctx, gcs.credentialJSON)
	if err != nil {
		return err
	}
	defer client.Close()

	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(15 * time.Minute),
	}

	u, err := client.Bucket(bucket).SignedURL(object, opts)
	if err != nil {
		return fmt.Errorf("Bucket(%q).SignedURL: %w", bucket, err)
	}

	fmt.Println("Generated GET signed URL:")
	fmt.Printf("%q\n", u)
	fmt.Println("You can use this URL with any user agent, for example:")
	fmt.Printf("curl %q\n", u)

	return nil
}
