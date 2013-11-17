//Package backend provides backend functions.
package backend

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/phack/robotmp/persistent"
	"github.com/phack/robotmp/markov"
	"github.com/NHTGD2013/twfy"
	"strings"
)

var apiKey string
type person twfy.PersonID

func (p person) key() (b []byte, err error) {
	var b bytes.Buffer
	_, err = binary.Write(b, binary.LittleEndian, p)
	return b, err
}


func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}


func MarkovMP(rw http.ResponseWriter, rq *http.Request) (err error) {
	//the last number of the URL is the MP's ID.
	var mp person
	s := strings.SplitN(reverse(rq.URL.String()), "/", 2)
	fmt.Sscan(" "+s[0]+" ", mpid)

	var chain *markov.Chain
	//find out if we already have the MP
	err = persistent.Retrieve(&chain, mp.key())
	if err != nil {
		if err == persistent.Inexistant{
			//TODO get the chain
		}
		return
	}

	rw.Header().Set("Content-Type", "text/plain")
	chain.Generate(rw, 100)
	return
}


func init() {
	flag.StringVar(&apiKey, "key", "", "API key.")
	flag.Parse()
	if apiKey == "" {
		glog.Fatal("API key cannot be empty!")
	}
	rand.Seed(time.Now().UnixNano())
}
