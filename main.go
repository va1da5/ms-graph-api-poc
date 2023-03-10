package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/users"

	"github.com/joho/godotenv"
)

func toInt32Ptr(i int32) *int32 { return &i }

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	tenantId := os.Getenv("TENANT_ID")
	clientId := os.Getenv("CLIENT_ID")
	fpxPassword := os.Getenv("FPX_PASSWORD")
	fpxCertificate := os.Getenv("FPX_CERTIFICATE_PATH")

	loadFile, err := os.ReadFile(fpxCertificate)

	if err != nil {
		log.Panicf("error reading certificate file from %s: %v\n", fpxCertificate, err)
	}

	certificates, privateKey, err := azidentity.ParseCertificates(loadFile, []byte(fpxPassword))
	if err != nil {
		log.Panicf("error parsing certificates: %v\n", err)
	}

	credentials, err := azidentity.NewClientCertificateCredential(
		tenantId, clientId, certificates, privateKey, &azidentity.ClientCertificateCredentialOptions{},
	)
	if err != nil {
		log.Panicf("error creating credentials: %v\n", err)
	}

	client, err := msgraphsdk.NewGraphServiceClientWithCredentials(credentials, []string{})
	if err != nil {
		log.Panicf("error creating client: %v\n", err)
	}

	result, err := client.Users().Get(
		context.Background(),
		&users.UsersRequestBuilderGetRequestConfiguration{
			QueryParameters: &users.UsersRequestBuilderGetQueryParameters{
				Top: toInt32Ptr(5),
			},
		},
	)
	if err != nil {
		log.Panicf("error getting users: %v\n", err)
	}

	list, err := userCreateListFromResult(client, result)
	if err != nil {
		log.Panicf(`failed to query MsGraph user: %v`, err)
	}

	switch len(list) {
	case 0:
		log.Panicf(`failed to query MsGraph user: %v`, err)
	case 1:
		fmt.Print(list[0])
	default:
		fmt.Printf("Found %d users. First %v", len(list), list[0])
	}

}

func userCreateListFromResult(client *msgraphsdk.GraphServiceClient, result models.UserCollectionResponseable) (list []interface{}, err error) {
	pageIterator, pageIteratorErr := msgraphcore.NewPageIterator(result, client.GetAdapter(), models.CreateUserCollectionResponseFromDiscriminatorValue)
	if pageIteratorErr != nil {
		return list, pageIteratorErr
	}

	iterateErr := pageIterator.Iterate(context.Background(), func(pageItem interface{}) bool {
		user := pageItem.(models.Userable)

		obj, serializeErr := serializeObject(client, user)
		if serializeErr != nil {
			err = serializeErr
			return false
		}

		list = append(list, obj)
		return true
	})
	if iterateErr != nil {
		return list, iterateErr
	}

	return
}

func serializeObject(client *msgraphsdk.GraphServiceClient, resultObj serialization.Parsable) (obj interface{}, err error) {
	writer, err := client.GetAdapter().GetSerializationWriterFactory().GetSerializationWriter("application/json")
	if err != nil {
		return nil, err
	}

	err = writer.WriteObjectValue("", resultObj)
	if err != nil {
		return nil, err
	}

	serializedValue, err := writer.GetSerializedContent()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(serializedValue, &obj)
	if err != nil {
		return nil, err
	}

	return
}
