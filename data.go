package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var privateKey8180 		*rsa.PrivateKey
var privateKey8181 		*rsa.PrivateKey
var privateKey8182 		*rsa.PrivateKey
var privateKey8183 		*rsa.PrivateKey
var privateKey8184 		*rsa.PrivateKey
var privateKey8185 		*rsa.PrivateKey
var privateKey8186 		*rsa.PrivateKey
var privateKeyClient  	*rsa.PrivateKey
var publicKey8180  		*rsa.PublicKey
var publicKey8181  		*rsa.PublicKey
var publicKey8182  		*rsa.PublicKey
var publicKey8183  		*rsa.PublicKey
var publicKey8184  		*rsa.PublicKey
var publicKey8185  		*rsa.PublicKey
var publicKey8186  		*rsa.PublicKey
var publicKeyClient 	*rsa.PublicKey
var KnownNodes []*KnownNode
var KeypairMap map[int]Keypair
var ClientNode *KnownNode
func init(){
	fmt.Println("data.go init running")   // temporary debug aid
	var err error
	generateKeyFiles()
	privateKey8180, publicKey8180, err = getKeyPairByFile(0)
	if err != nil {
		panic(err)
	}
	privateKey8181, publicKey8181, err = getKeyPairByFile(1)
	if err != nil {
		panic(err)
	}
	privateKey8182, publicKey8182, err = getKeyPairByFile(2)
	if err != nil {
		panic(err)
	}
	privateKey8183, publicKey8183, err = getKeyPairByFile(3)
	if err != nil {
		panic(err)
	}
	privateKey8184, publicKey8184, err = getKeyPairByFile(4)
	if err != nil {
		panic(err)
	}
	privateKey8185, publicKey8185, err = getKeyPairByFile(5)
	if err != nil {
		panic(err)
	}
	privateKey8186, publicKey8186, err = getKeyPairByFile(6)
	if err != nil {
		panic(err)
	}
	privateKeyClient, publicKeyClient, err = getKeyPairByFile(7)
	if err != nil {
		panic(err)
	}
	KnownNodes = []*KnownNode{
		{
			0,
			"localhost:8180",
			publicKey8180,
		},
		{
			1,
			"localhost:8181",
			publicKey8181,
		},
		{
			2,
			"localhost:8182",
			publicKey8182,
		},
		{
			3,
			"localhost:8183",
			publicKey8183,
		},
		{
			4,
			"localhost:8184",
			publicKey8184,
		},
		{
			5,
			"localhost:8185",
			publicKey8185,
		},
		{
			6,
			"localhost:8186",
			publicKey8186,
		},
	}
	KeypairMap = map[int]Keypair{
		0:{
			privateKey8180,
			publicKey8180,
		},
		1:{
			privateKey8181,
			publicKey8181,
		},
		2:{
			privateKey8182,
			publicKey8182,
		},
		3:{
			privateKey8183,
			publicKey8183,
		},
		4:{
			privateKey8184,
			publicKey8184,
		},
		5:{
			privateKey8185,
			publicKey8185,
		},
		6:{
			privateKey8186,
			publicKey8186,
		},
		7:{
			privateKeyClient,
			publicKeyClient,
		},
	}
	ClientNode = &KnownNode{
		7,
		"localhost:8187",
		publicKeyClient,
	}
}

func getKeyPairByFile(nodeID int) (*rsa.PrivateKey, *rsa.PublicKey, error){
	privFile, _ := filepath.Abs(fmt.Sprintf("./Keys/%d_priv",nodeID))
	pubFile, _ := filepath.Abs(fmt.Sprintf("./Keys/%d_pub",nodeID))
	fbytes, err := ioutil.ReadFile(privFile)
	if err != nil {
		return nil,nil, err
	}
	block, _:= pem.Decode(fbytes)
	if block == nil {
		return nil,nil, fmt.Errorf("parse block occured error")
	}
	privkey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil{
		return nil,nil, err
	}
	pubfbytes, err := ioutil.ReadFile(pubFile)
	if err != nil {
		return nil,nil, err
	}
	pubblock, _:= pem.Decode(pubfbytes)
	if pubblock == nil {
		return nil,nil, fmt.Errorf("parse block occured error")
	}
	pubkey, err := x509.ParsePKIXPublicKey(pubblock.Bytes)
	if err != nil{
		return nil,nil, err
	}
	return privkey, pubkey.(*rsa.PublicKey), nil
}

func generateKeyFiles(){
	if !FileExists("./Keys") {
		fmt.Println("hello")
		err := os.Mkdir("Keys", 0700)
		if err != nil {
			panic(err)
		}
		for i := 0; i<=7; i++ {
			filename, _ := filepath.Abs(fmt.Sprintf("./Keys/%d",i))
			if !FileExists(filename + "_priv") && !FileExists(filename + "_pub"){
				fmt.Println("creating keypair")
				priv, pub := generateKeyPair()
				err := ioutil.WriteFile(filename+"_priv", priv, 0644)
				if err != nil {
					panic(err)
				}
				ioutil.WriteFile(filename+"_pub", pub, 0644)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func generateKeyPair() ([]byte, []byte){
	privkey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	mprivkey := x509.MarshalPKCS1PrivateKey(privkey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: mprivkey,
	}
	bprivkey := pem.EncodeToMemory(block)
	pubkey := &privkey.PublicKey
	mpubkey, err := x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		panic(err)
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: mpubkey,
	}
	bpubkey := pem.EncodeToMemory(block)
	return bprivkey, bpubkey
}