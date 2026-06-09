# iqube

iQube is a puzzle game written in Go and Raylib.

## Description

Le but du jeu est de faire déplacer un pion sur un cube dont chaque face est divisé en une grille de cellule.
Au début de la partie, le pion est sur une cellule de départ et doit avancer jusqu'à une cellule d'arrivée, qui est l'objectif du niveau.
Les cellules peuvent être de plusieurs types : vide ou pleine. Les cellules vides permettent d'être traversée par le pion, contrairement aux cellules pleines.

Le joueur doit alors placer des marqueurs sur les cellules vides qui donnent des instructions au pion pour changer des directions.

Ces marqueurs sont limités pour chaque niveau et peuvent être de plusieurs types :
 - tourner à gauche
 - tourner à droite
 - faire demi-tour

Le pion ne peut avancer que dans un seul sens au départ. 
Dés que le pion passe sur une cellule contenant un marqueur