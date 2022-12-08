package bot

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"

	"blackout-bot/internal/checker"
	"blackout-bot/internal/db"
	"blackout-bot/internal/schedule"
	"blackout-bot/internal/servertime"
)

type bot struct {
	bot             *tgbotapi.BotAPI
	channelID       int64
	updateMessageID int
	db              *botDB
	sch             *schedule.Schedule
}

func InitBot(token string, channelID int64, updMsgID int, db *db.DB, sch *schedule.Schedule) error {
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	dbState := newBotDB(db)

	botState := &bot{
		bot:             botAPI,
		channelID:       channelID,
		updateMessageID: updMsgID,
		db:              dbState,
		sch:             sch,
	}
	go botState.InitUpdates()
	botState.Worker()

	return nil
}

func (b *bot) InitUpdates() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case getIDCommand:
					b.GetID(update)
				case getMessageIDCommand:
					b.GetMessageID(update)
				case getCurrentStatusCommand:
					b.GetCurrentStatus(update)
				case turnEmergencyCommand:
					b.TurnEmergency(update)
				case turnBotCommand:
					b.TurnBot(update)

				}

				continue
			}
		}
	}
}

func (b *bot) Worker() {
	for {
		time.Sleep(time.Second * 30)
		timeNow, err := servertime.GetKyivTimeNow()
		if err != nil {
			log.Error().Err(err).Send()
			continue
		}
		serverConfig := b.db.serverConfigDB().get()
		if serverConfig.IsBotOff {
			continue
		}

		lastSend := b.db.lastSendDB().get()
		totalSecsLastSend := timeNow.Unix() - lastSend.Unix()
		hoursLastSend, minutesLastSend := totalSecsLastSend/3600, (totalSecsLastSend%3600)/60

		lastScheduleSend := b.db.lastScheduleSendDB().get()
		isSoon, soonSchedule := b.sch.IsScheduleSoon()
		online := checker.Online()

		var nowStatus status
		if online {
			nowStatus = onlineStatus
		} else {
			nowStatus = offlineStatus
		}

		currentStatus := status(b.db.statusDB().get())
		if !currentStatus.Validate() && online {
			err := b.db.statusDB().set(onlineStatus.ToString())
			if err != nil {
				log.Error().Err(err).Send()
				continue
			}
			currentStatus = onlineStatus
		} else if !currentStatus.Validate() && !online {
			err := b.db.statusDB().set(offlineStatus.ToString())
			if err != nil {
				log.Error().Err(err).Send()
				continue
			}
			currentStatus = offlineStatus
		}

		if b.updateMessageID > 0 {
			if !nowStatus.Validate() {
				continue
			}

			var timeString string
			if currentStatus != nowStatus {
				timeString = getTimeString(0, 0)
			} else {
				timeString = getTimeString(int(hoursLastSend), int(minutesLastSend))
			}

			var text string
			if nowStatus == onlineStatus {
				text = fmt.Sprintf(
					"🟢 Електроенергія присутня\n💡 Світло вжє є %s\n\n⏳ Останнє оновлення: %s",
					timeString,
					timeNow.Format("2006-01-02 15:04:05"),
				)
			} else {
				text = fmt.Sprintf(
					"🔴 Електроенергія відсутня\n💔 Світла немає вже %s\n\n⏳ Останнє оновлення: %s",
					timeString,
					timeNow.Format("2006-01-02 15:04:05"),
				)
			}

			editMsg := tgbotapi.NewEditMessageText(b.channelID, b.updateMessageID, text)
			_, err = b.bot.Send(editMsg)
			if err != nil {
				log.Error().Err(err).Send()
				continue
			}
		}

		if currentStatus != nowStatus {
			if !nowStatus.Validate() {
				continue
			}

			var text string
			if nowStatus == onlineStatus {
				text = fmt.Sprintf(
					"💡 Електроенергію було увімкнено. Світла не було %s",
					getTimeString(int(hoursLastSend), int(minutesLastSend)),
				)
			} else {
				text = fmt.Sprintf(
					"💔 Електроенергію було вимкнено. Світло було присутнє %s",
					getTimeString(int(hoursLastSend), int(minutesLastSend)),
				)
			}

			msg := tgbotapi.NewMessage(b.channelID, text)
			if isLate() {
				msg.DisableNotification = true
			}
			_, err := b.bot.Send(msg)
			if err != nil {
				log.Error().Err(err).Send()
				continue
			}

			err = b.db.statusDB().set(nowStatus.ToString())
			if err != nil {
				log.Error().Err(err).Send()
				continue
			}
			err = b.db.lastSendDB().set(timeNow)
			if err != nil {
				log.Error().Err(err).Send()
				continue
			}

			continue
		}

		if isSoon &&
			timeNow.Unix()-lastScheduleSend.Unix() > 2400 && // 40 minutes
			timeNow.Unix()-lastSend.Unix() > 900 && // 15 minutes
			nowStatus == onlineStatus &&
			!serverConfig.IsEmergency {
			msg := tgbotapi.NewMessage(
				b.channelID,
				fmt.Sprintf(
					"⚠️ За розкладом планових відключень о %s відбудеться відключення електроенергії",
					GetTimeForOffWhereHour(soonSchedule.Start),
				),
			)
			if isLate() {
				msg.DisableNotification = true
			}
			_, err := b.bot.Send(msg)
			if err != nil {
				log.Error().Err(err).Send()
				continue
			}

			err = b.db.lastScheduleSendDB().set(timeNow)
			if err != nil {
				log.Error().Err(err).Send()
				continue
			}
		}
	}
}

func isLate() bool {
	timeNow, err := servertime.GetKyivTimeNow()
	if err != nil {
		return false
	}

	if timeNow.Hour() >= 1 && timeNow.Hour() < 7 {
		return true
	}

	return false
}

func GetTimeForOffWhereHour(hour int) string {
	formatTime := fmt.Sprintf("%d:00", hour)
	if len([]rune(formatTime)) == 4 {
		formatTime = "0" + formatTime
	}

	return formatTime
}

func getTimeString(hours, minutes int) string {
	var minutesWord, hoursWord = declOfNum(minutes, minutesForms), declOfNum(hours, hoursForms)
	if hours <= 0 {
		return fmt.Sprintf("%d %s", minutes, minutesWord)
	}

	return fmt.Sprintf("%d %s і %d %s", hours, hoursWord, minutes, minutesWord)
}

func declOfNum(n int, textForms []string) string {
	if len(textForms) != 3 {
		return ""
	}

	n = int(math.Abs(float64(n))) % 100
	n1 := n % 10
	if n > 10 && n < 20 {
		return textForms[2]
	}
	if n1 > 1 && n1 < 5 {
		return textForms[1]
	}
	if n1 == 1 {
		return textForms[0]
	}

	return textForms[2]
}

func GetAdminID() (int64, error) {
	return strconv.ParseInt(os.Getenv("ADMIN_ID"), 10, 64)
}

func HasAdmin(userID int64) bool {
	adminID, err := GetAdminID()
	if err != nil {
		return false
	}

	return userID == adminID
}
