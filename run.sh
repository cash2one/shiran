MASTERCONFIG=~/goyard/src/github.com/williammuji/shiran/master/masterconfig.json
SLAVECONFIG=~/goyard/src/github.com/williammuji/shiran/slave/slaveconfig.json
COMMANDERCONFIG=~/goyard/src/github.com/williammuji/shiran/commander/commanderconfig.json
CLIENTCONFIG=~/goyard/src/github.com/williammuji/shiran/client/client/clientconfig.json

./master --configFile=${MASTERCONFIG} --log_dir="./" &
sleep 1
./slave --configFile=${SLAVECONFIG} --log_dir="./" &
sleep 1

./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --add --slaveName="FJSlave" --appName="loginA"
./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --add --slaveName="FJSlave" --appName="loginB"

./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --start --slaveName="FJSlave" --appName="loginA"
./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --start --slaveName="FJSlave" --appName="loginB"
sleep 1


./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --add --slaveName="FJSlave" --appName="gateA_1001"
./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --add --slaveName="FJSlave" --appName="gateB_1001"
./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --add --slaveName="FJSlave" --appName="gateA_1002"
./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --add --slaveName="FJSlave" --appName="gateB_1002"

./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --start --slaveName="FJSlave" --appName="gateA_1001"
./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --start --slaveName="FJSlave" --appName="gateB_1001"
./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --start --slaveName="FJSlave" --appName="gateA_1002"
./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --start --slaveName="FJSlave" --appName="gateB_1002"
sleep 3

./client --configFile=${CLIENTCONFIG} --log_dir="./" &
