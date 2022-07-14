package tree

// Definition for a binary tree node.
type TreeNode struct {
	Left  *TreeNode
	Right *TreeNode
	Val   int
}

/*
Pre-order, NLR
1. Visit the current node.
2. Recursively traverse the current node's left subtree.
3. Recursively traverse the current node's right subtree.

The pre-order traversal is a topologically sorted one, because a parent node is processed before any of its child nodes
is done.

In-order, LNR
1. Recursively traverse the current node's left subtree.
2. Visit the current node.
3. Recursively traverse the current node's right subtree.

In a binary search tree ordered such that in each node the key is greater than all keys in its left subtree and less
than all keys in its right subtree, in-order traversal retrieves the keys in ascending sorted order.
*/

/*
With the pre-order traversal [𝑥1,…,𝑥𝑛] and the in-order traversal [𝑧1,…,𝑧𝑛], you can rebuild the tree as follows:

The root is the head of the pre-order traversal 𝑥1.

Let 𝑘 be the index such that 𝑧𝑘=𝑥1.

Then [𝑧1,…,𝑧𝑘−1] is the in-order traversal of the left child and [𝑧𝑘+1,…,𝑧𝑛] is the in-order traversal of the
right child.

Going by the number of elements, [𝑥2,…,𝑥𝑘] is the pre-order traversal of the left child and [𝑥𝑘+1,…,𝑥𝑛] that of
the right child.

Recurse to build the left and right subtrees.
*/

func BuildTree(preorder []int, inorder []int) *TreeNode {
	// With
	//   the pre-order traversal [𝑥1,…,𝑥𝑛] and
	//   the in-order traversal [𝑧1,…,𝑧𝑛],
	// you can rebuild the tree as follows:

	if len(preorder) < 1 || len(inorder) < 1 {
		return nil
	}

	// The root is the head of the pre-order traversal 𝑥1.
	root := &TreeNode{Val: preorder[0]}

	// Let 𝑘 be the index such that 𝑧𝑘=𝑥1.
	k := 0
	for a := range inorder {
		if inorder[a] == preorder[0] {
			k = a + 1 // k is set as a 1-based index, not a 0-based
			break
		}
	}

	// Then
	//   [𝑧1,…,𝑧𝑘−1] is the in-order traversal of the left child and
	//   [𝑧𝑘+1,…,𝑧𝑛] is the in-order traversal of the right child.
	leftInorder := make([]int, 0)
	if len(inorder) >= k-1 {
		leftInorder = append(leftInorder, inorder[:k-1]...)
	}

	rightInorder := make([]int, 0)
	if len(inorder) >= k {
		rightInorder = append(rightInorder, inorder[k:]...)
	}

	// Going by the number of elements,
	//   [𝑥2,…,𝑥𝑘] is the pre-order traversal of the left child and
	//   [𝑥𝑘+1,…,𝑥𝑛] is the pre-order traversal of the right child.
	leftPreorder := make([]int, 0)
	if len(preorder) >= k {
		leftPreorder = append(leftPreorder, preorder[1:k]...)
	}

	rightPreorder := make([]int, 0)
	if len(preorder) >= k {
		rightPreorder = append(rightPreorder, preorder[k:]...)
	}

	// Recurse to build the left and right subtrees.
	root.Left = BuildTree(leftPreorder, leftInorder)
	root.Right = BuildTree(rightPreorder, rightInorder)

	return root
}
