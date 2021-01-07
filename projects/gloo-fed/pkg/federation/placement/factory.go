package placement

type factory struct {
	podName string
}

var _ StatusBuilderFactory = factory{}

func NewFactory(podName string) StatusBuilderFactory {
	return factory{podName: podName}
}

func (f factory) GetBuilder() StatusBuilder {
	return &statusBuilder{podName: f.podName}
}
