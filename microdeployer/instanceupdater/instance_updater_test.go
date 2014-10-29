package instanceupdater_test

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	fakecmd "github.com/cloudfoundry/bosh-agent/platform/commands/fakes"
	fakesys "github.com/cloudfoundry/bosh-agent/system/fakes"
	fakeuuid "github.com/cloudfoundry/bosh-agent/uuid/fakes"
	fakebmagentclient "github.com/cloudfoundry/bosh-micro-cli/microdeployer/agentclient/fakes"
	fakebmas "github.com/cloudfoundry/bosh-micro-cli/microdeployer/applyspec/fakes"
	fakebmblobstore "github.com/cloudfoundry/bosh-micro-cli/microdeployer/blobstore/fakes"
	fakebmtemp "github.com/cloudfoundry/bosh-micro-cli/templatescompiler/fakes"

	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	bmdepl "github.com/cloudfoundry/bosh-micro-cli/deployment"
	bmas "github.com/cloudfoundry/bosh-micro-cli/microdeployer/applyspec"
	bmrel "github.com/cloudfoundry/bosh-micro-cli/release"
	bmstemcell "github.com/cloudfoundry/bosh-micro-cli/stemcell"

	. "github.com/cloudfoundry/bosh-micro-cli/microdeployer/instanceupdater"
)

