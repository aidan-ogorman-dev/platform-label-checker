package main

import (
	"log"
	"os"
	"strings"

	v1 "k8s.io/api/apps/v1"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/kubernetes/scheme"
)

const (
	ownerLabel   = "owner"
	ghVolumePath = "/github/workspace"
)

func main() {
	noFilesAddedModified := false
	noFilesRenamed := false
	filesAddedModified := os.Getenv("ADDED_MODIFIED_FILES")
	if filesAddedModified == "" {
		noFilesAddedModified = true
	}
	filesRenamed := os.Getenv("RENAMED_FILES")
	if filesRenamed == "" {
		noFilesRenamed = true
	}
	if noFilesAddedModified && noFilesRenamed {
		log.Printf("Check complete, good process.")
		return
	}
	files := []string{}
	filesAddedModifiedSplit := strings.Split(filesAddedModified, " ")
	for _, f := range filesAddedModifiedSplit {
		files = append(files, f)
	}
	filesRenamedSplit := strings.Split(filesRenamed, " ")
	for _, f := range filesRenamedSplit {
		files = append(files, f)
	}
	for _, file := range files {
		filePath := ghVolumePath + file
		log.Printf("Checking %s", filePath)
		buf, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("error reading %s: %v", file, err)
			return
		}
		decode := scheme.Codecs.UniversalDeserializer().Decode
		obj, gvk, err := decode(buf, nil, nil)
		if gvk == nil {
			log.Printf("Unmarshalled file: %s, but it's not a kubernetes manifest file", file)
			return
		}
		fileYAML := obj.(*v1.Deployment)
		if err != nil {
			log.Fatalf("Error while decoding YAML object. Err was: %s", err)
		}
		log.Printf("k8s manifest labels = %v", fileYAML.ObjectMeta.Labels)
		if _, ok := fileYAML.ObjectMeta.Labels[ownerLabel]; !ok {
			log.Printf("adding 'owner' label")
			fileYAML.ObjectMeta.Labels[ownerLabel] = "platform"
			path, err := os.Create(filePath)
			y := printers.YAMLPrinter{}
			err = y.PrintObj(fileYAML, path)
			if err != nil {
				log.Fatalf("Failed to write file: %v", err)
			}
		}
	}
}
