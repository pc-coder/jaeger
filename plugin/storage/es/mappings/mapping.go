// Copyright (c) 2021 The Jaeger Authors.
// SPDX-License-Identifier: Apache-2.0

package mappings

import (
	"bytes"
	"embed"
	"strings"

	"github.com/jaegertracing/jaeger/pkg/es"
)

// MAPPINGS contains embedded index templates.
//
//go:embed *.json
var MAPPINGS embed.FS

// MappingBuilder holds parameters required to render an elasticsearch index template
type MappingBuilder struct {
	TemplateBuilder              es.TemplateBuilder
	Shards                       int64
	Replicas                     int64
	PrioritySpanTemplate         int64
	PriorityServiceTemplate      int64
	PriorityDependenciesTemplate int64
	PrioritySamplingTemplate     int64
	EsVersion                    uint
	IndexPrefix                  string
	UseILM                       bool
	ILMPolicyName                string
	DisableLogsFieldSearch       bool
}

// GetMapping returns the rendered mapping based on elasticsearch version
func (mb *MappingBuilder) GetMapping(mapping string) (string, error) {
	if mb.EsVersion == 8 {
		return mb.fixMapping(mapping + "-8.json")
	} else if mb.EsVersion == 7 {
		return mb.fixMapping(mapping + "-7.json")
	}
	return mb.fixMapping(mapping + "-6.json")
}

// GetSpanServiceMappings returns span and service mappings
func (mb *MappingBuilder) GetSpanServiceMappings() (spanMapping string, serviceMapping string, err error) {
	spanMapping, err = mb.GetMapping("jaeger-span")
	if err != nil {
		return "", "", err
	}
	serviceMapping, err = mb.GetMapping("jaeger-service")
	if err != nil {
		return "", "", err
	}
	return spanMapping, serviceMapping, nil
}

// GetDependenciesMappings returns dependencies mappings
func (mb *MappingBuilder) GetDependenciesMappings() (string, error) {
	return mb.GetMapping("jaeger-dependencies")
}

// GetSamplingMappings returns sampling mappings
func (mb *MappingBuilder) GetSamplingMappings() (string, error) {
	return mb.GetMapping("jaeger-sampling")
}

func loadMapping(name string) string {
	s, _ := MAPPINGS.ReadFile(name)
	return string(s)
}

func (mb *MappingBuilder) fixMapping(mapping string) (string, error) {
	tmpl, err := mb.TemplateBuilder.Parse(loadMapping(mapping))
	if err != nil {
		return "", err
	}
	writer := new(bytes.Buffer)

	if mb.IndexPrefix != "" && !strings.HasSuffix(mb.IndexPrefix, "-") {
		mb.IndexPrefix += "-"
	}
	if err := tmpl.Execute(writer, mb); err != nil {
		return "", err
	}

	return writer.String(), nil
}
