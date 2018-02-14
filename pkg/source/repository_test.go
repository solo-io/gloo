package source

import "testing"

func TestAddUpstream(t *testing.T) {
	repo := newRepo()
	new := repo.set(Upstream{
		ID:   "test",
		Name: "us-east-1"})
	if !new {
		t.Error("expected upstream to be new")
	}
	upstreams := repo.upstreams()
	if len(upstreams) != 1 {
		t.Errorf("expected number of upstreams to be 1 but was %d", len(upstreams))
	}
	_, exists := repo.get("test")
	if !exists {
		t.Error("unable to get back saved upstream")
	}
}

func TestUpdateUpstream(t *testing.T) {
	repo := newRepo()
	new := repo.set(Upstream{ID: "test", Name: "us-east-1"})
	if !new {
		t.Error("expected upstream to be new")
	}
	new = repo.set(Upstream{ID: "test",
		Name: "us-east-1",
		Functions: []Function{
			Function{
				Name: "func1",
				Spec: map[string]interface{}{
					"Qualifier": "v1"},
			},
		},
	})
	if new {
		t.Error("expected upstream to be updated")
	}
	upstreams := repo.upstreams()
	if len(upstreams) != 1 {
		t.Errorf("expected number of upstreams to be 1 but was %d", len(upstreams))
	}
	r, exists := repo.get("test")
	if !exists {
		t.Error("unable to get back saved upstream")
	}
	if len(r.Functions) != 1 {
		t.Error("unable to see lambdas in updated upstream")
	}
}

func TestDeleteUpstream(t *testing.T) {
	repo := newRepo()
	new := repo.set(Upstream{ID: "test", Name: "us-east-1"})
	if !new {
		t.Error("expected upstream to be new")
	}
	repo.delete("test")

	upstreams := repo.upstreams()
	if len(upstreams) != 0 {
		t.Errorf("expected number of upstreams to be 0 but was %d", len(upstreams))
	}
	_, exists := repo.get("test")
	if exists {
		t.Error("able to get deleted upstream")
	}
}
