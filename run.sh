MASTERCONFIG=${GOPATH}/src/github.com/williammuji/shiran/master/masterconfig.json
SLAVECONFIG=${GOPATH}/src/github.com/williammuji/shiran/slave/slaveconfig.json
COMMANDERCONFIG=${GOPATH}/src/github.com/williammuji/shiran/commander/commanderconfig.json
CLIENTCONFIG=${GOPATH}/src/github.com/williammuji/shiran/client/client/clientconfig.json

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

echo ">>getHardware"
./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --getHardware --slaveName="FJSlave"
echo ">>getHardware --lshw"
./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --getHardware --slaveName="FJSlave" --lshw


echo ">>getFileContent noneexist"
./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --getFileContent --slaveName="FJSlave" --fileName="noneexist"
echo ">>getFileContent tmp"
./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --getFileContent --slaveName="FJSlave" --fileName="tmp"
echo ">>getFileContent tmp maxsize=10"
./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --getFileContent --slaveName="FJSlave" --fileName="tmp" --maxSize=10

echo ">>getFileChecksum"
./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --getFileChecksum --slaveName="FJSlave" --files="pingpongclient;pingpongserver"

echo ">>runCommand ls"
./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --runCommand --slaveName="FJSlave" --command="ls"
echo ">>runCommand /bin/ls"
./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --runCommand --slaveName="FJSlave" --command="/bin/ls"
echo ">>pingpongserver timeout=2"
./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --runCommand --slaveName="FJSlave" --command="pingpongserver" --timeout=2

echo -e "#!/bin/bash\necho 'i am script'" > testscript
./commander --configFile=${COMMANDERCONFIG} --log_dir="./" --runScript --slaveName="FJSlave" --script="testscript"



./client --configFile=${CLIENTCONFIG} --log_dir="./" &
