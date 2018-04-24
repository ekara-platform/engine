package engine

import (
	"net/url"
	"testing"

	"github.com/lagoon-platform/model"
	"github.com/stretchr/testify/assert"
)

func TestBuildDescriptorUrl(t *testing.T) {

	u := url.URL{Scheme: "https", Host: model.GitHubHost, Path: "blablabla"}
	cUrl := BuildDescriptorUrl(u)
	assert.Equal(t, u.Path+"/"+DescriptorFileName, cUrl.Path)

	u = url.URL{Scheme: "https", Host: model.GitHubHost, Path: "blablabla/"}
	cUrl = BuildDescriptorUrl(u)
	assert.Equal(t, u.Path+DescriptorFileName, cUrl.Path)

}
