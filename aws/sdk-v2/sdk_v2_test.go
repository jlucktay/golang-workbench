package sdk_v2_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	sdk_v2 "go.jlucktay.dev/golang-workbench/aws/sdk-v2"
)

// ...

type mockGetObjectAPI func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options),
) (*s3.GetObjectOutput, error)

func (m mockGetObjectAPI) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options),
) (*s3.GetObjectOutput, error) {
	return m(ctx, params, optFns...)
}

func TestGetObjectFromS3(t *testing.T) {
	const bodyBytes = `this is the body foo bar baz`

	cases := []struct {
		client func(t *testing.T) sdk_v2.S3GetObjectAPI
		bucket string
		key    string
		expect []byte
	}{
		{
			client: func(t *testing.T) sdk_v2.S3GetObjectAPI {
				return mockGetObjectAPI(func(_ context.Context, params *s3.GetObjectInput, _ ...func(*s3.Options),
				) (*s3.GetObjectOutput, error) {
					t.Helper()
					if params.Bucket == nil {
						t.Fatal("expect bucket to not be nil")
					}
					if e, a := "fooBucket", *params.Bucket; e != a {
						t.Errorf("expect %v, got %v", e, a)
					}
					if params.Key == nil {
						t.Fatal("expect key to not be nil")
					}
					if e, a := "barKey", *params.Key; e != a {
						t.Errorf("expect %v, got %v", e, a)
					}

					return &s3.GetObjectOutput{
						Body: io.NopCloser(bytes.NewReader([]byte(bodyBytes))),
					}, nil
				})
			},
			bucket: "fooBucket",
			key:    "barKey",
			expect: []byte(bodyBytes),
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			content, err := sdk_v2.GetObjectFromS3(ctx, tt.client(t), tt.bucket, tt.key)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			if e, a := tt.expect, content; !bytes.Equal(e, a) {
				t.Errorf("expect %v, got %v", e, a)
			}
		})
	}
}

// ...

type mockListObjectsV2Pager struct {
	PageNum int
	Pages   []*s3.ListObjectsV2Output
}

func (m *mockListObjectsV2Pager) HasMorePages() bool {
	return m.PageNum < len(m.Pages)
}

func (m *mockListObjectsV2Pager) NextPage(ctx context.Context, f ...func(*s3.Options)) (output *s3.ListObjectsV2Output, err error) {
	if m.PageNum >= len(m.Pages) {
		return nil, fmt.Errorf("no more pages")
	}
	output = m.Pages[m.PageNum]
	m.PageNum++
	return output, nil
}

func TestCountObjects(t *testing.T) {
	pager := &mockListObjectsV2Pager{
		Pages: []*s3.ListObjectsV2Output{
			{
				KeyCount: 5,
			},
			{
				KeyCount: 10,
			},
			{
				KeyCount: 15,
			},
		},
	}
	objects, err := sdk_v2.CountObjects(context.TODO(), pager)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if expect, actual := 30, objects; expect != actual {
		t.Errorf("expect %v, got %v", expect, actual)
	}
}
