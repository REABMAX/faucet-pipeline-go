package faucet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// IPipelineAdapter defines an interface for faucet-pipeline integrations
type IPipelineAdapter interface {
	TemplateFunc() (string, error) // TemplateFunc is a custom function that can be added to go templates as part of a template.FuncMap
	EnableHotReload()              // EnableHotReload enables hot reloading functionality for TemplateFunc
	DisableHotReload()             // DisableHotReload disables hot reloading functionality for TemplateFunc
	HotReloadIsEnabled() bool      // HotReloadIsEnabled checks if hot reload functionality is enabled and returns the fitting boolean value
}

// Pipeline defines a configuration for faucet-pipeline-go.
type PipelineAdapter struct {
	manifestPath string   // manifestPath: Path to your manifest.json file
	hotReload    bool     // hotReload: if true, manifest will be parsed on every call. Good for development purposes
	manifest     manifest // caches parsed manifest data internally
}

// manifest defines the go type for manifest data
type manifest map[string]string

// NewPipelineAdapter creates a new PipelineAdapter instance with hotReload functionality disabled by default. Accepts the
// manifestPath as only parameter
func NewPipelineAdapter(manifestPath string) *PipelineAdapter {
	return &PipelineAdapter{
		manifestPath: manifestPath,
		hotReload:    false,
	}
}

func (p *PipelineAdapter) TemplateFunc(assetPath string) (string, error) {
	manifest, err := p.loadManifest()
	if err != nil {
		return assetPath, err
	}

	fingerprintedAssetPath, ok := manifest[assetPath]
	if !ok {
		return assetPath, fmt.Errorf("could not find asset key %v in manifest", assetPath)
	}

	return fingerprintedAssetPath, nil
}

func (p *PipelineAdapter) EnableHotReload() {
	p.hotReload = true
}

func (p *PipelineAdapter) DisableHotReload() {
	p.hotReload = false
}

func (p *PipelineAdapter) HotReloadIsEnabled() bool {
	return p.hotReload
}

// loadManifests parses the manifest file and returns Manifest and error
func (p *PipelineAdapter) loadManifest() (manifest, error) {
	if !p.shouldUseManifestCache() {
		p.manifest = make(manifest)

		file, err := os.Open(p.manifestPath)
		if err != nil {
			return nil, err
		}

		defer file.Close()

		byteValue, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(byteValue, &p.manifest)
		if err != nil {
			return nil, err
		}
	}

	return p.manifest, nil
}

// shouldUseManifestCache determines whether the internal manifest cache should be used or not
func (p *PipelineAdapter) shouldUseManifestCache() bool {
	if p.manifest == nil {
		return false
	}

	if p.hotReload {
		return false
	}

	return true
}
