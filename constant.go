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

	// The name of the client requesting the action on the environment
	ClientEnvVariableKey string = "LAGOON_CLIENT"

	// The Name of the SSH Public key file
	SSHPuplicKeyFileName string = "ssh.pub"

	// The Name of the SSH Private key file
	SSHPrivateKeyFileName string = "ssh.pem"

	// The action to perform on the environment
	ActionEnvVariableKey string = "LAGOON_ACTION"

	// The environment variable used to get or store the http proxy definition
	HttpProxyEnvVariableKey string = "http_proxy"

	// The environment variable used to get or store the https proxy definition
	HttpsProxyEnvVariableKey string = "https_proxy"

	// The environment variable used to get or store the no proxy definition
	NoProxyEnvVariableKey string = "no_proxy"

	// The name of the generated proxy configuration file
	ProxyConfigFileName string = "proxy_env.yml"

	// The name of the file containing the map of all components locations
	ComponentPathsFileName string = "component_paths.yaml"

	// The name of the generated creation session file returned to the client
	CreationSessionFileName string = "create_session.json"

	// Prefix for the installer logs
	InstallerLogPrefix string = "Lagoon INSTALLER LOG:"

	// Volum location used by the installer
	InstallerVolume string = "/opt/lagoon/output"

	// The name of file containing params
	ParamYamlFileName string = "params.yaml"

	// The name of file containing extra vars
	ExtraVarYamlFileName string = "extra-vars.yaml"

	// // The name of file containing environment variables
	EnvYamlFileName string = "env.yaml"
)
