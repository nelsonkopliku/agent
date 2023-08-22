package gatherers

import (
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/trento-project/agent/pkg/factsengine/entities"
	utilsMocks "github.com/trento-project/agent/pkg/utils/mocks"
	"github.com/trento-project/agent/test/helpers"
)

type FactGatherer interface {
	Gather(factsRequests []entities.FactRequest) ([]entities.Fact, error)
}

func StandardGatherers(agentID string) map[string]FactGatherer {
	return map[string]FactGatherer{
		CibAdminGathererName:        demoCibAdmin(),
		CorosyncCmapCtlGathererName: demoCorosyncCmapCtl(),
		CorosyncConfGathererName:    demoCorosync(agentID),
		HostsFileGathererName:       NewDefaultHostsFileGatherer(),
		SystemDGathererName:         NewDefaultSystemDGatherer(),
		// PackageVersionGathererName:  NewDefaultPackageVersionGatherer(),
		PackageVersionGathererName: demoPackageVersion(),
		SBDConfigGathererName:      demoSbDConfig(agentID),
		SBDDumpGathererName:        demoSbdDump(),
		SapHostCtrlGathererName:    NewDefaultSapHostCtrlGatherer(),
		VerifyPasswordGathererName: NewDefaultPasswordGatherer(),
	}
}

func demoPackageVersion() *PackageVersionGatherer {
	mockExecutor := new(utilsMocks.CommandExecutor)

	multiversionsMockOutputFile, _ := os.Open(helpers.GetFixturePath("gatherers/rpm-query-multi-versions.variant-1.output"))
	multiversionsVersionMockOutput, _ := io.ReadAll(multiversionsMockOutputFile)
	mockExecutor.On("Exec", "/usr/bin/rpm", "-q", "--qf", "VERSION=%{VERSION}\nINSTALLTIME=%{INSTALLTIME}\n---\n", mock.AnythingOfType("string")).Return(
		[]byte(multiversionsVersionMockOutput), nil)

	mockExecutor.On("Exec", "/usr/bin/zypper", "--terse", "versioncmp", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(
		[]byte("Warning: The /etc/products.d/baseproduct symlink is dangling or missing!\nThe link must point to your core products .prod file in /etc/products.d.\n\n-1\n"), nil)
	// mockExecutor.On("Exec", "/usr/bin/rpm", "-q", "--qf", "VERSION=%{VERSION}\nINSTALLTIME=%{INSTALLTIME}\n---\n", "SLES_SAP-release").Return(
	// 	[]byte("package SLES_SAP-release is not installed"), errors.New(""))

	return NewPackageVersionGatherer(mockExecutor)
}

func demoCorosync(agentID string) *CorosyncConfGatherer {
	if agentID == "13e8c25c-3180-5a9a-95c8-51ec38e50cfc" {
		return NewCorosyncConfGatherer(helpers.GetFixturePath("gatherers/corosync.conf.one_node"))
	}
	return NewCorosyncConfGatherer(helpers.GetFixturePath("gatherers/corosync.conf.basic"))
}

func demoCorosyncCmapCtl() *CorosyncCmapctlGatherer {
	mockExecutor := new(utilsMocks.CommandExecutor)
	mockOutputFile, _ := os.Open(helpers.GetFixturePath("gatherers/corosynccmap-ctl.output"))
	mockOutput, _ := io.ReadAll(mockOutputFile)
	mockExecutor.On("Exec", "corosync-cmapctl", "-b").Return(mockOutput, nil)

	return NewCorosyncCmapctlGatherer(mockExecutor)
}

func RandBool() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2) == 1
}

func demoSbDConfig(agentID string) *SBDGatherer {
	if agentID == "13e8c25c-3180-5a9a-95c8-51ec38e50cfc" {
		return NewSBDGatherer(helpers.GetFixturePath("discovery/cluster/sbd/sbd_config"))
	}
	return NewSBDGatherer(helpers.GetFixturePath("discovery/cluster/sbd/sbd_config_argh"))
}

func demoSbdDump() *SBDDumpGatherer {
	mockExecutor := new(utilsMocks.CommandExecutor)

	deviceVDBMockOutputFile, _ := os.Open(helpers.GetFixturePath("gatherers/dev.vdb.sbddump.output"))
	deviceVDBMockOutput, _ := ioutil.ReadAll(deviceVDBMockOutputFile)

	deviceVDCMockOutputFile, _ := os.Open(helpers.GetFixturePath("gatherers/dev.vdc.sbddump.output"))
	deviceVDCMockOutput, _ := ioutil.ReadAll(deviceVDCMockOutputFile)

	mockExecutor.On("Exec", "sbd", "-d", "/dev/vdb", "dump").Return(deviceVDBMockOutput, nil)
	mockExecutor.On("Exec", "sbd", "-d", "/dev/vdc", "dump").Return(deviceVDCMockOutput, nil)

	return NewSBDDumpGatherer(
		mockExecutor,
		helpers.GetFixturePath("discovery/cluster/sbd/sbd_config"),
	)
}

func demoCibAdmin() *CibAdminGatherer {
	mockExecutor := new(utilsMocks.CommandExecutor)

	lFile, _ := os.Open(helpers.GetFixturePath("gatherers/cibadmin.xml"))
	content, _ := io.ReadAll(lFile)

	cibAdminOutput := content

	mockExecutor.On("Exec", "cibadmin", "--query", "--local").Return(
		cibAdminOutput, nil)

	return NewCibAdminGatherer(mockExecutor)
}
