package cmd_test

import (
	"testing"

	"github.com/matryer/is"

	"go.jlucktay.dev/golang-workbench/dtm/cmd"
	"go.jlucktay.dev/golang-workbench/dtm/docker"
)

func TestExecute(t *testing.T) {
	is := is.New(t)
	ctx := t.Context()

	dmc := docker.NewMockClient(t)
	creator := func() (docker.Client, error) { return dmc, nil }

	dmc.EXPECT().ListImages(ctx).Return([]string{"one", "two", "three"}, nil).Once()
	dmc.EXPECT().DeleteImage(ctx, "img string").Return(nil).Once()

	is.NoErr(cmd.Execute(ctx, creator))

	is.True(dmc.AssertNumberOfCalls(t, "ListImages", 1))
	is.True(dmc.AssertNumberOfCalls(t, "DeleteImage", 1))
}
