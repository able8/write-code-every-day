package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func main() {
	fmt.Println(CombinedOutputString("cal", nil, WithEnv("name", "able")))
}

func RunCommand(cmd string, arg []string, opts ...Option) error {
	c := exec.Command(cmd, arg...)
	applyOptions(c, opts)
	return c.Run()
}

func CombinedOutput(cmd string, arg []string, opts ...Option) ([]byte, error) {
	c := exec.Command(cmd, arg...)
	applyOptions(c, opts)
	return c.CombinedOutput()
}

func CombinedOutputString(cmd string, arg []string, opts ...Option) (string, error) {
	output, err := CombinedOutput(cmd, arg, opts...)
	return string(output), err
}

func Output(cmd string, arg []string, opts ...Option) ([]byte, error) {
	c := exec.Command(cmd, arg...)
	applyOptions(c, opts)
	return c.Output()
}

func OutputString(cmd string, arg []string, opts ...Option) (string, error) {
	output, err := Output(cmd, arg, opts...)
	return string(output), err
}

func SeparateOutput(cmd string, arg []string, opts ...Option) ([]byte, []byte, error) {
	var stdout, stderr bytes.Buffer
	c := exec.Command(cmd, arg...)
	applyOptions(c, opts)
	c.Stdout = &stdout
	c.Stderr = &stderr
	err := c.Run()
	return stdout.Bytes(), stderr.Bytes(), err
}

func SeparateOutputString(cmd string, arg []string, opts ...Option) (string, string, error) {
	stdout, stderr, err := SeparateOutput(cmd, arg, opts...)
	return string(stdout), string(stderr), err
}

type Option func(*exec.Cmd)

func WithStdin(stdin io.Reader) Option {
	return func(c *exec.Cmd) {
		c.Stdin = stdin
	}
}

func Without(stdout io.Writer) Option {
	return func(c *exec.Cmd) {
		c.Stdout = stdout
	}
}

func WithStderr(stderr io.Writer) Option {
	return func(c *exec.Cmd) {
		c.Stderr = stderr
	}
}

func WithOutWriter(out io.Writer) Option {
	return func(c *exec.Cmd) {
		c.Stdout = out
		c.Stderr = out
	}
}

func WithEnv(key, value string) Option {
	return func(c *exec.Cmd) {
		c.Env = append(os.Environ(), fmt.Sprintf("%s=%s", key, value))
	}
}

func applyOptions(cmd *exec.Cmd, opts []Option) {
	for _, opt := range opts {
		opt(cmd)
	}
}
