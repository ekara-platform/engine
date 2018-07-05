# Refresh "model"
rm -rf ./vendor/github.com/lagoon-platform/model/*.go
cp ../model/*.go  ./vendor/github.com/lagoon-platform/model/
rm -rf ./vendor/github.com/lagoon-platform/model/params/*.go
cp ../model/params/*.go  ./vendor/github.com/lagoon-platform/model/params/

