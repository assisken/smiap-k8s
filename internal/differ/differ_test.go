package differ

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func (someDiffObject) GetAnnotations() map[string]string            { return nil }
func (someDiffObject) SetAnnotations(annotations map[string]string) {}

type someDiffObject struct {
	Name string
}

func TestDiff(t *testing.T) {
	OldContainer := &someDiffObject{"old-container"}
	NewContainer := &someDiffObject{"new-container"}

	for _, tt := range []struct {
		name         string
		want         map[string]DiffObject
		have         map[string]DiffObject
		wantToCreate []string
		wantToDelete []string
		wantToUpdate []string
	}{
		{
			name:         "Empty list",
			want:         map[string]DiffObject{},
			have:         map[string]DiffObject{},
			wantToCreate: []string{},
			wantToDelete: []string{},
			wantToUpdate: []string{},
		},
		{
			name:         "Create some",
			want:         map[string]DiffObject{"container": NewContainer},
			have:         map[string]DiffObject{},
			wantToCreate: []string{"container"},
			wantToDelete: []string{},
			wantToUpdate: []string{},
		},
		{
			name:         "Delete some",
			want:         map[string]DiffObject{},
			have:         map[string]DiffObject{"container": OldContainer},
			wantToCreate: []string{},
			wantToDelete: []string{"container"},
			wantToUpdate: []string{},
		},
		{
			name:         "Update some",
			want:         map[string]DiffObject{"container": NewContainer},
			have:         map[string]DiffObject{"container": OldContainer},
			wantToCreate: []string{},
			wantToDelete: []string{},
			wantToUpdate: []string{"container"},
		},
		{
			name:         "Do nothing",
			want:         map[string]DiffObject{"container": OldContainer},
			have:         map[string]DiffObject{"container": OldContainer},
			wantToCreate: []string{},
			wantToDelete: []string{},
			wantToUpdate: []string{},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			gotToCreate, gotToDelete, gotToUpdate, err := Diff(tt.want, tt.have)

			assert.NoError(t, err)

			toCreate := make([]string, len(gotToCreate))
			for i := range gotToCreate {
				toCreate[i] = gotToCreate[i].Id
			}

			toDelete := make([]string, len(gotToDelete))
			for i := range gotToDelete {
				toDelete[i] = gotToDelete[i].Id
			}

			toUpdate := make([]string, len(gotToUpdate))
			for i := range gotToUpdate {
				toUpdate[i] = gotToUpdate[i].Id
			}

			assert.Equal(t, tt.wantToCreate, toCreate)
			assert.Equal(t, tt.wantToDelete, toDelete)
			assert.Equal(t, tt.wantToUpdate, toUpdate)
		})
	}
}
