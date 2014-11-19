package deployment_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	fakesys "github.com/cloudfoundry/bosh-agent/system/fakes"

	boshlog "github.com/cloudfoundry/bosh-agent/logger"

	. "github.com/cloudfoundry/bosh-micro-cli/deployment"
)

var _ = Describe("DeploymentRenderer", func() {
	var (
		deploymentPath string
		fakeFs         *fakesys.FakeFileSystem
		manifestParser ManifestParser
	)

	BeforeEach(func() {
		deploymentPath = "fake-deployment-path"
		fakeFs = fakesys.NewFakeFileSystem()
		logger := boshlog.NewLogger(boshlog.LevelNone)
		manifestParser = NewBoshDeploymentParser(fakeFs, logger)
	})

	Context("when deployment path does not exist", func() {
		It("returns an error", func() {
			_, err := manifestParser.Parse(deploymentPath)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when deployment path exists", func() {
		Context("when parser fails to read the deployment file", func() {
			BeforeEach(func() {
				fakeFs.ReadFileError = errors.New("fake-read-file-error")
			})

			It("returns an error", func() {
				_, err := manifestParser.Parse(deploymentPath)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when parser successfully reads the deployment file", func() {
			BeforeEach(func() {
				contents := `
---
name: fake-deployment-name
update:
  update_watch_time: 2000-7000
resource_pools:
- name: fake-resource-pool-name
  env:
    bosh:
      password: secret
networks:
- name: fake-network-name
  type: dynamic
  cloud_properties:
    subnet: fake-subnet
    a:
      b: value
- name: vip
  type: vip
disk_pools:
- name: fake-disk-pool-name
  disk_size: 2048
  cloud_properties:
    fake-disk-pool-cloud-property-key: fake-disk-pool-cloud-property-value
jobs:
- name: bosh
  networks:
  - name: vip
    static_ips: [1.2.3.4]
  persistent_disk: 1024
  persistent_disk_pool: fake-disk-pool-name
  properties:
    fake-prop-key:
      nested-prop-key: fake-prop-value
cloud_provider:
  properties:
    nested-property: fake-property-value
`
				fakeFs.WriteFileString(deploymentPath, contents)
			})

			It("parses deployment manifest", func() {
				deployment, err := manifestParser.Parse(deploymentPath)
				Expect(err).ToNot(HaveOccurred())

				Expect(deployment.Name).To(Equal("fake-deployment-name"))
				Expect(deployment.Update.UpdateWatchTime.Start).To(Equal(2000))
				Expect(deployment.Update.UpdateWatchTime.End).To(Equal(7000))

				networks := deployment.Networks
				Expect(networks).To(Equal([]Network{
					{
						Name: "fake-network-name",
						Type: Dynamic,
						RawCloudProperties: map[interface{}]interface{}{
							"subnet": "fake-subnet",
							"a": map[interface{}]interface{}{
								"b": "value",
							},
						},
					},
					{
						Name: "vip",
						Type: VIP,
					},
				}))
				resourcePools := deployment.ResourcePools
				Expect(resourcePools).To(Equal([]ResourcePool{
					{
						Name: "fake-resource-pool-name",
						RawEnv: map[interface{}]interface{}{
							"bosh": map[interface{}]interface{}{
								"password": "secret",
							},
						},
					},
				}))
				diskPools := deployment.DiskPools
				Expect(diskPools).To(Equal([]DiskPool{
					{
						Name:     "fake-disk-pool-name",
						DiskSize: 2048,
						RawCloudProperties: map[interface{}]interface{}{
							"fake-disk-pool-cloud-property-key": "fake-disk-pool-cloud-property-value",
						},
					},
				}))
				jobs := deployment.Jobs
				Expect(jobs).To(Equal([]Job{
					{
						Name: "bosh",
						Networks: []JobNetwork{
							{
								Name:      "vip",
								StaticIPs: []string{"1.2.3.4"},
							},
						},
						PersistentDisk:     1024,
						PersistentDiskPool: "fake-disk-pool-name",
						RawProperties: map[interface{}]interface{}{
							"fake-prop-key": map[interface{}]interface{}{
								"nested-prop-key": "fake-prop-value",
							},
						},
					},
				}))
			})

			Context("when update watch time is not set", func() {
				BeforeEach(func() {
					contents := `
---
name: fake-deployment-name
`
					fakeFs.WriteFileString(deploymentPath, contents)
				})

				It("uses default values", func() {
					deployment, err := manifestParser.Parse(deploymentPath)
					Expect(err).ToNot(HaveOccurred())

					Expect(deployment.Name).To(Equal("fake-deployment-name"))
					Expect(deployment.Update.UpdateWatchTime.Start).To(Equal(0))
					Expect(deployment.Update.UpdateWatchTime.End).To(Equal(300000))
				})
			})
		})
	})
})
