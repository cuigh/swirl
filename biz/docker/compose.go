package docker

//import (
//	"io/ioutil"
//
//	"github.com/docker/cli/cli/compose/loader"
//	composetypes "github.com/docker/cli/cli/compose/types"
//)
//
//func getConfigFile(filename string) (*composetypes.ConfigFile, error) {
//	bytes, err := ioutil.ReadFile(filename)
//	if err != nil {
//		return nil, err
//	}
//	config, err := loader.ParseYAML(bytes)
//	if err != nil {
//		return nil, err
//	}
//	return &composetypes.ConfigFile{
//		Filename: filename,
//		Config:   config,
//	}, nil
//}
