package commands

import "testing"

func TestCompleteFilter(t *testing.T) {
	t.Run("complete shows incomplete reminders", func(t *testing.T) {
		f := completeFilter(false, "", false)
		if f.Completed == nil || *f.Completed != false {
			t.Errorf("expected Completed=false, got %v", f.Completed)
		}
	})

	t.Run("uncomplete shows completed reminders", func(t *testing.T) {
		f := completeFilter(true, "", false)
		if f.Completed == nil || *f.Completed != true {
			t.Errorf("expected Completed=true, got %v", f.Completed)
		}
	})

	t.Run("list filter applied", func(t *testing.T) {
		f := completeFilter(false, "Work", false)
		if f.ListName != "Work" {
			t.Errorf("expected ListName=Work, got %q", f.ListName)
		}
	})

	t.Run("empty list not set", func(t *testing.T) {
		f := completeFilter(false, "", false)
		if f.ListName != "" {
			t.Errorf("expected empty ListName, got %q", f.ListName)
		}
	})

	t.Run("flagged applied on complete", func(t *testing.T) {
		f := completeFilter(false, "", true)
		if f.Flagged == nil || *f.Flagged != true {
			t.Errorf("expected Flagged=true, got %v", f.Flagged)
		}
	})

	t.Run("flagged ignored on uncomplete", func(t *testing.T) {
		f := completeFilter(true, "", true)
		if f.Flagged != nil {
			t.Errorf("expected Flagged=nil on uncomplete, got %v", *f.Flagged)
		}
	})
}

func TestDeleteFilter(t *testing.T) {
	t.Run("shows incomplete reminders", func(t *testing.T) {
		f := deleteFilter("", false)
		if f.Completed == nil || *f.Completed != false {
			t.Errorf("expected Completed=false, got %v", f.Completed)
		}
	})

	t.Run("list filter applied", func(t *testing.T) {
		f := deleteFilter("Personal", false)
		if f.ListName != "Personal" {
			t.Errorf("expected ListName=Personal, got %q", f.ListName)
		}
	})

	t.Run("flagged filter applied", func(t *testing.T) {
		f := deleteFilter("", true)
		if f.Flagged == nil || *f.Flagged != true {
			t.Errorf("expected Flagged=true, got %v", f.Flagged)
		}
	})

	t.Run("flagged not set when false", func(t *testing.T) {
		f := deleteFilter("", false)
		if f.Flagged != nil {
			t.Errorf("expected Flagged=nil, got %v", *f.Flagged)
		}
	})
}

func TestFlagFilter(t *testing.T) {
	t.Run("shows incomplete reminders", func(t *testing.T) {
		f := flagFilter("")
		if f.Completed == nil || *f.Completed != false {
			t.Errorf("expected Completed=false, got %v", f.Completed)
		}
	})

	t.Run("list filter applied", func(t *testing.T) {
		f := flagFilter("Work")
		if f.ListName != "Work" {
			t.Errorf("expected ListName=Work, got %q", f.ListName)
		}
	})

	t.Run("no flagged filter", func(t *testing.T) {
		f := flagFilter("")
		if f.Flagged != nil {
			t.Errorf("expected Flagged=nil, got %v", *f.Flagged)
		}
	})
}

func TestUnflagFilter(t *testing.T) {
	t.Run("shows incomplete reminders", func(t *testing.T) {
		f := unflagFilter("")
		if f.Completed == nil || *f.Completed != false {
			t.Errorf("expected Completed=false, got %v", f.Completed)
		}
	})

	t.Run("shows flagged reminders", func(t *testing.T) {
		f := unflagFilter("")
		if f.Flagged == nil || *f.Flagged != true {
			t.Errorf("expected Flagged=true, got %v", f.Flagged)
		}
	})

	t.Run("list filter applied", func(t *testing.T) {
		f := unflagFilter("Personal")
		if f.ListName != "Personal" {
			t.Errorf("expected ListName=Personal, got %q", f.ListName)
		}
	})
}
