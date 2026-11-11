package importer

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	domain "github.com/HolyArrs/holyarrs/app/backend/internal/core/model"
)

func TestImport__usenetDownloadWithNoFileList__placesTheVideoFromItsFolder(t *testing.T) {
	// given
	w := setup(t)
	mapExec(t, w.s, "UPDATE downloaders SET adapter = 'nzbget', protocol = 'usenet'")
	mapExec(t, w.s, "UPDATE downloads SET files_seen = '[]' WHERE id = ?", w.dlID)

	// when
	imp, err := (&Importer{Store: w.s}).Import(context.Background(), w.dlID, false, domain.System)

	// then
	if err != nil || imp.State != "done" {
		t.Fatalf("import %+v %v", imp, err)
	}
	if n := count(t, w.s, "files"); n != 1 {
		t.Fatalf("%d files placed", n)
	}
}

func TestImport__fileListNamesAFolder__placesTheVideoInsideIt(t *testing.T) {
	// given
	w := setup(t)
	job := filepath.Join(w.dl, "Dune.Part.Two.2024.1080p.WEB.H264-HOLYARRS.TAoE")
	if err := os.Mkdir(job, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(filepath.Join(w.dl, relName+".mkv"), filepath.Join(job, relName+".mkv")); err != nil {
		t.Fatal(err)
	}
	mapExec(t, w.s, "UPDATE downloaders SET adapter = 'nzbget', protocol = 'usenet'")
	mapExec(t, w.s, "UPDATE downloads SET files_seen = json_array(json_object('name', ?, 'size', 2048)) WHERE id = ?", filepath.Base(job), w.dlID)

	// when
	imp, err := (&Importer{Store: w.s}).Import(context.Background(), w.dlID, false, domain.System)

	// then
	if err != nil || imp.State != "done" {
		t.Fatalf("import %+v %v", imp, err)
	}
	if n := count(t, w.s, "files"); n != 1 {
		t.Fatalf("%d files placed", n)
	}
}

func TestImport__torrentDownloadWithNoFileList__placesNothing(t *testing.T) {
	// given
	w := setup(t)
	mapExec(t, w.s, "UPDATE downloads SET files_seen = '[]' WHERE id = ?", w.dlID)

	// when
	imp, err := (&Importer{Store: w.s}).Import(context.Background(), w.dlID, false, domain.System)

	// then
	if err != nil || imp.State == "done" {
		t.Fatalf("import %+v %v", imp, err)
	}
	if n := count(t, w.s, "files"); n != 0 {
		t.Fatalf("%d files placed", n)
	}
}
