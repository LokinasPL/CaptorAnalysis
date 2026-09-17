Auteur : Lokinas Peter Lemire

Date : Mars 2026

Tâche demandé : Vous travaillez pour une entreprise qui surveille des capteurs déployés dans une usine intelligente. Chaque capteur envoie des données sous forme d’entiers 32 bits, où :

    Bits 0 à 6 (7 bits les moins significatifs) : Identifiant du capteur (valeur entre 0 et 127).
    Bit 7 : Bit de validation (doit être 0, sinon il y a une erreur).
    Bits 8 à 31 (24 bits les plus significatifs) : Valeur du capteur. Exactement un de ces bits doit être à 1 (indiquant une mesure spécifique), ou tous peuvent être à 0 (indiquant aucune mesure).

Votre tâche est d’écrire un programme Go avec une fonction Analyse(data []uint32, capteur uint8) ([24]int, error) qui analyse ces données pour un capteur donné, compte les occurrences de chaque mesure, et valide les données.
