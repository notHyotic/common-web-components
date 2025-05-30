package lib

import (
    "context"
    "log"
    "os"
    "path/filepath"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
)

// UploadFolderToS3 uploads the contents of the ./www folder to the specified S3 bucket.
func UploadFolderToS3(bucketName string, folderPath string) error {
    // Load AWS configuration
    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
        log.Fatal(err)
    }

    client := s3.NewFromConfig(cfg)
    uploader := manager.NewUploader(client)

    // Walk through the folder and upload each file
    err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        // Skip directories
        if info.IsDir() {
            return nil
        }

        // Open the file
        file, err := os.Open(path)
        if err != nil {
            return err
        }
        defer file.Close()

        // Calculate the relative path for the S3 key
        relativePath, err := filepath.Rel(folderPath, path)
        if err != nil {
            return err
        }

        // Upload the file to S3
        _, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
            Bucket: &bucketName,
            Key:    aws.String(filepath.ToSlash(relativePath)), // Use relative path for the S3 key
            Body:   file,
        })
        if err != nil {
            return err
        }

        log.Printf("Uploaded %s to bucket %s\n", path, bucketName)
        return nil
    })

    return err
}

func ClearS3Bucket(bucketName string) error {
    // Load AWS configuration
    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
        log.Fatal(err)
    }

    client := s3.NewFromConfig(cfg)

    // List objects in the bucket
    listObjectsOutput, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
        Bucket: aws.String(bucketName),
    })
    if err != nil {
        return err
    }

    // Delete each object
    for _, object := range listObjectsOutput.Contents {
        _, err := client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
            Bucket: aws.String(bucketName),
            Key:    object.Key,
        })
        if err != nil {
            return err
        }
        log.Printf("Deleted object: %s\n", *object.Key)
    }

    // Check if there are more objects to delete (pagination)
    for listObjectsOutput.IsTruncated != nil && *listObjectsOutput.IsTruncated {
        listObjectsOutput, err = client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
            Bucket: aws.String(bucketName),
            ContinuationToken: listObjectsOutput.NextContinuationToken,
        })
        if err != nil {
            return err
        }

        for _, object := range listObjectsOutput.Contents {
            _, err := client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
                Bucket: aws.String(bucketName),
                Key:    object.Key,
            })
            if err != nil {
                return err
            }
            log.Printf("Deleted object: %s\n", *object.Key)
        }
    }

    log.Printf("Cleared all objects from bucket: %s\n", bucketName)
    return nil
}
