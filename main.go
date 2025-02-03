package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"watcharis/go-poc-gcs/config"
	"watcharis/go-poc-gcs/gcs"
)

const (
	BUCKET_NAME          = "go-poc-gcs"
	PATHFILE_DESTINATION = "./data/%s"
)

func main() {
	log.Println("start project go-poc-gcs")

	ctx := context.Background()

	mainConfigs := config.InitConfig(ctx)
	log.Printf("mainConfigs : %+v\n", mainConfigs)

	// credentialByte, err := os.ReadFile("./fleet-impact-331903-e6c38ecca990.json")
	// if err != nil {
	// 	log.Panic(err)
	// }

	credentialByte, err := json.Marshal(mainConfigs.Secret.GcsCredential)
	if err != nil {
		log.Panic(err)
	}

	storageClient, err := gcs.NewGoogleCloudStorageClient(ctx, credentialByte)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("storageClient : %+v\n", storageClient)

	googleCloudStorageRepository := gcs.NewGoogleCloudStorageRepository(credentialByte)
	// fmt.Printf("googleCloudStorageRepository : %+v\n", googleCloudStorageRepository)

	destination := fmt.Sprintf(PATHFILE_DESTINATION, "profiles_202406222320.sql")
	if err := googleCloudStorageRepository.Download(ctx, BUCKET_NAME, "profiles_202406222320.sql", destination); err != nil {
		fmt.Println("[ERROR] cannot download file from gcs :", err)
		return
	}

	if err := googleCloudStorageRepository.Upload(ctx, BUCKET_NAME, "google-create-credential.txt"); err != nil {
		fmt.Println("[ERROR] cannot upload file to gcs :", err)
		return
	}

	if err := googleCloudStorageRepository.ReadFile(ctx, BUCKET_NAME, "profiles_202406222320.sql"); err != nil {
		fmt.Println("[ERROR] cannot readfile from gcs :", err)
		return
	}

	if err := googleCloudStorageRepository.GenerateV4GetObjectSignedURL(ctx, BUCKET_NAME, "profiles_202406222320.sql"); err != nil {
		fmt.Println("[ERROR] cannot sign url from gcs :", err)
		return
	}
}
