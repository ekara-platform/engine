echo "ENGINE copydep..."

# Refresh "model"
rm -rf ./vendor/github.com/ekara-platform/model/*.go
cp ../model/*.go  ./vendor/github.com/ekara-platform/model/

