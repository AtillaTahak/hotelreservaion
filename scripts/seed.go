package main

import "go/types"

func main() {

	hotel := types.Hotel{
		Name: "Hotel California",
		Location: "California",
	}
	room := types.Room{
		Type: types.Single,
		BasePrice: 100,
	}

}