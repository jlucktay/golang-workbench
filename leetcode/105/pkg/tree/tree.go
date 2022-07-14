package tree

// Definition for a binary tree node.
type TreeNode struct {
	Left  *TreeNode
	Right *TreeNode
	Val   int
}

func BuildTree(preorder []int, inorder []int) *TreeNode {
	return &TreeNode{Val: -1}
}
