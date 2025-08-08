// main is the package for the C API.
package main

// Make sure cgo is enabled `$env:CGO_ENABLED="1"`.
// generate shared library and header, this requires setting up a gcc compiler on windows.
//go:generate go build -o ../generated/libinsights.dll -buildmode=c-shared libinsights.go

// Copy insights_types.h to the generated folder
//go:generate sh -c "cp insights_types.h ../generated/types.h"
