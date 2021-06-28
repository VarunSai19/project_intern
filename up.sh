rm -rf project_intern
nvm install v10.15.3 
sleep 3
docker stop $(docker ps -aq) && docker rm $(docker ps -aq)
sleep 3
chmod a+w project_intern
cd project_intern
cd artifacts/channel/create-certificate-with-ca
docker-compose up -d
sleep 3
./create-certificate-with-ca.sh 
cd ..
./create-artifacts.sh 
cd ..
docker-compose up -d
sleep 3
cd ..
./createChannel.sh 
sleep 3
./deployChaincode.sh 
sleep 3
cd api-2.0/config/
./generate-ccp.sh 
sleep 3
cd ..
sleep 3
npm install --unsafe-perm
