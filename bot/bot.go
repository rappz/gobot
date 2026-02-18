// Package bot implements a Discord bot with slash commands and message handlers.
package bot

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
)

// --- Constants and global state ---

const CHRIS_ID = "245751144209448961"     // Chris's Discord user ID
var crissyMode = false                    // When true, Chris gets timed out if he @-mentions someone
var BotToken string                       // Discord bot token, set from main
var punsihedUsers = make(map[string]bool) // User IDs currently "punished" (in-memory only)
var defaultMemberPermissions = int64(discordgo.PermissionAdministrator)

// PUNISH_MESSAGES are random messages sent when a punished user speaks; %s is replaced with their mention.
var PUNISH_MESSAGES = []string{
	"%sYou are punished!😡\n Bad kitten!🤪",
	"%s Daddy is mad now 😤\nNaughty little kitten 🐾",
	"%s Uh-oh… daddy saw that 🙃\nBad baby cat 🐱💢",
	"%s Love muffin misbehaved again 🧁\nDaddy's disappointed 😑",
	"%s Tiny rat energy detected 🐀⚡\nDaddy is not amused 😠",
	"%s Who scratched the couch?? 😡\nWas it my chaotic kitten? 🐈‍⬛",
	"%s No treats for spicy kitty 😼🚫\nDaddy says behave.",
	"%s You hiss at daddy?? 😾\nBold move, little gremlin.",
	"%s That's it. Jail for kitten. 🚔🐱\nDaddy is furious.",
	"%s You adorable menace 😤💘\nWhy is daddy always stressed.",
	"%s Rat behavior. Absolute rat behavior. 🐀\nDaddy is shaking his head.",
	"%s Love muffin turned into chaos muffin 🧁🔥\nDaddy needs a minute.",
	"%s Tiny paws, big crimes 🐾\nDaddy witnessed everything 👀",
	"%s Don't blink at me like that 😒\nYou know daddy is mad.",
	"%s Sweet baby kitten by day 😇\nCertified rat by night 🐀🌙",
	"%s Who knocked over the water?? 💦\nConfess, fuzzy criminal.",
	"%s Daddy gave you one job 😐\nYou chose violence, kitten.",
	"%s Stop being cute while guilty 😤💕\nIt's manipulative.",
	"%s You bit daddy?? 😡\nThat's betrayal, little fang gremlin 🐱🩸",
	"%s Suspicious whiskers detected 🕵️‍♂️\nDaddy knows.",
	"%s Love muffin revoked. Now just muffin. 😑🧁",
	"%s Tiny toe beans, massive audacity 🐾\nDaddy is stunned.",
	"%s Why are you staring like that 👁️\nYou absolutely did something.",
	"%s Menace in a fur coat 😼\nDaddy demands order.",
	"%s One more zoomie and it's over 😤💨\nDaddy said calm.",
	"%s Come here, you chaotic rat kitten 🐀🐱\nDaddy is mad… but also holding you anyway.",
}

var KING_IMAGE = []string{
	"https://media1.tenor.com/m/zzh5EGMb8KcAAAAd/yes-king.gif",
	"https://media1.tenor.com/m/wzBvSvmdhdMAAAAd/yes-king-yes.gif",
	"https://media1.tenor.com/m/1exE1H-iGGsAAAAd/martene3-yesking.gif",
	"https://media1.tenor.com/m/psMStUrhCp4AAAAd/burger-king-yes-sir.gif",
	"https://media2.giphy.com/media/v1.Y2lkPTc5MGI3NjExbm1zNGNwbTFtbjFvbXlmMzByY3p3azhmNzN3cDN5d2FhMjRrcTF0eSZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/L0nhTaYf038ZqDMZpY/giphy.gif",
}

func getKingImage() (*bytes.Reader, error) {
	resp, err := http.Get(KING_IMAGE[rand.Intn(len(KING_IMAGE))])
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return bytes.NewReader(body), err
}

