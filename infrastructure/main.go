package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create an AWS resource (S3 Bucket)
		stackName := ctx.Stack()
		bucket, err := s3.NewBucketV2(ctx, "common-web-components-bucket", &s3.BucketV2Args{
			Bucket: pulumi.String("common-web-components-bucket-"+stackName),
		})
		if err != nil {
			return err
		}

		// Export the name of the bucket
		ctx.Export("bucketName", bucket.ID())
		return nil
	})
}
