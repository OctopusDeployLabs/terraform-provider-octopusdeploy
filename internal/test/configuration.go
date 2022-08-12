package test

func GetConfiguration(configurations []string) string {
	output := ""
	for _, configuration := range configurations {
		output += configuration + "\n"
	}
	return output
}