// --- Slash command definitions and handlers ---

var (
	RemoveCommands = flag.Bool("rm-cmd", false, "Remove commands after execution")
	//dmPermission   = false

	// commands are registered with Discord as application (slash) commands.
	commands = []*discordgo.ApplicationCommand{
		{
			Name:                     "punish",
			Description:              "Punish a bad kitten",
			DefaultMemberPermissions: &defaultMemberPermissions,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "the-kitten",
					Description: "The user to absolve of sin",
					Required:    true,
				},
			},
		},
		{
			Name:                     "absolve",
			Description:              "Absolve a kitten of sin",
			DefaultMemberPermissions: &defaultMemberPermissions,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "the-kitten",
					Description: "The user to punish for sin",
					Required:    true,
				},
			},
		},
		{
			Name:                     "crissy",
			Description:              "Toggle Crissy mode",
			DefaultMemberPermissions: &defaultMemberPermissions,
		},
	}

	// commandHandlers dispatch slash command interactions; key is command name.
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		// punish: mark a user as "punished"; only in allowed guild.
		"punish": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if i.GuildID == "698693159261306921" {
				userID := i.ApplicationCommandData().Options[0].UserValue(nil).ID
				member, err := s.GuildMember(i.GuildID, userID)
				if err != nil {
					log.Printf("Error getting guild member: %v", err)
					return
				}
				if member.User.Bot {
					err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseType(discordgo.InteractionResponseChannelMessageWithSource),
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("Only god can punish %s", member.Mention()),
						},
					})
					if err != nil {
						log.Printf("Error responding to interaction: %v", err)
					}
					return
				}
				content := ""
				switch punsihedUsers[userID] {
				case true:
					content = fmt.Sprintf("%s already punished", member.Mention())
					//s.ChannelMessageSend(i.ChannelID, fmt.Sprintf("%s already punished", member.Mention()))
				default:
					content = fmt.Sprintf("Punishing user %s", member.Mention())
					//s.ChannelMessageSend(i.ChannelID, fmt.Sprintf("Punishing user %s", member.Mention()))
					punsihedUsers[userID] = true
				}
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseType(discordgo.InteractionResponseChannelMessageWithSource),
					Data: &discordgo.InteractionResponseData{
						Content: content,
					},
				})
				if err != nil {
					log.Printf("Error responding to interaction: %v", err)
					s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
						Content: "Something went wrong",
					})
					return
				}
				time.AfterFunc(time.Second*5, func() {
					// delete the message after 5 seconds
					s.InteractionResponseDelete(i.Interaction)
				})
			} else {
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseType(discordgo.InteractionResponseChannelMessageWithSource),
					Data: &discordgo.InteractionResponseData{
						Content: "You are not allowed to use this command in this server",
					},
				})
				if err != nil {
					s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
						Content: "Something went wrong",
					})
					return
				}
				time.AfterFunc(time.Second*5, func() {
					// delete the message after 5 seconds
					s.InteractionResponseDelete(i.Interaction)
				})
			}
		},
		// absolve: remove a user from the punished set; only in allowed guild.
		"absolve": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if i.GuildID == "698693159261306921" {
				userID := i.ApplicationCommandData().Options[0].UserValue(nil).ID
				member, err := s.GuildMember(i.GuildID, userID)
				if err != nil {
					log.Printf("Error getting guild member: %v", err)
					return
				}
				if member.User.Bot {
					err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseType(discordgo.InteractionResponseChannelMessageWithSource),
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("Only god can absolve %s", member.Mention()),
						},
					})
					if err != nil {
						log.Printf("Error responding to interaction: %v", err)
					}
					return
				}
				content := ""
				switch punsihedUsers[userID] {
				case true:
					content = fmt.Sprintf("Absolving %s of sin", member.Mention())
					delete(punsihedUsers, userID)
				default:
					content = fmt.Sprintf("%s already free of sin", member.Mention())

				}
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseType(discordgo.InteractionResponseChannelMessageWithSource),
					Data: &discordgo.InteractionResponseData{
						Content: content,
					},
				})
				if err != nil {
					log.Printf("Error responding to interaction: %v", err)
					s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
						Content: "Something went wrong",
					})
					return
				}
				//time.AfterFunc(time.Second*5, func() {
				// delete the message after 5 seconds
				//	s.InteractionResponseDelete(i.Interaction)
				//})

			} else {
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseType(discordgo.InteractionResponseChannelMessageWithSource),
					Data: &discordgo.InteractionResponseData{
						Content: "You are not allowed to use this command in this server",
					},
				})
				if err != nil {
					s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
						Content: "Something went wrong",
					})
					return
				}
				time.AfterFunc(time.Second*5, func() {
					// delete the message after 5 seconds
					s.InteractionResponseDelete(i.Interaction)
				})
			}
		},
		// crissy: toggle Crissy mode (timeout Chris on @-mention); only in allowed guild.
		"crissy": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if i.GuildID == "698693159261306921" {
				member, err := s.GuildMember(i.GuildID, CHRIS_ID)
				if err != nil {
					log.Printf("Error getting guild member: %v", err)
					return
				}
				crissyContent := ""
				switch crissyMode {
				case true:
					crissyContent = fmt.Sprintf("%s is now less of a sinner💦", member.Mention())
					s.UpdateCustomStatus("Crissy Punisher is now off duty 💦")
					crissyMode = false

				case false:
					crissyContent = fmt.Sprintf("%s will now get timed out for 1 minute if he @'s💦", member.Mention())
					s.UpdateCustomStatus("Crissy Punisher is now set to punish 💦💦💦")
					crissyMode = true
				}
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseType(discordgo.InteractionResponseChannelMessageWithSource),
					Data: &discordgo.InteractionResponseData{
						Content: crissyContent,
					},
				})
				if err != nil {
					log.Printf("Error responding to interaction: %v", err)
					return
				}
				time.AfterFunc(time.Second*5, func() {
					// delete the message after 5 seconds
					s.InteractionResponseDelete(i.Interaction)
				})
			} else {
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseType(discordgo.InteractionResponseChannelMessageWithSource),
					Data: &discordgo.InteractionResponseData{
						Content: "You are not allowed to use this command in this server",
					},
				})
				if err != nil {
					s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
						Content: "Something went wrong",
					})
					return
				}
				time.AfterFunc(time.Second*5, func() {
					// delete the message after 5 seconds
					s.InteractionResponseDelete(i.Interaction)
				})
			}
		},
	}
)

