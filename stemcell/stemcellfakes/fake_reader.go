// This file was generated by counterfeiter
package stemcellfakes

import (
	"sync"

	"github.com/cloudfoundry/bosh-cli/stemcell"
)

type FakeReader struct {
	ReadStub        func(stemcellTarballPath string, extractedPath string) (stemcell.ExtractedStemcell, error)
	readMutex       sync.RWMutex
	readArgsForCall []struct {
		stemcellTarballPath string
		extractedPath       string
	}
	readReturns struct {
		result1 stemcell.ExtractedStemcell
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeReader) Read(stemcellTarballPath string, extractedPath string) (stemcell.ExtractedStemcell, error) {
	fake.readMutex.Lock()
	fake.readArgsForCall = append(fake.readArgsForCall, struct {
		stemcellTarballPath string
		extractedPath       string
	}{stemcellTarballPath, extractedPath})
	fake.recordInvocation("Read", []interface{}{stemcellTarballPath, extractedPath})
	fake.readMutex.Unlock()
	if fake.ReadStub != nil {
		return fake.ReadStub(stemcellTarballPath, extractedPath)
	} else {
		return fake.readReturns.result1, fake.readReturns.result2
	}
}

func (fake *FakeReader) ReadCallCount() int {
	fake.readMutex.RLock()
	defer fake.readMutex.RUnlock()
	return len(fake.readArgsForCall)
}

func (fake *FakeReader) ReadArgsForCall(i int) (string, string) {
	fake.readMutex.RLock()
	defer fake.readMutex.RUnlock()
	return fake.readArgsForCall[i].stemcellTarballPath, fake.readArgsForCall[i].extractedPath
}

func (fake *FakeReader) ReadReturns(result1 stemcell.ExtractedStemcell, result2 error) {
	fake.ReadStub = nil
	fake.readReturns = struct {
		result1 stemcell.ExtractedStemcell
		result2 error
	}{result1, result2}
}

func (fake *FakeReader) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.readMutex.RLock()
	defer fake.readMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeReader) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ stemcell.Reader = new(FakeReader)