document.addEventListener("DOMContentLoaded", async function () {
    let steamid = new URLSearchParams(window.location.search).get('steamid64');

    if (steamid != null) {
        getData(steamid);
    }
});

async function getData(modName) {
    // send mod name to back-end
    let str = modName // what the user entered into the text field
    const data = { str }
    const options = {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(data)
    };
  
    var response = await fetch('/api', options);
    let info = await response.json();
    console.log(info);
  
    if (response.status == 200) {
        const html = `<img src="./output.png"></img>`
        document.body.innerHTML = html;
    }
}
    //const html = `<h1>${JSON.stringify(info)}</h1>`;
    //document.body.innerHTML = html;