package xcmd

import (
	"bytes"
	"os/exec"
)

//在指定目录执行命令，stdout、stderr返回string
func ExecCmdInDir(dir, name string, args ...string) (string, string, error) {
	stdOutBytes, stdErrBytes, err := ExecCmdInDirBytes(dir, name, args...)
	return string(stdOutBytes), string(stdErrBytes), err
}

//在指定脚本执行命令，stdout、stderr返回[]byte
func ExecCmdInDirBytes(dir, name string, args ...string) ([]byte, []byte, error) {
	bufOut := &bytes.Buffer{}
	bufErr := &bytes.Buffer{}

	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = bufOut
	cmd.Stderr = bufErr

	err := cmd.Run()
	return bufOut.Bytes(), bufErr.Bytes(), err
}
