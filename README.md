# blackout-bot
[Читати Українською](README.uk.md)

Bot for notification of actual power outages and notification of planned outages

## Features:
- Send a message about actual outages
- Send messages about scheduled outages
- Displays the amount of time electricity has been present or absent
- The ability to disable the bot
- Ability to disable scheduled messages
- The ability to update real-time notifications about the status of electricity

### Admin Commands:
| Command          | Description                  |
|------------------|------------------------------|
| `turn_bot`       | Turn on/off the Bot          |
| `turn_emergency` | Turn on/off planned messages |

### Commands:
| Command          | Description                                     |
|------------------|-------------------------------------------------|
| `get_id`         | Return current chat id                          |
| `get_message_id` | Return messageID for reply message (only chats) |
| `curr_status`    | Return current online status                    |

### How Bot Works:
1. You configure the config
2. You add the bot to the channel or chat and give it the necessary rights
3. If there is no response within 15 seconds from the set endpoint, the bot will send information that the electricity has been turned off
4. If there were no shutdowns during the last 40 minutes, the bot will send a shutdown message 30 minutes before the scheduled shutdown
5. Also, in one message, if you specify which bot will update information in real time

## Requirements:
- Docker `(or golang)`
- Docker-Compose `(or golang)`
- Makefile

## Installation:
1. Clone this repository `git clone https://github.com/sergikoio/blackout-bot.git`
2. Create .env file `cp  .env.example .env`
3. Edit `.env` file
4. Run Bot:
    - With Docker: `make build`
    - Local with Go:
        - `go mod download`
        - `go run main.go`

## FAQ
- Where to get the endpoint?
  - You should start a server at home that will respond to any request. Or use your router (see something like: router in global). But for this you must have either a permanent IP address or DDNS on the router
- The shutdown schedule has changed, where can I get a new one?
   - Go to https://yasno.com.ua/schedule-turn-off-electricity-kyiv, check Page Source Code, `CTRL + F` -> `group_1`. Copy this line and reformat JS to JSON
- For which city is the schedule of scheduled outages made?
   - Kyiv
- Where can I launch the bot?
   - VPS, Qovery, Heroku, etc...

## Contributions:
Any contributions are welcome, also if there are problems in the process, then create an issue