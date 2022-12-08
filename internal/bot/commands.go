package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

const (
	getIDCommand            = "get_id"
	getMessageIDCommand     = "get_message_id"
	getCurrentStatusCommand = "curr_status"

	turnEmergencyCommand = "turn_emergency"
	turnBotCommand       = "turn_bot"
)

func (b *bot) GetID(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprintf("ID: %d", update.FromChat().ID))
	_, err := b.bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}
}

func (b *bot) GetMessageID(update tgbotapi.Update) {
	if update.Message.ReplyToMessage == nil {
		return
	}

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		fmt.Sprintf("Message ID: %d", update.Message.ReplyToMessage.MessageID),
	)
	_, err := b.bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}
}

func (b *bot) GetCurrentStatus(update tgbotapi.Update) {
	currStatus := b.db.statusDB().get()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Curr Status: %s", currStatus))
	_, err := b.bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}
}

func (b *bot) TurnEmergency(update tgbotapi.Update) {
	if !HasAdmin(update.SentFrom().ID) {
		return
	}

	currConfig := b.db.serverConfigDB().get()

	var msg tgbotapi.MessageConfig
	if currConfig.IsEmergency {
		msg = tgbotapi.NewMessage(update.SentFrom().ID, "Emergency turn off")
		currConfig.IsEmergency = !currConfig.IsEmergency
	} else {
		msg = tgbotapi.NewMessage(update.SentFrom().ID, "Emergency turn on")
		currConfig.IsEmergency = !currConfig.IsEmergency
	}

	_, err := b.bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	err = b.db.serverConfigDB().set(currConfig)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}
}

func (b *bot) TurnBot(update tgbotapi.Update) {
	if !HasAdmin(update.SentFrom().ID) {
		return
	}

	currConfig := b.db.serverConfigDB().get()

	var msg tgbotapi.MessageConfig
	if currConfig.IsBotOff {
		msg = tgbotapi.NewMessage(update.SentFrom().ID, "Bot turn on")
		currConfig.IsBotOff = !currConfig.IsBotOff
	} else {
		msg = tgbotapi.NewMessage(update.SentFrom().ID, "Bot turn off")
		currConfig.IsBotOff = !currConfig.IsBotOff
	}

	_, err := b.bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	err = b.db.serverConfigDB().set(currConfig)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}
}
