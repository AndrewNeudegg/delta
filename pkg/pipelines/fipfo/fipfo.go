package fipfo

import (
	"context"
	"sync"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/pipelines"
	"github.com/andrewneudegg/delta/pkg/pipelines/definitions"
	resources "github.com/andrewneudegg/delta/pkg/resources"
	resourceMapping "github.com/andrewneudegg/delta/pkg/resources/builder"
	resourceDefinitions "github.com/andrewneudegg/delta/pkg/resources/definitions"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

// Pipeline is a type of pipeline that fans in inputs,
// applies processing, then fans out to the various outputs.
type Pipeline struct {
	Inputs    []definitions.PipelineNode `mapstructure:"input"`
	Processes []definitions.PipelineNode `mapstructure:"process"`
	Outputs   []definitions.PipelineNode `mapstructure:"output"`
}

// ID returns a human readable identifier for this pipeline.
func (f Pipeline) ID() string {
	return "pipelines/fipfo"
}

func (f Pipeline) buildInput(n definitions.PipelineNode) (resourceDefinitions.Input, error) {
	// if this node has sub processes build a meta input, otherwise return the process.
	if len(n.Nodes) > 0 {
		sub := make([]resourceDefinitions.Input, 0)
		for _, sn := range n.Nodes {
			snR, err := f.buildInput(sn)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to build subnode for '%s'", n.ID)
			}
			sub = append(sub, snR)
		}

		meta, err := resources.BuildMetaResource(n.ID, n.Config, resourceMapping.MetaMapping())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to build meta resource for '%s'", n.ID)
		}

		return meta.I(sub)
	}

	return resources.BuildInputResource(n.ID, n.Config, resourceMapping.InputMapping())
}

func (f Pipeline) inputs(ctx context.Context) (chan events.Collection, error) {
	sourceChannels := make([]chan events.Collection, 0)
	outputCh := make(chan events.Collection)

	for _, inputDef := range f.Inputs {
		thisSourceChan := make(chan events.Collection)

		input, err := f.buildInput(inputDef)

		if err != nil {
			return nil, errors.Wrapf(err, "failed to BuildInputResource, for input with ID '%s'", inputDef.ID)
		}

		go func(i resourceDefinitions.Input, ch chan events.Collection) {
			log.Infof("launching input '%s'", i.ID())

			err := i.DoInput(context.TODO(), ch)
			if err != nil {
				log.Error(errors.Wrap(err, "failed to do input"))
			}
			log.Warnf("input '%s' has exited", i.ID())
		}(input, thisSourceChan)

		sourceChannels = append(sourceChannels, thisSourceChan)
	}
	log.Infof("found '%d' sources", len(f.Inputs))

	go pipelines.FanIn(ctx, sourceChannels, outputCh)

	return outputCh, nil
}

func (f Pipeline) buildProcess(n definitions.PipelineNode) (resourceDefinitions.Process, error) {
	// if this node has sub processes build a meta input, otherwise return the process.
	if len(n.Nodes) > 0 {
		sub := make([]resourceDefinitions.Process, 0)
		for _, sn := range n.Nodes {
			snR, err := f.buildProcess(sn)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to build subnode for '%s'", n.ID)
			}
			sub = append(sub, snR)
		}

		meta, err := resources.BuildMetaResource(n.ID, n.Config, resourceMapping.MetaMapping())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to build meta resource for '%s'", n.ID)
		}

		return meta.P(sub)
	}

	return resources.BuildProcessResource(n.ID, n.Config, resourceMapping.ProcessMapping())
}

func (f Pipeline) processes(ch chan events.Collection) (chan events.Collection, error) {
	var previousSourceOutput *chan events.Collection
	previousSourceOutput = &ch

	for _, processDef := range f.Processes {

		process, err := f.buildProcess(processDef)

		if err != nil {
			return nil, errors.Wrapf(err, "failed to BuildProcessResource, for process with ID '%s'", processDef.ID)
		}

		thisRelayOutputChan := make(chan events.Collection)

		go func(p resourceDefinitions.Process, inCh <-chan events.Collection, outCh chan<- events.Collection) {
			log.Infof("launching process '%s'", p.ID())
			err := p.DoProcess(context.TODO(), inCh, outCh)
			if err != nil {
				log.Error(errors.Wrap(err, "failed to do process"))
			}
			log.Warnf("process '%s' has exited", p.ID())

		}(process, *previousSourceOutput, thisRelayOutputChan)
		previousSourceOutput = &thisRelayOutputChan
	}
	log.Infof("found '%d' processes", len(f.Processes))

	return *previousSourceOutput, nil
}

func (f Pipeline) buildOutputs(n definitions.PipelineNode) (resourceDefinitions.Output, error) {
	// if this node has sub processes build a meta input, otherwise return the process.
	if len(n.Nodes) > 0 {
		sub := make([]resourceDefinitions.Output, 0)
		for _, sn := range n.Nodes {
			snR, err := f.buildOutputs(sn)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to build subnode for '%s'", n.ID)
			}
			sub = append(sub, snR)
		}

		meta, err := resources.BuildMetaResource(n.ID, n.Config, resourceMapping.MetaMapping())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to build meta resource for '%s'", n.ID)
		}

		return meta.O(sub)
	}

	return resources.BuildOutputResource(n.ID, n.Config, resourceMapping.OutputMapping())
}

func (f Pipeline) outputs(ctx context.Context, ch chan events.Collection) error {
	distributorChannels := make([]chan events.Collection, len(f.Processes))

	for _, outputDef := range f.Outputs {
		distributorInputChannel := make(chan events.Collection)

		output, err := f.buildOutputs(outputDef)
		if err != nil {
			return errors.Wrapf(err, "failed to BuildOutputResource, for output with ID '%s'", outputDef.ID)
		}

		go func(o resourceDefinitions.Output, ch chan events.Collection) {
			log.Infof("launching output '%s'", o.ID())
			err := o.DoOutput(context.TODO(), ch)
			if err != nil {
				log.Error(errors.Wrap(err, "failed to do distributor"))
			}
			log.Warnf("output '%s' has exited", o.ID())
		}(output, distributorInputChannel)
		distributorChannels = append(distributorChannels, distributorInputChannel)
	}
	log.Infof("found '%d' output", len(f.Outputs))

	go pipelines.FanOut(ctx, ch, distributorChannels)

	return nil
}

// Do constructs and initialises all elements of this pipeline.
func (f Pipeline) Do(ctx context.Context) error {
	wg := sync.WaitGroup{}
	wg.Add(1)

	ch1, err := f.inputs(ctx)
	if err != nil {
		return errors.Wrap(err, "failed inputs")
	}

	ch2, err := f.processes(ch1)
	if err != nil {
		return errors.Wrap(err, "failed processes")
	}

	err = f.outputs(ctx, ch2)
	if err != nil {
		return errors.Wrap(err, "failed outputs")
	}

	<-ctx.Done()
	return ctx.Err()
}
