package model

//ProxyAware represents the common behavior of proxy definition holders
type ProxyAware interface {
	Proxy() Proxy
}

//Proxy represents the proxy definition
type Proxy struct {
	Http    string `yaml:"http_proxy,omitempty" json:",omitempty"`
	Https   string `yaml:"https_proxy,omitempty" json:",omitempty"`
	NoProxy string `yaml:"no_proxy,omitempty" json:",omitempty"`
}

func createProxy(yamlRef yamlProxy) Proxy {
	return Proxy{
		Http:    yamlRef.Http,
		Https:   yamlRef.Https,
		NoProxy: yamlRef.NoProxy,
	}
}

func (r Proxy) override(with Proxy) Proxy {
	if with.Http != "" {
		r.Http = with.Http
	}
	if with.Https != "" {
		r.Https = with.Https
	}
	if with.NoProxy != "" {
		r.NoProxy = with.NoProxy
	}
	return r
}
