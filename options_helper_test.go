package bazilfuse

// for TestMountOptionCommaError
func ForTestSetMountOption(conf *MountConfig, k, v string) {
	conf.options[k] = v
}
