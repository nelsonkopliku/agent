package mocks

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/trento-project/agent/internal/cloud"
)

func NewDiscoveredCloudMock() cloud.CloudInstance {
	metadata := &cloud.AzureMetadata{}

	jsonFile, err := os.Open("./test/fixtures/discovery/azure/azure_metadata.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, metadata)

	return cloud.CloudInstance{
		Provider: cloud.Azure,
		Metadata: metadata,
	}
}
