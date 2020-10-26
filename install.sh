sudo apt-get install -y curl;
curl -sL https://deb.nodesource.com/setup_12.x | sudo -E bash -;
sudo apt-get install -y nodejs;
sudo npm i -g pm2;
echo "const child = require('child_process')
const process = child.exec('./RoleBot')
process.stdout.on('data', (content) => {
  console.log(content)
})
process.stderr.on('data', (content) => {
  console.log(content)
})">>rolebot.js;
pm2 start rolebot.js --name=RoleBot;
pm2 save;
pm2 startup;