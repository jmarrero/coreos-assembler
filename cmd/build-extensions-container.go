package main

import (
	"fmt"
	"github.com/coreos/coreos-assembler/internal/pkg/cosash"
	"path/filepath"
	"io/ioutil"
	"crypto/sha256"
	"io"
	"log"
	"os"
	"strings"
	"encoding/json"
	"os/exec"
)

type MetaJSON struct {
	Images Images `json:"images"`
}

type Images struct {
	ExtensionsContainer ExtensionsContainer `json:"extensions-container"`
}

type ExtensionsContainer struct {
    Path string `json:"path"`
    Sha256 string `json:"sha256"`
    Size int64 `json:"size"`
	SkipCompression bool `json:"skip-compression"`
}

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
	ociarchive := strings.TrimSpace(string(content))
	fmt.Printf("calculating size & sha for: %s \n", ociarchive)


        // file.Size())
	
		f, err := os.Open(ociarchive)
	    if err != nil {
		  log.Fatal(err)
	    }
	    defer f.Close()

    	h := sha256.New()
	    if _, err := io.Copy(h, f); err != nil {
		   log.Fatal(err)
	   }

	   file, err := f.Stat()
	   if err != nil {
		 //file no here?
	   }

	   sha256 := fmt.Sprintf("%x", h.Sum(nil))

       meta := MetaJSON{Images{ExtensionsContainer{filepath.Base(ociarchive), sha256, file.Size(), false}}}
	   json, err := json.Marshal(meta)
	   if err != nil {
		// Could not marshall meta.json
		fmt.Println("Could not marshall meta.json")
	   }
	   fmt.Println(string(json))


       //dumpCurrentJSON := 
	   //cosaAddJSON := exec.Command("cosa", "--workdir", filepath.root(ociarchive), "--build", build, "--artifact", "extensions-container", "--artifact-json", json)	  
	   //finalizeCmd := exec.Command("/usr/lib/coreos-assembler/finalize-artifact", filepath.Base(ociarchive), ociarchive)
	return nil
}
