docker build --build-arg APP=cloud --build-arg TYPE=local --build-arg PORT=8080 --build-arg HOST=0.0.0.0 -t fehmi-cloud-connector .
docker create --name temp-extractor fehmi-cloud-connector
docker cp temp-extractor:/output/fehmi-connector-windows.exe ./app/fehmi-connector-windows.exe
docker cp temp-extractor:/output/fehmi-connector-linux ./app/fehmi-connector-linux
docker rm -f temp-extractor

docker build -t windows-connector -f Dockerfile.app .
