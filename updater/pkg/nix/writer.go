package nix

import (
	"fmt"
	"os"
	"text/template"
)

var nixTemplate = template.Must(template.New("nix").Parse(`{
  image = "{{ .Image }}";
  followTag = "{{ .FollowTag }}";
  hash = "{{ .Hash }}";
}
`))

func WriteNix(imagePath string, info ImageInfo) error {
	file, err := os.Create(imagePath)
	if err != nil {
		return fmt.Errorf("failed creating file: %w", err)
	}
	defer file.Close()

	// Execute the template, writing the output to the file
	err = nixTemplate.Execute(file, info)
	if err != nil {
		return fmt.Errorf("failed executing template: %w", err)
	}
	return nil
}
