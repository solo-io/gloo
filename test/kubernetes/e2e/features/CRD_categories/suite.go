package crd_categories

import (
	"bytes"
	"context"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	"github.com/stretchr/testify/suite"
	"io"
	"strings"
)

var _ e2e.NewSuiteFunc = NewTestingSuite

type testingSuite struct {
	suite.Suite
	ctx context.Context
	ti  *e2e.TestInstallation
}

func NewTestingSuite(
	ctx context.Context,
	testInst *e2e.TestInstallation,
) suite.TestingSuite {
	return &testingSuite{
		ctx: ctx,
		ti:  testInst,
	}
}

func (s *testingSuite) SetupSuite() {
	err := s.ti.Actions.Kubectl().ApplyFile(s.ctx, enterpriseCRsManifest)
	s.NoError(err, "can apply manifest "+enterpriseCRsManifest)

	err = s.ti.Actions.Kubectl().ApplyFile(s.ctx, ossCRsManifest)
	s.NoError(err, "can apply manifest "+ossCRsManifest)
}

func (s *testingSuite) TearDownSuite() {
	output, err := s.ti.Actions.Kubectl().DeleteFileWithOutput(s.ctx, enterpriseCRsManifest)
	s.ti.Assertions.ExpectObjectDeleted(enterpriseCRsManifest, err, output)

	output, err = s.ti.Actions.Kubectl().DeleteFileWithOutput(s.ctx, ossCRsManifest)
	s.ti.Assertions.ExpectObjectDeleted(ossCRsManifest, err, output)
}

func (s *testingSuite) TestSoloIoCategory() {
	cmd := s.ti.Actions.Kubectl().Command(s.ctx, "get", "solo-io", "-o", "name")

	var out bytes.Buffer
	err := cmd.WithStdout(io.Writer(&out)).Run().Cause()
	s.NoError(err)

	// output should contain all installed enterprise CRs and no additional CRs
	s.ElementsMatch(strings.Split(strings.TrimSpace(out.String()), "\n"), installedEnterpriseCRs)
}
