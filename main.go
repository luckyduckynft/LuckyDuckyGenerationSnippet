package main

import (
	"encoding/json"
	"fmt"
	wr "github.com/mroth/weightedrand"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var config Config

func init() {
	if !Exists("config.json") {
		f, err := os.OpenFile("config.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			log.Fatal(err)
		}
		writeme, _ := json.Marshal(config)
		if _, err := f.WriteString(string(writeme)); err != nil {
			log.Fatal(err)
		}
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	} else {
		jsonFile, err := os.Open("config.json")
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()
		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &config)
	}
}

func main() {
	generateMetatata()
	fmt.Println("All Done! Check the output folder and you can cross verify everything matches up :)")
	fmt.Scanln()
}

func generateMetatata() {
	fmt.Println("Loading CSV...")

	makeDirectoryIfNotExists(filepath.FromSlash("./output"))
	makeDirectoryIfNotExists(filepath.FromSlash("./output/metadata"))

	f, _ := os.Open("assets.csv")
	r := NewScanner(f)
	var allAssets []Asset
	for r.Scan() {
		weight := r.Text("Weight")
		weightInt, _ := strconv.Atoi(weight)
		asset := Asset{
			Name:   r.Text("Name"),
			Type:   r.Text("Type"),
			Weight: weightInt,
		}
		allAssets = append(allAssets, asset)
	}
	var allConflicts []Conflict
	if Exists("conflicts.csv") {
		f3, _ := os.Open("conflicts.csv")
		r3 := NewScanner(f3)

		for r3.Scan() {
			conflict := Conflict{
				Name: r3.Text("Name"),
				Type: r3.Text("Type"),
			}
			allConflicts = append(allConflicts, conflict)
		}
	}

	fmt.Println("Loading Legendaries...")

	var legendaryNames []string
	filepath.Walk("legendaries", func(path string, info os.FileInfo, err error) error {
		//fmt.Println(niceName)
		if info.IsDir() && strings.Contains(path, string(os.PathSeparator)) {
			fmt.Println(path)
			niceName := strings.Split(path, filepath.FromSlash("legendaries/"))[1]
			fmt.Println(niceName)
			//static
			legendaryNames = append(legendaryNames, niceName)
			//allAssets[category]
		}
		return nil
	})

	fmt.Println("Loading Asset Order...")
	var allTypes []string
	f2, _ := os.Open("types.csv")
	r2 := NewScanner(f2)
	m := make(map[string]int)
	for r2.Scan() {
		typeName := r2.Text("Type")
		order, _ := strconv.Atoi(r2.Text("Order"))
		m[typeName] = order
	}
	for i := 0; i < len(m); i++ {
		for k, v := range m {
			if v == i {
				if !strings.Contains(k, string(os.PathSeparator)) {
					//fmt.Println(k)
					allTypes = append(allTypes, k)
				}
			}
		}
	}

	fmt.Println("Generating Metadata...")
	var generatedMetatata []Metatata
	var generatedMetadataStrings []string
	for i := config.StartIndex; i <= config.EndIndex; i++ {
		complete := false
		var body string
		for !complete {
			var meta Metatata
			fileName := strconv.Itoa(i)
			meta.Name = strings.Replace(config.AssetName, "ID", fileName, -1)
			meta.Image = strings.Replace(config.ImagePath, "ID", fileName, -1)
			meta.Description = config.Description
			for _, attribute := range allTypes {
				if !strings.Contains(attribute, string(os.PathSeparator)) {
					choices := make([]wr.Choice, 0)
					for _, asset := range allAssets {
						if asset.Type == attribute {
							if strings.Contains(body, "WAVE") || strings.Contains(body, "STRIPE") {
								conflicts := []string{
									"SHIRT_BBALL",
									"DETECTIVE",
									"SHIRT_BUSINESS",
									"DRESS",
									"DRESS_LOCKETT",
									"TURTLE_GREEN",
									"TURTLE_BLUE",
									"VEST_BLACK",
									"VEST_BROWN",
									"BANDANA",
									"SCARF",
								}
								if !contains(conflicts, asset.Name) {
									c := wr.NewChoice(asset.Name, uint(asset.Weight))
									choices = append(choices, c)
								}
							} else {
								if asset.Name == "NONE" && asset.Type == "ACCESSORY" {
									//re roll if naked

								} else {
									c := wr.NewChoice(asset.Name, uint(asset.Weight))
									choices = append(choices, c)
								}
							}
						}
					}
					chs, _ := wr.NewChooser(choices...)
					chosenMeta := chs.Pick()
					chosenTypeString := strings.Replace(fmt.Sprintf("%v", chosenMeta), "_", " ", -1)
					conficts := false
					for _, metaD := range meta.Attributes {
						for _, conflictD := range allConflicts {
							if metaD.Value == conflictD.Name {
								if attribute == conflictD.Type {
									conficts = true
								}
							}
						}
					}
					if chosenTypeString != "N/A" && !conficts {
						metaAttr := MetadataAttribute{
							TraitType: attribute,
							Value:     strings.Replace(fmt.Sprintf("%v", chosenMeta), "_", " ", -1),
						}
						meta.Attributes = append(meta.Attributes, metaAttr)
					}
					if strings.ToLower(attribute) == "body" {
						body = chosenTypeString
					}
				}
			}
			//sort further conflicts
			hasOutfit := false
			outfitName := ""
			for _, attr := range meta.Attributes {
				if attr.Value == "ASTRONAUT" || attr.Value == "CYCLOPS" {
					hasOutfit = true
					outfitName = attr.Value
				}
				if attr.Value == "FIRE" {
					for xxx, attr := range meta.Attributes {
						if attr.TraitType == "GLASSES" {
							meta.Attributes[xxx].Value = "NONE"
						}
					}
				}

				if strings.Contains(attr.Value, "SILLY") {
					for xxx, attr2 := range meta.Attributes {
						conflicts := []string{
							"DEFAULT",
							"FIRE",
							"GRIN",
							"WORRIED",
						}
						if contains(conflicts, attr2.Value) {
							choices := []string{"DOOT", "SMILE"}
							meta.Attributes[xxx].Value = choices[rand.Intn(len(choices))]
						}
					}
				}
				if attr.Value == "VSHAPE" {
					for xxx, attr := range meta.Attributes {
						if attr.TraitType == "HATS" {
							meta.Attributes[xxx].Value = "NONE"
						}
					}
				}
			}
			//outfit fixups
			if hasOutfit {
				if outfitName == "ASTRONAUT" {
					for xxx, attr := range meta.Attributes {
						if attr.TraitType == "GLASSES" {
							meta.Attributes[xxx].Value = "NONE"
						}
						if attr.TraitType == "ACCESSORY" {
							meta.Attributes[xxx].Value = "NONE"
						}
						if attr.TraitType == "HATS" {
							meta.Attributes[xxx].Value = "NONE"
						}
						if attr.TraitType == "BACKGROUND" {
							meta.Attributes[xxx].Value = "SPACE ASTRONAUT"
						}
					}
				}
				if outfitName == "CYCLOPS" {
					for xxx, attr := range meta.Attributes {
						if attr.TraitType == "GLASSES" {
							meta.Attributes[xxx].Value = "NONE"
						}
						if attr.TraitType == "HATS" {
							meta.Attributes[xxx].Value = "NONE"
						}
					}
				}
			}
			if config.Duplicates {
				parsed, _ := json.Marshal(&meta)
				f, err := os.OpenFile("./output/metadata/"+fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
				if err != nil {
					log.Fatal(err)
				}
				if _, err := f.WriteString(string(parsed)); err != nil {
					log.Fatal(err)
				}
				if err := f.Close(); err != nil {
					log.Fatal(err)
				}
				generatedMetatata = append(generatedMetatata, meta)
				complete = true
			} else {
				//re-roll if duplicate
				metaString, _ := json.Marshal(meta.Attributes)
				if !contains(generatedMetadataStrings, string(metaString)) {
					generatedMetatata = append(generatedMetatata, meta)
					generatedMetadataStrings = append(generatedMetadataStrings, string(metaString))
					complete = true
				} else {
					fmt.Println("Duplicate Prevented!")
				}
			}
		}
	}
	var pickedLegendaries []int
	for _, lgnd := range legendaryNames {
		done := false
		for !done {
			min := config.LegendaryBlock + 1
			max := config.EndIndex - config.StartIndex
			randness := rand.Intn(max-min) + min
			if !containsInt(pickedLegendaries, randness) {
				pickedLegendaries = append(pickedLegendaries, randness)
				var md []MetadataAttribute
				md = append(md, MetadataAttribute{TraitType: "Legendary", Value: lgnd})
				generatedMetatata[randness] = Metatata{Attributes: md}
				done = true
			}
		}
	}
	fmt.Println("Legendaries Replaced!")
	rand.Shuffle(len(generatedMetatata), func(i, j int) {
		generatedMetatata[i], generatedMetatata[j] = generatedMetatata[j], generatedMetatata[i]
	})

	//we added shuffling of data due to an error with the hashToDecimal function (which was identified to be the cause after reveal) which resulted in the int64 being maxxed out and producing the same data as the testnet hash.
	//whilst shuffling wasn't the most optimal fix, it prevented the ability to work out which rares are which during the pre-reveal stage based off our testnet run.
	//a fixed up version of this will be posted down the line alongside our image generation scripts :) <3

	fmt.Println("Shuffled!")
	startCount := config.StartIndex
	for _, meta := range generatedMetatata {
		meta.Name = strings.Replace(config.AssetName, "ID", strconv.Itoa(startCount), -1)
		meta.Image = strings.Replace(config.ImagePath, "ID", strconv.Itoa(startCount), -1)
		parsed, _ := json.Marshal(&meta)
		f, err := os.OpenFile(filepath.FromSlash("./output/metadata/"+strconv.Itoa(startCount)), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := f.WriteString(string(parsed)); err != nil {
			log.Fatal(err)
		}
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
		startCount = startCount + 1
	}
	fmt.Println("Done!")
}
