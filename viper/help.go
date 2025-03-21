package viper

import (
	"flag"
	"github.com/gone-io/gone/v2"
	"github.com/sagikazarmark/locafero"
	"github.com/spf13/afero"
	"os"
	"path"
	"path/filepath"
)

const testDataDir = "testdata"
const defaultEnv = "local"

const EEnv = "ENV"
const EConf = "CONF"
const ConPath = "config"

const TestSuffix = "_test"
const DefaultConf = "default"

var envFlag = flag.String("env", "", "environment")
var confFlag = flag.String("conf", "", "config directory")

// GetEnv get environment, fetch value from command line flag(-env) first, then from environment variable(ENV), then use default value
func GetEnv() string {
	flag.Parse()
	if *envFlag != "" {
		return *envFlag
	}

	env := os.Getenv(EEnv)
	if env != "" {
		return env
	}
	return defaultEnv
}

// GetConfDir get config directory, fetch value from command line flag(-conf) first, then from environment variable(CONF)
func GetConfDir() string {
	flag.Parse()
	if *confFlag != "" {
		return *confFlag
	}
	return os.Getenv(EConf)
}

func MustGetExecutableConfDir() string {
	dir, err := os.Executable()
	if err != nil {
		panic(gone.ToError(err))
	}
	return path.Join(filepath.Dir(dir), ConPath)
}

func MustGetWorkDir() string {
	workDir, err := os.Getwd()
	if err != nil {
		panic(gone.ToError(err))
	}
	return workDir
}

func lookForModDir(workDir string) string {
	if workDir == "/" {
		return ""
	}
	modFile := path.Join(workDir, "go.mod")
	if _, err := os.Stat(modFile); os.IsNotExist(err) {
		return lookForModDir(path.Dir(workDir))
	}
	return workDir
}

var SupportedExts = []string{"json", "toml", "yaml", "yml", "properties"}

func getConfigPaths(isInTestKit bool) []string {
	configPaths := []string{
		MustGetExecutableConfDir(),
		path.Join(MustGetExecutableConfDir(), ConPath),
		MustGetWorkDir(),
		path.Join(MustGetWorkDir(), ConPath),
	}
	if isInTestKit {
		modDir := lookForModDir(MustGetWorkDir())
		if modDir != "" {
			configPaths = append(configPaths, path.Join(modDir, ConPath))
		}
		configPaths = append(configPaths,
			path.Join(MustGetWorkDir(), testDataDir),
			path.Join(MustGetWorkDir(), testDataDir, ConPath),
		)
	}
	settingConfPath := GetConfDir()
	if settingConfPath != "" {
		configPaths = append(configPaths, settingConfPath)
	}
	return configPaths
}

func findConfigFiles(env string, isInTestKit bool, paths []string, fsys afero.Fs) ([]string, error) {
	filenames := locafero.NameWithExtensions(DefaultConf, SupportedExts...)
	if isInTestKit {
		filenames = append(filenames, locafero.NameWithExtensions(DefaultConf+TestSuffix, SupportedExts...)...)
	}

	if env != "" {
		filenames = append(filenames, locafero.NameWithExtensions(env, SupportedExts...)...)
		if isInTestKit {
			filenames = append(filenames, locafero.NameWithExtensions(env+TestSuffix, SupportedExts...)...)
		}
	}

	finder := locafero.Finder{
		Paths: paths,
		Names: filenames,
		Type:  locafero.FileTypeFile,
	}

	return finder.Find(fsys)
}

func getConfigFiles(isInTestKit bool, fsys afero.Fs) ([]string, error) {
	paths := getConfigPaths(isInTestKit)
	return findConfigFiles(GetEnv(), isInTestKit, paths, fsys)
}

func fileExt(cf string) string {
	ext := filepath.Ext(cf)
	if len(ext) > 1 {
		return ext[1:]
	}
	return ""
}
