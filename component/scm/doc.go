// Package scm is the abstraction of all available SCM Manager allowing to fetch components.
//
// Currently available SCM Manager are :
//   - git
//   - file (local files, only for test purpose)
//
// The SCM Manager used to fetch a repository will depend on the github.com/ekara-platform/model.SCMType
// of the desired repository.
//   - GitScm
//   - SvnScm
//   - UnknownScm
//
// A repository with an UnknownScm but with an URL but with a file scheme will
// be intended to be fetched as a local files repository.
package scm
