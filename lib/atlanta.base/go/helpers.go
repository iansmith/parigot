package lib

import "github.com/iansmith/parigot/g/pb/kernel"

// Export1 is a wrapper around Export which makes it easy to say you export a single
// service. It does not change any of the Export behavior.
func Export1(packagePath, service string) (*kernel.ExportResponse, error) {
	fqSvc := &kernel.FullyQualifiedService{
		PackagePath: packagePath, Service: service}
	req := &kernel.ExportRequest{}
	req.Service = []*kernel.FullyQualifiedService{fqSvc}
	return Export(req)
}

// Require1 is a wrapper around Require which makes it easy to say you require a single
// service. It does not change any of the Require behavior.
func Require1(packagePath, service string) (*kernel.RequireResponse, error) {
	fqSvc := &kernel.FullyQualifiedService{
		PackagePath: packagePath, Service: service}
	req := &kernel.RequireRequest{}
	req.Service = []*kernel.FullyQualifiedService{fqSvc}
	return Require(req)
}
