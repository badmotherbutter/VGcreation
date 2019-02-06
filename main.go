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
	VGID    *Variations
	Color   *CustomAttrs
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
	var products = make([][]Product, len(records), len(records))

	//Loop all file rows
	for i := 0; i < len(records); i += 2 {

		products[i] = []Product{
			Product{
				ID: records[i][1] + "_" + records[i][3],
				VGID: &Variations{
					VariationGroups: &VariationGroups{
						VariationGroup: &VariationGroup{
							ID: records[i][1] + "_" + records[i][3] + "_" + records[i][12],
						},
					},
				},
			},
			Product{
				ID: records[i][1] + "_" + records[i][3] + "_" + records[i][12],
				Color: &CustomAttrs{
					CustomAttr: &CustomAttr{
						ID:    "color",
						Value: records[i][12],
					},
				},
			},
		}

	}

	output, err := xml.MarshalIndent(products, "  ", "    ")

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
