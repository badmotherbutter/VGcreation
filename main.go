package main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"os"
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
	reader.Comma = ';'

	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	//var masters = make([]Product, len(records), len(records))
	var products = make([][]Product, 0, len(records))

	//Loop all file rows
	for i := 0; i < len(records); i++ {

		masterID := records[i][1] + "_" + records[i][3]

		if !valueInSlice(masterID, products) {
			variationID := records[i][1] + "_" + records[i][3] + "_" + records[i][12]
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
							Value: records[i][12],
						},
					},
				},
			})
		}

	}

	output, err := xml.MarshalIndent(products, "  ", "    ")

	//output = append(output, "</catalog>"...)

	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

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

func valueInSlice(val string, slice [][]Product) bool {

	if len(slice) > 0 {
		for _, v := range slice {
			if val == v[0].ID || val == v[1].ID {
				return true
			}
		}
	}
	return false
}
