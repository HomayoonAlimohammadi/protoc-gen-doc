package main

import (
	"flag"
	"fmt"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/HomayoonAlimohammadi/protoc-gen-doc/doc"
)

const version = "1.0.0"

func main() {

	var rootCMD cobra.Command

	var versionCMD = cobra.Command{
		Use:   "version",
		Short: "print version",
		Run: func(*cobra.Command, []string) {
			fmt.Printf("protoc-gen-divar-doc %s\n", version)
		},
	}

	rootCMD.AddCommand(&versionCMD)

	var flags flag.FlagSet
	exclude := flags.String("exclude", "", "All widgets that are excluded from doc "+
		"validation should be listed in a dash-separated fashion, like: SELECTOR_ROW-MY_RANDOM_WIDGET-SOMETHING_ELSE "+
		"This will make 'SELECTOR_ROW' , 'MY_RANDOM_WIDGET' and 'SOMETHING_ELSE' to "+
		"bypass validation. Note that these names SHOULD be formatted like the widget "+
		"'Type' enum inside widgets.proto")

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		d, err := doc.Generate(gen, *exclude)
		if err != nil {
			return err
		}
		return doc.Export(d, gen)
	})
}
