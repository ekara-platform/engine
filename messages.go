package engine

const (

	//ErrorParsingDescriptor indicates than an error occured parsing the environment descriptor
	ErrorParsingDescriptor string = "Error parsing the descriptor %s, run the \"cli check\" command to get the details"
	//ErrorParsingEnvironment indicates than an error occured parsing the environment
	ErrorParsingEnvironment string = "Error parsing the environment: %s"
	//ErrorCreatingReportFile indicates than an error occured creating the execution report file
	ErrorCreatingReportFile string = "Error creating the report file  %s"
	//ErrorCreatingDumpFile indicates than an error occured creating the environment dump file
	ErrorCreatingDumpFile string = "Error creating the dump file  %s"

	//ErrorAddingExchangeFolder indicates than an error occured adding an exchange folder
	ErrorAddingExchangeFolder string = "An error occurred adding the exchange folder %s: %s"
	//ErrorCreatingExchangeFolder indicates than an error occured creating an exchange folder
	ErrorCreatingExchangeFolder string = "An error occurred creating the exchange folder %s: %s"
	//ErrorReadingReport indicates than an error occured reading the execution report file
	ErrorReadingReport string = "error reading the report file \"%s\", error \"%s\""
	//ErrorUnmarshallingReport indicates than an error occured unmarshalling the execution report file
	ErrorUnmarshallingReport string = "error Unmarshalling the report file \"%s\", error \"%s\""
	//ErrorGeneric  indicates than an error occured
	ErrorGeneric string = "An error occurred  %s:"

	//LogValidationSuccessful indicates that the descriptor validation was sucessful
	LogValidationSuccessful string = "The envinronment descriptor validation is successful!"
	//LogProcessingNode  indicates that a node is being processed
	LogProcessingNode string = "Processing node: %s\n"
	//LogRunningSetupFor  indicates that a setup step is being processes
	LogRunningSetupFor string = "Running setup for provider %s"
	//LogLaunchingAction  indicates that an action is being launched
	LogLaunchingAction string = "Engine, launching action %s\n"
	//LogRunningAction  indicates that an action is running
	LogRunningAction string = "Engine, running action %s\n"
	//LogProcessingStackPlaybook  indicates that a stack is being deploy using a playbook
	LogProcessingStackPlaybook string = "Processing playbook for stack: %s on node: %s"
	//LogProcessingStackCompose  indicates that a stack is being deploy using a docker compose file
	LogProcessingStackCompose string = "Processing Docker Compose for stack: %s on node: %s"
	//LogReportWritten   indicates where the execution report file has been written
	LogReportWritten string = "The execution report file has been written in %s\n"
)
