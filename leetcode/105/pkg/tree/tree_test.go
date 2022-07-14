package tree_test

import (
	"testing"

	"github.com/matryer/is"

	"go.jlucktay.dev/golang-workbench/leetcode/105/pkg/tree"
)

// Given two integer arrays 'preorder' and 'inorder' where 'preorder' is the preorder traversal of a binary tree and
// 'inorder' is the inorder traversal of the same tree, construct and return _the binary tree_.

func TestBuildTree(t *testing.T) {
	testCases := map[string]struct {
		expected          *tree.TreeNode
		inorder, preorder []int
	}{
		"Example 1": {
			inorder:  []int{9, 3, 15, 20, 7},
			preorder: []int{3, 9, 20, 15, 7},
			expected: &tree.TreeNode{
				Val:  3,
				Left: &tree.TreeNode{Val: 9},
				Right: &tree.TreeNode{
					Val:   20,
					Left:  &tree.TreeNode{Val: 15},
					Right: &tree.TreeNode{Val: 7},
				},
			},
		},
		"Example 2": {
			inorder:  []int{-1},
			preorder: []int{-1},
			expected: &tree.TreeNode{Val: -1},
		},
	}
	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			is := is.New(t)

			is.Equal(
				tree.BuildTree(tc.preorder, tc.inorder),
				tc.expected,
			)
		})
	}
}
