package node

import (
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMerkleNode(t *testing.T) {
	datas := [][]byte{
		[]byte("node1"),
		[]byte("node2"),
		[]byte("node3"),
	}

	// level 1
	n1 := NewMerkleNode(nil, nil, datas[0])
	n2 := NewMerkleNode(nil, nil, datas[1])
	n3 := NewMerkleNode(nil, nil, datas[2])
	n4 := &MerkleNode{Hash: n3.Hash}

	// level 2
	n5 := NewMerkleNode(n1, n2, nil)
	n6 := NewMerkleNode(n3, n4, nil)

	// level 3
	root := NewMerkleNode(n5, n6, nil)

	assert.Equal(
		t,
		"64b04b718d8b7c5b6fd17f7ec221945c034cfce3be4118da33244966150c4bd4",
		hex.EncodeToString(n5.Hash),
		"Node 5 hash is correct",
	)

	assert.Equal(
		t,
		"08bd0d1426f87a78bfc2f0b13eccdf6f5b58dac6b37a7b9441c1a2fab415d76c",
		fmt.Sprintf("%x", n6.Hash),
		"Node 6 hash is correct",
	)

	assert.Equal(
		t,
		"4e3e44e55926330ab6c31892f980f8bfd1a6e910ff1ebc3f778211377f35227e",
		hex.EncodeToString(root.Hash),
		"Root hash is correct",
	)
}

func TestNewMerkleTree0(t *testing.T) {
	datas := [][]byte{}

	// merkle tree
	tree := NewMerkleTree(datas)

	assert.Equal(
		t,
		true,
		tree == nil,
	)
}

func TestNewMerkleTree1(t *testing.T) {
	datas := [][]byte{
		[]byte("node1"),
	}

	// level 1
	root := NewMerkleNode(nil, nil, datas[0])

	// merkle tree
	tree := NewMerkleTree(datas)

	assert.Equal(
		t,
		root.Hash,
		tree.Root.Hash,
	)
}

func TestNewMerkleTree2(t *testing.T) {
	datas := [][]byte{
		[]byte("node1"),
		[]byte("node2"),
	}

	// level 1
	n1 := NewMerkleNode(nil, nil, datas[0])
	n2 := NewMerkleNode(nil, nil, datas[1])

	// level 2
	root := NewMerkleNode(n1, n2, nil)

	// merkle tree
	tree := NewMerkleTree(datas)

	assert.Equal(
		t,
		root.Hash,
		tree.Root.Hash,
	)
}

func TestNewMerkleTree3(t *testing.T) {
	datas := [][]byte{
		[]byte("node1"),
		[]byte("node2"),
		[]byte("node3"),
	}

	// level 1
	n1 := NewMerkleNode(nil, nil, datas[0])
	n2 := NewMerkleNode(nil, nil, datas[1])
	n3 := NewMerkleNode(nil, nil, datas[2])
	n4 := &MerkleNode{Hash: n3.Hash}

	// level 2
	n5 := NewMerkleNode(n1, n2, nil)
	n6 := NewMerkleNode(n3, n4, nil)

	// level 3
	root := NewMerkleNode(n5, n6, nil)

	// merkle tree
	tree := NewMerkleTree(datas)

	assert.Equal(
		t,
		root.Hash,
		tree.Root.Hash,
	)
}

func TestNewMerkleTree4(t *testing.T) {
	datas := [][]byte{
		[]byte("node1"),
		[]byte("node2"),
		[]byte("node3"),
		[]byte("node4"),
	}

	// level 1
	n1 := NewMerkleNode(nil, nil, datas[0])
	n2 := NewMerkleNode(nil, nil, datas[1])
	n3 := NewMerkleNode(nil, nil, datas[2])
	n4 := NewMerkleNode(nil, nil, datas[3])

	// level 2
	n5 := NewMerkleNode(n1, n2, nil)
	n6 := NewMerkleNode(n3, n4, nil)

	// level 3
	root := NewMerkleNode(n5, n6, nil)

	// merkle tree
	tree := NewMerkleTree(datas)

	assert.Equal(
		t,
		root.Hash,
		tree.Root.Hash,
	)
}

func TestNewMerkleTree5(t *testing.T) {
	datas := [][]byte{
		[]byte("node1"),
		[]byte("node2"),
		[]byte("node3"),
		[]byte("node4"),
		[]byte("node5"),
	}

	// level 1
	n1 := NewMerkleNode(nil, nil, datas[0])
	n2 := NewMerkleNode(nil, nil, datas[1])
	n3 := NewMerkleNode(nil, nil, datas[2])
	n4 := NewMerkleNode(nil, nil, datas[3])
	n5 := NewMerkleNode(nil, nil, datas[4])
	n6 := &MerkleNode{Hash: n5.Hash}

	// level 2
	n7 := NewMerkleNode(n1, n2, nil)
	n8 := NewMerkleNode(n3, n4, nil)
	n9 := NewMerkleNode(n5, n6, nil)
	n10 := &MerkleNode{Hash: n9.Hash}

	// level 3
	n11 := NewMerkleNode(n7, n8, nil)
	n12 := NewMerkleNode(n9, n10, nil)

	// level 4
	root := NewMerkleNode(n11, n12, nil)

	// merkle tree
	tree := NewMerkleTree(datas)

	assert.Equal(
		t,
		root.Hash,
		tree.Root.Hash,
	)
}

func TestNewMerkleTree6(t *testing.T) {
	datas := [][]byte{
		[]byte("node1"),
		[]byte("node2"),
		[]byte("node3"),
		[]byte("node4"),
		[]byte("node5"),
		[]byte("node6"),
	}

	// level 1
	n1 := NewMerkleNode(nil, nil, datas[0])
	n2 := NewMerkleNode(nil, nil, datas[1])
	n3 := NewMerkleNode(nil, nil, datas[2])
	n4 := NewMerkleNode(nil, nil, datas[3])
	n5 := NewMerkleNode(nil, nil, datas[4])
	n6 := NewMerkleNode(nil, nil, datas[5])

	// level 2
	n7 := NewMerkleNode(n1, n2, nil)
	n8 := NewMerkleNode(n3, n4, nil)
	n9 := NewMerkleNode(n5, n6, nil)
	n10 := &MerkleNode{Hash: n9.Hash}

	// level 3
	n11 := NewMerkleNode(n7, n8, nil)
	n12 := NewMerkleNode(n9, n10, nil)

	// level 4
	root := NewMerkleNode(n11, n12, nil)

	// merkle tree
	tree := NewMerkleTree(datas)

	assert.Equal(
		t,
		root.Hash,
		tree.Root.Hash,
	)
}

func TestNewMerkleTree7(t *testing.T) {
	datas := [][]byte{
		[]byte("node1"),
		[]byte("node2"),
		[]byte("node3"),
		[]byte("node4"),
		[]byte("node5"),
		[]byte("node6"),
		[]byte("node7"),
	}

	// level 1
	n1 := NewMerkleNode(nil, nil, datas[0])
	n2 := NewMerkleNode(nil, nil, datas[1])
	n3 := NewMerkleNode(nil, nil, datas[2])
	n4 := NewMerkleNode(nil, nil, datas[3])
	n5 := NewMerkleNode(nil, nil, datas[4])
	n6 := NewMerkleNode(nil, nil, datas[5])
	n7 := NewMerkleNode(nil, nil, datas[6])
	n8 := &MerkleNode{Hash: n7.Hash}

	// level 2
	n9 := NewMerkleNode(n1, n2, nil)
	n10 := NewMerkleNode(n3, n4, nil)
	n11 := NewMerkleNode(n5, n6, nil)
	n12 := NewMerkleNode(n7, n8, nil)

	// level 3
	n13 := NewMerkleNode(n9, n10, nil)
	n14 := NewMerkleNode(n11, n12, nil)

	// level 4
	root := NewMerkleNode(n13, n14, nil)

	// merkle tree
	tree := NewMerkleTree(datas)

	assert.Equal(
		t,
		root.Hash,
		tree.Root.Hash,
	)
}

func TestNewMerkleTree8(t *testing.T) {
	datas := [][]byte{
		[]byte("node1"),
		[]byte("node2"),
		[]byte("node3"),
		[]byte("node4"),
		[]byte("node5"),
		[]byte("node6"),
		[]byte("node7"),
		[]byte("node8"),
	}

	// level 1
	n1 := NewMerkleNode(nil, nil, datas[0])
	n2 := NewMerkleNode(nil, nil, datas[1])
	n3 := NewMerkleNode(nil, nil, datas[2])
	n4 := NewMerkleNode(nil, nil, datas[3])
	n5 := NewMerkleNode(nil, nil, datas[4])
	n6 := NewMerkleNode(nil, nil, datas[5])
	n7 := NewMerkleNode(nil, nil, datas[6])
	n8 := NewMerkleNode(nil, nil, datas[7])

	// level 2
	n9 := NewMerkleNode(n1, n2, nil)
	n10 := NewMerkleNode(n3, n4, nil)
	n11 := NewMerkleNode(n5, n6, nil)
	n12 := NewMerkleNode(n7, n8, nil)

	// level 3
	n13 := NewMerkleNode(n9, n10, nil)
	n14 := NewMerkleNode(n11, n12, nil)

	// level 4
	root := NewMerkleNode(n13, n14, nil)

	// merkle tree
	tree := NewMerkleTree(datas)

	assert.Equal(
		t,
		root.Hash,
		tree.Root.Hash,
	)
}
