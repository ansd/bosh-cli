package cmd_test

import (
	"github.com/cppforlife/go-patch/patch"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry/bosh-cli/cmd"
	boshtpl "github.com/cloudfoundry/bosh-cli/director/template"
	fakeui "github.com/cloudfoundry/bosh-cli/ui/fakes"
)

var _ = Describe("BuildManifestCmd", func() {
	var (
		ui      *fakeui.FakeUI
		command BuildManifestCmd
	)

	BeforeEach(func() {
		ui = &fakeui.FakeUI{}
		command = NewBuildManifestCmd(ui)
	})

	Describe("Run", func() {
		var (
			opts BuildManifestOpts
		)

		BeforeEach(func() {
			opts = BuildManifestOpts{}
		})

		act := func() error { return command.Run(opts) }

		It("shows templated manifest", func() {
			opts.Args.Manifest = FileBytesArg{
				Bytes: []byte("name1: ((name1))\nname2: ((name2))"),
			}

			opts.VarKVs = []boshtpl.VarKV{
				{Name: "name1", Value: "val1-from-kv"},
			}

			opts.VarsFiles = []boshtpl.VarsFileArg{
				{Vars: boshtpl.Variables(map[string]interface{}{"name1": "val1-from-file"})},
				{Vars: boshtpl.Variables(map[string]interface{}{"name2": "val2-from-file"})},
			}

			opts.OpsFiles = []OpsFileArg{
				{
					Ops: patch.Ops([]patch.Op{
						patch.ReplaceOp{Path: patch.MustNewPointerFromString("/xyz?"), Value: "val"},
					}),
				},
			}

			err := act()
			Expect(err).ToNot(HaveOccurred())

			bytes := "name1: val1-from-kv\nname2: val2-from-file\nxyz: val\n"
			Expect(ui.Blocks).To(Equal([]string{bytes}))
		})

		It("returns portion of the template if out path is given", func() {
			opts.Args.Manifest = FileBytesArg{
				Bytes: []byte("name1: ((name1))\nname2: ((name2))"),
			}

			opts.VarKVs = []boshtpl.VarKV{
				{Name: "name1", Value: "val1-from-kv"},
			}

			opts.VarsFiles = []boshtpl.VarsFileArg{
				{Vars: boshtpl.Variables(map[string]interface{}{"var": "var-val"})},
			}

			opts.OpsFiles = []OpsFileArg{
				{
					Ops: patch.Ops([]patch.Op{
						patch.ReplaceOp{Path: patch.MustNewPointerFromString("/name2"), Value: "((var))"},
					}),
				},
			}

			ptr := patch.MustNewPointerFromString("/name2")
			opts.OutPath = &ptr

			err := act()
			Expect(err).ToNot(HaveOccurred())
			Expect(ui.Blocks).To(Equal([]string{"var-val\n"}))
		})

		It("returns portion of the template formatting multiline string without YAML indent", func() {
			opts.Args.Manifest = FileBytesArg{
				Bytes: []byte(`key: "line1\nline2"`),
			}

			ptr := patch.MustNewPointerFromString("/key")
			opts.OutPath = &ptr

			err := act()
			Expect(err).ToNot(HaveOccurred())
			Expect(ui.Blocks).To(Equal([]string{"line1\nline2\n"}))
		})

		It("returns portion of the template formatting result as regular YAML", func() {
			opts.Args.Manifest = FileBytesArg{
				Bytes: []byte("key:\n  subkey:\n    subsubkey: key"),
			}

			ptr := patch.MustNewPointerFromString("/key")
			opts.OutPath = &ptr

			err := act()
			Expect(err).ToNot(HaveOccurred())
			Expect(ui.Blocks).To(Equal([]string{"subkey:\n  subsubkey: key\n"}))
		})

		It("returns error if variables are not found in templated manifest if var-errors flag is set", func() {
			opts.Args.Manifest = FileBytesArg{
				Bytes: []byte("name1: ((name1))\nname2: ((name2))"),
			}

			opts.VarKVs = []boshtpl.VarKV{
				{Name: "name1", Value: "val1-from-kv"},
			}

			opts.VarErrors = true

			err := act()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Expected to find variables: name2"))
		})
	})
})
