package bazilfuse

// for TestMountOptionCommaError
func ForTestSetMountOption(conf *mountConfig, k, v string) {
	conf.options[k] = v
}