var _ = Describe("InstanceUpdater", func() {
	var (
		fakeAgentClient      *fakebmagentclient.FakeAgentClient
		instanceUpdater      InstanceUpdater
		applySpec            bmstemcell.ApplySpec
		fakeJobRenderer      *fakebmtemp.FakeJobRenderer
		fakeCompressor       *fakecmd.FakeCompressor
		fakeBlobstore        *fakebmblobstore.FakeBlobstore
		fakeUUIDGenerator    *fakeuuid.FakeGenerator
		fakeApplySpecFactory *fakebmas.FakeApplySpecFactory
		job                  bmrel.Job
		fs                   *fakesys.FakeFileSystem
		tempFile             *os.File
		compileDir           string
		extractDir           string
		logger               boshlog.Logger
	)

	BeforeEach(func() {
		fakeJobRenderer = fakebmtemp.NewFakeJobRenderer()
		fakeCompressor = fakecmd.NewFakeCompressor()
		fakeCompressor.CompressFilesInDirTarballPath = "fake-tarball-path"
		fakeAgentClient = fakebmagentclient.NewFakeAgentClient()
		applySpec = bmstemcell.ApplySpec{
			Packages: map[string]bmstemcell.Blob{
				"first-package-name": bmstemcell.Blob{
					Name:        "first-package-name",
					Version:     "first-package-version",
					SHA1:        "first-package-sha1",
					BlobstoreID: "first-package-blobstore-id",
				},
				"second-package-name": bmstemcell.Blob{
					Name:        "second-package-name",
					Version:     "second-package-version",
					SHA1:        "second-package-sha1",
					BlobstoreID: "second-package-blobstore-id",
				},
			},
			Job: bmstemcell.Job{
				Name: "fake-job-name",
				Templates: []bmstemcell.Blob{
					{
						Name:        "first-job-name",
						Version:     "first-job-version",
						SHA1:        "first-job-sha1",
						BlobstoreID: "first-job-blobstore-id",
					},
					{
						Name:        "second-job-name",
						Version:     "second-job-version",
						SHA1:        "second-job-sha1",
						BlobstoreID: "second-job-blobstore-id",
					},
					{
						Name:        "third-job-name",
						Version:     "third-job-version",
						SHA1:        "third-job-sha1",
						BlobstoreID: "third-job-blobstore-id",
					},
				},
			},
		}
		logger = boshlog.NewLogger(boshlog.LevelNone)
		fakeBlobstore = fakebmblobstore.NewFakeBlobstore()
		fs = fakesys.NewFakeFileSystem()
		deployment := bmdepl.Deployment{
			Name: "fake-deployment-name",
			Jobs: []bmdepl.Job{
				{
					Name: "fake-manifest-job-name",
					Templates: []bmdepl.ReleaseJobRef{
						{Name: "first-job-name"},
						{Name: "third-job-name"},
					},
					RawProperties: map[interface{}]interface{}{
						"fake-property-key": "fake-property-value",
					},
					Networks: []bmdepl.JobNetwork{
						{
							Name:      "fake-network-name",
							StaticIPs: []string{"fake-network-ip"},
						},
					},
				},
			},
			Networks: []bmdepl.Network{
				{
					Name: "fake-network-name",
					Type: "fake-network-type",
				},
			},
		}
		fakeUUIDGenerator = &fakeuuid.FakeGenerator{
			GeneratedUuid: "fake-blob-id",
		}
		fakeApplySpecFactory = fakebmas.NewFakeApplySpecFactory()
		instanceUpdater = NewInstanceUpdater(
			fakeAgentClient,
			applySpec,
			deployment,
			fakeBlobstore,
			fakeCompressor,
			fakeJobRenderer,
			fakeUUIDGenerator,
			fakeApplySpecFactory,
			fs,
			logger,
		)

		var err error
		tempFile, err = fs.TempFile("fake-blob-temp-file")
		Expect(err).ToNot(HaveOccurred())

		fs.ReturnTempFile = tempFile

		fs.TempDirDir = "/fake-tmp-dir"
		// fake file system only supports one temp dir
		compileDir = "/fake-tmp-dir"
		extractDir = "/fake-tmp-dir"
		job = bmrel.Job{
			Templates: map[string]string{
				"director.yml.erb": "config/director.yml",
			},
			ExtractedPath: extractDir,
		}
		blobJobJSON, err := json.Marshal(job)
		Expect(err).ToNot(HaveOccurred())

		fakeCompressor.DecompressFileToDirCallBack = func() {
			fs.WriteFile("/fake-tmp-dir/job.MF", blobJobJSON)
			fs.WriteFile("/fake-tmp-dir/monit", []byte("fake-monit-contents"))
		}

		fakeCompressor.CompressFilesInDirTarballPath = "fake-tarball-path"
	})

	Describe("Update", func() {
		It("stops the agent", func() {
			err := instanceUpdater.Update()
			Expect(err).ToNot(HaveOccurred())
			Expect(fakeAgentClient.StopCalled).To(BeTrue())
		})

		It("downloads only job template blobs from the blobstore that are specified in the manifest", func() {
			err := instanceUpdater.Update()
			Expect(err).ToNot(HaveOccurred())
			Expect(fakeBlobstore.GetInputs).To(Equal([]fakebmblobstore.GetInput{
				{
					BlobID:          "first-job-blobstore-id",
					DestinationPath: tempFile.Name(),
				},
				{
					BlobID:          "third-job-blobstore-id",
					DestinationPath: tempFile.Name(),
				},
			}))

		})

		It("removes the tempfile for downloaded blobs", func() {
			tempFile, err := fs.TempFile("fake-blob-temp-file")
			Expect(err).ToNot(HaveOccurred())

			fs.ReturnTempFile = tempFile
			err = instanceUpdater.Update()
			Expect(err).ToNot(HaveOccurred())

			Expect(fs.FileExists(tempFile.Name())).To(BeFalse())
		})

		It("decompressed job templates", func() {
			err := instanceUpdater.Update()
			Expect(err).ToNot(HaveOccurred())

			Expect(fakeCompressor.DecompressFileToDirTarballPaths[0]).To(Equal(tempFile.Name()))
			Expect(fakeCompressor.DecompressFileToDirDirs[0]).To(Equal(extractDir))
		})

		It("renders job templates", func() {
			err := instanceUpdater.Update()
			Expect(err).ToNot(HaveOccurred())

			Expect(fakeJobRenderer.RenderInputs).To(Equal([]fakebmtemp.RenderInput{
				{
					SourcePath:      extractDir,
					DestinationPath: filepath.Join(compileDir, "first-job-name"),
					Job:             job,
					Properties: map[string]interface{}{
						"fake-property-key": "fake-property-value",
					},
					DeploymentName: "fake-deployment-name",
				},
				{
					SourcePath:      extractDir,
					DestinationPath: filepath.Join(compileDir, "third-job-name"),
					Job:             job,
					Properties: map[string]interface{}{
						"fake-property-key": "fake-property-value",
					},
					DeploymentName: "fake-deployment-name",
				},
			}))
		})

		It("compresses rendered templates", func() {
			err := instanceUpdater.Update()
			Expect(err).ToNot(HaveOccurred())

			Expect(fakeCompressor.CompressFilesInDirDir).To(Equal(compileDir))
		})

		It("cleans up rendered tarball", func() {
			err := instanceUpdater.Update()
			Expect(err).ToNot(HaveOccurred())

			Expect(fs.FileExists("fake-tarball-path")).To(BeFalse())
		})

		It("uploads rendered jobs to the blobstore", func() {
			err := instanceUpdater.Update()
			Expect(err).ToNot(HaveOccurred())

			Expect(fakeBlobstore.SaveInputs).To(Equal([]fakebmblobstore.SaveInput{
				{
					BlobID:     "fake-blob-id",
					SourcePath: "fake-tarball-path",
				},
			}))
		})

		It("creates apply spec", func() {
			err := instanceUpdater.Update()
			Expect(err).ToNot(HaveOccurred())

			Expect(fakeApplySpecFactory.CreateInput).To(Equal(
				fakebmas.CreateInput{
					ApplySpec:      applySpec,
					DeploymentName: "fake-deployment-name",
					JobName:        "fake-manifest-job-name",
					NetworksSpec: map[string]interface{}{
						"fake-network-name": map[string]interface{}{
							"type":             "fake-network-type",
							"ip":               "fake-network-ip",
							"cloud_properties": map[string]interface{}{},
						},
					},
					ArchivedTemplatesBlobID: "fake-blob-id",
					ArchivedTemplatesPath:   "fake-tarball-path",
					TemplatesDir:            compileDir,
				},
			))
		})

		It("sends apply spec to the agent", func() {
			applySpec := bmas.ApplySpec{
				Deployment: "fake-deployment-name",
			}
			fakeApplySpecFactory.CreateApplySpec = applySpec
			err := instanceUpdater.Update()
			Expect(err).ToNot(HaveOccurred())
			Expect(fakeAgentClient.ApplyApplySpec).To(Equal(applySpec))
		})

		Context("when sending apply spec to the agent fails", func() {
			BeforeEach(func() {
				fakeAgentClient.ApplyErr = errors.New("fake-agent-apply-err")
			})

			It("returns an error", func() {
				err := instanceUpdater.Update()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-agent-apply-err"))
			})
		})

		Context("when stopping an agent fails", func() {
			BeforeEach(func() {
				fakeAgentClient.SetStopBehavior(errors.New("fake-stop-error"))
			})

			It("returns an error", func() {
				err := instanceUpdater.Update()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-stop-error"))
			})
		})

		Context("when downloading a blob fails", func() {
			BeforeEach(func() {
				fakeBlobstore.GetErr = errors.New("fake-blobstore-get-error")
			})

			It("returns an error", func() {
				err := instanceUpdater.Update()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-blobstore-get-error"))
			})
		})

		Context("when rendering jobs fails", func() {
			BeforeEach(func() {
				fakeJobRenderer.SetRenderBehavior(extractDir, errors.New("fake-render-error"))
			})

			It("returns an error", func() {
				err := instanceUpdater.Update()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-render-error"))
			})
		})

		Context("when compressing rendered templates fails", func() {
			BeforeEach(func() {
				fakeCompressor.CompressFilesInDirErr = errors.New("fake-compress-tarball-error")
			})

			It("returns an error", func() {
				err := instanceUpdater.Update()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-compress-tarball-error"))
			})
		})
	})

	Describe("Start", func() {
		It("starts agent services", func() {
			err := instanceUpdater.Start()
			Expect(err).ToNot(HaveOccurred())
			Expect(fakeAgentClient.StartCalled).To(BeTrue())
		})

		Context("when starting an agent fails", func() {
			BeforeEach(func() {
				fakeAgentClient.SetStartBehavior(errors.New("fake-start-error"))
			})

			It("returns an error", func() {
				err := instanceUpdater.Start()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-start-error"))
			})
		})
	})
})
