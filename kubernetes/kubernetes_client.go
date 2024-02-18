package kubernetes

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
)

func check(err error, message string) {
	if err != nil {
		log.Fatalf("%v: %v", message, err)
	}
}

type KubernetesClient struct {
	kubeConfigPath         string
	caCertPath             string
	apiServierHost         string
	kubectlCommandExecutor CommandExecutor
}

func NewKubernetesClient(
	kubeConfigPath string,
	caCertPath string,
	apiServierHost string,
) KubernetesClient {
	return KubernetesClient{
		kubeConfigPath:         kubeConfigPath,
		caCertPath:             caCertPath,
		apiServierHost:         apiServierHost,
		kubectlCommandExecutor: RealCommandExecutor{},
	}
}

func (client KubernetesClient) Version() KubernetesVersionInfo {
	caCert, err := os.ReadFile(client.caCertPath)
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

	host := fmt.Sprintf("%v:6443/version", client.apiServierHost)
	resp, err := httpClient.Get(host)
	check(err, "")
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	check(err, "")

	var kubernetesVersionInfo KubernetesVersionInfo

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
		client.kubeConfigPath,
		"--no-headers",
	}

	output, err := client.kubectlCommandExecutor.executeCommand(binary, args)
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

func (client KubernetesClient) Workers() ([]Worker, error) {
	got, err := client.kubectlCommand("get nodes")
	if err != nil {
		return []Worker{}, err
	}

	var workers []Worker
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
		worker := Worker{
			Name:   result["Name"],
			Status: result["Status"],
		}
		workers = append(workers, worker)
	}

	return workers, nil
}
