Color_Off='\033[0m'
Black='\033[0;30m'
Red='\033[0;31m'
Green='\033[0;32m'
Yellow='\033[0;33m'
Blue='\033[0;34m'
Purple='\033[0;35m'
Cyan='\033[0;36m'
White='\033[0;37m'

BASEDIR=$(dirname $0)

rm -rf ${BASEDIR}/webServer/App/front &>/dev/null
mkdir -p ${BASEDIR}/webServer/App/front &>/dev/null
export FRONT_DOMAIN_NAME=http://localhost:8080
if [ "$1" = "build" ]; then
	echo ${Red}Please don\'t forget to update vuejs prod env, \'API_DOMAIN_NAME\' to \'localhost:4001\'${Color_Off}
	export FRONT_DOMAIN_NAME=http://localhost:4001
	echo ${Purple}Install vuejs dependency...${Color_Off}
	npm i  --prefix ${BASEDIR}/front/ &>/dev/null

	echo ${Purple}Build single page application...${Color_Off}
	npm run build  --prefix ${BASEDIR}/front/ &>/dev/null

	echo ${Blue}Copy app to public folder...${Color_Off}
	mv ${BASEDIR}/front/dist/* ${BASEDIR}/webServer/App/front &>/dev/null
fi
echo ${Cyan}Build docker-compose...${Color_Off}
rm -rf ./mongodbData/.* mongodbData/*/.* &>/dev/null
docker-compose build &>/dev/null
echo ${Green}Launch containers...${Color_Off}
docker-compose up

#echo ${Purple}Build docker-compose...${Color_Off}
#docker-compose create &>/dev/null
# echo ${Purple}Lauch and build all containers in detach mode...${Color_Off}
# if [ "$(whoami)" != "hypertube" ]; then
# 	docker-compose -f ${BASEDIR}/docker-compose.yml up -d
# else
# 	echo "Running on prod launch docker-compose with sudo"
# 	sudo docker-compose -f ${BASEDIR}/docker-compose.yml up -d
# fi
