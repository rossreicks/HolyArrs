package downloads

import (
	"testing"

	"github.com/HolyArrs/holyarrs/app/backend/internal/mediafiles"
)

func TestRetryImport__finishedPackWithStuckEpisodes__queuesOneImportForEveryStuckEpisode(t *testing.T) {
	t.Parallel()
	// given
	x, ids := seedPack(t)

	// when
	do(t, "POST", x.srv.URL+"/api/v1/needs-decision/"+itoa(ids[0])+"/decide", `{"action":"retry_import"}`, 202)

	// then
	want := map[int64]string{ids[0]: "importing", ids[1]: "in_library", ids[2]: "importing"}
	for id, st := range want {
		var got string
		_ = x.s.Read().QueryRow("SELECT state FROM slots WHERE id = ?", id).Scan(&got)
		if got != st {
			t.Errorf("slot %d is %s, want %s", id, got, st)
		}
	}
	var jobs int
	_ = x.s.Read().QueryRow("SELECT count(*) FROM jobs WHERE kind = 'import_download' AND json_extract(subject, '$.id') = ?", x.dlID).Scan(&jobs)
	if jobs != 1 {
		t.Fatalf("%d import jobs", jobs)
	}
}

func TestRetryImport__downloadNotFinished__isRefusedAndChangesNothing(t *testing.T) {
	t.Parallel()
	// given
	x := seed41(t, mediafiles.SlotNeedsADecision, "downloading")

	// when
	do(t, "POST", x.srv.URL+"/api/v1/needs-decision/"+itoa(x.slotID)+"/decide", `{"action":"retry_import"}`, 409)

	// then
	var st string
	_ = x.s.Read().QueryRow("SELECT state FROM slots WHERE id = ?", x.slotID).Scan(&st)
	if st != "needs_a_decision" {
		t.Fatalf("slot is %s", st)
	}
}
