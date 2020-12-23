package faucet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// PipelineInterface defines an interface for faucet-pipeline integrations
type PipelineInterface interface {
	TemplateFunc() (string, error)
}

// Pipeline defines a configuration for faucet-pipeline-go.
type Pipeline struct {
	ManifestPath string   // ManifestPath: Path to your manifest.json file
	HotReload    bool     // HotReload: if true, manifest will be parsed on every call. Good for development purposes
	manifest     Manifest // caches parsed manifest data internally
}

// Manifest defines the go type for manifest data
type Manifest map[string]string

// TemplateFunc is a custom function that can be added to go templates as part of a template.FuncMap
func (pipeline *Pipeline) TemplateFunc(assetPath string) (string, error) {
	manifest, err := pipeline.loadManifest()
	if err != nil {
		return assetPath, err
	}

	fingerprintedAssetPath, ok := manifest[assetPath]
	if !ok {
		return assetPath, fmt.Errorf("could not find asset key %v in manifest", assetPath)
	}

	return fingerprintedAssetPath, nil
}

// loadManifests parses the manifest file and returns Manifest and error
func (pipeline *Pipeline) loadManifest() (Manifest, error) {
	if !pipeline.shouldUseManifestCache() {
		pipeline.manifest = make(Manifest)

		file, err := os.Open(pipeline.ManifestPath)
		if err != nil {
			return nil, err
		}

		defer file.Close()

		byteValue, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(byteValue, &pipeline.manifest)
		if err != nil {
			return nil, err
		}
	}

	return pipeline.manifest, nil
}

// shouldUseManifestCache determines whether the internal manifest cache should be used or not
func (pipeline *Pipeline) shouldUseManifestCache() bool {
	if pipeline.manifest == nil {
		return false
	}

	if pipeline.HotReload {
		return false
	}

	return true
}
