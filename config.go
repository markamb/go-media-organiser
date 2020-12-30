package main

var defaultConfig = Config{
	SourceConfig: []Source{
		{
			Source: "D:/Home/Mark/OneDrive/Pictures/Camera Roll",
			Destinations: []string{
				"D:/Drive/Photo Uploads",
				"D:/Home/Mark/OneDrive/1. Shared [AM]/Photo Album",
			},
		},
		{
			Source: "D:/Home/Mark/OneDrive/z Alix Camera Roll",
			Destinations: []string{
				"D:/Home/Mark/OneDrive/1. Shared [AM]/Photo Album",
			},
		},
	},
	ArchiveDirectory: "D:/Backups/Photos",
	LogToFile:        true,
}
