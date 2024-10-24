package kyma

import (
	"path"
	"testing"

	"github.com/kyma-project/test-infra/development/tools/jobs/tester"
	"github.com/kyma-project/test-infra/development/tools/jobs/tester/jobsuite"
)

var tests = []struct {
	path              string
	image             string
	suite             func(config *jobsuite.Config) jobsuite.Suite
	additionalOptions []jobsuite.Option
}{
	{path: "service-catalog", image: tester.ImageGolangKubebuilder2BuildpackLatest, suite: tester.NewGenericComponentSuite,
		additionalOptions: []jobsuite.Option{
			jobsuite.JobFileSuffix("tests-generic"),
		},
	},
	{path: "application-connector-tests", image: tester.ImageGolangKubebuilder2BuildpackLatest, suite: tester.NewGenericComponentSuite,
		additionalOptions: []jobsuite.Option{
			jobsuite.JobFileSuffix("generic"),
		},
	},
	{path: "application-gateway-tests", image: tester.ImageGolangKubebuilder2BuildpackLatest, suite: tester.NewGenericComponentSuite,
		additionalOptions: []jobsuite.Option{
			jobsuite.JobFileSuffix("generic"),
		},
	},
	{path: "application-gateway-legacy-tests", image: tester.ImageGolangKubebuilder2BuildpackLatest, suite: tester.NewGenericComponentSuite,
		additionalOptions: []jobsuite.Option{
			jobsuite.JobFileSuffix("generic"),
		},
	},
	{path: "application-operator-tests", image: tester.ImageGolangKubebuilder2BuildpackLatest, suite: tester.NewGenericComponentSuite,
		additionalOptions: []jobsuite.Option{
			jobsuite.JobFileSuffix("generic"),
		},
	},
	{path: "application-registry-tests", image: tester.ImageGolangKubebuilder2BuildpackLatest, suite: tester.NewGenericComponentSuite,
		additionalOptions: []jobsuite.Option{
			jobsuite.JobFileSuffix("generic"),
		},
	},
	{path: "compass-runtime-agent", image: tester.ImageGolangKubebuilder2BuildpackLatest, suite: tester.NewGenericComponentSuite,
		additionalOptions: []jobsuite.Option{
			jobsuite.JobFileSuffix("tests-generic"),
		},
	},
	{path: "connection-token-handler-tests", image: tester.ImageGolangKubebuilder2BuildpackLatest, suite: tester.NewGenericComponentSuite,
		additionalOptions: []jobsuite.Option{
			jobsuite.JobFileSuffix("generic"),
		},
	},
	{path: "connector-service-tests", image: tester.ImageGolangKubebuilder2BuildpackLatest, suite: tester.NewGenericComponentSuite,
		additionalOptions: []jobsuite.Option{
			jobsuite.JobFileSuffix("generic"),
		},
	},
	{path: "console-backend-service", image: tester.ImageGolangKubebuilder2BuildpackLatest, suite: tester.NewGenericComponentSuite,
		additionalOptions: []jobsuite.Option{
			jobsuite.JobFileSuffix("tests-generic"),
			jobsuite.RunIfChanged("components/console-backend-service/main.go", "scripts/go-dep.mk"),
		},
	},
	{path: "end-to-end/upgrade", image: tester.ImageGolangKubebuilder2BuildpackLatest, suite: tester.NewGenericComponentSuite, additionalOptions: []jobsuite.Option{
		jobsuite.RunIfChanged("^tests/end-to-end/upgrade/[^chart]", "tests/end-to-end/upgrade/fix"),
		jobsuite.JobFileSuffix("tests-generic"),
	}},
	{path: "end-to-end/external-solution-integration", image: tester.ImageGolangKubebuilder2BuildpackLatest, suite: tester.NewGenericComponentSuite,
		additionalOptions: []jobsuite.Option{
			jobsuite.JobFileSuffix("tests-generic"),
		},
	},
	{path: "integration/event-service", image: tester.ImageGolangKubebuilder2BuildpackLatest, suite: tester.NewGenericComponentSuite,
		additionalOptions: []jobsuite.Option{
			jobsuite.JobFileSuffix("tests-generic"),
			jobsuite.AllReleases(),
		},
	},
	{path: "integration/apiserver-proxy", image: tester.ImageGolangKubebuilder2BuildpackLatest, suite: tester.NewGenericComponentSuite,
		additionalOptions: []jobsuite.Option{
			jobsuite.JobFileSuffix("tests-generic"),
		},
	},
	{path: "rafter", image: tester.ImageGolangKubebuilder2BuildpackLatest, suite: tester.NewGenericComponentSuite,
		additionalOptions: []jobsuite.Option{
			jobsuite.JobFileSuffix("tests-generic"),
		},
	},
	{path: "integration/api-gateway", image: tester.ImageGolangKubebuilder2BuildpackLatest, suite: tester.NewGenericComponentSuite,
		additionalOptions: []jobsuite.Option{
			jobsuite.JobFileSuffix("tests-generic"),
		},
	},
	{path: "integration/dex", image: tester.ImageGolangKubebuilder2BuildpackLatest, suite: tester.NewGenericComponentSuite,
		additionalOptions: []jobsuite.Option{
			jobsuite.JobFileSuffix("tests-generic"),
		},
	},
	{path: "integration/logging", image: tester.ImageGolangKubebuilder2BuildpackLatest, suite: tester.NewGenericComponentSuite,
		additionalOptions: []jobsuite.Option{
			jobsuite.JobFileSuffix("generic"),
		},
	},
	{path: "contract/knative-channel-kafka", image: tester.ImageGolangKubebuilder2BuildpackLatest, suite: tester.NewGenericComponentSuite,
		additionalOptions: []jobsuite.Option{
			jobsuite.JobFileSuffix("tests-generic"),
		},
	},
	{path: "function-controller", image: tester.ImageGolangKubebuilder2BuildpackLatest, suite: tester.NewGenericComponentSuite,
		additionalOptions: []jobsuite.Option{
			jobsuite.JobFileSuffix("tests-generic"),
		},
	},
}

func TestTestJobs(t *testing.T) {
	testedConfigurations := make(map[string]struct{})
	repos := map[string]struct{}{}
	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			opts := []jobsuite.Option{
				jobsuite.Test(test.path, test.image),
				jobsuite.KymaRepo(),
				jobsuite.AllReleases(),
			}
			opts = append(opts, test.additionalOptions...)
			cfg := jobsuite.NewConfig(opts...)
			suite := test.suite
			if suite == nil {
				suite = tester.NewComponentSuite
			}
			ts := suite(cfg)
			if pathProvider, ok := ts.(jobsuite.JobConfigPathProvider); ok {
				testedConfigurations[path.Clean(pathProvider.JobConfigPath())] = struct{}{}
			}
			repos[cfg.Repository] = struct{}{}
			ts.Run(t)
		})
	}
	t.Run("All Files covered by test", jobsuite.CheckFilesAreTested(repos, testedConfigurations, jobBasePath, []string{"tests"}))
}
