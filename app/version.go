package app

// Stamp is the version number placed by
var Stamp string

// Name is the name of the application
const Name = "pod"

// Version returns the application version as a properly formed string per the semantic versioning 2.0.0 spec (http://semver.org/).
func Version() string {
	return Name + "-" + Stamp
}
