package bot

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
)

const CHRIS_ID = "245751144209448961"     // Chris's ID
var crissyMode = false                    // Crissy mode flag
var BotToken string                       // Token for the bot, set in main.go
var punsihedUsers = make(map[string]bool) // map to keep track of punished users
var (
	RemoveCommands = flag.Bool("rm-cmd", false, "Remove commands after execution")
	//dmPermission   = false
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "punish",
			Description: "Punish a bad kitten",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user-option",
					Description: "User option",
					Required:    true,
				},
			},
		},
		{
			Name:        "absolve",
			Description: "Absolve a kitten of sin",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user-option",
					Description: "User option",
					Required:    true,
				},
			},
		},
		{
			Name:        "crissy",
			Description: "Toggle Crissy mode",
		},
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"punish": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			userID := i.ApplicationCommandData().Options[0].UserValue(nil).ID
			member, err := s.GuildMember(i.GuildID, userID)
			if err != nil {
				log.Fatal(err)
			}
			if member.User.Bot {
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseType(discordgo.InteractionResponseChannelMessageWithSource),
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("Only god can punish %s", member.Mention()),
					},
				})
				if err != nil {
					log.Fatal(err)
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
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "Something went wrong",
				})
				return
			}
			time.AfterFunc(time.Second*5, func() {
				// delete the message after 5 seconds
				s.InteractionResponseDelete(i.Interaction)
			})
		},
		"absolve": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			userID := i.ApplicationCommandData().Options[0].UserValue(nil).ID
			member, err := s.GuildMember(i.GuildID, userID)
			if err != nil {
				log.Fatal(err)
			}

			if member.User.Bot {
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseType(discordgo.InteractionResponseChannelMessageWithSource),
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("Only god can absolve %s", member.Mention()),
					},
				})
				if err != nil {
					log.Fatal(err)
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
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "Something went wrong",
				})
				return
			}
			time.AfterFunc(time.Second*5, func() {
				// delete the message after 5 seconds
				s.InteractionResponseDelete(i.Interaction)
			})
			//s.ChannelMessageSend(i.ChannelID, fmt.Sprintf("Absolving user %s", member.Mention()))
			//delete(punsihedUsers, userID)
		},
		"crissy": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if i.GuildID == "698693159261306921" {
				member, err := s.GuildMember(i.GuildID, CHRIS_ID)
				if err != nil {
					log.Fatal(err)
				}
				crissyContent := ""
				switch crissyMode {
				case true:
					crissyContent = fmt.Sprintf("%s is now less of a sinner💦", member.Mention())
					crissyMode = false

				case false:
					crissyContent = fmt.Sprintf("%s will now get timed out for 1 minute if he @'s💦", member.Mention())
					crissyMode = true
				}
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseType(discordgo.InteractionResponseChannelMessageWithSource),
					Data: &discordgo.InteractionResponseData{
						Content: crissyContent,
					},
				})
				if err != nil {
					log.Fatal(err)
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

func checkNilErr(e error) {
	if e != nil {
		log.Fatal("Error message")
	}
}

func Run() {
	flag.Parse()
	// create a session
	discord, err := discordgo.New("Bot " + BotToken)
	checkNilErr(err)

	// add a event handler
	discord.AddHandler(newMessage)
	discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
	// open session
	err = discord.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	discord.UpdateCustomStatus("Daddy I'm Coming!")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := discord.ApplicationCommandCreate(discord.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer discord.Close() // close session, after function termination

	// keep bot running untill there is NO os interruption (ctrl + C)
	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	if *RemoveCommands {
		log.Println("Removing commands...")
		// // We need to fetch the commands, since deleting requires the command ID.
		// // We are doing this from the returned commands on line 375, because using
		// // this will delete all the commands, which might not be desirable, so we
		// // are deleting only the commands that we added.
		// registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
		// if err != nil {
		// 	log.Fatalf("Could not fetch registered commands: %v", err)
		// }

		for _, v := range registeredCommands {
			err := discord.ApplicationCommandDelete(discord.State.User.ID, "", v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}
}
func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {

	/* prevent bot responding to its own message
	this is achived by looking into the message author id
	if message.author.id is same as bot.author.id then just return
	*/
	if message.Author.ID == discord.State.User.ID {
		return
	}

	if message.Author.Bot {
		return
	}
	member, err := discord.GuildMember(message.GuildID, message.Author.ID)
	if err != nil {
		log.Fatal(err)
	}

	if crissyMode {
		if message.Author.ID == CHRIS_ID && len(message.Mentions) == 0 {
			currentTime := time.Now()
			oneMinute := currentTime.Add(time.Minute * 1)
			// if the message is from Chris, then punish him for 1 minute
			err = discord.GuildMemberTimeout(message.GuildID, message.Author.ID, &oneMinute)
			if err != nil {
				log.Fatal(err)
			}
			discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%s get fucked pussy fart! ", member.Mention()))

		}
	}

	// respond to user message if it contains `!help` or `!bye`
	switch {
	case punsihedUsers[message.Author.ID]:
		pun, err := discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%s \nYou are punished!😡\n Bad kitten!🤪", member.Mention()))
		checkNilErr(err)
		discord.ChannelMessageDelete(message.ChannelID, message.ID) // delete the message
		time.Sleep(5 * time.Second)                                 // wait for 5 seconds
		discord.ChannelMessageDelete(message.ChannelID, pun.ID)     // delete the punishment message
	}

}
