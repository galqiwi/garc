package ssh

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/gliderlabs/ssh"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
	"net"
	"os"
	"path"
	"strings"
	"sync"
	"testing"
	"time"
)

func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}

	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func runSingleResponseServer(t *testing.T, wg *sync.WaitGroup) int {
	var server *ssh.Server

	port, err := getFreePort()
	require.NoError(t, err)

	server = &ssh.Server{Addr: fmt.Sprintf(":%v", port), Handler: func(session ssh.Session) {
		_, err := fmt.Fprint(session, "test")
		require.NoError(t, err)
		go func() {
			time.Sleep(time.Second)
			err := server.Shutdown(context.Background())
			require.NoError(t, err)
		}()
	}}

	wg.Add(1)
	go func() {
		_ = server.ListenAndServe()
		wg.Done()
	}()

	return port
}

func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {

	block := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(privateKey),
	}

	privatePEM := pem.EncodeToMemory(&block)

	return privatePEM
}

func createSSHKey(t *testing.T, keyPath string) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	require.NoError(t, err)
	require.NoError(t, privateKey.Validate())
	privateKeyPEM := encodePrivateKeyToPEM(privateKey)
	require.NoError(t, os.WriteFile(keyPath, privateKeyPEM, 0600))
}

func TestClient(t *testing.T) {
	defer goleak.VerifyNone(t)

	dirname, err := os.MkdirTemp("", "")
	require.NoError(t, err)
	defer func() {
		err := os.RemoveAll(dirname)
		require.NoError(t, err)
	}()

	keyPath := path.Join(dirname, "id_rsa")

	createSSHKey(t, keyPath)

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	port := runSingleResponseServer(t, wg)
	cmd := NewClient(&ClientConfig{KeyPath: keyPath}).
		Host("localhost", port).
		Command("username", "echo test")

	stdout := strings.Builder{}
	cmd.SetStdout(&stdout)
	require.NoError(t, cmd.Run())
	require.Equal(t, "test", stdout.String())
}
