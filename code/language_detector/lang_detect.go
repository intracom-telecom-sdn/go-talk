package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
        "os"
)



type PairSI struct {
	S string
	I int
}


func fileToString(filename string) string {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return string(file)
}

func compressedSize(s string) int {
	var b bytes.Buffer
	compressor, _ := gzip.NewWriterLevel(&b, gzip.BestCompression)

	_, err := compressor.Write([]byte(s))
	if err != nil {
		panic(err)
	}
	compressor.Close()
	return b.Len()
}

func main() {
	languages := []string{"EN", "DE"}
	corpus := make(map[string]string)
	corpus_sz := make(map[string]int)

	done := make(chan bool)

	// Preprocessing: Create ground truth
	for _, lang := range languages {
		go func(lang string) { 
			corpus[lang] = fileToString("corpora/" + lang)
			corpus_sz[lang] = compressedSize(corpus[lang])
			done <- true
		}(lang)
	}

	for i := 0; i < len(languages); i++ { <-done }	

	
	// Testing: Evaluate phrases


	phrase_len := len(os.Args) - 1

	phrases := make([]string, phrase_len)

	for i := range phrases {
		phrases[i] = fileToString(os.Args[i+1])
	}

	

	classification := make(chan PairSI)

	//fine-grain start OMIT
	for i := range phrases {
		go func(idx int, phrase string) { // HL
			results := make(chan PairSI)
			for _, lang := range languages {
				go func(lang string, results *chan PairSI) { // HL
					eval := phrase + corpus[lang]
					eval_sz := compressedSize(eval) - corpus_sz[lang]
					*results <- PairSI{S: lang, I: eval_sz}
				}(lang, &results)
			}

			class := PairSI{S: "", I: 1<<32 - 1}
			for i := 0; i < len(languages); i++ {
				curr_class := <- results
				if curr_class.I < class.I  {
					class.S = curr_class.S
					class.I = curr_class.I
				}
			}
			class.I = idx
			classification <- class
  		}(i, phrases[i])
	}
	//fine-grain end OMIT

	for i := 0; i < len(phrases); i++ {
		class := <-classification 
		fmt.Println("PHRASE:")   	
		fmt.Println(phrases[class.I])
		fmt.Println("CLASSIFICATION:", class.S)
		fmt.Println()
		fmt.Println("-------------------------------------------------------------------------------------------------------------------")
		fmt.Println()
	}	
	
}
