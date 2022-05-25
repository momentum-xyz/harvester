package example_harvester

type ExampleHarvester struct{}

func NewExampleHarvester() (*ExampleHarvester, error) {
	return &ExampleHarvester{}, nil
}

func (eh *ExampleHarvester) Start() {

}
func (eh *ExampleHarvester) Stop() {}
