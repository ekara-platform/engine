package util

const (
	//StarterEnvVariableKey is the environment variable key used to pass
	// the environment descriptor location between several Ekara components.
	StarterEnvVariableKey string = "EKARA_ENV_DESCR"

	//StarterEnvQualifiedVariableKey is the environment variable key used to pass
	// the environment qualified name between several Ekara components.
	StarterEnvQualifiedVariableKey string = "EKARA_QUALIFIED_NAME"

	//StarterEnvNameVariableKey is the environment variable key used to pass
	// the environment descriptor name between several Ekara components.
	StarterEnvNameVariableKey string = "EKARA_ENV_DESCR_NAME"

	//StarterEnvLoginVariableKey is the environment variable key used to pass
	//the environment descriptor user
	StarterEnvLoginVariableKey string = "EKARA_ENV_DESCR_LOGIN"

	//StarterEnvPasswordVariableKey is the environment variable key used
	//to pass the environment descriptor password
	StarterEnvPasswordVariableKey string = "EKARA_ENV_DESCR_PASSWORD"

	//SSHPuplicKeyFileName is the Name of the SSH Public key file
	SSHPuplicKeyFileName string = "ssh.pub"

	//SSHPrivateKeyFileName is the Name of the SSH Private key file
	SSHPrivateKeyFileName string = "ssh.pem"

	//ActionEnvVariableKey is the environment variable key used to pass
	// the ekera action to perform between several Ekara components.
	ActionEnvVariableKey string = "EKARA_ACTION"

	//HttpProxyEnvVariableKey is the environment variable key used to
	// get or store the http proxy definition
	HttpProxyEnvVariableKey string = "http_proxy"

	//HttpsProxyEnvVariableKey  is the environment variable key used to
	// get or store the https proxy definition
	HttpsProxyEnvVariableKey string = "https_proxy"

	//NoProxyEnvVariableKey is the environment variable key used to
	// get or store the http noproxy definition
	NoProxyEnvVariableKey string = "no_proxy"

	//ProxyConfigFileName is the name of the generated proxy configuration file
	ProxyConfigFileName string = "proxy_env.yml"

	//ComponentPathsFileName is the name of the file containing the map of all components locations
	ComponentPathsFileName string = "component_paths.yaml"

	//CliParametersFileName is the name of the file containing the map of all components locations
	CliParametersFileName string = "cli_parameters.yaml"

	//CreationSessionFileName is the name of the generated creation session file returned to the client
	CreationSessionFileName string = "create_session.json"

	//InstallerLogPrefix is the prefix for the installer logs
	InstallerLogPrefix string = "Ekara INSTALLER LOG:"

	//InstallerVolume is the volume location used by the installer
	InstallerVolume string = "/opt/ekara/output"

	//ParamYamlFileName is the name of any file containing params
	ParamYamlFileName string = "params.yaml"

	//ExtraVarYamlFileName is the name of any file containing extra vars
	ExtraVarYamlFileName string = "extra-vars.yaml"

	//EnvYamlFileName is the name of any file containing environment variables
	EnvYamlFileName string = "env.yaml"

	//InventoryYamlFileName is the name of any file containing host inventories
	InventoryYamlFileName string = "inventory.yaml"

	//ComponentModuleFolder is the name of any folder containing component modules
	ComponentModuleFolder string = "modules"

	//InventoryModuleFolder is the name of any folder containing inventories
	InventoryModuleFolder string = "inventory"
)
