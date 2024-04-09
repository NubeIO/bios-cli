package commander

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (bt *BuildTool) handleFiles(params interface{}) (interface{}, error) {
	var paramList []string
	switch p := params.(type) {
	case string:
		paramList = strings.Fields(p)
	case []string:
		paramList = p
	case []interface{}:
		for _, v := range p {
			if str, ok := v.(string); ok {
				paramList = append(paramList, str)
			} else {
				return nil, fmt.Errorf("invalid element type in []interface{}")
			}
		}
	default:
		return nil, fmt.Errorf("invalid params type for file operations")
	}
	if len(paramList) < 2 {
		return nil, fmt.Errorf("invalid params for file operations")
	}

	operation := paramList[0]
	filePath := paramList[1]
	fileOps := fileImpl{}

	switch operation {
	case "mkdir":
		err := os.MkdirAll(filePath, 0755)
		if err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %v", filePath, err)
		}

	case "delete":
		if filePath == "/" {
			return nil, fmt.Errorf("cannot delete root directory")
		}
		// Check if the filePath is one level down from the root directory
		rootDir := filepath.Dir(filepath.Clean(filePath))
		if rootDir == "/" {
			return nil, fmt.Errorf("cannot delete one level down from the root directory")
		}
		// Check if the filePath is the user's home directory
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user's home directory: %v", err)
		}
		if filePath == userHomeDir {
			return nil, fmt.Errorf("cannot delete user's home directory")
		}
		// Check if the filePath is one level down from the user's home directory
		userHomeDirClean := filepath.Clean(userHomeDir)
		if strings.HasPrefix(filePath, userHomeDirClean+string(os.PathSeparator)) {
			return nil, fmt.Errorf("cannot delete one level down from the user's home directory")
		}
		if err != nil {
			return nil, fmt.Errorf("failed to delete %s: %v", filePath, err)
		}

	case "unzip":
		if len(paramList) < 3 {
			return nil, fmt.Errorf("unzip requires source and destination")
		}
		dest := paramList[2]

		err := unzip(filePath, dest)
		if err != nil {
			return nil, fmt.Errorf("failed to unzip %s: %v", filePath, err)
		}

	case "mv":
		if len(paramList) < 3 {
			return nil, fmt.Errorf("move requires source and destination")
		}
		dest := paramList[2]
		if _, err := os.Stat(dest); err == nil {
			// Remove the existing file or directory at the destination
			if err := os.RemoveAll(dest); err != nil {
				return nil, fmt.Errorf("failed to move %s to %s: %v", filePath, dest, err)
			}
		}
		// Move the file or directory
		err := os.Rename(filePath, dest)
		if err != nil {
			return nil, fmt.Errorf("failed to move %s to %s: %v", filePath, dest, err)
		}

	case "rename":
		if len(paramList) < 3 {
			return nil, fmt.Errorf("rename requires path and new name")
		}
		newName := paramList[2]
		err := fileOps.Rename(filePath, newName)
		if err != nil {
			return nil, fmt.Errorf("failed to rename %s to %s: %v", filePath, newName, err)
		}

	case "walkup":
		walkedPaths, err := fileOps.WalkUpTree(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to walk up from %s: %v", filePath, err)
		}
		return walkedPaths, nil

	case "walkdown":
		walkedPaths, err := fileOps.WalkDownTree(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to walk down from %s: %v", filePath, err)
		}
		return walkedPaths, nil

	case "listfiles":
		files, err := fileOps.ListAllFiles(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to list files in %s: %v", filePath, err)
		}
		return files, nil

	default:
		return nil, fmt.Errorf("unsupported file operation: %s", operation)
	}

	return nil, nil
}

type fileImpl struct {
	permissions os.FileMode
}

func (f *fileImpl) Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (f *fileImpl) Rename(path string, newName string) error {
	parentDir := filepath.Dir(path)
	newFolderPath := filepath.Join(parentDir, newName)
	err := os.Rename(path, newFolderPath)
	if err != nil {
		return err
	}
	return nil
}

// WalkUpTree walks up the folder tree from the current folder.
func (f *fileImpl) WalkUpTree(path string) ([]string, error) {
	var walkedPaths []string
	currentDir := path
	for currentDir != "" {
		walkedPaths = append(walkedPaths, currentDir)
		currentDir = filepath.Dir(currentDir)
	}
	return walkedPaths, nil
}

// WalkDownTree walks down the folder tree starting from the current folder.
func (f *fileImpl) WalkDownTree(path string) ([]string, error) {
	var walkedPaths []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		walkedPaths = append(walkedPaths, path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return walkedPaths, nil
}

func (f *fileImpl) ListAllFiles(path string) ([]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var fileNames []string
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}

	return fileNames, nil
}
