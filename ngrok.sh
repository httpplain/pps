rm -rf ngrok ngrok.zip ngrok.sh > /dev/null 2>&1

echo "======================="

echo "Download ngrok"

echo "======================="

wget -O ngrok.zip https://bin.equinox.io/c/4VmDzA7iaHb/ngrok-stable-windows-amd64.zip > /dev/null 2>&1

unzip ngrok.zip > /dev/null 2>&1

read -p "Paste You Ngrok Authtoken: " CRP

./ngrok authtoken $CRP 
