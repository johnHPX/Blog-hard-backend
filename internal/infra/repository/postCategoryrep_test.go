package repository

import (
	"fmt"
	"testing"
)

func TestFind(t *testing.T) {
	t.Run("teste positivo", func(t *testing.T) {
		rep := NewPostCategoryRepository()
		entity, err := rep.Find("08990115-1540-4bde-9ed7-d916dcac1730", "e34c2924-fadd-4cde-9554-69e99ddcdd3c")
		if err != nil {
			t.Errorf("Error: %s", err.Error())
		}
		fmt.Println(entity.PostId)
	})
}
