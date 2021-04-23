package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Name of APIs that the dapla-cli communicates with
const (
	APINameDataMaintenanceSvc = "data-maintenance"
	APINamePseudoSvc          = "dapla-pseudo-service"
)

func apiURLOf(apiName string) string {
	var apiURL, err = apiURLOrError(apiName)
	cobra.CheckErr(err)
	return apiURL
}

func apiURLOrError(apiName string) (string, error) {
	apiURLs := viper.GetStringMapString("apis")
	if apiURLs == nil {
		return "", fmt.Errorf("unable to determine API URLs from config")
	}

	apiURL := apiURLs[apiName]
	if apiURL == "" {
		return "", fmt.Errorf("unable to determine API URL for %v", apiName)
	}

	if strings.HasPrefix(apiURL, "$") {
		if resolvedURL := os.Getenv(apiURL[1:]); resolvedURL != "" {
			return resolvedURL, nil
		}

		return "", fmt.Errorf("unable to resolve %v API URL for environment variable %v", apiName, apiURL)
	}

	return apiURL, nil
}

func allAPIUrls() map[string]string {
	return map[string]string{
		APINameDataMaintenanceSvc: apiURLOf(APINameDataMaintenanceSvc),
		APINamePseudoSvc:          apiURLOf(APINamePseudoSvc),
	}
}

func allAPIUrlsString() string {
	apis, _ := json.MarshalIndent(allAPIUrls(), "", "\t")
	return string(apis)
}
