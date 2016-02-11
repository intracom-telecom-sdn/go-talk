package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
)


///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var EN_PHRASE string = `A very hard, brittle, silvery-white transition metal of the platinum family, iridium is the second-densest element (after osmium) and is the most corrosion-resistant metal, even at temperatures as high as 2000 °C. Although only certain molten salts and halogens are corrosive to solid iridium, finely divided iridium dust is much more reactive and can be flammable. Iridium was discovered in 1803 among insoluble impurities in natural platinum. Smithson Tennant, the primary discoverer, named the iridium for the goddess Iris, personification of the rainbow, because of the striking and diverse colors of its salts. Iridium is one of the rarest elements in the Earth's crust, with annual production and consumption of only three tonnes.`


var DE_PHRASE string = `Iridium ist ein chemisches Element mit dem Symbol Ir und der Ordnungszahl 77. Es zählt zu den Übergangsmetallen, im Periodensystem steht es in der Gruppe 9 (in der älteren Zählung Teil der 8. Nebengruppe) oder Cobaltgruppe. Das sehr schwere, harte, spröde, silber-weiß glänzende Edelmetall aus der Gruppe der Platinmetalle gilt als das korrosionsbeständigste Element. Unter 0,11 Kelvin wechselt es in den supraleitfähigen Zustand über.`

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////


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

	phrases := []string{ EN_PHRASE, DE_PHRASE }
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
	   	fmt.Println(phrases[class.I], class.S)
		fmt.Println()
	}	
	
}
