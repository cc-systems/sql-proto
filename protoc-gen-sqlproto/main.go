package main

import (
	"github.com/cc-systems/sql-proto/gensql"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	protogen.Options{}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}

			err := gensql.GenerateFile(gen, f)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
