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

		bucketOwnershipControls, err := s3.NewBucketOwnershipControls(ctx, "bucket-ownership-control", &s3.BucketOwnershipControlsArgs{
			Bucket: bucket.ID(),
			Rule: &s3.BucketOwnershipControlsRuleArgs{
				ObjectOwnership: pulumi.String("BucketOwnerPreferred"),
			},
		})
		if err != nil {
			return err
		}

		bucketPublicAccessBlock, err := s3.NewBucketPublicAccessBlock(ctx, "public-access", &s3.BucketPublicAccessBlockArgs{
			Bucket:                bucket.ID(),
			BlockPublicAcls:       pulumi.Bool(false),
			BlockPublicPolicy:     pulumi.Bool(false),
			IgnorePublicAcls:      pulumi.Bool(false),
			RestrictPublicBuckets: pulumi.Bool(false),
		})
		if err != nil {
			return err
		}

		_, err = s3.NewBucketCorsConfigurationV2(ctx, "example", &s3.BucketCorsConfigurationV2Args{
					Bucket: bucket.ID(),
					CorsRules: s3.BucketCorsConfigurationV2CorsRuleArray{
						&s3.BucketCorsConfigurationV2CorsRuleArgs{
							AllowedHeaders: pulumi.StringArray{
								pulumi.String("*"),
							},
							AllowedMethods: pulumi.StringArray{
								pulumi.String("GET"),
								pulumi.String("HEAD"),
							},
							AllowedOrigins: pulumi.StringArray{
								// Replace this with your dev origin or "*" for testing
								pulumi.String("http://127.0.0.1:5500"),
							},
							ExposeHeaders: pulumi.StringArray{},
							MaxAgeSeconds: pulumi.Int(3000),
						},
					},
				})
		if err != nil {
			return err
		}

		_, err = s3.NewBucketAclV2(ctx, "public-read-acl", &s3.BucketAclV2Args{
			Bucket: bucket.ID(),
			Acl:    pulumi.String("public-read"),
		}, pulumi.DependsOn([]pulumi.Resource{
			bucketOwnershipControls,
			bucketPublicAccessBlock,
		}))
		if err != nil {
			return err
		}

        // Create a bucket policy to make the bucket public read-only
        _, err = s3.NewBucketPolicy(ctx, "common-web-components-bucket-policy", &s3.BucketPolicyArgs{
            Bucket: bucket.ID(),
            Policy: pulumi.String(`{
                "Version": "2012-10-17",
                "Statement": [
                    {
                        "Effect": "Allow",
                        "Principal": "*",
                        "Action": "s3:GetObject",
                        "Resource": "arn:aws:s3:::common-web-components-bucket-` + stackName + `/*"
                    }
                ]
            }`),
        })
        if err != nil {
            return err
        }

		// Export the name of the bucket
		ctx.Export("bucketName", bucket.ID())
		region := pulumi.String("us-east-1") // Replace with your bucket's region
		bucketURL := pulumi.Sprintf("https://%s.s3.%s.amazonaws.com", bucket.Bucket, region)
		ctx.Export("bucketURL", bucketURL)
		return nil
	})
}
