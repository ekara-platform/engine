package engine

const (
	//ValidationOutputFile is the name of the file resulting of the environment descriptor validation
	ValidationOutputFile = "validation_output.json"
	//DumpSourceYamlOutputFile is the name of the file where the complete yaml, coming from
	// the environment descriptor parsing, will be written
	DumpSourceYamlOutputFile = "dump_output_source.yaml"
	//DumpYamlOutputFile is the name of the file where a lighter  yaml, coming from
	// the environment descriptor parsing, will be written. The lighter  version of the dump contains
	// only the non empty stuff really used to define an environment
	DumpYamlOutputFile = "dump_output.yaml"
	//ReportOutputFile is the name of the file reporting the engine execution details
	ReportOutputFile = "report_output.json"
)
