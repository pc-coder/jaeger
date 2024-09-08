// Copyright (c) 2021 The Jaeger Authors.
// SPDX-License-Identifier: Apache-2.0

package renderer

import (
	"strconv"

	"github.com/jaegertracing/jaeger/cmd/esmapping-generator/app"
	"github.com/jaegertracing/jaeger/pkg/es"
	"github.com/jaegertracing/jaeger/plugin/storage/es/mappings"
)

var supportedMappings = map[string]struct{}{
	"jaeger-span":         {},
	"jaeger-service":      {},
	"jaeger-dependencies": {},
	"jaeger-sampling":     {},
}

// GetMappingAsString returns rendered index templates as string
func GetMappingAsString(builder es.TemplateBuilder, opt *app.Options) (string, error) {
	enableILM, err := strconv.ParseBool(opt.UseILM)
	if err != nil {
		return "", err
	}
	disableLogsFieldSearch, err := strconv.ParseBool(opt.DisableLogsFieldSearch)
	if err != nil {
		return "", err
	}

	mappingBuilder := mappings.MappingBuilder{
		TemplateBuilder:        builder,
		Shards:                 opt.Shards,
		Replicas:               opt.Replicas,
		EsVersion:              opt.EsVersion,
		IndexPrefix:            opt.IndexPrefix,
		UseILM:                 enableILM,
		ILMPolicyName:          opt.ILMPolicyName,
		DisableLogsFieldSearch: disableLogsFieldSearch,
	}
	return mappingBuilder.GetMapping(opt.Mapping)
}

// IsValidOption checks if passed option is a valid index template.
func IsValidOption(val string) bool {
	_, ok := supportedMappings[val]
	return ok
}
