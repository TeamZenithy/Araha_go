sudo apt-get install -y curl;
curl -sL https://deb.nodesource.com/setup_12.x | sudo -E bash -;
sudo apt-get install -y nodejs;
sudo npm i -g pm2;
echo "const child = require('child_process')
const process = child.exec('./Araha')
process.stdout.on('data', (content) => {
  console.log(content)
})
process.stderr.on('data', (content) => {
  console.log(content)
})">>araha.js;
pm2 start araha.js --name=Araha;
pm2 save;
pm2 startup;
