package faucet

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FaucetPipeline", func() {
	Context("HotReload", func() {
		When("hot reload is enabled and manifest.json changes", func() {
			It("serves new content", func() {
				faucetPipeline := NewPipelineAdapter("test_fixtures/working/dist/manifest.json")
				faucetPipeline.EnableHotReload()

				asset, err := faucetPipeline.TemplateFunc("dist/bundle.css")

				Expect(err).To(BeNil())
				Expect(asset).To(Equal("/dist/bundle-68b329da9893e34099c7d8ad5cb9c940.css"))

				faucetPipeline.manifestPath = "test_fixtures/working/dist/manifest2.json"

				asset, err = faucetPipeline.TemplateFunc("dist/bundle.css")

				Expect(err).To(BeNil())
				Expect(asset).To(Equal("/dist/bundle-68b329da9893e34099c7d8ad5cb9c941.css"))
			})
		})

		When("hot reload is disabled and manifest.json changes", func() {
			It("serves the same content as before", func() {
				faucetPipeline := NewPipelineAdapter("test_fixtures/working/dist/manifest.json")

				asset, err := faucetPipeline.TemplateFunc("dist/bundle.css")

				Expect(err).To(BeNil())
				Expect(asset).To(Equal("/dist/bundle-68b329da9893e34099c7d8ad5cb9c940.css"))

				faucetPipeline.manifestPath = "test_fixtures/working/dist/manifest2.json"

				asset, err = faucetPipeline.TemplateFunc("dist/bundle.css")

				Expect(err).To(BeNil())
				Expect(asset).To(Equal("/dist/bundle-68b329da9893e34099c7d8ad5cb9c940.css"))
			})
		})
	})

	Context("Execution of generated TemplateFunc by GetTemplateFunc()", func() {
		When("requested asset is available", func() {
			It("returns correct path", func() {
				faucetPipeline := NewPipelineAdapter("test_fixtures/working/dist/manifest.json")

				asset, err := faucetPipeline.TemplateFunc("dist/bundle.css")

				Expect(err).To(BeNil())
				Expect(asset).To(Equal("/dist/bundle-68b329da9893e34099c7d8ad5cb9c940.css"))
			})
		})

		When("requested asset is not available", func() {
			It("returns an error", func() {
				faucetPipeline := NewPipelineAdapter("test_fixtures/working/dist/manifest.json")

				_, err := faucetPipeline.TemplateFunc("dist/nobundle.css")

				Expect(err).NotTo(BeNil())
			})
		})
	})

	Context("loadManifest()", func() {
		When("manifest.json is available at the given path", func() {
			It("returns map with correct contents", func() {
				faucetPipeline := NewPipelineAdapter("test_fixtures/working/dist/manifest.json")
				manifest, err := faucetPipeline.loadManifest()

				Expect(err).To(BeNil())
				Expect(manifest).To(HaveLen(2))
				Expect(manifest["dist/bundle.css"]).To(Equal("/dist/bundle-68b329da9893e34099c7d8ad5cb9c940.css"))
			})
		})

		When("manifest.json is not available at the given path", func() {
			It("returns an error", func() {
				faucetPipeline := NewPipelineAdapter("test_fixtures/not_available/dist/manifest.json")
				manifest, err := faucetPipeline.loadManifest()

				Expect(err).NotTo(BeNil())
				Expect(manifest).To(BeNil())
			})
		})
	})
})
