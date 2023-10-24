package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/galqiwi/garc/internal/utils/ssh"
	"github.com/galqiwi/garc/internal/utils/ssh/ssh_utils"
	"golang.org/x/sync/errgroup"
	"io"
	"path"
)

type LimboClient struct {
	cfg *LimboConfig
}

func (c *LimboClient) getSshHost() ssh.Host {
	sshClient := ssh.NewClient(&ssh.ClientConfig{KeyPath: c.cfg.Client.KeyPath})
	return sshClient.Host(c.cfg.Hostname, c.cfg.Port)
}

func (c *LimboClient) ListRemotes() ([]string, error) {
	sshHost := c.getSshHost()
	return ssh_utils.ListRemoteDir(sshHost, c.cfg.Username, c.cfg.Path)
}

func (c *LimboClient) DoesRemoteExist(remote string) (bool, error) {
	remotes, err := c.ListRemotes()
	if err != nil {
		return false, err
	}

	for _, existingRemote := range remotes {
		if remote == existingRemote {
			return true, nil
		}
	}
	return false, nil
}

func (c *LimboClient) GetRemoteMeta(remoteUUID string) (*ArchiveMeta, error) {
	sshHost := c.getSshHost()

	buf := &bytes.Buffer{}
	err := ssh_utils.ReadRemoteFile(
		sshHost, c.cfg.Username,
		path.Join(c.cfg.Path, remoteUUID, LimboMetadataFile), buf)
	if err != nil {
		return nil, err
	}

	var output ArchiveMeta
	err = json.Unmarshal(buf.Bytes(), &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (c *LimboClient) GetRemoteSize(remoteUUID string) (int64, error) {
	sshHost := c.getSshHost()

	return ssh_utils.GetFileSize(
		sshHost,
		c.cfg.Username,
		path.Join(c.cfg.Path, remoteUUID, "tarball"),
	)
}

func (c *LimboClient) ReadRemoteTarball(remoteUUID string, to io.Writer) error {
	sshHost := c.getSshHost()
	remotePath := path.Join(c.cfg.Path, remoteUUID)

	return ssh_utils.ReadRemoteFile(
		sshHost,
		c.cfg.Username, path.Join(remotePath, "tarball"),
		to,
	)
}

func (c *LimboClient) RemoveRemote(remoteUUID string) error {
	sshHost := c.getSshHost()
	remotePath := path.Join(c.cfg.Path, remoteUUID)

	return ssh_utils.RemoveRemoteDir(
		sshHost,
		c.cfg.Username,
		remotePath,
	)
}

func (c *LimboClient) createRemoteDir(meta *ArchiveMeta) error {
	sshHost := c.getSshHost()
	remotePath := path.Join(c.cfg.Path, fmt.Sprint(meta.Id))

	err := ssh_utils.CreateRemoteDir(sshHost, c.cfg.Username, remotePath)
	if err != nil {
		return err
	}
	return nil
}

func (c *LimboClient) createRemoteMetaFile(meta *ArchiveMeta) error {
	sshHost := c.getSshHost()
	remotePath := path.Join(c.cfg.Path, fmt.Sprint(meta.Id))
	remoteMetaPath := path.Join(remotePath, LimboMetadataFile)

	metaData, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	err = ssh_utils.CreateRemoteFile(sshHost, c.cfg.Username, remoteMetaPath, bytes.NewReader(metaData))
	if err != nil {
		return err
	}
	return nil
}

func (c *LimboClient) createRemoteTarFile(meta *ArchiveMeta, tarballSrc io.Reader) error {
	sshHost := c.getSshHost()
	remotePath := path.Join(c.cfg.Path, fmt.Sprint(meta.Id))
	remoteTarballPath := path.Join(remotePath, "tarball")

	return ssh_utils.CreateRemoteFile(sshHost, c.cfg.Username, remoteTarballPath, tarballSrc)
}

func (c *LimboClient) CreateRemote(meta *ArchiveMeta, tarballSrc io.Reader) error {
	err := c.createRemoteDir(meta)
	if err != nil {
		return err
	}
	fmt.Println("created remote dir")

	err = c.createRemoteMetaFile(meta)
	if err != nil {
		return err
	}
	fmt.Println("created remote meta file")

	tarReader, tarWriter := io.Pipe()

	g := new(errgroup.Group)
	g.Go(func() error {
		_, err := io.Copy(tarWriter, tarballSrc)
		if err != nil {
			return err
		}
		return tarWriter.Close()
	})

	g.Go(func() error {
		err = c.createRemoteTarFile(meta, tarReader)
		if err != nil {
			return err
		}
		fmt.Println("created remote tarball")
		return err
	})

	return g.Wait()
}

func NewLimboClient(cfg *LimboConfig) *LimboClient {
	return &LimboClient{
		cfg: cfg,
	}
}
