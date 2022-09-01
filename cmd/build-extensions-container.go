package main

import (
	"fmt"
	"github.com/coreos/coreos-assembler-schema/cosa"
	"github.com/coreos/coreos-assembler/internal/pkg/cosash"

	"crypto/sha256"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func buildExtensionContainer() error {
	sh, err := cosash.NewCosaSh()
	if err != nil {
		return err
	}
	if _, err := sh.PrepareBuild(); err != nil {
		return err
	}
	sh.Process("runvm -- /usr/lib/coreos-assembler/build-extensions-oscontainer.sh $tmp_builddir/output.txt")
	tmpdir, err := sh.ProcessWithReply("echo $tmp_builddir>&3\n")
	if err != nil {
		return err
	}
	content, err := ioutil.ReadFile(filepath.Join(tmpdir, "output.txt"))
	if err != nil {
		return err
	}
	arch := cosa.BuilderArch()
	ociarchive := strings.TrimSpace(string(content))
	workdir := getWorkDir(ociarchive)
	lastBuild, _, err := cosa.ReadBuild(workdir+"/builds", "latest", arch)
	if err != nil {
		return err
	}
	buildID := lastBuild.BuildID
	renamedArchive := filepath.Join(filepath.Dir(ociarchive), "extensions-container-"+buildID+"."+arch+".ociarchive")
	err = os.Rename(ociarchive, renamedArchive)
	if err != nil {
		return err
	}
	file, err := os.Open(renamedArchive)
	if err != nil {
		return err
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return err
	}
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	sha256 := fmt.Sprintf("%x", hash.Sum(nil))
	builddir := filepath.Join(workdir, "builds", "latest", arch)
	metapath := filepath.Join(builddir, "meta.json")

	jsonFile, err := os.Open(metapath)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	jsonBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	var cosaBuild cosa.Build
	err = json.Unmarshal(jsonBytes, &cosaBuild)
	if err != nil {
		return err
	}

	cosaBuild.BuildArtifacts.ExtensionsContainer = &cosa.Artifact{
		Path:            filepath.Base(renamedArchive),
		Sha256:          sha256,
		SizeInBytes:     float64(stat.Size()),
		SkipCompression: false,
	}

	newBytes, err := json.MarshalIndent(cosaBuild, "", "    ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(metapath, newBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func getWorkDir(path string) string {
	directories := strings.Split(path, "/")
	//expects path starts with /.
	return "/" + directories[1]
}
