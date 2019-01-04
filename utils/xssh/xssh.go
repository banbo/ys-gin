package xssh

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/banbo/ys-gin/errors"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

//ssh连接
type SshConn struct {
	addr       string
	user       string
	key        string
	client     *ssh.Client
	sftpClient *sftp.Client
}

func NewSshConn(ip string, port int32, user, key string) *SshConn {
	return &SshConn{
		addr: fmt.Sprintf("%s:%d", ip, port),
		user: user,
		key:  os.ExpandEnv(key),
	}
}

//测试连接
func (s *SshConn) TryConnect() error {
	_, err := s.getSshClient()
	if err != nil {
		return err
	}

	s.Close()

	return nil
}

//在远程主机执行命令
func (s *SshConn) RunCmd(cmd string) (string, string, error) {
	client, err := s.getSshClient()
	if err != nil {
		return "", "", err
	}

	//创建会话
	session, err := client.NewSession()
	if err != nil {
		return "", "", errors.NewNormal(fmt.Sprintf("创建会话失败: %v", err))
	}
	defer session.Close()

	//输出
	outBuf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}

	session.Stdout = outBuf
	session.Stderr = errBuf

	//执行命令
	if err := session.Run(cmd); err != nil {
		return "", "", errors.NewNormal(fmt.Sprintf("执行命令失败: %v", err))
	}

	return outBuf.String(), errBuf.String(), nil
}

//拷贝本地文件到远程服务器
func (s *SshConn) CopyFile(srcFile, dstFile string) error {
	sftpClient, err := s.getSftpClient()
	if err != nil {
		return err
	}

	//创建远程目录
	toPath := path.Dir(dstFile)
	_, _, err = s.RunCmd("mkdir -p " + toPath)
	if err != nil {
		return errors.NewNormal(fmt.Sprintf("创建远程目录失败: %v", err))
	}

	//打开本地文件
	srcF, err := os.Open(srcFile)
	if err != nil {
		return errors.NewNormal(fmt.Sprintf("打开本地文件失败: %v", err))
	}
	defer srcF.Close()

	//创建远程文件
	dstF, err := sftpClient.Create(dstFile)
	if err != nil {
		return errors.NewNormal(fmt.Sprintf("创建文件失败 [%s]: %v", dstFile, err))
	}
	defer dstF.Close()

	//拷贝
	n, err := io.Copy(dstF, srcF)
	if err != nil {
		return errors.NewNormal(fmt.Sprintf("拷贝文件失败: %v", err))
	}

	//判断文件大小
	fStat, _ := srcF.Stat()
	if fStat.Size() != n {
		return errors.NewNormal(fmt.Sprintf("写入文件大小错误，源文件大小：%d, 写入大小：%d", fStat.Size(), n))
	}

	return nil
}

//关闭连接
func (s *SshConn) Close() {
	if s.client != nil {
		s.client.Close()
		s.client = nil
	}

	if s.sftpClient != nil {
		s.sftpClient.Close()
		s.sftpClient = nil
	}
}

//创建ssh连接
func (s *SshConn) getSshClient() (*ssh.Client, error) {
	if s.client != nil {
		return s.client, nil
	}

	//config
	config := ssh.ClientConfig{
		User:            s.user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	//私钥
	signers := make([]ssh.Signer, 0)

	pk, err := s.readPrivateKey(s.key)
	if err != nil {
		return nil, err
	}
	signers = append(signers, pk)

	config.Auth = append(config.Auth, ssh.PublicKeys(signers...))

	//建立连接
	s.client, err = ssh.Dial("tcp", s.addr, &config)
	if err != nil {
		return nil, errors.NewNormal(fmt.Sprintf("无法连接到服务器[%s]: %v", s.addr, err))
	}
	return s.client, nil
}

//创建sftp客户端
func (s *SshConn) getSftpClient() (*sftp.Client, error) {
	if s.sftpClient != nil {
		return s.sftpClient, nil
	}

	sshClient, err := s.getSshClient()
	if err != nil {
		return nil, err
	}

	s.sftpClient, err = sftp.NewClient(sshClient, sftp.MaxPacket(1<<15))
	if err != nil {
		return nil, errors.NewNormal(fmt.Sprintf("无法创建sftp客户端: %v", err))
	}

	return s.sftpClient, nil
}

//读取私钥
func (s *SshConn) readPrivateKey(path string) (ssh.Signer, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.NewNormal(fmt.Sprintf("打开私钥文件出错: %v", err))
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, errors.NewNormal(fmt.Sprintf("读取私钥文件出错: %v", err))
	}

	signer, err := ssh.ParsePrivateKey(b)
	if err != nil {
		return nil, errors.NewNormal(fmt.Sprintf("解析私钥文件出错: %v", err))
	}

	return signer, nil
}
