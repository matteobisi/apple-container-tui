package main

import (
	"bytes"
	"testing"
)

func TestVersionOutput(t *testing.T) {
	originalVersion := version
	t.Cleanup(func() { version = originalVersion })

	for _, testCase := range []struct {
		name    string
		version string
		want    string
	}{
		{name: "development build", version: "dev", want: "actui version dev\n"},
		{name: "linker injected release build", version: "0.1.13", want: "actui version 0.1.13\n"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			version = testCase.version
			command := newRootCmd()
			output := &bytes.Buffer{}
			command.SetOut(output)
			command.SetArgs([]string{"--version"})

			if err := command.Execute(); err != nil {
				t.Fatalf("execute --version: %v", err)
			}
			if got := output.String(); got != testCase.want {
				t.Errorf("version output = %q, want %q", got, testCase.want)
			}
		})
	}
}
