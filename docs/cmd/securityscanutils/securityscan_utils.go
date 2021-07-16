package securityscanutils

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"github.com/imroc/req"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/rotisserie/eris"
	"github.com/solo-io/go-utils/log"
	"k8s.io/apimachinery/pkg/util/version"
)

func BuildSecurityScanReportGloo(tags []string) error {
	// tags are sorted by minor version
	latestTag := tags[0]
	prevMinorVersion, _ := version.ParseSemantic(latestTag)
	for ix, tag := range tags {
		semver, err := version.ParseSemantic(tag)
		if err != nil {
			return err
		}
		if ix == 0 || semver.Minor() != prevMinorVersion.Minor() {
			fmt.Printf("\n***Latest %d.%d.x Gloo Open Source Release: %s***\n\n", semver.Major(), semver.Minor(), tag)
			err = printImageReportGloo(tag)
			if err != nil {
				return err
			}
			prevMinorVersion = semver
		} else {
			fmt.Printf("<details><summary> Release %s </summary>\n\n", tag)
			err = printImageReportGloo(tag)
			if err != nil {
				return err
			}
			fmt.Println("</details>")
		}
	}

	return nil
}

func BuildSecurityScanReportGlooE(tags []string) error {
	// tags are sorted by minor version
	latestTag := tags[0]
	prevMinorVersion, _ := version.ParseSemantic(latestTag)
	for ix, tag := range tags {
		semver, err := version.ParseSemantic(tag)
		if err != nil {
			return err
		}
		if ix == 0 || semver.Minor() != prevMinorVersion.Minor() {
			fmt.Printf("\n***Latest %d.%d.x Gloo Enterprise Release: %s***\n\n", semver.Major(), semver.Minor(), tag)
			err = printImageReportGlooE(semver)
			if err != nil {
				return err
			}
			prevMinorVersion = semver
		} else {
			fmt.Printf("<details><summary>Release %s </summary>\n\n", tag)
			err = printImageReportGlooE(semver)
			if err != nil {
				return err
			}
			fmt.Println("</details>")
		}
	}

	return nil
}

func printImageReportGloo(tag string) error {
	images := []string{"gateway", "discovery", "gloo", "gloo-envoy-wrapper", "ingress", "access-logger", "sds", "certgen"}
	for _, image := range images {
		fmt.Printf("**Gloo %s image**\n\n", image)
		// remove `v` from tag
		url := "https://storage.googleapis.com/solo-gloo-security-scans/gloo/" + tag[1:] + "/" + image + "_cve_report.docgen"
		report, err := GetSecurityScanReport(url)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n\n", report)
	}
	return nil
}

func printImageReportGlooE(semver *version.Version) error {
	tag := semver.String()
	images := []string{"rate-limit-ee", "gloo-ee", "gloo-ee-envoy-wrapper", "observability-ee", "extauth-ee", "ext-auth-plugins"}

	// gloo-fed images replaced grpcserver images in 1.7+
	grpcserverImages := []string{"grpcserver-ee", "grpcserver-envoy", "grpcserver-ui"}
	fedImages := []string{"gloo-fed", "gloo-fed-apiserver", "gloo-fed-apiserver-envoy", "gloo-federation-console", "gloo-fed-rbac-validating-webhook"}
	hasFedVersion, _ := semver.Compare("1.7.0")
	if hasFedVersion >= 0 {
		images = append(images, fedImages...)
	} else {
		images = append(images, grpcserverImages...)
	}

	for _, image := range images {
		fmt.Printf("**Gloo Enterprise %s image**\n\n", image)
		// remove `v` from tag
		url := "https://storage.googleapis.com/solo-gloo-security-scans/glooe/" + tag + "/" + image + "_cve_report.docgen"
		report, err := GetSecurityScanReport(url)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n\n", report)
	}
	return nil
}

func GetSecurityScanReport(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	var report string
	if resp.StatusCode == http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		report = string(bodyBytes)
	} else if resp.StatusCode == http.StatusNotFound {
		// Older releases may be missing scan results
		report = "No scan found\n"
	}
	resp.Body.Close()

	return report, nil
}

func RunSecurityScanOnVersions(tags, images []string, outputDirSuffix string) error {
	// Create temporary dir to output security scans to
	// e.g. $SCAN_DIR/gloo/markdown_results, $SCAN_DIR/gloo/sarif_results
	scanDir := path.Join(os.Getenv("SCAN_DIR"), outputDirSuffix)
	markdownDir, sarifDir := path.Join(scanDir, "markdown_results"), path.Join(scanDir, "sarif_results")
	if _, err := os.Stat(markdownDir); os.IsNotExist(err) {
		err := os.MkdirAll(markdownDir, os.ModePerm)
		if err != nil {
			return eris.Wrapf(err, "error creating dir %s", markdownDir)
		}
	}
	if _, err := os.Stat(sarifDir); os.IsNotExist(err) {
		err := os.MkdirAll(sarifDir, os.ModePerm)
		if err != nil {
			return eris.Wrapf(err, "error creating dir %s", sarifDir)
		}
	}

	for _, tag := range tags {
		fmt.Printf("Scanning version %s\n", tag)
		err := ScanVersionMd(images, tag, markdownDir)
		if err != nil {
			log.Fatalf(err.Error())
		}
		err = ScanVersionSarif(images, tag, sarifDir)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
	return nil
}

func ScanVersionMd(images []string, version, outputDir string) error {
	// Create a separate directory for each version we are scanning
	outputDir = path.Join(outputDir, version)
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.Mkdir(outputDir, os.ModePerm)
	}
	for _, image := range images {
		fileName := fmt.Sprintf("%s_cve_report.docgen", image)
		templateFile := "@hack/utils/security_scan_report/markdown.tpl"
		err := CreateScanFileAndRunScan(outputDir, fileName, image, version, templateFile)
		if err != nil {
			return err
		}
	}
	return nil
}

