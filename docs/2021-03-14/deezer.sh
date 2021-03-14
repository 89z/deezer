function gw(){ curl -s "https://www.deezer.com/ajax/gw-light.php?method=$1&input=3&api_version=1.0&api_token=$3" -H "cookie: sid=$2" $@;}

TRK_IDS=75498415
DZR_API="4V...."
DZR_URL="www.deezer.com/ajax/gw-light.php?method=deezer.ping&api_version=1.0&api_token"
DZR_SID=$(curl -s "$DZR_URL" | jq -r .results.SESSION)
USR_NFO=$(gw deezer.getUserData $DZR_SID)
USR_TOK=$(jq -r .results.USER_TOKEN <<< $USR_NFO)
USR_LIC=$(jq -r .results.USER.OPTIONS.license_token <<< $USR_NFO)
API_TOK=$(jq -r .results.checkForm <<< $USR_NFO)
TRK_NFO=$(gw song.getListData $DZR_SID $API_TOK --data-binary '{"sng_ids":['"$TRK_IDS"']}')
TRK_TOK=$(jq -r .results.data[].TRACK_TOKEN <<< "$TRK_NFO")

curl 'https://media.deezer.com/v1/get_url' --data-binary '{"license_token":"'$USR_LIC'","media":[{"type":"FULL","formats":[{"cipher":"BF_CBC_STRIPE","format":"MP3_128"}]}],"track_tokens":["'$TRK_TOK'"]}'
