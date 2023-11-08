package sdk_v2

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// ...

type S3GetObjectAPI interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

func GetObjectFromS3(ctx context.Context, api S3GetObjectAPI, bucket, key string) ([]byte, error) {
	object, err := api.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	defer object.Body.Close()

	return io.ReadAll(object.Body)
}

// ...

type ListObjectsV2Pager interface {
	HasMorePages() bool
	NextPage(context.Context, ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

func CountObjects(ctx context.Context, pager ListObjectsV2Pager) (count int, err error) {
	for pager.HasMorePages() {
		var output *s3.ListObjectsV2Output
		output, err = pager.NextPage(ctx)
		if err != nil {
			return count, err
		}
		count += int(output.KeyCount)
	}
	return count, nil
}
