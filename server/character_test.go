package server_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/xackery/eqcp/pb"
)

func TestCharacterSCRUD(t *testing.T) {
	ctx := context.Background()
	s := serverSetup(t)

	//create
	respC, err := s.CharacterCreate(ctx, &pb.CharacterCreateRequest{Values: map[string]string{"name": "test"}})
	if err != nil {
		t.Fatal(err)
	}

	if respC.Id < 1 {
		t.Fatal("response invalid")
	}
	id := respC.Id

	//read
	respR, err := s.CharacterRead(ctx, &pb.CharacterReadRequest{Id: id})
	if err != nil {
		t.Fatal(err)
	}
	if respR.Character.Id != id {
		t.Fatalf("expected %d, got %d", respR.Character.Id, id)
	}

	//search
	respS, err := s.CharacterSearch(ctx, &pb.CharacterSearchRequest{Values: map[string]string{
		"name": respR.Character.Name,
		"id":   fmt.Sprintf("%d", respR.Character.Id),
	}})
	if err != nil {
		t.Fatal(err)
	}
	if respS == nil || len(respS.Characters) < 1 {
		t.Fatal("search failed to get any results")
	}

	//patch
	respP, err := s.CharacterPatch(ctx, &pb.CharacterPatchRequest{Id: id, Key: "name", Value: "test2"})
	if err != nil {
		t.Fatal(err)
	}
	if respP == nil || respP.Rowsaffected < 1 {
		t.Fatal("patch failed to get any results")
	}

	//update
	respU, err := s.CharacterUpdate(ctx, &pb.CharacterUpdateRequest{Values: map[string]string{"name": "test3"}})
	if err != nil {
		t.Fatal(err)
	}
	if respU == nil || respP.Rowsaffected < 1 {
		t.Fatal("updated failed to get any results")
	}

	//delete
	respD, err := s.CharacterDelete(ctx, &pb.CharacterDeleteRequest{Id: id})
	if err != nil {
		t.Fatal(err)
	}
	if respD == nil || respD.Rowsaffected < 1 {
		t.Fatal("updated failed to get any results")
	}
}
