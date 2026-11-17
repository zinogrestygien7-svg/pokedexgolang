package main

import (
	"fmt"
	"net/http"
	"encoding/json"
)
type PokemonList struct{
	Results []struct {
		Nom string `json:"name"`
	} `json:"results"`
}
type Pokemon struct {
	Nom   string `json:"name"`
	ID     int    `json:"id"`
	Taille int    `json:"height"`
	Poids int    `json:"weight"`
	Types  []struct {
		Type struct {
			Nom string `json:"name"`
		} `json:"type"`
	} `json:"types"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Nom string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Sprites struct {
		FrontDefault string `json:"front_default"`
	} `json:"sprites"`
}
type Region struct {
	Nom    string
	Offset int
	Limit  int
}

var regions = []Region{
	{"Kanto", 0, 151},
	{"Johto", 151, 100},
	{"Hoenn", 251, 135},
	{"Sinnoh", 386, 107},
	{"Unys", 493, 156},
	{"Kalos", 649, 72},
	{"Alola", 721, 88},
	{"Galar", 809, 96},
	{"Paldea", 905, 120},
}
func GetPokemonList(r Region) (PokemonList, error) {
	var list PokemonList

	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon?limit=%d&offset=%d", r.Limit, r.Offset)
	resp, err := http.Get(url)
	if err != nil {
		return list, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return list, fmt.Errorf("erreur API: statut %d", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return list, err
	}
	return list, nil
}
func GetPokemon(name string) (Pokemon, error) {
	var p Pokemon

	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + name)
	if err != nil {
		return p, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return p, fmt.Errorf("pokémon %q introuvable (404)", name)
	}
	if resp.StatusCode != http.StatusOK {
		return p, fmt.Errorf("erreur API : statut %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return p, err
	}
	return p, nil
}
func ChoisirRegion() (Region, bool) {
	fmt.Println("\n=== Choisis une région ===")
	for i, r := range regions {
		fmt.Printf("%d. %s\n", i+1, r.Nom)
	}
	fmt.Print("Numéro de la région (0 pour quitter) : ")

	var n int
	if _, err := fmt.Scan(&n); err != nil || n == 0 {
		return Region{}, false // quitter
	}
	if n < 1 || n > len(regions) {
		fmt.Println("Région inconnue, réessaie.")
		return ChoisirRegion()
	}
	return regions[n-1], true
}
func main() {
	for {
		region, ok := ChoisirRegion()
		if !ok {
			fmt.Println("Au revoir !")
			return
		}

		list, err := GetPokemonList(region)
		if err != nil {
			fmt.Println("Erreur :", err)
			continue
		}

		fmt.Printf("\n=== Pokémons de %s ===\n", region.Nom)
		for i, p := range list.Results {
			// numéro national = offset + i + 1
			fmt.Printf("%d. %s\n", region.Offset+i+1, p.Nom)
		}

	choixPokemon:
		for {
			fmt.Print("Numéro du pokémon (0 = changer de région) : ")
			var n int
			if _, err := fmt.Scan(&n); err != nil {
				fmt.Println("Entrée invalide.")
				return
			}
			if n == 0 {
				break choixPokemon // retour au choix de région
			}

			idx := n - region.Offset - 1 // numéro national -> index du tableau
			if idx < 0 || idx >= len(list.Results) {
				fmt.Printf("Ce pokémon n'est pas dans %s.\n", region.Nom)
				continue
			}

			p, err := GetPokemon(list.Results[idx].Nom)
			if err != nil {
				fmt.Println("Erreur :", err)
				continue
			}
			PrintPokemon(p)
		}
	}
}
func PrintPokemon(p Pokemon) {
	fmt.Println("=================================")
	fmt.Printf("Nom        : %s\n", p.Nom)
	fmt.Printf("Identifiant: %d\n", p.ID)
	fmt.Printf("Image      : %s\n", p.Sprites.FrontDefault)
	fmt.Printf("Taille     : %d\n", p.Taille)
	fmt.Printf("Poids      : %d\n", p.Poids)

	fmt.Print("Type       : ")
	for i, t := range p.Types {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(t.Type.Nom)
		
	}
	fmt.Println()

	for _, s := range p.Stats {
		switch s.Stat.Nom {
		case "hp", "attack", "defense":
			fmt.Printf("%-11s: %d\n", s.Stat.Nom, s.BaseStat) 
		}
	}
	fmt.Println("=================================")
}