document.addEventListener("DOMContentLoaded", async function () {
    let steamid = new URLSearchParams(window.location.search).get('steamid64');

    if (steamid != null) {
        getData(steamid);
    }
});

async function getData(steamid64) {
    // send mod name to back-end
    let str = steamid64
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
        document.body.style = "margin: 0px; background: #0e0e0e; height: 100%";
        const html = `<img -webkit-user-select: none;margin: auto;background-color: hsl(0, 0%, 90%);transition: background-color 300ms; src="${info.url}"></img>`;
        document.body.innerHTML = html;
    }
}
    //const html = `<h1>${JSON.stringify(info)}</h1>`;
    //document.body.innerHTML = html;