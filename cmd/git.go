package commander

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func (bt *BuildTool) handleGitHubDownload(params interface{}) (interface{}, error) {
	paramMap, ok := params.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid params for GitHub download")
	}

	owner, _ := paramMap["owner"].(string)
	repo, _ := paramMap["repo"].(string)
	tag, _ := paramMap["tag"].(string)
	arch, _ := paramMap["arch"].(string)
	token, _ := paramMap["token"].(string)
	downloadDir, _ := paramMap["location"].(string)
	if downloadDir == "" {
		downloadDir = "./"
	}
	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", fmt.Sprintf("token %s", token)).
		Get(fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", owner, repo, tag))
	if err != nil {
		return nil, fmt.Errorf("failed to get release information: %v", err)
	}
	if resp.StatusCode() > 300 {
		return nil, fmt.Errorf("http reponse: %v", resp.String())
	}

	// Parse the release information to find the correct zip file based on architecture
	var releaseInfo map[string]interface{}
	err = json.Unmarshal(resp.Body(), &releaseInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to parse release information: %v", err)
	}

	if assets, ok := releaseInfo["assets"].([]interface{}); ok {
		for _, asset := range assets {
			assetInfo := asset.(map[string]interface{})
			name := assetInfo["name"].(string)
			if strings.Contains(name, arch) {
				// Download the release zip file
				zipFilePath := filepath.Join(downloadDir, name)
				resp, err := client.R().
					SetHeader("Authorization", fmt.Sprintf("token %s", token)).
					SetOutput(zipFilePath).
					Get(assetInfo["browser_download_url"].(string))
				if err != nil {
					return nil, fmt.Errorf("failed to download release zip: %v", err)
				}
				defer resp.RawResponse.Body.Close()
				bt.UpdateVar("zipName", zipFilePath)
				fmt.Printf("Release successfully downloaded to: %s\n", zipFilePath)
				return nil, nil
			}
		}
	} else {
		return nil, fmt.Errorf("no assets found in the release information")
	}

	return nil, fmt.Errorf("no matching zip file found for architecture: %s", arch)
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	os.MkdirAll(dest, 0755)

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}

		path := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			f, err := os.OpenFile(
				path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
			f.Close()
		}
		rc.Close()
	}
	return nil
}
