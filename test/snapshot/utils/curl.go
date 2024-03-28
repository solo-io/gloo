package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/rotisserie/eris"
	"github.com/solo-io/go-utils/log"
	"github.com/solo-io/go-utils/threadsafe"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type CurlCall interface {
	SetHeaders([]string)
	Execute() (string, error)
}

type KubeContext struct {
	Context     string
	ClusterName string

	KubernetesClients *kubernetes.Clientset
}

func (k *KubeContext) GetPod(ctx context.Context, ns string, labelSelector map[string]string) (corev1.Pod, error) {
	pods, err := k.GetPods(ctx, ns, labelSelector)
	if err != nil {
		// will return error if no pods are found for given selector
		return corev1.Pod{}, err
	}
	return pods[0], nil
}

func (k *KubeContext) GetPods(ctx context.Context, ns string, labelSelector map[string]string) ([]corev1.Pod, error) {
	var buildSelector []string
	for lKey, lValue := range labelSelector {
		buildSelector = append(buildSelector, lKey+"="+lValue)
	}
	pl, err := k.KubernetesClients.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{LabelSelector: strings.Join(buildSelector, ","), FieldSelector: "status.phase==Running"})
	if err != nil {
		return nil, err
	}
	if len(pl.Items) == 0 {
		return nil, eris.New(fmt.Sprintf("No pods found for given selector %v", buildSelector))
	}

	var runningPods []corev1.Pod
	for _, pod := range pl.Items {
		// We should only be returning pods that have not been marked for deletion
		if pod.DeletionTimestamp == nil {
			runningPods = append(runningPods, pod)
		}
	}

	if len(runningPods) == 0 {
		return nil, eris.New("Only found pods marked for deletion")
	}

	return runningPods, nil
}

type CurlFromPod struct {
	Url             string
	Cluster         *KubeContext
	App             string
	Namespace       string
	Headers         []string
	TimeoutSeconds  float32
	DisableVerbose  bool
	DisableWriteOut bool
	Insecure        bool
	DataUrlEncode   string
	ExtraArgs       []string
	// Container from which curl is executed, defaults to "curl"
	Container string
}

func (c *CurlFromPod) Execute(out io.Writer) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute/2)
	defer cancel()

	args := []string{
		c.Url,
		"-s",
		"--connect-timeout",
		"5",
	}

	if !c.DisableVerbose {
		args = append(args, "-v")
	}

	if !c.DisableWriteOut {
		args = append(args, "-w", "time_total: %{time_total}s\n")
	}

	for _, header := range c.Headers {
		args = append(args, "-H", header)
	}

	if c.TimeoutSeconds > 0 {
		args = append(args, "--max-time", fmt.Sprintf("%v", c.TimeoutSeconds))
	}

	if c.Insecure {
		args = append([]string{"-k"}, args...)
	}

	if len(c.ExtraArgs) > 0 {
		args = append(args, c.ExtraArgs...)
	}

	if c.DataUrlEncode != "" {
		args = append(args, "--data-urlencode", c.DataUrlEncode)
	}

	pod, err := c.Cluster.GetPod(ctx, c.Namespace, map[string]string{"app": c.App})
	if err != nil {
		return "", err
	}

	container := "curl"
	if c.Container != "" {
		container = c.Container
	}

	output, err := Curl(ctx, out, c.Cluster.Context, pod, container, args...)
	if err != nil {
		return "", err

	}
	if _, err := out.Write([]byte(output)); err != nil {
		fmt.Fprintf(out, "failed to write output with the following error: %v", err)
	}
	return output, nil
}

func (c *CurlFromPod) SetHeaders(headers []string) {
	c.Headers = headers
}

// this curl requires a container named 'curl' to exist on the deployment
func Curl(ctx context.Context, out io.Writer, kubeContext string, curlPod corev1.Pod, container string, args ...string) (string, error) {
	args = append(
		[]string{
			"-n", curlPod.Namespace,
			"exec", curlPod.Name,
			"-c", container,
			"--", "curl",
		}, args...,
	)
	result, err := execute(ctx, out, kubeContext, args...)
	if err != nil {
		return "", err
	}
	return result, nil
}

func execute(ctx context.Context, out io.Writer, kubeContext string, args ...string) (string, error) {
	args = append([]string{"--context", kubeContext}, args...)
	fmt.Fprintf(out, "Executing: kubectl %v \n", args)
	readerChan, done, err := KubectlOutChan(&bytes.Buffer{}, args...)
	if err != nil {
		return "", err
	}
	defer close(done)
	select {
	case <-ctx.Done():
		return "", nil
	case reader := <-readerChan:
		data, err := ioutil.ReadAll(reader)
		if err != nil {
			return "", err
		}
		fmt.Fprintf(out, "<kubectl %v> output: %v\n", args, string(data))
		return string(data), nil
	}
}

func KubectlOutChan(r *bytes.Buffer, args ...string) (<-chan io.Reader, chan struct{}, error) {
	cmd := kubectl(args...)
	buf := &threadsafe.Buffer{}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, nil, err
	}
	if r != nil && r.Len() > 0 {
		go func() {
			defer stdin.Close()
			if _, err := io.Copy(stdin, r); err != nil {
				panic(err)
			}
		}()
	}

	cmd.Stdout = buf
	cmd.Stderr = buf
	log.Debugf("async running: %s", strings.Join(cmd.Args, " "))
	err = cmd.Start()
	if err != nil {
		return nil, nil, err
	}
	done := make(chan struct{})
	go func() {
		<-done // wait until done
		_ = cmd.Process.Kill()
	}()

	result := make(chan io.Reader)
	go func() {
		for {
			time.Sleep(time.Second)
			select {
			case result <- buf:
				continue
			case <-done:
				return
			default:
				continue
			}
		}
	}()

	return result, done, err
}

func kubectl(args ...string) *exec.Cmd {
	cmd := exec.Command("kubectl", args...)
	cmd.Env = os.Environ()
	// disable DEBUG=1 from getting through to kube
	for i, pair := range cmd.Env {
		if strings.HasPrefix(pair, "DEBUG") {
			cmd.Env = append(cmd.Env[:i], cmd.Env[i+1:]...)
			break
		}
	}
	return cmd
}
