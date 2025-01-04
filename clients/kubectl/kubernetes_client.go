package kubectl

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/living-etc/sightseer.go/entities/kubernetes"
)

func check(err error, message string) {
	if err != nil {
		log.Fatalf("%v: %v", message, err)
	}
}

type KubernetesClient struct {
	KubeConfigPath         string
	CaCertPath             string
	ApiServierHost         string
	KubectlCommandExecutor CommandExecutor
}

func NewKubernetesClient(
	kubeConfigPath string,
	caCertPath string,
	apiServierHost string,
) KubernetesClient {
	return KubernetesClient{
		KubeConfigPath:         kubeConfigPath,
		CaCertPath:             caCertPath,
		ApiServierHost:         apiServierHost,
		KubectlCommandExecutor: RealCommandExecutor{},
	}
}

func (client KubernetesClient) Version() kubernetes.KubernetesVersionInfo {
	caCert, err := os.ReadFile(client.CaCertPath)
	check(err, "")

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	host := fmt.Sprintf("%v:6443/version", client.ApiServierHost)
	resp, err := httpClient.Get(host)
	check(err, "")
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	check(err, "")

	var kubernetesVersionInfo kubernetes.KubernetesVersionInfo

	err = json.Unmarshal([]byte(body), &kubernetesVersionInfo)
	kubernetesVersionInfo.Full = fmt.Sprintf(
		"%v.%v",
		kubernetesVersionInfo.Major,
		kubernetesVersionInfo.Minor,
	)
	check(err, "")

	return kubernetesVersionInfo
}

func (client KubernetesClient) kubectlCommand(command string) (string, error) {
	binary := "kubectl"
	args := []string{
		"get",
		"nodes",
		"--kubeconfig",
		client.KubeConfigPath,
		"--no-headers",
	}

	output, err := client.KubectlCommandExecutor.ExecuteCommand(binary, args)
	if err != nil {
		return "", errors.New(
			fmt.Sprintf(
				"Error running command...\nCommand: %v %v\n%v",
				binary,
				args,
				output,
			),
		)
	}

	o := strings.TrimSuffix(string(output), "\n")

	return o, nil
}

func (client KubernetesClient) Workers() ([]kubernetes.Worker, error) {
	got, err := client.kubectlCommand("get nodes")
	if err != nil {
		return []kubernetes.Worker{}, err
	}

	var workers []kubernetes.Worker
	scanner := bufio.NewScanner(strings.NewReader(got))
	for scanner.Scan() {
		pattern := `(?P<Name>.*?)\s+(?P<Status>\w+)\s+(?P<Roles>.*?)\s+(?P<Age>.*?)\s+(?P<Version>.+)`
		line := scanner.Text()

		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(line)

		result := make(map[string]string)
		for i, name := range re.SubexpNames() {
			if i > 0 {
				result[name] = matches[i]
			}
		}
		worker := kubernetes.Worker{
			Name:   result["Name"],
			Status: result["Status"],
		}
		workers = append(workers, worker)
	}

	return workers, nil
}
