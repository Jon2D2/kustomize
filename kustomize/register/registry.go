package register

func RegisterCoreKinds() {
	if err := NewArmadaCRDRegisterPlugin().Config(nil, []byte{}); err != nil {
		panic(err)
	}
	if err := NewKfDefCRDRegisterPlugin().Config(nil, []byte{}); err != nil {
		panic(err)
	}
	if err := NewNetworkingCRDRegisterPlugin().Config(nil, []byte{}); err != nil {
		panic(err)
	}
	if err := NewOpenShiftCRDRegisterPlugin().Config(nil, []byte{}); err != nil {
		panic(err)
	}
	if err := NewRolloutCRDRegisterPlugin().Config(nil, []byte{}); err != nil {
		panic(err)
	}
	if err := NewWorkflowCRDRegisterPlugin().Config(nil, []byte{}); err != nil {
		panic(err)
	}
	if err := NewMonitoringCRDRegisterPlugin().Config(nil, []byte{}); err != nil {
		panic(err)
	}
}
