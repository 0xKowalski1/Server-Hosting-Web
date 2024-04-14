package seeds

import "0xKowalski1/server-hosting-web/models"

func SeedGames() []models.Game {
	var games []models.Game

	// Minecraft
	minecraft := models.Game{
		Name:             "Minecraft",
		ShortDescription: "Craft your world in this groundbreaking sandbox game.",
		GridImage:        "/images/minecraft-grid.jpg",
	}
	games = append(games, minecraft)

	// Rust
	rust := models.Game{
		Name:             "Rust",
		ShortDescription: "Struggle to survive in a harsh environment where the only goal is to outlast others.",
		GridImage:        "/images/rust-grid.jpg",
	}
	games = append(games, rust)

	// ARK: Survival Evolved
	ark := models.Game{
		Name:             "ARK: Survival Evolved",
		ShortDescription: "Survive and tame prehistoric creatures in a mysterious island ecosystem.",
		GridImage:        "/images/ark-grid.jpg",
	}
	games = append(games, ark)

	// Counter Strike 2
	cs2 := models.Game{
		Name:             "Counter Strike 2",
		ShortDescription: "Engage in intense, fast-paced battles in the latest installment of this competitive first-person shooter series.",
		GridImage:        "/images/cs2-grid.jpg",
	}
	games = append(games, cs2)

	// Eco
	eco := models.Game{
		Name:             "Eco",
		ShortDescription: "Build, craft, and create a sustainable civilization within a living ecosystem.",
		GridImage:        "/images/eco-grid.jpg",
	}
	games = append(games, eco)

	// Garry's Mod
	garrysMod := models.Game{
		Name:             "Garry's Mod",
		ShortDescription: "Manipulate objects and experiment with physics in this strange sandbox world.",
		GridImage:        "/images/garrysmod-grid.jpg",
	}
	games = append(games, garrysMod)

	// Battlebit Remastered
	battlebit := models.Game{
		Name:             "Battlebit Remastered",
		ShortDescription: "Dive into large-scale warfare with minimalistic graphics and tactical gameplay.",
		GridImage:        "/images/battlebit-grid.jpg",
	}
	games = append(games, battlebit)

	// Valheim
	valheim := models.Game{
		Name:             "Valheim",
		ShortDescription: "Embark on a mythical adventure in 'Valheim,' a brutal exploration and survival game for warriors fallen in battle.",
		GridImage:        "/images/valheim-grid.jpg",
	}
	games = append(games, valheim)

	return games
}
