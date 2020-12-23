package faucet_test

import (
	"bytes"
	"github.com/REABMAX/faucet-pipeline-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"html/template"
)

var _ = Describe("Acceptance", func() {
	Context("working", func() {
		It("returns correctly parsed template", func() {
			rendered, err := execute("test_fixtures/working/dist/manifest.json", "test_fixtures/working/index.html")

			Expect(err).To(BeNil())
			Expect(rendered).To(Equal("<script src=\"/dist/bundle-6f9cdfdf5d45a70ad818b45090761ede.js\"></script><link href=\"/dist/bundle-68b329da9893e34099c7d8ad5cb9c940.css\" />"))
		})
	})

	Context("invalid", func() {
		It("returns error", func() {
			_, err := execute("test_fixtures/invalid/dist/manifest.json", "test_fixtures/invalid/index.html")
			Expect(err).NotTo(BeNil())
		})
	})
})

func execute(manifest string, html string) (string, error) {
	faucetPipeline := faucet.Pipeline{
		ManifestPath: manifest,
	}

	tpl, err := template.
		New("index.html").
		Funcs(template.FuncMap{
			"asset": faucetPipeline.TemplateFunc,
		}).
		ParseFiles(html)

	if err != nil {
		return "", err
	}

	out := &bytes.Buffer{}
	err = tpl.Execute(out, nil)
	if err != nil {
		return "", err
	}

	return out.String(), nil
}
