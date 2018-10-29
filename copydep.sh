echo "ENGINE copydep..."

# Refresh "model"
rm -rf ./vendor/github.com/ekara-platform/model/*.go
mkdir ./vendor/github.com/ekara-platform/model/
cp ../model/*.go  ./vendor/github.com/ekara-platform/model/

