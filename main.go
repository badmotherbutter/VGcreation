package main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"os"
)

type Product struct {
	XMLName xml.Name `xml:"product"`
	ID      string   `xml:"product-id,attr"`
	VGID    *Variation
	Color   CustomAttr `xml:"custom-attributes>custom-attribute,omitempty"`
}

type Variation struct {
	XMLName xml.Name `xml:"variations,omitempty"`
	ID      string   `xml:"product-id,attr"`
}

type CustomAttr struct {
	XMLName xml.Name `xml:"custom-attribute"`
	ID      string   `xml:"attribute-id,attr"`
	Value   string   `xml:",chardata"`
}

func main() {

	csvFile, err := os.Open("source/products.csv")

	if err != nil {
		fmt.Println("Err: ", err)
		return
	}

	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.Comma = ';'

	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	//var masters = make([]Product, len(records), len(records))
	var vGs = make([]Product, len(records), len(records))

	//Loop al file rows
	for i := 0; i < len(records); i++ {

		//Loop single csv element

		/*	masters[i] = Product{
			ID: records[i][1] + "_" + records[i][3],
			VGID: Variation{
				ID: records[i][1] + "_" + records[i][3] + "_" + records[i][12],
			},
		}*/

		vGs[i] = Product{
			ID: records[i][1] + "_" + records[i][3] + "_" + records[i][12],
			VGID: &Variation{
				ID: "ook",
			},
			Color: CustomAttr{
				ID:    "color",
				Value: records[i][12],
			},
		}

	}

	output, err := xml.MarshalIndent(vGs, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	os.Stdout.Write(output)

}
