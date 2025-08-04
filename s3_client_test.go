/*
 * Copyright 2022-present Kuei-chun Chen. All rights reserved.
 * s3_client_test.go
 */

package hatchet

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type mockS3Client struct {
	s3iface.S3API
	err    error
	bucket string
	key    string
	body   []byte
}

func (m *mockS3Client) CreateBucket(input *s3.CreateBucketInput) (*s3.CreateBucketOutput, error) {
	m.bucket = *input.Bucket
	return &s3.CreateBucketOutput{}, m.err
}

func (m *mockS3Client) DeleteBucket(input *s3.DeleteBucketInput) (*s3.DeleteBucketOutput, error) {
	m.bucket = ""
	return &s3.DeleteBucketOutput{}, m.err
}

func (m *mockS3Client) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	m.key = *input.Key
	buf := new(bytes.Buffer)
	buf.ReadFrom(input.Body)
	m.body = buf.Bytes()
	return &s3.PutObjectOutput{}, m.err
}

func (m *mockS3Client) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return &s3.GetObjectOutput{
		Body: ioutil.NopCloser(bytes.NewReader(m.body)),
	}, m.err
}

func (m *mockS3Client) DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	m.key = ""
	m.body = nil
	return &s3.DeleteObjectOutput{}, m.err
}

func TestS3Client(t *testing.T) {
	mockSvc := &mockS3Client{}
	s3client := &S3Client{
		service: mockSvc,
	}

	// create a new S3 bucket
	bucketName := "test-bucket"
	err := s3client.CreateBucket(bucketName)
	if err != nil {
		t.Fatalf("failed to create S3 bucket: %v", err)
	}

	// upload a file to S3
	fileName := "test-file.txt"
	testDataDir := "./testdata"
	if err := os.MkdirAll(testDataDir, 0755); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(testDataDir) })
	filePath := testDataDir + "/" + fileName
	f, err := os.Create(filePath)
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString("hello s3")
	f.Close()
	err = s3client.PutObject(bucketName, fileName, filePath)
	if err != nil {
		t.Fatalf("failed to upload file to S3: %v", err)
	}

	// download the file from S3
	var buf []byte
	buf, err = s3client.GetObject(fmt.Sprintf("%v/%v", bucketName, fileName))
	if err != nil {
		t.Fatalf("failed to download file from S3: %v", err)
	}

	// check that the contents of the downloaded file match the original file
	originalBytes, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read original file: %v", err)
	}
	if !bytes.Equal(buf, originalBytes) {
		t.Fatalf("file contents do not match: expected '%s', got '%s'", string(originalBytes), string(buf))
	}
	t.Log(string(buf))

	// delete the file from S3
	err = s3client.DeleteObject(bucketName, fileName)
	if err != nil {
		t.Fatalf("failed to delete file from S3: %v", err)
	}

	err = s3client.DeleteBucket(bucketName)
	if err != nil {
		t.Fatalf("failed to delete S3 bucket: %v", err)
	}
}
