package main

import (
	"fmt"
	"math/bits"
)

// Analyse traite les donnees d'un capteur specifique et retourne un tableau de compteurs pour les bits actifs
// dans les bits 31 a 8 des entrees valides, sachant que les entrees valides sont celles qui ont le 7e bit a 0
// et dont les bits 0 a 6 correspondent a l'identifiant du capteur specifique
func Analyse(data []uint32, capteur uint8) ([24]int, error) {
	var result [24]int
	//verifie que les donnees ne sont pas vides
	if len(data) == 0 {
		return [24]int{}, fmt.Errorf("Aucune donnee a analyser")
	}
	//verifie que la valeur capteur est valide
	if capteur > 127 {
		return [24]int{}, fmt.Errorf("Capteur invalide, doit etre entre 0 et 127")
	}
	for index, value := range data {
		// verifie l'identifiant du capteur dans les bits 0 a 6
		if value&0x7F != uint32(capteur) {
			continue //identifiant du capteur ne correspond pas, on passe a l'entree suivante
		}
		//verifie que le 7e bit est a 0
		if value&(1<<7) != 0 {
			return [24]int{}, fmt.Errorf("Bit de validation invalide, le 7e bit est a 1 dans l'entree a l'index %d", index)
		}
		//note les positions des bits activees dans les bits 31 a 8
		index := Indexes(value >> 8)
		//verifie que juste un bit est active dans les bits 31 a 8
		if len(index) > 1 {
			return [24]int{}, fmt.Errorf("Plus d'un bit actif dans les bits 31 a 8 dans l'entree a l'index %d", index)
		} else if len(index) == 0 {
			continue //si aucun bit active, on passe a l'entree suivante, pourrait declancher une erreur selon les besoins
		}
		//ajoute un a la position correspondante dans le resultat
		result[index[0]]++
	}
	return result, nil
}

// tel que disponnible dans le manuel du cours p108, modifier de 64 bits a 32 bits pour correspondre au type de données utilisé dans ce projet
func Indexes(x uint32) []int {
	var ind = make([]int, bits.OnesCount32(x))
	pos := 0
	for x != 0 {
		ind[pos] = bits.TrailingZeros32(x)
		x &= x - 1
		pos++
	}
	return ind
}

func main() {
	data := []uint32{
		0x00000105, // ID=5, bit 7=0, bit 8=1 (valide)
		0x00000205, // ID=5, bit 7=0, bit 9=1 (valide)
		0x00000080, // ID=0, bit 7=1 (invalide)
		0x00000305, // ID=5, bit 7=0, bits 8 et 9=1 (invalide)
	}
	counts, err := Analyse(data, 5)
	if err != nil {
		// Résultat attendu : err != nil (à cause des entrées invalides)
		fmt.Println("Erreur lors de l'analyse:", err)
	} else {
		// Si aucune erreur : counts = [1, 1, 0, 0, ..., 0] (bit 8: 1 fois, bit 9: 1 fois)
		fmt.Println("Analyse réussie, counts:", counts)
	}
}
