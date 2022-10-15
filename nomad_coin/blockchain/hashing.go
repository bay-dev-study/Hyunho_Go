package blockchain

import "crypto/sha256"

func Hashing(data []byte) [32]byte {
	return sha256.Sum256(data)
}
