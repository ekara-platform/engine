package engine

const (

	// Error messages
	ERROR_PARSING_DESCRIPTOR       string = "Error parsing the descriptor %s, run the \"cli check\" command to get the details"
	ERROR_PARSING_ENVIRONMENT      string = "Error parsing the environment: %s"
	ERROR_CREATING_REPORT_FILE     string = "Error creating the report file  %s"
	ERROR_ADDING_EXCHANGE_FOLDER   string = "An error occurred adding the exchange folder %s: %s"
	ERROR_CREATING_EXCHANGE_FOLDER string = "An error occurred creating the exchange folder %s: %s"
	ERROR_READING_REPORT           string = "error reading the report file \"%s\", error \"%s\""
	ERROR_UNMARSHALLING_REPORT     string = "error Unmarshalling the report file \"%s\", error \"%s\""
	ERROR_GENERIC                  string = "An error occurred  %s:"

	// Log messages
	LOG_VALIDATION_SUCCESSFUL     string = "The envinronment descriptor validation is successful!"
	LOG_PROCESSING_NODE           string = "Processing node: %s"
	LOG_RUNNING_SETUP_FOR         string = "Running setup for provider %s"
	LOG_LAUNCHING_ACTION          string = "Engine, launching action %s"
	LOG_RUNNING_ACTION            string = "Engine, running action %s"
	LOG_PROCESSING_STACK_PLAYBOOK string = "Processing playbook for stack: %s on node: %s"
	LOG_PROCESSING_STACK_COMPOSE  string = "Processing Docker Compose for stack: %s on node: %s"
	LOG_REPORT_WRITTEN            string = "The execution report file has been written in %s\n"
)
