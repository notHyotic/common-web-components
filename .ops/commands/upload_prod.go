package commands

import (
	"log"

	"lesiw.io/cmdio/sys"
	s3 "ops/lib"
)


func (Ops) Uploadprod() {
	var rnr = sys.Runner().WithEnv(map[string]string{
		"PWD": ".",
	})
	defer rnr.Close()
	var err error
	bucketName := "common-web-components-bucket-prod"

	err = rnr.Run("npm", "install")
	if err != nil {	
		log.Fatal(err)
	}	

	err = rnr.Run("npm", "run", "build")
	if err != nil {
		log.Fatal(err)
	}

	err = s3.ClearS3Bucket(bucketName)
	if err != nil {	
		log.Fatal(err)
	}

	err = s3.UploadFolderToS3(bucketName, "./www")
	if err != nil {
		log.Fatal(err)
	}
}
