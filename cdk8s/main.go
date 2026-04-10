package main

import (
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"

	cdk8splus33 "github.com/cdk8s-team/cdk8s-plus-go/cdk8splus33/v2"
)

func main() {
	app := cdk8s.NewApp(&cdk8s.AppProps{
		YamlOutputType: cdk8s.YamlOutputType_FILE_PER_RESOURCE,
	})

	chart := cdk8s.NewChart(app, jsii.String("ubuntu"), &cdk8s.ChartProps{
		DisableResourceNameHashes: jsii.Bool(true),
	})

	cdk8splus33.NewDeployment(chart, jsii.String("Deployment"), &cdk8splus33.DeploymentProps{
		Replicas: jsii.Number(3),
		Containers: &[]*cdk8splus33.ContainerProps{{
			Image: jsii.String("ubuntu"),
		}},
	})

	app.Synth()
}
