package reference

import (
	"fmt"
	"strings"

	distreference "github.com/docker/distribution/reference"
)

// SubstituteReferenceName creates a new image reference from given ref with
// its *name* part substituted for reposName.
func SubstituteReferenceName(ref distreference.Named, reposName string) (newRef distreference.Named, err error) {
	reposNameRef, err := distreference.WithName(reposName)
	if err != nil {
		return nil, err
	}
	if tagged, isTagged := ref.(distreference.Tagged); isTagged {
		newRef, err = distreference.WithTag(reposNameRef, tagged.Tag())
		if err != nil {
			return nil, err
		}
	} else if digested, isDigested := ref.(distreference.Digested); isDigested {
		newRef, err = distreference.WithDigest(reposNameRef, digested.Digest())
		if err != nil {
			return nil, err
		}
	} else {
		newRef = reposNameRef
	}
	return
}

// UnqualifyReference ...
func UnqualifyReference(ref distreference.Named) distreference.Named {
	_, remoteName, err := SplitReposName(ref)
	if err != nil {
		return ref
	}
	newRef, err := SubstituteReferenceName(ref, remoteName.Name())
	if err != nil {
		return ref
	}
	return newRef
}

// QualifyUnqualifiedReference ...
func QualifyUnqualifiedReference(ref distreference.Named, indexName string) (distreference.Named, error) {
	if !isValidHostname(indexName) {
		return nil, fmt.Errorf("Invalid hostname %q", indexName)
	}
	orig, remoteName, err := SplitReposName(ref)
	if err != nil {
		return nil, err
	}
	if orig == "" {
		return SubstituteReferenceName(ref, indexName+"/"+remoteName.Name())
	}
	return ref, nil
}

// IsReferenceFullyQualified determines whether the given reposName has prepended
// name of index.
func IsReferenceFullyQualified(reposName distreference.Named) bool {
	indexName, _, _ := SplitReposName(reposName)
	return indexName != ""
}

// SplitReposName breaks a reposName into an index name and remote name
func SplitReposName(reposName distreference.Named) (indexName string, remoteName distreference.Named, err error) {
	var remoteNameStr string
	indexName, remoteNameStr = distreference.SplitHostname(reposName)
	if !isValidHostname(indexName) {
		// This is a Docker Index repos (ex: samalba/hipache or ubuntu)
		// 'docker.io'
		indexName = ""
		remoteName = reposName
	} else {
		remoteName, err = distreference.WithName(remoteNameStr)
	}
	return
}

func isValidHostname(hostname string) bool {
	return hostname != "" && !strings.Contains(hostname, "/") &&
		(strings.Contains(hostname, ".") ||
			strings.Contains(hostname, ":") || hostname == "localhost")
}
