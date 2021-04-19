// setup
const express = require('express');
const path = require('path');
const app = express();
const axios = require('axios');
const cheerio = require('cheerio');
const Jimp = require('jimp');
const imgur = require('imgur');

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

//
app.post('/api', async (request, response) => {
  let id = request.body.str;
  console.log('Got a request: ' + id);

  const url = 'http://javid.ddns.net/tModLoader/tools/ranksbysteamid.php?steamid64=' + id;

  // scrape data
  const mods = await scrapeData(url)

  console.log(mods);

  // generate image using JIMP
  await generateImage(mods);

  // generate and get imgur link
  const link = await uploadImage();

  response.redirect(link);
  // send data to frontend
  //response.status(200).send(mods);
});

// returns a json array
async function scrapeData(url) {
  let mods = [];

  await axios(url)
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
    }).catch(console.error);

  return mods;
}

async function generateImage(mods) {
  Jimp.loadFont("Andy32.fnt")
    .then(font => Jimp.read('card.png')
    .then(card => {
      card = card.print(font, 0, 0, mods[0].DisplayName).write('output.png');
      return card;
    })
    .catch(err => console.error(err)))
    .catch(err => console.error(err));
}

// returns a imgur link to the image
async function uploadImage() {
  const json = await imgur.uploadFile('output.png');
  const link = json.link;

  console.log(link);
  return link;
}

//stuff to do on exit
process.on('exit', function () {
  console.log('About to close');
});