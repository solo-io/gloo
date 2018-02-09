package aws

import "testing"

func TestAddRegion(t *testing.T) {
	repo := newRepo()
	new := repo.set(Region{
		ID:   "test",
		Name: "us-east-1"})
	if !new {
		t.Error("expected region to be new")
	}
	regions := repo.regions()
	if len(regions) != 1 {
		t.Errorf("expected number of regions to be 1 but was %d", len(regions))
	}
	_, exists := repo.get("test")
	if !exists {
		t.Error("unable to get back saved region")
	}
}

func TestUpdateRegion(t *testing.T) {
	repo := newRepo()
	new := repo.set(Region{ID: "test", Name: "us-east-1"})
	if !new {
		t.Error("expected region to be new")
	}
	new = repo.set(Region{ID: "test",
		Name:    "us-east-1",
		Lambdas: []Lambda{Lambda{Name: "func1", Qualifier: "v1"}}})
	if new {
		t.Error("expected region to be updated")
	}
	regions := repo.regions()
	if len(regions) != 1 {
		t.Errorf("expected number of regions to be 1 but was %d", len(regions))
	}
	r, exists := repo.get("test")
	if !exists {
		t.Error("unable to get back saved region")
	}
	if len(r.Lambdas) != 1 {
		t.Error("unable to see lambdas in updated region")
	}
}

func TestDeleteRegion(t *testing.T) {
	repo := newRepo()
	new := repo.set(Region{ID: "test", Name: "us-east-1"})
	if !new {
		t.Error("expected region to be new")
	}
	repo.delete("test")

	regions := repo.regions()
	if len(regions) != 0 {
		t.Errorf("expected number of regions to be 0 but was %d", len(regions))
	}
	_, exists := repo.get("test")
	if exists {
		t.Error("able to get deleted region")
	}
}
