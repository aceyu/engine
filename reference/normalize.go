package reference

import (
	"errors"
	"fmt"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/registry"
	"strings"
)

type normalizedNamed interface {
	reference.Named
	Familiar() reference.Named
}

func ParseNormalizedNamed(s string) (reference.Named, error){
	if ok := anchoredIdentifierRegexp.MatchString(s); ok {
		return nil, fmt.Errorf("invalid repository name (%s), cannot specify 64-byte hexadecimal strings", s)
	}
	domain, remainder := splitDockerDomain(s)

	var remoteName string
	if tagSep := strings.IndexRune(remainder, ':'); tagSep > -1 {
		remoteName = remainder[:tagSep]
	} else {
		remoteName = remainder
	}
	if strings.ToLower(remoteName) != remoteName {
		return nil, errors.New("invalid reference format: repository name must be lowercase")
	}
	sn :=""
	if domain == "" {
		sn = remainder
	} else {
		sn = domain + "/" + remainder
	}
	ref, err := reference.Parse(sn)
	if err != nil {
		return nil, err
	}
	named, isNamed := ref.(reference.Named)
	if !isNamed {
		return nil, fmt.Errorf("reference %s has no name", ref.String())
	}
	return named, nil
}

func splitDockerDomain(name string) (domain, remainder string) {
	i := strings.IndexRune(name, '/')
	if i == -1 || (!strings.ContainsAny(name[:i], ".:") && name[:i] != "localhost") {
		domain, remainder = "", name
	} else {
		domain, remainder = name[:i], name[i+1:]
	}
	for _, r := range registry.DefaultRegistries {
		if (domain == r || domain=="") && !strings.ContainsRune(remainder, '/') {
			remainder = "library/" + remainder
			return
		}
	}
	return
}

func trimDefaultRegistry(s string) string {
	domain, _:= splitDockerDomain(s)
	for _, r := range registry.DefaultRegistries {
		if domain == r {
			if strings.Index(s, domain +"/library") != -1 {
				return strings.Replace(s, r + "/library/", "", 1)
			} else {
				return strings.Replace(s, r + "/", "", 1)
			}
		}
	}
	return s
}

func FamiliarName(ref reference.Named) (s string) {
	s = reference.FamiliarName(ref)
	s = trimDefaultRegistry(s)
	return
}

func FamiliarString(ref reference.Named)  (s string) {
	s =  reference.FamiliarString(ref)
	s = trimDefaultRegistry(s)
	return
}