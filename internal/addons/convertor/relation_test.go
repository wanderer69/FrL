package convertor

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadRelation(t *testing.T) {
	relations, err := LoadRelation("/home/user/Go_projects/SemanticNet/data/relation.txt") // s_short
	require.NoError(t, err)
	for i := range relations {
		fmt.Printf("%v\r\n", relations[i].RelationType)
		for j := range relations[i].RelationItem {
			fmt.Printf("\t%#v\r\n", relations[i].RelationItem[j])
		}
	}

	require.True(t, false)
}
