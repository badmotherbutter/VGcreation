package main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
)

type Product struct {
	XMLName    xml.Name `xml:"product"`
	ID         string   `xml:"product-id,attr"`
	Searchable bool     `xml:"searchable-flag,omitempty"`
	Brand      string   `xml:"brand,omitempty"`
	VGID       *Variations
	Color      *CustomAttrs
}

type Variations struct {
	XMLName         xml.Name `xml:"variations,omitempty"`
	VariationGroups *VariationGroups
}

type VariationGroups struct {
	XMLName        xml.Name `xml:"variation-groups,omitempty"`
	VariationGroup *VariationGroup
}

type VariationGroup struct {
	XMLName xml.Name `xml:"variation-group,omitempty"`
	ID      string   `xml:"product-id,attr"`
}

type CustomAttrs struct {
	XMLName    xml.Name `xml:"custom-attributes,omitempty"`
	CustomAttr *CustomAttr
}

type CustomAttr struct {
	XMLName xml.Name `xml:"custom-attribute,omitempty"`
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
	//Set ; as CSV separator
	reader.Comma = ';'

	//Read all CSV records
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	//Create a bidimensional slice to store both master product and variations
	var products = make([][]Product, 0, len(records))

	//Loop all file rows
	for i := 0; i < len(records); i++ {

		colorID := records[i][12]
		if nb, _ := strconv.Atoi(colorID); nb < 10 {
			colorID = "0" + colorID
		}

		variationID := records[i][1] + "_" + records[i][3] + "_" + colorID
		//Check if master is
		if !valueInSlice(variationID, products) {
			masterID := records[i][1] + "_" + records[i][3]

			products = append(products, []Product{
				Product{
					ID: masterID,
					VGID: &Variations{
						VariationGroups: &VariationGroups{
							VariationGroup: &VariationGroup{
								ID: variationID,
							},
						},
					},
				},
				Product{
					ID:         variationID,
					Searchable: true,
					Brand:      records[i][40],
					Color: &CustomAttrs{
						CustomAttr: &CustomAttr{
							ID:    "color",
							Value: colorID,
						},
					},
				},
			})
		}

	}

	//Marshal products to XML
	output, err := xml.MarshalIndent(products, "  ", "    ")

	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	//Create target file
	fXML, err := os.Create("dest/vgcatalog.xml")

	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	_, err1 := fXML.Write(output)

	if err1 != nil {
		fmt.Printf("error: %v\n", err1)
	}

	fXML.Sync()

}

//Check if a Variation group ID isinside the products slice
func valueInSlice(val string, slice [][]Product) bool {

	if len(slice) > 0 {
		for _, v := range slice {
			if val == v[1].ID {
				return true
			}
		}
	}
	return false
}
