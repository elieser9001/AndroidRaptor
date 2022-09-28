package whatsapp

import (
	"strings"

	"github.com/elieser9001/AndroidRaptor/android"
)

const (
	WHATSAPP_MEDIA_PATH       = "/sdcard/WhatsApp/Media"
	WHATSAPP_IMAGES_PATH      = WHATSAPP_MEDIA_PATH + "/WhatsApp Images/"
	WHATSAPP_IMAGES_SENT_PATH = WHATSAPP_MEDIA_PATH + "/WhatsApp Images/Sent/"
	WHATSAPP_VIDEO_PATH       = WHATSAPP_MEDIA_PATH + "/WhatsApp Video/"
	WHATSAPP_STATUSES_PATH    = WHATSAPP_MEDIA_PATH + "/.Statuses/"
	WHATSAPP_VIDEO_SENT_PATH  = WHATSAPP_MEDIA_PATH + "/WhatsApp Video/Sent/"
	WHATSAPP_VOICE_NOTES_PATH = WHATSAPP_MEDIA_PATH + "/WhatsApp Voice Notes/"
)

func GetWhatsAppAllMediaFiles(fLogPath string) error {
	path := strings.ReplaceAll(WHATSAPP_MEDIA_PATH, " ", "\\ ")
	cmd := "ls -laR " + path

	_, err := android.ExecuteCommand(cmd + " > " + fLogPath)
	if err != nil {
		return err
	}

	return nil
}

func getWhatsAppFiles(wsFilesPath string, fLogPath string) error {
	path := strings.ReplaceAll(wsFilesPath, " ", "\\ ")
	cmd := "find " + path + " '*/\\.*' | sort > " + fLogPath

	_, err := android.ExecuteCommand(cmd + " > " + fLogPath)
	if err != nil {
		return err
	}

	return nil
}

func GetWhatsAppImages(fLogPath string) error {
	err := getWhatsAppFiles(WHATSAPP_IMAGES_PATH, fLogPath)
	return err
}

func GetWhatsAppImgSent(fLogPath string) error {
	err := getWhatsAppFiles(WHATSAPP_IMAGES_SENT_PATH, fLogPath)
	return err
}

func GetWhatsAppVideoSent(fLogPath string) error {
	err := getWhatsAppFiles(WHATSAPP_VIDEO_SENT_PATH, fLogPath)
	return err
}

func GetWhatsAppVideos(fLogPath string) error {
	err := getWhatsAppFiles(WHATSAPP_VIDEO_PATH, fLogPath)
	return err
}

func GetWhatsAppStatuses(fLogPath string) error {
	err := getWhatsAppFiles(WHATSAPP_STATUSES_PATH, fLogPath)
	return err
}

func GetWhatsappVoiceNotes(fLogPath string) error {
	err := getWhatsAppFiles(WHATSAPP_VOICE_NOTES_PATH, fLogPath)
	return err
}
