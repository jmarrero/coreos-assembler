package main

import (
	"fmt"
	"github.com/coreos/coreos-assembler/internal/pkg/cosash"
	"path/filepath"
	"io/ioutil"
)

func buildExtensionContainer()  error {
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
	//TODO add sha256 and size to meta.json here.
	//Modify cmd-push-container to add the extensions-container.
	fmt.Printf("File contents: %s", content)

	return nil
}