// --- Helpers ---

// getRandomPunishMessage returns a random entry from PUNISH_MESSAGES.
func getRandomPunishMessage() string {
	return PUNISH_MESSAGES[rand.Intn(len(PUNISH_MESSAGES))]
}

// checkNilErr logs the error if non-nil (no panic).
func checkNilErr(e error) {
	if e != nil {
		log.Print("Error: ", e)
	}
}

// --- Main bot lifecycle ---

// Run starts the Discord session, registers slash commands, and blocks until OS interrupt (e.g. Ctrl+C).
// If -rm-cmd is set, slash commands are removed on shutdown.
func Run() {
	flag.Parse()

	discord, err := discordgo.New("Bot " + BotToken)
	checkNilErr(err)

	discord.AddHandler(newMessage)
	discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	err = discord.Open()
	if err != nil {
		log.Printf("Cannot open the session: %v", err)
		return
	}
	discord.UpdateCustomStatus("Daddy I'm Coming!")

	// Register each slash command globally (empty guild ID = global).
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := discord.ApplicationCommandCreate(discord.State.User.ID, "", v)
		if err != nil {
			log.Printf("Cannot create '%v' command: %v", v.Name, err)
			continue
		}
		registeredCommands[i] = cmd
	}

	defer discord.Close()

	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	if *RemoveCommands {
		log.Println("Removing commands...")
		for _, v := range registeredCommands {
			if v == nil {
				continue
			}
			err := discord.ApplicationCommandDelete(discord.State.User.ID, "", v.ID)
			if err != nil {
				log.Printf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}
}

// --- Message handler ---
// newMessage handles every message: applies Crissy mode (timeout Chris on @-mention) and
// punishes already-marked users by deleting their message and posting a random punish message.
func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	if message == nil || message.Author == nil {
		return
	}
	if message.Type == discordgo.MessageTypeReply {
		return
	}
	if message.GuildID == "" {
		return
	}
	if message.Author.ID == discord.State.User.ID {
		return
	}
	if message.Author.Bot {
		return
	}
	member, err := discord.GuildMember(message.GuildID, message.Author.ID)
	if err != nil {
		log.Printf("Error getting guild member: %v", err)
		return
	}
	if member == nil {
		return
	}

	// Crissy mode: if Chris sends a message with mentions, timeout him 1 minute and reply.
	if crissyMode {
		if message.Author.ID == CHRIS_ID && len(message.Mentions) != 0 {
			currentTime := time.Now()
			oneMinute := currentTime.Add(time.Minute * 1)
			err = discord.GuildMemberTimeout(message.GuildID, message.Author.ID, &oneMinute)
			if err != nil {
				log.Printf("Error setting timeout: %v", err)
				return
			}
			userDM, err := discord.UserChannelCreate(message.Author.ID)
			if err != nil {
				log.Printf("Error creating user channel: %v", err)
				return
			}
			kingImg, err := getKingImage()
			if err != nil {
				log.Printf("Error getting king image: %v", err)
				discord.ChannelMessageSend(userDM.ID, fmt.Sprintf("%s get fucked pussy fart! ", member.Mention()))
			}
			punishMsg := &discordgo.MessageSend{
				Content: fmt.Sprintf(getRandomPunishMessage(), member.Mention()),
				File:    &discordgo.File{Name: "king.gif", Reader: kingImg},
			}
			discord.ChannelMessageSendComplex(userDM.ID, punishMsg)
			discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%s get fucked pussy fart! ", member.Mention()))
		}

	}

	// If author is punished: send random punish message, delete their message, then delete punish message after 5s.
	if punsihedUsers[message.Author.ID] {
		kingImg, err := getKingImage()
		if err != nil {
			log.Printf("Error getting king image: %v", err)
		} else {
			punishMsg := &discordgo.MessageSend{
				Content: fmt.Sprintf(getRandomPunishMessage(), member.Mention()),
				File:    &discordgo.File{Name: "king.gif", Reader: kingImg},
			}
			_, err := discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf(getRandomPunishMessage(), member.Mention()))
			checkNilErr(err)
			userDM, err := discord.UserChannelCreate(message.Author.ID)
			checkNilErr(err)
			discord.ChannelMessageSendComplex(userDM.ID, punishMsg)

			err = discord.ChannelMessageDelete(message.ChannelID, message.ID)
			if err != nil {
				log.Printf("Failed to delete punished user message: channel=%s message=%s author=%s: %v", message.ChannelID, message.ID, message.Author.ID, err)
			} else {
				log.Printf("Deleted punished user message: channel=%s message=%s author=%s", message.ChannelID, message.ID, message.Author.ID)
			}
			/*
				time.Sleep(5 * time.Second)
				err = discord.ChannelMessageDelete(message.ChannelID, pun.ID)
				if err != nil {
					log.Printf("Failed to delete punishment reply: channel=%s message=%s: %v", message.ChannelID, pun.ID, err)
				} else {
					log.Printf("Deleted punishment reply: channel=%s message=%s", message.ChannelID, pun.ID)
				}*/
		}
	}
}