func ScanVersionSarif(images []string, version, outputDir string) error {
	for _, image := range images {
		fileName := fmt.Sprintf("%s_%s_cve_report.sarif", version, image)
		templateFile := "@hack/utils/security_scan_report/sarif.tpl"
		err := CreateScanFileAndRunScan(outputDir, fileName, image, version, templateFile)
		if err != nil {
			return err
		}
	}
	return nil
}

// Creates File for trivy scan and runs the scan
func CreateScanFileAndRunScan(fileDir, fileName, image, versionTag, templateFile string) error {
	file, err := os.Create(path.Join(fileDir, fileName))
	if err != nil {
		return eris.Wrap(err, "unable to create temp file for sarif file")
	}
	imageRepo := os.Getenv("IMAGE_REPO")
	if imageRepo == "" {
		return fmt.Errorf("no IMAGE_REPO environment variable specified")
	}
	imageToScan := fmt.Sprintf("%s/%s:%s", imageRepo, image, versionTag)
	err = RunTrivyScan(imageToScan, versionTag, templateFile, file.Name())
	if err != nil {
		return err
	}
	err = UploadSecurityScanToGithub(fileName, versionTag)
	return err
}

type SarifMetadata struct {
	Ref       string `json:"ref"`
	CommitSha string `json:"commit_sha"`
	Sarif     string `json:"sarif"`
}

type Response struct {
	ShaObject `json:"object"`
}

type ShaObject struct {
	Sha string `json:"sha"`
}

func UploadSecurityScanToGithub(fileName, versionTag string) error {
	resp, err := req.Get("https://api.github.com/repos/solo-io/gloo/git/refs/tags/v" + versionTag)
	if err != nil {
		return eris.Wrapf(err, "Unable to get commit for version v%s", versionTag)
	}
	shaResp := &Response{}
	resp.ToJSON(shaResp)
	fmt.Printf("%+v\n", shaResp)
	//cmd := exec.Command("git", "rev-parse", fmt.Sprintf("refs/tags/v%s", versionTag))
	//out, err := cmd.Output()
	//if err != nil {
	//	fmt.Println(string(out))
	//	return eris.Wrapf(err, "error getting commit sha for tag %s", versionTag)
	//}
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		return fmt.Errorf("no GITHUB_TOKEN environment variable specified")
	}
	b, err := ioutil.ReadFile(fileName)
	var r bytes.Buffer
	w := gzip.NewWriter(&r)
	_, err = w.Write(b)
	if err != nil {
		return eris.Wrap(err, "error writing gzip file")
	}
	w.Close()
	sarifMetadata := SarifMetadata{
		Ref:       fmt.Sprintf("refs/tags/v%s", versionTag),
		CommitSha: shaResp.Sha,
		Sarif:     base64.StdEncoding.EncodeToString(r.Bytes()),
	}
	header := req.Header{
		"Authorization": fmt.Sprintf("token %s", githubToken),
		"Content-Type":  "application/json",
	}
	res, err := req.Post("https://api.github.com/repos/solo-io/gloo/code-scanning/sarifs", req.BodyJSON(sarifMetadata), header)
	fmt.Println(res.String())
	if err != nil || res.Response().StatusCode != 200 {
		return eris.Wrapf(err, "error uploading sarif file to github, response: \n%s", res)
	}
	return nil
}

// Runs trivy scan command
func RunTrivyScan(image, version, templateFile, output string) error {
	args := []string{"image", "--severity", "HIGH,CRITICAL", "--format", "template", "--template", templateFile, "--output", output, image}
	cmd := exec.Command("trivy", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		// delete empty trivy output file that may have been created
		_ = os.Remove(output)
		// swallow error if image is not found error, so that we can continue scanning releases
		// even if some releases failed and we didn't publish images for those releases
		if IsImageNotFoundErr(string(out)) {
			return nil
		}
		return eris.Wrapf(err, "error running trivy scan on image %s, version %s, Logs: \n%s", image, version, string(out))
	}
	return nil
}

func IsImageNotFoundErr(logs string) bool {
	if strings.Contains(logs, "No such image: ") {
		return true
	}
	return false
}
