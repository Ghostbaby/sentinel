package work

import (
	"os"
	"strings"

	"github.com/go-logr/logr"

	"github.com/ghostbaby/sentinel/pkg/models"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/ghostbaby/sentinel/pkg/g"

	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
)

const (
	DefaultRemoteUser      = "root"
	DefaultRemoteFilePath  = "/tmp/test"
	DefaultRemoteFileChmod = "0655"
)

func Scp() []*models.ScpResult {

	var list []*models.ScpResult

	config := g.Config()
	log := ctrl.Log.WithName("controllers").WithName("sentinel").WithName("scp")

	chJobs := make(chan *models.Job, len(config.IPS))
	chResults := make(chan *models.ScpResult, len(config.IPS))

	for w := 1; w <= config.Routines; w++ {
		go ScpFileWork(chJobs, chResults)
	}

	for ip, _ := range config.IPS {
		job := &models.Job{
			Log:        log,
			IP:         ip,
			RetryTimes: config.RetryTimes,
		}
		chJobs <- job
	}
	close(chJobs)

	for i := 1; i <= len(config.IPS); i++ {
		res := <-chResults
		list = append(list, res)
	}

	return list

}

func ScpFileWork(jobs <-chan *models.Job, results chan<- *models.ScpResult) {
	for j := range jobs {
		rt := j.RetryTimes

		for {

			result, err := ScpFile(j.Log, j.IP)
			if err != nil {
				j.Log.Error(err, "fail to exec PinCnTask")
				rt--
				if rt == 0 {
					results <- result
					break
				}
				continue
			}

			results <- result
			break
		}

	}
}

func ScpFile(log logr.Logger, ip string) (*models.ScpResult, error) {
	config := g.Config()

	result := &models.ScpResult{
		IP:      ip,
		IsReady: false,
	}
	log.Info("start to scp file", "host", ip)

	// Use SSH key authentication from the auth package
	// we ignore the host key in this example, please change this if you use this library
	clientConfig, err := auth.PrivateKey(DefaultRemoteUser, config.IdRsa, ssh.InsecureIgnoreHostKey())
	if err != nil {
		log.Error(err, "fail to create ssh client config")
		return result, err
	}

	// For other authentication methods see ssh.ClientConfig and ssh.AuthMethod
	// Create a new SCP client
	client := scp.NewClient(strings.Join([]string{ip, config.IPS[ip].Port}, ":"), &clientConfig)

	// Connect to the remote server
	if err := client.Connect(); err != nil {
		log.Error(err, "couldn't establish a connection to the remote server ")
		return result, err
	}

	// Open a file
	f, err := os.Open(config.ScpFile)
	if err != nil {
		log.Error(err, "fail to open scp source file")
		return result, err
	}

	// Close client connection after the file has been copied
	defer client.Close()

	// Close the file after it has been copied
	defer f.Close()

	// Finaly, copy the file over
	// Usage: CopyFile(fileReader, remotePath, permission)
	if err := client.CopyFile(f, DefaultRemoteFilePath, DefaultRemoteFileChmod); err != nil {
		log.Error(err, "fail to copying file to remote server", "host", ip)
		return result, err
	}

	log.Info("success to scp file", "host", ip)

	result.IsReady = true
	result.Name = config.IPS[ip].Describe
	return result, nil
}
