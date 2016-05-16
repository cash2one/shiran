MASTERCONFIG=${GOPATH}/src/github.com/williammuji/shiran/master/masterconfig.json
SLAVECONFIG=${GOPATH}/src/github.com/williammuji/shiran/master/slaveconfig.json
COMMANDERCONFIG=${GOPATH}/src/github.com/williammuji/shiran/master/masterconfig_pingpong.json

./master --configFile=${MASTERCONFIG} --log_dir="./" &
sleep 1
./slave --configFile=${SLAVECONFIG} --log_dir="./" &
sleep 1

./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --add --slaveName="FJSlave" --appName="ppserver"
./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --add --slaveName="FJSlave" --appName="ppclient"

./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --start --slaveName="FJSlave" --appName="ppserver"
sleep 1
./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --start --slaveName="FJSlave" --appName="ppclient"

sleep 50
./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --list --slaveName="FJSlave"
./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --get --slaveName="FJSlave" --appName="ppserver"
./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --stop --slaveName="FJSlave" --appName="ppserver"
./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --get --slaveName="FJSlave" --appName="ppserver"
./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --list --slaveName="FJSlave"

sleep 1
./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --start --slaveName="FJSlave" --appName="ppserver"
sleep 1
./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --start --slaveName="FJSlave" --appName="ppclient"
sleep 50
./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --stop --slaveName="FJSlave" --appName="ppserver"
./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --get --slaveName="FJSlave" --appName="ppserver;ppclient"

sleep 1
./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --remove --slaveName="FJSlave" --appName="ppclient"
./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --list --slaveName="FJSlave"
./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --remove --slaveName="FJSlave" --appName="ppserver"
./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --list --slaveName="FJSlave"
