// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build darwin && arm64 && !go1.24 && !datadog.no_waf && (cgo || appsec)

package lib

// THIS FILE IS AUTOGENERATED. DO NOT EDIT.

import _ "embed" // Needed for go:embed

//go:embed libddwaf-darwin-arm64.dylib.gz
var libddwaf []byte

const embedNamePattern = "libddwaf-*.dylib"
