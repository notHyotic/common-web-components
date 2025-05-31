package lib

import (
    "context"
    "log"
    "mime"
    "os"
    "path/filepath"
    "strings"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
)

// UploadFolderToS3 uploads the contents of a folder to the specified S3 bucket with MIME types and public-read access.
func UploadFolderToS3(bucketName string, folderPath string) error {
    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
        log.Fatal(err)
    }

    client := s3.NewFromConfig(cfg)
    uploader := manager.NewUploader(client)

    err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if info.IsDir() {
            return nil
        }

        file, err := os.Open(path)
        if err != nil {
            return err
        }
        defer file.Close()

        relativePath, err := filepath.Rel(folderPath, path)
        if err != nil {
            return err
        }

        ext := strings.ToLower(filepath.Ext(relativePath))
        mimeType := mime.TypeByExtension(ext)

        // Normalize common web types if mime.TypeByExtension fails or is incomplete
        if mimeType == "" {
            switch ext {
            case ".js":
                mimeType = "application/javascript"
            case ".css":
                mimeType = "text/css"
            case ".html", ".htm":
                mimeType = "text/html"
            case ".json":
                mimeType = "application/json"
            case ".svg":
                mimeType = "image/svg+xml"
            case ".woff":
                mimeType = "font/woff"
            case ".woff2":
                mimeType = "font/woff2"
            default:
                mimeType = "application/octet-stream"
            }
        }

        // Upload the file with MIME type and public-read access
        _, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
            Bucket:      aws.String(bucketName),
            Key:         aws.String(filepath.ToSlash(relativePath)),
            Body:        file,
            ContentType: aws.String(mimeType),
            ACL:         "public-read",
        })

        if err != nil {
            return err
        }

        log.Printf("Uploaded %s to %s with Content-Type: %s\n", path, bucketName, mimeType)
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
