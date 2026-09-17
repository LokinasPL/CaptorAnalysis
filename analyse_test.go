package main

import (
	"testing"
)

func TestAnalyse(t *testing.T) {
	tests := []struct {
		name     string
		data     []uint32
		capteur  uint8
		expected [24]int
		wantErr  bool
	}{
		{
			name:     "une entree valide, bit 8 active",
			data:     []uint32{0x00000105}, // bit 8 active
			capteur:  5,
			expected: [24]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr:  false,
		},
		{
			name:     "une entree valide, bit 31 active",
			data:     []uint32{0x80000005}, // bit 31 active
			capteur:  5,
			expected: [24]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			wantErr:  false,
		},
		{
			name:     "une entree valide, bit 16 active",
			data:     []uint32{0x00010005}, // bit 16 active
			capteur:  5,
			expected: [24]int{0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr:  false,
		},
		{
			name:     "plusieurs entrees valides, bits differents actives",
			data:     []uint32{0x00000105, 0x00010005, 0x80000005}, // bits 8, 16, 31 actives
			capteur:  5,
			expected: [24]int{1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			wantErr:  false,
		},
		{
			name:     "plusieurs entrees valides, memes bits actives",
			data:     []uint32{0x00010005, 0x00010005, 0x00010005}, // bit 16 active 3 fois
			capteur:  5,
			expected: [24]int{0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr:  false,
		},
		{
			name:     "plusieurs entrees valides, memes bits et bits differents actives",
			data:     []uint32{0x00010005, 0x00010005, 0x80000005}, // bit 16 active 2 fois, bit 31 active 1 fois
			capteur:  5,
			expected: [24]int{0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			wantErr:  false,
		},
		{
			name:     "aucun bit activé mais entrées valides",
			data:     []uint32{0x00000005}, // aucun bit active
			capteur:  5,
			expected: [24]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr:  false,
		},
		{
			name:     "aucun bit activé mais entrées valides, parmis plusieurs entrees valides",
			data:     []uint32{0x00000005, 0x00010005, 0x80000005}, // aucun bit active dasn premier, bits 16 et 31 actives
			capteur:  5,
			expected: [24]int{0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			wantErr:  false,
		},
		{
			name:     "valeur de capteur non-presente dans la donnee",
			data:     []uint32{0x00000105}, // bits 8 actives
			capteur:  1,                    // capteur pas utilise dans les donnees
			expected: [24]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr:  false,
		},
		{
			name:     "valeur de capteur non-presente dans les donnees",
			data:     []uint32{0x00000105, 0x00010005, 0x80000005}, // bits 8, 16, 31 actives
			capteur:  1,                                            // capteur pas utilise dans les donnees
			expected: [24]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr:  false,
		},
		{
			name:     "valeur de capteur presente dans seulement deux des quatre donnees",
			data:     []uint32{0x00010001, 0x00000105, 0x00020001, 0x80000005}, // bits 16, 8, 17, 31 actives
			capteur:  1,                                                        // capteur utilise dans la 1e et 3e donnée
			expected: [24]int{0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr:  false,
		},
		{
			name:     "trop de bits actives dans une entree",
			data:     []uint32{0x80010005}, // bit 16 et 31 active
			capteur:  5,
			expected: [24]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr:  true,
		},
		{
			name:     "trops de bits actives dans une entree parmis plusieurs entrees valides",
			data:     []uint32{0x00010005, 0x00010005, 0x80010005}, // bit 16 active 2 fois, bit 31 active 1 fois
			capteur:  5,
			expected: [24]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr:  true,
		},
		{
			name:     "trops de bits actives dans plusieurs entrees parmis plusieurs entrees valides",
			data:     []uint32{0x00010005, 0x00010005, 0x00010105, 0x80010005}, // bit 16 active 2 fois, puis 8 et 16 dans une entree, et 31 et 16 dans une autre entree
			capteur:  5,
			expected: [24]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr:  true,
		},
		{
			name:     "un bit de validation invalide pour une entree",
			data:     []uint32{0x80000085}, // bit 7 active
			capteur:  5,
			expected: [24]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr:  true,
		},
		{
			name:     "un bit de validation invalide pour plusieurs entree",
			data:     []uint32{0x00010005, 0x00010085, 0x80000005}, // bit 7 active dans deuxieme entree
			capteur:  5,
			expected: [24]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr:  true,
		},
		{
			name:     "plusieurs bits de validation invalides pour plusieurs entree",
			data:     []uint32{0x00010085, 0x00010085, 0x80000005, 0x00000105}, // bit 7 active dans premiere et deuxieme entree
			capteur:  5,
			expected: [24]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr:  true,
		},
		{
			name:     "aucune donnees",
			data:     []uint32{},
			capteur:  5,
			expected: [24]int{},
			wantErr:  true,
		},
		{
			name:     "valeur de capteur invalide (un peu plus grand que 127)",
			data:     []uint32{0x00000105, 0x00010005, 0x80000005}, // bits 8, 16, 31 actives
			capteur:  128,                                          // capteur invalide
			expected: [24]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr:  true,
		},
		{
			name:     "valeur de capteur invalide (beaucoup plus grand que 127)",
			data:     []uint32{0x00000105, 0x00010005, 0x80000005}, // bits 8, 16, 31 actives
			capteur:  250,                                          // capteur invalide
			expected: [24]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Analyse(tt.data, tt.capteur)
			if (err != nil) != tt.wantErr {
				t.Errorf("Analyse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("Analyse() = %v, want %v", result, tt.expected)
			}
		})
	}
}
