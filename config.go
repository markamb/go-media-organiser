package main

var defaultConfig = Config{
	SourceConfig: []Source{
		{
			Source: "D:/Dropbox/Camera Uploads",
			Destinations: []string{
				"D:/Dropbox/Shared/Pictures/Album",
			},
		},
		{
			Source: "D:/Home/Mark/OneDrive/Pictures/Camera Roll",
			Destinations: []string{
				"D:/Drive/Photo Uploads",
				"D:/Home/Mark/OneDrive/1. Shared [AM]/Photo Album",
			},
		},
		{
			Source: "D:/Home/Mark/OneDrive/zzz Alix Photo Uploads",
			Destinations: []string{
				"D:/Home/Mark/OneDrive/1. Shared [AM]/Photo Album",
			},
		},
	},
	ArchiveDirectory: "D:/Backups/Photos",
}
