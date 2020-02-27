package node

import "crypto/sha256"

type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Hash  []byte
}

func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	node := MerkleNode{}

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		node.Hash = hash[:]
	} else {
		prevHashes := append(left.Hash, right.Hash...)
		hash := sha256.Sum256(prevHashes)
		node.Hash = hash[:]
	}

	node.Left = left
	node.Right = right
	return &node
}

type MerkleTree struct {
	Root *MerkleNode
}

func NewMerkleTree(datas [][]byte) *MerkleTree {
	if len(datas) < 1 {
		return nil
	}

	var nodes []MerkleNode

	for _, data := range datas {
		node := NewMerkleNode(nil, nil, data)
		nodes = append(nodes, *node)
	}

	for len(nodes) > 1 {
		var newLevel []MerkleNode

		if len(nodes) % 2 != 0 {
			node := MerkleNode{Hash: nodes[len(nodes) - 1].Hash}
			nodes = append(nodes, node)
		}

		for j, end := 0, len(nodes) - 1; j < end; j += 2 {
			node := NewMerkleNode(&nodes[j], &nodes[j + 1], nil)
			newLevel = append(newLevel, *node)
		}
		nodes = newLevel
	}

	return &MerkleTree{&nodes[0]}
}
