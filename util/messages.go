package util

const (

	// Errors messages
	ERROR                         string = "Error: %s"
	ERROR_MISSING                 string = "%s is missing"
	ERROR_CREATING_CONFIG_FILE    string = "Error creating the configuration file: %s, %s"
	ERROR_GENERATING_PROXY_CONFIG string = "Error generating the proxy configuration: %s"
	ERROR_REQUIRED_ENV            string = "the environment variable \"%s\" should be defined"

	// Log messages
	LOG_SAVING             string = "Saving: %s"
	LOG_CREATING_FILE      string = "Creation of the file %s"
	LOG_LAUNCHING_PLAYBOOK string = "Launching playbook %s"
	LOG_STARTING_PLAYBOOK  string = "Starting playbook %s"
)
