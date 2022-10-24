#!/usr/bin/env bash
# @describe: build script

set -e

function usage() {
    echo "usage:    ./build.sh program work-env tag"
    echo "example:  ./build.sh api dev v1.0.0"
    echo "program: api|task|cmd"
    echo "work-env: test|prod"
    exit 0
}

workspace=$(cd $(dirname $0) && pwd)
cd ${workspace}

if [[ $# -lt 3 ]]; then
    usage
fi

binary_program=$1
work_env=$2
tag=$3

image="tinypro-api"
if [[ ${binary_program} == "task" ]]; then
    image="tinypro-task"
elif [[ ${binary_program} == "cmd" ]]; then
    image="tinypro-cmd"
fi


if [[ "$work_env" != "test" ]] && [[ "$work_env" != "prod" ]];then
    echo "work_env is invalid"
    echo ""
    usage
fi

rm -rf build
mkdir -p build/data build/logs

git checkout . && git pull origin tinypro
if [[ $? != 0 ]]
then
    echo "[`date +"%Y-%m-%d %H:%M:%S"`] update failed, please check"
    exit 110
fi

echo "[`date +"%Y-%m-%d %H:%M:%S"`] git pull finished"
git_hash=`git rev-parse HEAD`
echo "[`date +"%Y-%m-%d %H:%M:%S"`] git latest hash: $git_hash"
echo "$tag $git_hash" > "api/conf/git-rev-hash"
cp -Rfp "api/conf"  "build/"
cp /opt/data/GeoLite2-City.mmdb build/data/

dk_user="user"
dk_pwd="pass"
dk_repo="registry.cn-beijing.aliyuncs.com"
registry="$dk_repo/tinypro-test/$image"

if [[ "$work_env" == "prod" ]];then
    # transfer configuration file
    cp -fp build/conf/app.prod.conf build/conf/app.conf

    # prod docker
    dk_user="user"
    dk_pwd="pass"
    dk_repo="registry.cn-beijing.aliyuncs.com"
    registry="$dk_repo/tinypro-prod/$image"
else
    #sed -i "" 's#mysql.db.rds#'$DEV_HOST_IP'#g' build/conf/app.conf
    #sed -i "" 's#cache.redis#'$DEV_HOST_IP'#g' build/conf/app.conf
    #sed -i "" 's#storage.redis#'$DEV_HOST_IP'#g' build/conf/app.conf
    #cp -fp build/conf/app.conf build/conf/app.conf
    #sed -i "" 's#runmode = "dev"#runmode = "prod"#g' build/conf/app.conf
    echo "this is dev"
fi

export GO111MODULE=on

distributeOs=`uname  -a`

b="Darwin"
c="centos"
d="ubuntu"

currentOs="mac"
if [[ $distributeOs =~ $b ]];then
    currentOs="mac"
else
    currentOs="linux"
fi

if [[ ${binary_program} == "api" ]]; then
    # h5 files
    cp -Rfp api/views build/views

    if [[ $currentOs == "mac" ]]; then
      GOOS=linux GOARCH=amd64 go build -o build/api api/main.go
    else
      # fix standard_init_linux.go:228: exec user process caused: no such file or direct
      CGO_ENABLED=0 go build -o build/api api/main.go
    fi
elif [[ ${binary_program} == "task" ]]; then
    if [[ $currentOs == "mac" ]]; then
      GOOS=linux GOARCH=amd64 go build -o build/task api/cmd/task/main.go
    else
      CGO_ENABLED=0 go build -o build/task api/cmd/task/main.go
    fi
else
    if [[ $currentOs == "mac" ]]; then
      GOOS=linux GOARCH=amd64 go build -o build/task api/cmd/main.go
    else
      CGO_ENABLED=0 go build -o build/task api/cmd/main.go
    fi
fi

if [[ $? != 0 ]]
then
    echo "[`date +"%Y-%m-%d %H:%M:%S"`] compile failed,please check"
    exit 1
else
    echo "[`date +"%Y-%m-%d %H:%M:%S"`] compile finished!"
fi

if [[ ${binary_program} == "api" ]]; then
    docker build -f Dockerfile -t ${image}:${tag} .
else
    docker build -f DockerTaskFile -t ${image}:${tag} .
fi

echo "[`date +"%Y-%m-%d %H:%M:%S"`] docker build finished!"

image_id=`docker images | grep ${image} | grep ${tag} | awk '{print $3}'`
echo ${image_id}

# 登陆
docker login -u ${dk_user} -p ${dk_pwd} ${dk_repo}

img_name=${registry}:${tag}
echo "docker name: $img_name"

# prod docker tag, push to cloud
if [[ "$work_env" == "prod" ]];then
  docker tag ${image_id} ${img_name}
  docker push ${img_name}
fi

echo "[`date +"%Y-%m-%d %H:%M:%S"`] release successfully"
