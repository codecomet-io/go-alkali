package cache

import (
	"github.com/moby/buildkit/client"
)

type cacheType string

const (
	typeRegistry cacheType = "registry"
	typeLocal    cacheType = "local"
	typeGHA      cacheType = "gha"
	typeS3       cacheType = "s3"
	typeAzure    cacheType = "azure"
)

type Entry interface {
	ToExport() client.CacheOptionsEntry
	ToImport() client.CacheOptionsEntry
}

type Options struct {
	NoCache bool
	Import  []Entry
	Export  []Entry
}

func (o *Options) ToClientImport() []client.CacheOptionsEntry {
	ret := []client.CacheOptionsEntry{}
	for _, v := range o.Import {
		ret = append(ret, v.ToImport())
	}

	return ret
}

func (o *Options) ToClientExport() []client.CacheOptionsEntry {
	ret := []client.CacheOptionsEntry{}
	for _, v := range o.Import {
		ret = append(ret, v.ToExport())
	}

	return ret
}

type Local struct {
	Path string
}

func (o *Local) ToImport() client.CacheOptionsEntry {
	return client.CacheOptionsEntry{
		Type: string(typeLocal),
		Attrs: map[string]string{
			"src": o.Path,
		},
	}
}

func (o *Local) ToExport() client.CacheOptionsEntry {
	return client.CacheOptionsEntry{
		Type: string(typeLocal),
		Attrs: map[string]string{
			"dest":           o.Path,
			"compression":    "uncompressed",
			"oci-mediatypes": "true",
			"mode":           "max",
		},
	}
}

/*
func loadGithubEnv(cache client.CacheOptionsEntry) (client.CacheOptionsEntry, error) {
	if _, ok := cache.Attrs["url"]; !ok {
		url, ok := os.LookupEnv("ACTIONS_CACHE_URL")
		if !ok {
			return cache, errors.New("cache with type gha requires url parameter or $ACTIONS_CACHE_URL")
		}
		cache.Attrs["url"] = url
	}

	if _, ok := cache.Attrs["token"]; !ok {
		token, ok := os.LookupEnv("ACTIONS_RUNTIME_TOKEN")
		if !ok {
			return cache, errors.New("cache with type gha requires token parameter or $ACTIONS_RUNTIME_TOKEN")
		}
		cache.Attrs["token"] = token
	}
	return cache, nil
}
*/
