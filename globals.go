package main

import "os"

type Asset struct {
	Type   string
	Name   string
	Weight int
}

type Config struct {
	BlockHash      string `json:"block_hash"`
	ImagePath      string `json:"image_path"`
	AssetName      string `json:"asset_name"`
	Description    string `json:"description"`
	LegendaryBlock int    `json:"legendary_block"`
	StartIndex     int    `json:"start_index"`
	EndIndex       int    `json:"end_index"`
	Duplicates     bool   `json:"duplicates"`
}

type Metatata struct {
	Image       string              `json:"image"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Attributes  []MetadataAttribute `json:"attributes"`
}

type MetadataAttribute struct {
	TraitType string `json:"trait_type"`
	Value     string `json:"value"`
}

type Conflict struct {
	Name string
	Type string
}

func makeDirectoryIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, os.ModeDir|0755)
	}
	return nil
}

func Exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func contains(s []string, searchterm string) bool {
	for _, value := range s {
		if value == searchterm {
			return true
		}
	}
	return false
}

func containsInt(s []int, searchterm int) bool {
	for _, value := range s {
		if value == searchterm {
			return true
		}
	}
	return false
}
