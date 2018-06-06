package engine

const (
	// The name of the environment file descriptor
	//
	// All requests of environment creation or update are made passing as parameter
	// the root folder of the environment specification. Then the Laggon Engine
	// will look for this specific file into the given folder.
	DescriptorFileName string = "lagoon.yaml"

	// The environment variable used to pass the environment descriptor
	// location between several Lagoon components.
	StarterEnvVariableKey string = "LAGOON_ENV_DESCR"

	// The name of the client requesting the environment creation or update
	ClientEnvVariableKey string = "LAGOON_CLIENT"

	// The environment variable used to get or store the http proxy definition
	HttpProxyEnvVariableKey string = "http_proxy"

	// The environment variable used to get or store the https proxy definition
	HttpsProxyEnvVariableKey string = "https_proxy"

	// The environment variable used to get or store the no proxy definition
	NoProxyEnvVariableKey string = "no_proxy"

	// The name of the generated proxy configuration file
	ProxyConfigFileName string = "proxy_env.yml"

	// The name of the generated configuration file for a node
	NodeConfigFileName string = "conf.yml"

	// The name of the generated docker config file for a node
	NodeDockerFileName string = "docker.yml"

	// The name of the generated creation session file returned to the client
	CreationSessionFileName string = "create_session.json"

	// Prefix for the installer logs
	InstallerLogPrefix string = "Lagoon INSTALLER LOG:"

	// Volum location used by the installer
	InstallerVolume string = "/opt/lagoon/output"
)
