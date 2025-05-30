package commands

import (
	"log"

	"lesiw.io/cmdio/sys"
	lib "ops/lib"
)

func (Ops) Upload() {
	var rnr = sys.Runner().WithEnv(map[string]string{
		"PWD": ".",
	})
	defer rnr.Close()
	var err error

	err = rnr.Run("npm", "install")
	if err != nil {	
		log.Fatal(err)
	}	

	err = rnr.Run("npm", "run", "build")
	if err != nil {
		log.Fatal(err)
	}

	err = lib.ClearS3Bucket("common-web-components-bucket-dev")
	if err != nil {	
		log.Fatal(err)
	}

	err = lib.UploadFolderToS3("common-web-components-bucket-dev", "./www")
	if err != nil {
		log.Fatal(err)
	}
}
