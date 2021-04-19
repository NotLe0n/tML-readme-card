// setup
const express = require('express');
const path = require('path');
const app = express();
const axios = require('axios');
const cheerio = require('cheerio');
var Jimp = require('jimp');

// fixes file paths
app.use(express.static(__dirname));
//enables json use
app.use(express.json());

//get request for the homepage
app.get('/', async (request, response) => {
  await response.sendFile(path.join(process.cwd(), 'index.html'));
});

// listening for requests
app.listen(3000, () => {
  console.log('server started');
});

// get data from database and send it to front-end
app.post('/api', async (request, response) => {
  let id = request.body.str;
  console.log('Got a request: ' + id);

  const url = 'http://javid.ddns.net/tModLoader/tools/ranksbysteamid.php?steamid64=' + id;
  let mods = [];

  axios(url)
      .then(site => {
        const html = site.data;
        const $ = cheerio.load(html);

        // find the *first* table
        const table = $('.primary')[0];
        // firstChild is the <tbody>, get its children
        const rows = table.firstChild.children;

        let RankTotal;
        let DisplayName;
        let DownloadsTotal;
        let DownloadsYesterday;

        // go trough all rows and grab the data
        for (let i = 1; i < rows.length; i++) {
          RankTotal = rows[i].children[0].children[0].data;
          DisplayName = rows[i].children[1].children[0].data;
          DownloadsTotal = rows[i].children[2].children[0].data;
          DownloadsYesterday = rows[i].children[3].children[0].data;
        
          // generate json
          mods.push({
            "DisplayName": DisplayName,
            "RankTotal": RankTotal,
            "DownloadsTotal": DownloadsTotal,
            "DownloadsYesterday": DownloadsYesterday
          });
        }
        console.log(mods);

        generateImage(JSON.stringify(mods));
      }).catch(console.error);

      // send data to frontend
      response.status(200).send(JSON.stringify(mods));
});

async function generateImage(mod) {
  Jimp.loadFont("Andy32.fnt")
    .then(font => Jimp.read('card.png')
    .then(card => {
      return card
        .print(font, 0, 0, "amongus")
        .write('output.png'); // save
    })
    .catch(err => console.error(err)))
    .catch(err => console.error(err));
  
}

//stuff to do on exit
process.on('exit', function () {
  console.log('About to close');
});