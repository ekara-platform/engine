package util

const (
	// The name of the environment file descriptor
	DescriptorFileName string = "ekara.yaml"

	// The environment variable used to pass the environment descriptor
	// location between several Ekara components.
	StarterEnvVariableKey string = "EKARA_ENV_DESCR"

	// The environment variable used to pass the environment descriptor
	// name between several Ekara components.
	StarterEnvNameVariableKey string = "EKARA_ENV_DESCR_NAME"

	// The Name of the SSH Public key file
	SSHPuplicKeyFileName string = "ssh.pub"

	// The Name of the SSH Private key file
	SSHPrivateKeyFileName string = "ssh.pem"

	// The action to perform on the environment
	ActionEnvVariableKey string = "EKARA_ACTION"

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

	// The name of the file containing the map of all components locations
	CliParametersFileName string = "cli_parameters.yaml"

	// The name of the generated creation session file returned to the client
	CreationSessionFileName string = "create_session.json"

	// Prefix for the installer logs
	InstallerLogPrefix string = "Ekara INSTALLER LOG:"

	// Volum location used by the installer
	InstallerVolume string = "/opt/ekara/output"

	// The name of file containing params
	ParamYamlFileName string = "params.yaml"

	// The name of file containing extra vars
	ExtraVarYamlFileName string = "extra-vars.yaml"

	//  The name of file containing environment variables
	EnvYamlFileName string = "env.yaml"

	//  The name of file containing host inventories
	InventoryYamlFileName string = "inventory.yaml"

	// Component modules folder
	ComponentModuleFolder string = "modules"

	// Inventory modules folder
	InventoryModuleFolder string = "inventory"
)
