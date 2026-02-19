package bot

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

func getVoiceClips() ([]string, error) {
	files, err := os.ReadDir("media")
	if err != nil {
		log.Printf("Error reading media directory: %v", err)
		return nil, err
	}
	var names []string
	for _, file := range files {
		if !file.IsDir() {
			names = append(names, file.Name())
		}
	}
	return names, nil

}

func voiceKing(s *discordgo.Session, guildID, channelID string) {
	voice, err := s.ChannelVoiceJoin(guildID, channelID, false, false)
	if err != nil {
		log.Printf("Error joining voice channel: %v", err)
		return
	}
	defer voice.Disconnect()
	voice.Speaking(true)
	defer voice.Speaking(false)
	voice.Speaking(true)
}
