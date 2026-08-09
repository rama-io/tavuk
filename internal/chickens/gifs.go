package chickens

import "math/rand"

var List = []string{
	"https://media.tenor.com/EMQV-OcGUIQAAAAC/cody.gif",
	"https://media.tenor.com/poQ7zWY_rOcAAAAC/chicken-werk-chicken-dance.gif",
	"https://media.tenor.com/sIaCtfXMB1kAAAAC/fried-chicken-with-shoes.gif",
	"https://media.tenor.com/zme4VU2jjpQAAAAC/goobin-novoura.gif",
	"https://media.tenor.com/El2t8T9WOq4AAAAC/poker.gif",
}

func Random() string {
	return List[rand.Intn(len(List))]
}
