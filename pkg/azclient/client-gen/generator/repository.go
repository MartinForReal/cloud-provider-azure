/*
Copyright 2023 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package generator
package generator

import (
	"go/ast"

	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

type RepositoryGenerator struct {
	HeaderFile string `marker:",optional"`
}

func (RepositoryGenerator) RegisterMarkers(into *markers.Registry) error {
	return markers.RegisterAll(into, clientGenMarker, enableClientGenMarker)
}

func (generator RepositoryGenerator) Generate(ctx *genall.GenerationContext) error {
	for _, root := range ctx.Roots {
		pkgMakers, err := markers.PackageMarkers(ctx.Collector, root)
		if err != nil {
			root.AddError(err)
			break
		}
		if _, markedForGeneration := pkgMakers[enableClientGenMarker.Name]; !markedForGeneration {
			continue
		}

		//visit each type
		root.NeedTypesInfo()

		err = markers.EachType(ctx.Collector, root, func(typeInfo *markers.TypeInfo) {
			marker := typeInfo.Markers.Get(clientGenMarker.Name)
			if marker == nil {
				return
			}
			markerConf := marker.(ClientGenConfig)
			if markerConf.OutOfSubscriptionScope {
				return
			}

			// generte repository code

		})
		if err != nil {
			break
		}
	}
	return err
}

func (generator RepositoryGenerator) Help() *markers.DefinitionHelp {
	return &markers.DefinitionHelp{
		DetailedHelp: markers.DetailedHelp{
			Summary: "Generate mock for the given package",
			Details: `Generate mock for the given package`,
		},
		FieldHelp: map[string]markers.DetailedHelp{
			"HeaderFile": {
				Summary: "header file path",
				Details: `header file path`,
			},
		},
	}
}

func (RepositoryGenerator) CheckFilter() loader.NodeFilter {
	return func(node ast.Node) bool {
		// ignore structs
		_, isIface := node.(*ast.InterfaceType)
		return isIface
	}
}
