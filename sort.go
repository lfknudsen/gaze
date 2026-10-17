package main

import "strings"

func sortByName(a *VideoInfo, b *VideoInfo) int {
	return strings.Compare(a.Name, b.Name)
}

func sortByWidth(a *VideoInfo, b *VideoInfo) int {
	if a.Width > b.Width {
		return 1
	}
	if a.Width < b.Width {
		return -1
	}
	return 0
}

func sortByHeight(a *VideoInfo, b *VideoInfo) int {
	if a.Height > b.Height {
		return 1
	}
	if a.Height < b.Height {
		return -1
	}
	return 0
}

func sortByResolution(a *VideoInfo, b *VideoInfo) int {
	if a.Width*a.Height > b.Width*b.Height {
		return 1
	}
	if a.Width*a.Height < b.Width*b.Height {
		return -1
	}
	return 0
}

func sortByBitrate(a *VideoInfo, b *VideoInfo) int {
	if a.Bitrate > b.Bitrate {
		return 1
	}
	if a.Bitrate < b.Bitrate {
		return -1
	}
	return 0
}

func sortByFramerate(a *VideoInfo, b *VideoInfo) int {
	if a.Framerate > b.Framerate {
		return 1
	}
	if a.Framerate < b.Framerate {
		return -1
	}
	return 0
}

func sortByNameDesc(a *VideoInfo, b *VideoInfo) int {
	return strings.Compare(b.Name, a.Name)
}

func sortByWidthDesc(a *VideoInfo, b *VideoInfo) int {
	return sortByWidth(b, a)
}

func sortByHeightDesc(a *VideoInfo, b *VideoInfo) int {
	return sortByHeight(b, a)
}

func sortByResolutionDesc(a *VideoInfo, b *VideoInfo) int {
	return sortByResolution(b, a)
}

func sortByBitrateDesc(a *VideoInfo, b *VideoInfo) int {
	return sortByBitrate(b, a)
}

func sortByFramerateDesc(a *VideoInfo, b *VideoInfo) int {
	return sortByFramerate(b, a)
}
