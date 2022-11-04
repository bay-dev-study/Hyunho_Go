package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"
	"nomad_coin/utils"
	"os"
	"sync"
)

var walletOnce = sync.Once{}

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string
}

var loadedWallet *wallet

const privateKeyFileName string = "key.wallet"

func privateKeyFileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	if err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		utils.ErrHandler(err)
	}
	return false
}

func createPrivateKey() *ecdsa.PrivateKey {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.ErrHandler(err)
	return privKey
}

func getPublicKeyInString(privateKey *ecdsa.PrivateKey) string {
	return encodeBigIntsToHexString(privateKey.X, privateKey.Y)
}

func encodeBigIntsToHexString(a, b *big.Int) string {
	z := append(a.Bytes(), b.Bytes()...)
	return fmt.Sprintf("%x", z)
}

func decodeHexStringToBigInts(hexString string) (*big.Int, *big.Int) {
	bytes, err := hex.DecodeString(hexString)
	utils.ErrHandler(err)
	bytesLength := len(bytes)
	firstHalf := bytes[:int(bytesLength/2)]
	secondHalf := bytes[int(bytesLength/2):]

	var a, b = &big.Int{}, &big.Int{}
	a.SetBytes(firstHalf)
	b.SetBytes(secondHalf)
	return a, b
}

func Sign(payloadInHexString string, w *wallet) string {
	payloadAsBytes, err := hex.DecodeString(payloadInHexString)
	utils.ErrHandler(err)
	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, payloadAsBytes)
	utils.ErrHandler(err)
	return encodeBigIntsToHexString(r, s)
}

func Verify(signature, payloadInHexString, address string) bool {
	r, s := decodeHexStringToBigInts(signature)
	x, y := decodeHexStringToBigInts(address)
	publicKey := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}
	payloadInBytes, err := hex.DecodeString(payloadInHexString)
	utils.ErrHandler(err)
	return ecdsa.Verify(publicKey, payloadInBytes, r, s)
}

func loadPrivateKeyFile(fileName string) *ecdsa.PrivateKey {
	privateKeyInBytes, err := os.ReadFile(fileName)
	utils.ErrHandler(err)
	privateKey, err := x509.ParseECPrivateKey(privateKeyInBytes)
	utils.ErrHandler(err)
	return privateKey
}

func savePrivateKeyFile(fileName string, privateKey *ecdsa.PrivateKey) {
	privateKeyInBytes, err := x509.MarshalECPrivateKey(privateKey)
	utils.ErrHandler(err)
	err = os.WriteFile(fileName, privateKeyInBytes, 0644)
	utils.ErrHandler(err)
}

func GetWallet() *wallet {
	walletOnce.Do(func() {
		loadedWallet = &wallet{}
		if privateKeyFileExist(privateKeyFileName) {
			loadedWallet.privateKey = loadPrivateKeyFile(privateKeyFileName)
		} else {
			loadedWallet.privateKey = createPrivateKey()
			savePrivateKeyFile(privateKeyFileName, loadedWallet.privateKey)
		}
		loadedWallet.Address = getPublicKeyInString(loadedWallet.privateKey)
	})
	return loadedWallet
}
